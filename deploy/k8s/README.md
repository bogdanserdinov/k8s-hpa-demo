This folder contains deployments and services for both the gateway and addition services.

### Notes:

1) Each service is deployed in its own independent namespace.
2) We use Deployments instead of ReplicaSets or Pod definitions directly.

### How to Apply All K8s Changes
To apply changes to your Kubernetes cluster, use the following commands in the order specified to ensure that dependencies are applied correctly:

```bash
kubectl apply -f namespace.yaml
kubectl apply -f config.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

### How to Rollback Changes
To remove the resources you applied, use the following commands:

```bash
kubectl delete deployment <deployment-name> -n <namespace>
kubectl delete service <service-name> -n <namespace>
kubectl delete configmap <configmap-name> -n <namespace>
kubectl delete namespace <namespace>
```
