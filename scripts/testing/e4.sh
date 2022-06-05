#!/usr/bin/env sh
export KUBERNETES_DIR=$1
export N_CLIENTS=$2
export N_THREADS=$3
export READ_RATE=$4
export SCENE=$5
export PAYLOAD_SIZE=1
export QTY_ITERATION=100000
export THINKING_TIME=0.2
export PERCENTAGE_SAMPLING=90
export SERVICE_NAME=hermes-leader

TEST="qty-clients"-$(expr $N_CLIENTS \* $N_THREADS)
EXPERIMENT_NAME=$SCENE/$(expr 100 \- $READ_RATE)/$TEST
mkdir -p logs/$EXPERIMENT_NAME/throughput
mkdir -p logs/$EXPERIMENT_NAME/logs
mkdir -p logs/$EXPERIMENT_NAME/operations
mkdir -p logs/$EXPERIMENT_NAME/operations
mkdir -p logs/$EXPERIMENT_NAME/latency

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

echo "collecting hermes throughput..."
echo $(kubectl get pods -l app=hermes-leader -o=jsonpath='{.items[0].metadata.name}')
kubectl logs $(kubectl get pods -l app=hermes-leader -o=jsonpath='{.items[0].metadata.name}') http-log-server > logs/$EXPERIMENT_NAME/throughput/leader.log

for i in $(seq 0 1)
do
echo "collecting followers throughput..."
echo $(kubectl get pods -l app=hermes-followers -o=jsonpath="{.items[$i].metadata.name}")
kubectl logs $(kubectl get pods -l app=hermes-followers -o=jsonpath="{.items[$i].metadata.name}") http-log-server > logs/$EXPERIMENT_NAME/http-log-server/throughput/follower-$i.log
done

echo "collecting hermes logs..."
echo $(kubectl get pods -l app=hermes-leader -o=jsonpath='{.items[0].metadata.name}')
kubectl logs $(kubectl get pods -l app=hermes-leader -o=jsonpath='{.items[0].metadata.name}') hermes > logs/$EXPERIMENT_NAME/logs/leader.log

for i in $(seq 0 1)
do
echo "collecting followers logs..."
echo $(kubectl get pods -l app=hermes-followers -o=jsonpath="{.items[$i].metadata.name}")
kubectl logs $(kubectl get pods -l app=hermes-followers -o=jsonpath="{.items[$i].metadata.name}") hermes > logs/$EXPERIMENT_NAME/logs/follower-$i.log
done

echo "collecting data..."
# Defaulted container "http-log-server" out of: http-log-server, hermes
echo $(kubectl get pods -l app=hermes-leader -o=jsonpath="{.items[0].metadata.name}")
kubectl cp $(kubectl get pods -l app=hermes-leader -o=jsonpath="{.items[0].metadata.name}"):/tmp/logs/operations.log logs/$EXPERIMENT_NAME/operations/leader.log

for i in $(seq 0 1)
do
echo "collecting data..."
kubectl cp $(kubectl get pods -l app=hermes-followers -o=jsonpath="{.items[$i].metadata.name}"):/tmp/logs/operations.log logs/$EXPERIMENT_NAME/operations/follower-$i.log
done

echo "collecting latency log..."
for i in $(seq $(expr $N_CLIENTS - 1))
do
echo $(kubectl get pods -l app=http-log-client -o=jsonpath="{.items[$i].metadata.name}")
kubectl logs $(kubectl get pods -l app=http-log-client -o=jsonpath="{.items[$i].metadata.name}") > logs/$EXPERIMENT_NAME/latency/client-$i.log
done

echo "deleting client..."
kubectl delete -f $KUBERNETES_DIR/http-log-client.yml

echo "deleting server..."
kubectl delete -f $KUBERNETES_DIR/hermes-leader.yml
kubectl delete -f $KUBERNETES_DIR/hermes-followers.yml
