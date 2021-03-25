#!/usr/bin/env sh
KUBERNETES_DIR=$1
export N_CLIENTS=$2
export N_THREADS=$3
SCENE=$4
TEST=$5

echo "apply server..."
kubectl apply -f $KUBERNETES_DIR/tcp-kv-server.yml

sleep 5

echo "wait all replicas to be ready..."
until [ "$(kubectl get deployments -l app=tcp-kv-server -o jsonpath="{.items[0].status.replicas}")" = "$(kubectl get deployments -l app=tcp-kv-server -o jsonpath="{.items[0].status.readyReplicas}")" ]
do
  sleep 5;
done

echo "wait server to be running..."
until [ "$(kubectl get pods -l app=tcp-kv-server -o jsonpath="{.items[0].status.phase}")" = "Running" ]
do
  sleep 5;
done

echo "apply clients..."
envsubst < $KUBERNETES_DIR/tcp-kv-client.yml | kubectl apply -f -

echo "wait job to complete..."
kubectl wait --for=condition=complete --timeout=1h job.batch/tcp-kv-client

echo "collecting throughput log..."
kubectl cp $(kubectl get pods -l app=tcp-kv-server -o=jsonpath='{.items[0].metadata.name}'):/tmp/throughput.log logs/$SCENE/throughput/$TEST.log

echo "collecting latency log..."
kubectl logs $(kubectl get pods -l app=tcp-kv-client -o=jsonpath='{.items[0].metadata.name}') > logs/$SCENE/latency/$TEST.log

echo "deleting client..."
kubectl delete -f $KUBERNETES_DIR/tcp-kv-client.yml

echo "deleting server..."
kubectl delete -f $KUBERNETES_DIR/tcp-kv-server.yml