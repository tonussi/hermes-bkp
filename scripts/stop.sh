#!/usr/bin/env sh
KUBERNETES_DIR=$1

kubectl delete -f $KUBERNETES_DIR/http-log-client.yml
kubectl delete -f $KUBERNETES_DIR/hermes-leader.yml
kubectl delete -f $KUBERNETES_DIR/hermes-followers.yml
