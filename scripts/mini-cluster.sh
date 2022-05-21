#!/usr/bin/env sh
#  minikube start --nodes 3 -p 10.10.1.2 10.10.1.1 10.10.1.3
kubectl label nodes multinode-demo kubernetes.io/role=landowner --overwrite
kubectl label nodes multinode-demo-m02 kubernetes.io/role=hardworker --overwrite
kubectl label nodes multinode-demo-m03 kubernetes.io/role=hardworker --overwrite
kubectl get nodes
