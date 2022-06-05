#!/usr/bin/env sh
export KUBERNETES_DIR=kubernetes

kubectl delete -f $KUBERNETES_DIR/http-log-client.yml
kubectl delete -f $KUBERNETES_DIR/hermes-leader.yml
kubectl delete -f $KUBERNETES_DIR/hermes-followers.yml
