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
To remove the manifests you applied, use the following commands:

```bash
kubectl delete deployment <deployment-name> -n <namespace>
kubectl delete service <service-name> -n <namespace>
kubectl delete configmap <configmap-name> -n <namespace>
kubectl delete namespace <namespace>
```

### How to expose the gateway service (testing purpose only)
```bash
kubectl port-forward -n gateway svc/gateway 8081:80
```

### Metrics server

Apply:

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

Delete:

```bash
kubectl delete -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

Patch metrics-server to use insecure tls:

```bash
kubectl patch deployment metrics-server -n kube-system \
  --type='json' \
  -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'
```
