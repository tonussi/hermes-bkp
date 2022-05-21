#!/usr/bin/env sh
KUBERNETES_DIR=$1

echo "deleting client..."
kubectl delete -f $KUBERNETES_DIR/tcp-kv-client.yml

echo "deleting server..."
kubectl delete -f $KUBERNETES_DIR/tcp-kv-hashicorp-raft-leader.yml
kubectl delete -f $KUBERNETES_DIR/tcp-kv-hashicorp-raft-followers.yml
