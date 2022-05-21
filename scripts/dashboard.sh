#!/usr/bin/env sh
KUBERNETES_DIR=$1

kubectl delete -f $KUBERNETES_DIR/dashboard-adminuser.yaml
kubectl delete -f $KUBERNETES_DIR/dashboard.yaml
kubectl delete -f $KUBERNETES_DIR/liveness.yml

kubectl apply -f $KUBERNETES_DIR/dashboard-adminuser.yaml
kubectl apply -f $KUBERNETES_DIR/dashboard.yaml
kubectl apply -f $KUBERNETES_DIR/liveness.yml

kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep eks-admin | awk '{print $1}')

PORT=8081

kubectl proxy --port=$PORT

echo "http://127.0.0.1:$PORT/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/#/workloads?namespace=default"
