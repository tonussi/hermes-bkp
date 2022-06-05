#!/usr/bin/env sh
MINIKUBE_NODE_PREFIX=$1

# minikube

kubectl label nodes $MINIKUBE_NODE_PREFIX kubernetes.io/role=server --overwrite
kubectl label nodes $MINIKUBE_NODE_PREFIX-m02 kubernetes.io/role=client --overwrite
kubectl label nodes $MINIKUBE_NODE_PREFIX-m03 kubernetes.io/role=client --overwrite

kubectl label nodes $MINIKUBE_NODE_PREFIX role=server --overwrite
kubectl label nodes $MINIKUBE_NODE_PREFIX-m02 role=client --overwrite
kubectl label nodes $MINIKUBE_NODE_PREFIX-m03 role=client --overwrite

# GKE

kubectl label nodes gke-cluster-1-default-pool-c6be4ebb-12np kubernetes.io/role=server --overwrite
kubectl label nodes gke-cluster-1-default-pool-c6be4ebb-1p0h kubernetes.io/role=client --overwrite
kubectl label nodes gke-cluster-1-default-pool-c6be4ebb-jgsd kubernetes.io/role=client --overwrite

kubectl label nodes gke-cluster-1-default-pool-c6be4ebb-12np role=server --overwrite
kubectl label nodes gke-cluster-1-default-pool-c6be4ebb-1p0h role=client --overwrite
kubectl label nodes gke-cluster-1-default-pool-c6be4ebb-jgsd role=client --overwrite

# k3d

kubectl label nodes k3d-hermes-server-0 role=server --overwrite
kubectl label nodes k3d-hermes-server-1 role=client --overwrite
kubectl label nodes k3d-hermes-server-2 role=client --overwrite

kubectl label nodes k3d-hermes-server-0 kubernetes.io/role=server --overwrite
kubectl label nodes k3d-hermes-server-1 kubernetes.io/role=client --overwrite
kubectl label nodes k3d-hermes-server-2 kubernetes.io/role=client --overwrite

# emulab

kubectl label nodes node0.hermes.lptonussi.emulab.net kubernetes.io/role=admin --overwrite
kubectl label nodes node1.hermes.lptonussi.emulab.net kubernetes.io/role=server --overwrite
kubectl label nodes node2.hermes.lptonussi.emulab.net kubernetes.io/role=server --overwrite
kubectl label nodes node3.hermes.lptonussi.emulab.net kubernetes.io/role=server --overwrite
kubectl label nodes node4.hermes.lptonussi.emulab.net kubernetes.io/role=client --overwrite

kubectl label nodes node0.hermes.lptonussi.emulab.net role=admin --overwrite
kubectl label nodes node1.hermes.lptonussi.emulab.net role=server --overwrite
kubectl label nodes node2.hermes.lptonussi.emulab.net role=server --overwrite
kubectl label nodes node3.hermes.lptonussi.emulab.net role=server --overwrite
kubectl label nodes node4.hermes.lptonussi.emulab.net role=client --overwrite
