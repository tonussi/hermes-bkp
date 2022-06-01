#!/usr/bin/env sh
export KUBERNETES_DIR=$1
export N_CLIENTS=$2
export N_THREADS=$3
export READ_RATE=$4
export SCENE=$5
export PAYLOAD_SIZE=2
export QTY_ITERATION=100000
export THINKING_TIME=0.2
export PERCENTAGE_SAMPLING=90
export SERVICE_NAME=hermes-leader

echo "apply leader..."
kubectl apply -f $KUBERNETES_DIR/hermes-leader.yml

sleep 5

echo "wait all replicas to be ready..."
until [ "$(kubectl get deployments -l app=hermes-leader -o jsonpath="{.items[0].status.replicas}")" = "$(kubectl get deployments -l app=hermes-leader -o jsonpath="{.items[0].status.readyReplicas}")" ]
do
sleep 5;
done

echo "wait server to be running..."
until [ "$(kubectl get pods -l app=hermes-leader -o jsonpath="{.items[0].status.phase}")" = "Running" ]
do
sleep 5;
done

echo "apply followers..."
kubectl apply -f $KUBERNETES_DIR/hermes-followers.yml

sleep 10

echo "wait all replicas to be ready..."
until [ "$(kubectl get deployments -l app=hermes-followers -o jsonpath="{.items[0].status.replicas}")" = "$(kubectl get deployments -l app=hermes-followers -o jsonpath="{.items[0].status.readyReplicas}")" ]
do
sleep 5;
done

echo "wait server to be running..."
until [ "$(kubectl get pods -l app=hermes-followers -o jsonpath="{.items[0].status.phase}")" = "Running" ]
do
sleep 5;
done

sleep 10

echo "apply clients..."
envsubst < $KUBERNETES_DIR/http-log-client.yml | kubectl apply -f -

echo "wait job to complete..."
kubectl wait --for=condition=complete --timeout=100s job.batch/http-log-client

TEST=$(expr $N_CLIENTS \* $N_THREADS)-$N_CLIENTS

echo "collecting hermes logs..."
mkdir -p logs/$SCENE/throughput/http-log-server/
echo $(kubectl get pods -l app=hermes-leader -o=jsonpath='{.items[0].metadata.name}')
kubectl logs $(kubectl get pods -l app=hermes-leader -o=jsonpath='{.items[0].metadata.name}') http-log-server > logs/$SCENE/throughput/http-log-server/$TEST.log

echo "collecting hermes logs..."
mkdir -p logs/$SCENE/hermes-output/
echo $(kubectl get pods -l app=hermes-leader -o=jsonpath='{.items[0].metadata.name}')
kubectl logs $(kubectl get pods -l app=hermes-leader -o=jsonpath='{.items[0].metadata.name}') hermes > logs/$SCENE/hermes-output/hermes-leader.log

for i in $(seq 0 1)
do
echo "collecting throughput log..."
mkdir -p logs/$SCENE/throughput/hermes-follower-$i/
echo $(kubectl get pods -l app=hermes-followers -o=jsonpath="{.items[$i].metadata.name}")
kubectl logs $(kubectl get pods -l app=hermes-followers -o=jsonpath="{.items[$i].metadata.name}") http-log-server > logs/$SCENE/throughput/http-log-server-$i/$TEST.log

echo "collecting hermes logs..."
mkdir -p logs/$SCENE/hermes-output/
echo $(kubectl get pods -l app=hermes-followers -o=jsonpath="{.items[$i].metadata.name}")
kubectl logs $(kubectl get pods -l app=hermes-followers -o=jsonpath="{.items[$i].metadata.name}") hermes > logs/$SCENE/hermes-output/hermes-follower-$i.log
done

echo "collecting data..."
# Defaulted container "http-log-server" out of: http-log-server, hermes
echo $(kubectl get pods -l app=hermes-leader -o=jsonpath="{.items[0].metadata.name}"):tmp/logs/operations.log
mkdir -p logs/$SCENE/operations
kubectl cp $(kubectl get pods -l app=hermes-leader -o=jsonpath="{.items[0].metadata.name}"):/tmp/logs/operations.log logs/$SCENE/operations/operations.log

for i in $(seq 0 1)
do
echo "collecting data..."
mkdir -p logs/$SCENE/operations
kubectl cp $(kubectl get pods -l app=hermes-followers -o=jsonpath="{.items[$i].metadata.name}"):/tmp/logs/operations.log logs/$SCENE/operations/operations-$i.log
done

echo "collecting latency log..."
for i in $(seq $(expr $N_CLIENTS - 1))
do
mkdir -p logs/$SCENE/latency/client-$i
echo $(kubectl get pods -l app=http-log-client -o=jsonpath="{.items[$i].metadata.name}")
kubectl logs $(kubectl get pods -l app=http-log-client -o=jsonpath="{.items[$i].metadata.name}") > logs/$SCENE/latency/client-$i/$TEST.log
done

echo "deleting client..."
kubectl delete -f $KUBERNETES_DIR/http-log-client.yml

echo "deleting server..."
kubectl delete -f $KUBERNETES_DIR/hermes-leader.yml
kubectl delete -f $KUBERNETES_DIR/hermes-followers.yml
