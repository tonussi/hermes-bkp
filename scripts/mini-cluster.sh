#!/usr/bin/env sh
MINIKUBE_NODE_PREFIX=$1
#  minikube start --nodes 3 -p 10.10.1.2 10.10.1.1 10.10.1.3
kubectl label nodes $MINIKUBE_NODE_PREFIX kubernetes.io/role=landowner --overwrite
kubectl label nodes $MINIKUBE_NODE_PREFIX-m02 kubernetes.io/role=hardworker --overwrite
kubectl label nodes $MINIKUBE_NODE_PREFIX-m03 kubernetes.io/role=hardworker --overwrite
kubectl get nodes
