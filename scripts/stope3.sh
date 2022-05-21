#!/usr/bin/env sh
export KUBERNETES_DIR=$1

kubectl delete -f $KUBERNETES_DIR/tcp-kv-client.yml
kubectl delete -f $KUBERNETES_DIR/hermes-leader.yml
kubectl delete -f $KUBERNETES_DIR/hermes-followers.yml

kubectl delete -f $KUBERNETES_DIR/testing/hermes-leader.yml
kubectl delete -f $KUBERNETES_DIR/testing/hermes-followers.yml
kubectl delete -f $KUBERNETES_DIR/testing/http-log-client.yml
