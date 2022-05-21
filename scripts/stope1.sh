#!/usr/bin/env sh
export KUBERNETES_DIR=kubernetes

kubectl delete -f $KUBERNETES_DIR/tcp-kv-client.yml
kubectl delete -f $KUBERNETES_DIR/tcp-kv-server.yml
