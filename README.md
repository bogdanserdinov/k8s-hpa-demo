# k8s-hpa-demo

This project consists of several Go services that communicates with each other using gRPC.
Project contains k8s manifests that can be used to deploy the services to a Kubernetes cluster along with HPA.

How to deploy the application to the Kubernetes cluster:

1. Create metrics server
```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

2. Update metrics policy to use insecure tls:
```bash
kubectl patch deployment metrics-server -n kube-system \
  --type='json' \
  -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'
```

3. Deploy the services
```bash
mage -v k8s:apply
```

4. Expose the gateway service (testing purpose only)
```bash
kubectl port-forward -n gateway svc/gateway 8083:80
```

5. Load testing
```bash
mage -v LoadTest http://localhost:8083 500 100
```

6. Observe the HPA or list the pods
```bash
mage -v k8s:getHpa <service-name>
mage -v k8s:getPods <service-name>
```

7. Rollback the deployment
```bash
kubectl delete -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

mage -v k8s:delete
```

The deployment manifests for all the services are managed within the `deploy/k8s/deployment` folder
