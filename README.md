# annotation-operator
Demonstrates reconciliation on the top of ingress only 
- No CRD
- No custom type

## Getting Started
- https://book.kubebuilder.io/quick-start.html

```shell
 kubebuilder init --domain cloud.example.com --repo cloud.example.com/annotation-operator
 kubebuilder create api --group annov1 --version v1 --kind Anno
 make generate
 make deploy
```
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

# annotation-operator
