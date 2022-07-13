#!/usr/bin/env sh
KUBERNETES_DIR=$1
APP=$2

kubectl delete -f $KUBERNETES_DIR/http-$APP-client.yml
kubectl delete -f $KUBERNETES_DIR/hermes-leader.yml
kubectl delete -f $KUBERNETES_DIR/hermes-followers.yml
