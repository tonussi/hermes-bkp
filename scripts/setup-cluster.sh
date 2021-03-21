#!/usr/bin/env sh
ANSIBLE_PB_DIR=$1
KUBERNETES_DIR=$2
shift 1

DOCKERD_JSON_FILE="$(cd ansible; pwd)/daemon.json"

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
