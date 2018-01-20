# poc-argo-fromcode

PoC showing how we can run an [Argo](https://github.com/argoproj/argo) workflow from code.

## Setup

```
# Setup Kubernetes
minikube start --vm-driver=xhyve
eval $(minikube docker-env)

# Argo
kubectl create namespace poc-argo-fromcode
argo install --install-namespace poc-argo-fromcode

# Dependencies
dep ensure
```

There is also an Argo UI, which you can access at http://127.0.0.1:8080 after port forwarding:

```
kubectl get --all-namespaces services | grep argo-ui
kubectl get --all-namespaces pods | grep argo-ui
kubectl port-forward -n kube-system argo-ui-77f8ff9588-qcgp5 8080:8001
```

## Usage

Normally you would run the argo workflow with:

```
argo submit --namespace poc-argo-fromcode example.yaml
```

To run this example from code (and the manually hardcoded equivalent of `example.yaml`):

```
go run main.go
```

Then you can check the status/logs/etc in the normal Argo ways:

```
argo list --namespace poc-argo-fromcode
argo get --namespace poc-argo-fromcode poc-argo-fromcode-7p9sp
argo logs --namespace poc-argo-fromcode poc-argo-fromcode-7p9sp
argo delete --namespace poc-argo-fromcode poc-argo-fromcode-7p9sp 
```

These could also be implemented in code, but that goes beyond this initial PoC.
