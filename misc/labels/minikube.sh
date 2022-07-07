#!/usr/bin/env sh
MINIKUBE_NODE_PREFIX=$1

# minikube

kubectl label nodes $MINIKUBE_NODE_PREFIX kubernetes.io/role=server --overwrite
kubectl label nodes $MINIKUBE_NODE_PREFIX-m02 kubernetes.io/role=client --overwrite
kubectl label nodes $MINIKUBE_NODE_PREFIX-m03 kubernetes.io/role=client --overwrite

kubectl label nodes $MINIKUBE_NODE_PREFIX role=server --overwrite
kubectl label nodes $MINIKUBE_NODE_PREFIX-m02 role=client --overwrite
kubectl label nodes $MINIKUBE_NODE_PREFIX-m03 role=client --overwrite
