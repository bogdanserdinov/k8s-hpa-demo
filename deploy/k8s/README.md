This folder contains k8s manifests for both the gateway and addition services.

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
