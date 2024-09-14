This folder contains deployments and services for both gateway and addition service.

Notes:

1) Every service has independent namespace.
2) We use deployments instead of replicaset or pod definition
3) Namespace declaration file named as 0_namespace.yaml because of order, namespace should be created firstly.

Files:
- config.yaml 
- deployment.yaml
- service.yaml
- namespace.yaml

## How to apply all k8s changes:
``
kubectl apply -f namespace.yaml
kubectl apply -f config.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
``

## How to rollback changes:

``
kubectl delete deployment <deployment-name> -n <namespace>
kubectl delete service <service-name> -n <namespace>
kubectl delete configmap <configmap-name> -n <namespace>
kubectl delete namespace <namespace>
``
