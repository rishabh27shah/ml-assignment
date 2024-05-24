# MailerSend

Email Operetor

## Description

This is an email operetor to send emails and manage email sender configuration.

Written in GO Lang

![Alt text](./Images/Arch.png?raw=true "Architecture")

## Getting Started

### Prerequisites

- Go (v1.21.0+)
- Docker (17.03+)
- Operator-sdk (v1.4.0+)
- Kubectl (v1.11.3+)
- Kubernetes cluster (v1.11.3+)

**Install the CRDs into the cluster:**

```sh
make install
```

**Check CRDs:**

```sh
kubectl get crd
```

### To Deploy on the cluster

**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/mailersend:tag
```

**Default image name:** controller:latest

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/mailersend:tag
```

**To deploy the CRD with image:**
We have the file for [deployment file](./CRD-deployment.yaml) in which we can update the configuration and image name.

```sh
kubectl apply -f CRD-deployment.yaml
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

**To run in Development mode:**
Using this it will genrate manifest files, formating, installting and running CRD in development mode and we can see the logs.

```sh
make run
```

### To Uninstall

**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

Please check the screenshots for reference [here](./Images/).

## Contributing

**Rishabh Shah**

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

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
