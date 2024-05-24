/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	"github.com/mailersend/mailersend-go"

	tutorialv1 "my.domain/tutorial/api/v1"
)

// EmailReconciler reconciles a Email object
type EmailReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tutorial.my.domain,resources=emails,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tutorial.my.domain,resources=emails/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tutorial.my.domain,resources=emails/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Email object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile

func (r *EmailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("email", req.NamespacedName)
	log.Info("Reconciling Email")

	// Fetch the Email instance
	var email tutorialv1.Email
	if err := r.Get(ctx, req.NamespacedName, &email); err != nil {
		log.Error(err, "unable to fetch Email")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("Fetched Email instance", "email", email)

	// Check if the email has already been sent
	if email.Status.DeliveryStatus == "Sent" || email.Status.DeliveryStatus == "Failed" {
		log.Info("Email already processed", "deliveryStatus", email.Status.DeliveryStatus)
		return ctrl.Result{}, nil
	}

	// Fetch the referenced EmailSenderConfig
	var senderConfig tutorialv1.EmailSenderConfig
	senderConfigRef := email.Spec.SenderConfigRef
	if err := r.Get(ctx, client.ObjectKey{Name: senderConfigRef, Namespace: req.Namespace}, &senderConfig); err != nil {
		log.Error(err, "unable to fetch EmailSenderConfig")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("Fetched EmailSenderConfig instance", "emailSenderConfig", senderConfig)

	// Initialize MailerSend client
	apiToken := senderConfig.Spec.ApiToken
	ms := mailersend.NewMailersend(apiToken)

	// Prepare email details
	subject := email.Spec.Subject
	text := email.Spec.Body
	html := "<p>" + email.Spec.Body + "</p>"

	from := mailersend.From{
		Name:  senderConfig.Spec.SenderName,
		Email: senderConfig.Spec.SenderEmail,
	}

	recipients := []mailersend.Recipient{
		{
			Name:  "Recipient",
			Email: email.Spec.RecipientEmail,
		},
	}

	// Create and send the email
	message := ms.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)
	message.SetText(text)

	log.Info("Sending email...")
	res, err := ms.Email.Send(ctx, message)
	if err != nil {
		// Check for rate limiting (HTTP 429)
		if res != nil && res.StatusCode == 429 {
			log.Error(err, "Rate limit exceeded, retrying...")
			// Retry with exponential backoff
			retryAfter := time.Duration(1) * time.Minute
			return ctrl.Result{RequeueAfter: retryAfter}, nil
		}

		log.Error(err, "Failed to send email")
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = err.Error()
		if updateErr := r.Status().Update(ctx, &email); updateErr != nil {
			log.Error(updateErr, "Failed to update email status")
			return ctrl.Result{}, updateErr
		}
		return ctrl.Result{}, err
	}

	log.Info("Email sent successfully", "X-Message-Id", res.Header.Get("X-Message-Id"))

	// Update status
	email.Status.DeliveryStatus = "Sent"
	email.Status.MessageId = res.Header.Get("X-Message-Id")
	if err := r.Status().Update(ctx, &email); err != nil {
		log.Error(err, "Failed to update email status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EmailReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tutorialv1.Email{}).
		Complete(r)
}
