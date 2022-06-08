#!/usr/bin/env sh
KUBERNETES_DIR=$1
PORT=$2

kubectl delete -f $KUBERNETES_DIR/recommended.yaml
kubectl delete -f $KUBERNETES_DIR/dashboard.admin-user.yml
kubectl delete -f $KUBERNETES_DIR/dashboard.admin-user-role.yml

kubectl apply -f $KUBERNETES_DIR/recommended.yaml
kubectl apply -f $KUBERNETES_DIR/dashboard.admin-user.yml
kubectl apply -f $KUBERNETES_DIR/dashboard.admin-user-role.yml

kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep eks-admin | awk '{print $1}')

echo "http://localhost:$PORT/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/#/workloads?namespace=default"

kubectl proxy --port=$PORT
