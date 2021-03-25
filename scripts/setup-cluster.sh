#!/usr/bin/env sh
ANSIBLE_PB_DIR=$1
KUBERNETES_DIR=$2
EXPERIMENT_NAME=$3
N_SERVER_NODES=$4

DOCKERD_JSON_FILE=`echo $(cd $ANSIBLE_PB_DIR; pwd)/daemon.json`

echo "ping all nodes..."
ansible all -i $ANSIBLE_PB_DIR/hosts -m ping

echo "initial cluster setup..."
ansible-playbook -i $ANSIBLE_PB_DIR/hosts --extra-vars daemon_json_file=$DOCKERD_JSON_FILE $ANSIBLE_PB_DIR/initial-setup.yml

echo "initial master setup..."
ansible-playbook -i $ANSIBLE_PB_DIR/hosts $ANSIBLE_PB_DIR/init-master.yml

echo "apply flannel pod network..."
kubectl apply -f $KUBERNETES_DIR/kube-flannel.yml

echo "join worker nodes to cluster..."
ansible-playbook -i $ANSIBLE_PB_DIR/hosts $ANSIBLE_PB_DIR/join-workers.yml

echo "label server nodes..."
for i in $(seq 1 $N_SERVER_NODES)
do
  kubectl label nodes node$i.$EXPERIMENT_NAME.scalablesmr.emulab.net role=server
done

echo "label client nodes..."
N_NODES=$(kubectl get nodes -o go-template="{{len .items}}")

for i in $(seq $(expr $N_SERVER_NODES + 1) $N_NODES)
do
  kubectl label nodes node$i.$EXPERIMENT_NAME.scalablesmr.emulab.net role=client
done
