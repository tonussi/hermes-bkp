#!/usr/bin/env sh
MINIKUBE_NODE_PREFIX=$1
kubectl label nodes $MINIKUBE_NODE_PREFIX kubernetes.io/role=server --overwrite
kubectl label nodes $MINIKUBE_NODE_PREFIX-m02 kubernetes.io/role=client --overwrite
kubectl label nodes $MINIKUBE_NODE_PREFIX-m03 kubernetes.io/role=client --overwrite
kubectl get nodes

# GKE labels

# kubectl label nodes gke-cluster-1-default-pool-c6be4ebb-12np role=server --overwrite
# kubectl label nodes gke-cluster-1-default-pool-c6be4ebb-1p0h role=client --overwrite
# kubectl label nodes gke-cluster-1-default-pool-c6be4ebb-jgsd role=client --overwrite
