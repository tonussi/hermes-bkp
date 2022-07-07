#!/usr/bin/env sh

CLUSTER_NAME=$1
TYPE=$2

echo "label the worker nodes (where the experiment happens)..."
N_NODES=$(kubectl get nodes -o go-template="{{len .items}}")

echo "server roles"
for i in $(seq 0 2)
do
  kubectl label nodes k3d-$CLUSTER_NAME-$TYPE-$i role=server --overwrite
  kubectl label nodes k3d-$CLUSTER_NAME-$TYPE-$i kubernetes.io/role=server --overwrite
done

echo "client roles"
for i in $(seq 3 $(expr $N_NODES \- 1))
do
  kubectl label nodes k3d-$CLUSTER_NAME-$TYPE-$i role=client --overwrite
  kubectl label nodes k3d-$CLUSTER_NAME-$TYPE-$i kubernetes.io/role=client --overwrite
done

kubectl get nodes --show-labels
