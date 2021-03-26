#!/usr/bin/env sh
KUBERNETES_DIR=$1
SCENE=$2
TEST=$3
APP=tcp-kv-hashicorp-raft-leader

echo "collecting throughput log..."
kubectl cp $(kubectl get pods -l app=$APP -o=jsonpath='{.items[0].metadata.name}'):/tmp/throughput.log logs/$SCENE/throughput/$TEST.log

echo "collecting latency log..."
kubectl logs $(kubectl get pods -l app=tcp-kv-client -o=jsonpath='{.items[0].metadata.name}') > logs/$SCENE/latency/$TEST.log

echo "deleting client..."
kubectl delete -f $KUBERNETES_DIR/tcp-kv-client.yml

echo "deleting server..."
kubectl delete -f $KUBERNETES_DIR/$APP.yml
kubectl delete -f $KUBERNETES_DIR/tcp-kv-hashicorp-raft-followers.yml
