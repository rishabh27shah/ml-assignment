# Objective

Create a Kubernetes operator that manages custom resources for configuring email sending and sending of emails via a transactional email provider like MailerSend. The operator should work cross namespace, and demonstrate sending from multiple providers.

## Specifications

### Custom Resource Definitions (CRDs)

#### EmailSenderConfig

- Define a CRD named `EmailSenderConfig`.
- The CRD should have the following specifications:

```yaml
spec:
  apiTokenSecretRef: string
  senderEmail: string
```

The token should be retrieved from a secret object.

#### Email

- Define a CRD named `Email`.
- The CRD should have the following specifications:

```yaml
spec:
  senderConfigRef: string # Reference to EmailSenderConfig
  recipientEmail: string
  subject: string
  body: string
```

```yaml
status:
  deliveryStatus: string
  messageId: string
  error: string
```

### Controller Deployment

Create a deployment for the controller with the following specifications:

```yaml
name: email-operator
namespace: default
service_account: email-operator
resources: 0.8 CPU, 256Mi memory
replicas: 1
```

### Operator Implementation

Implement the operator using the Go programming language.
The operator should watch for changes to resources of kind `Email` and `EmailSenderConfig`.

When a new `EmailSenderConfig` is created or updated, the operator should confirm the email sending settings, you should be able to create multiple configs, the operator should log these actions (Created / Updated)

When a new `Email` is created, the operator should send the email using the provided sender configuration and MailerSend API.

- Use the provided MailerSend API token from the referenced EmailSenderConfig to authenticate and send emails.
- Configure the sender email, recipient email, subject, and body according to the specifications provided in the `Email` CRD.
- Update the status of the `Email` resource to reflect the delivery status and message ID.
- If the request fails, update the status to failed and provide a error response.

### Testing

Test the operator by creating instances of `EmailSenderConfig` and `Email` with different configurations.

Verify that the operator properly sends emails via MailerSend, and another provider.
Check that the status of the `Email` resource reflects the delivery status and message ID.
Check that email errors are recorded

### Deliverables

- Source code for the operator implementation, on Github.
- Container and Deployment manifests for the operator.
- Documentation on deploying and using the operator.

### Bonus items

Ensure that the operator securely handles sensitive information like API tokens (secret refs)
GitOps deployment
