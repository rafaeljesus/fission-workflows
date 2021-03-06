# Installation

This document covers the installation of Fission Workflows.

### Prerequisites

Fission Workflows requires the following to be installed on your host machine:

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [helm](https://github.com/kubernetes/helm)

Additionally, Fission Workflows requires a [Fission](https://github.com/fission/fission) 
deployment on your Kubernetes cluster. If you do not have a Fission deployment, follow
[Fission's installation guide](http://fission.io/docs/0.4.0/install/).

**(Note that Fission Workflows requires Fission 0.4.1 or higher, with the NATS component installed!)**

### Installing Fission Workflows

Fission Workflows is an add-on to Fission. You can install both
Fission and Fission Workflows using helm charts.

Assuming you have a Kubernetes cluster, run the following commands:

```bash
# Add the Fission charts repo
helm repo add fission-charts https://fission.github.io/fission-charts/
helm repo update

# Install Fission 
# This assumes that you do not have a Fission deployment yet, and are installing on a standard Minikube deployment.
# Otherwise see http://fission.io/docs/0.4.0/install/ for more detailed instructions
helm install --wait -n fission-all --namespace fission --set serviceType=NodePort --set analytics=false fission-charts/fission-all --version 0.4.1

# Install Fission Workflows
helm install --wait -n fission-workflows fission-charts/fission-workflows --version 0.2.0
```

### Creating your first workflow

After installing Fission and Workflows, you're ready to run a simple
test workflow.  Clone this repository, and from its root directory, run:

```bash
#
# Add binary environment and create two test functions on your Fission setup:
#
fission env create --name binary --image fission/binary-env
fission function create --name whalesay --env binary --deploy examples/whales/whalesay.sh
fission function create --name fortune --env binary --deploy examples/whales/fortune.sh

#
# Create a workflow that uses those two functions. A workflow is just
# a function that uses the special "workflow" environment.
#
fission function create --name fortunewhale --env workflow --src examples/whales/fortunewhale.wf.yaml

#
# Map an HTTP GET to your new workflow function:
#
fission route create --method GET --url /fortunewhale --function fortunewhale

#
# Invoke the workflow with an HTTP request:
#
curl ${FISSION_ROUTER}/fortunewhale
``` 
