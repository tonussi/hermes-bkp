#!/usr/bin/env sh
ANSIBLE_PB_DIR=$1
KUBERNETES_DIR=$2
EXPERIMENT_NAME=$3
EMULAB_GROUP_NAME=$4
N_SERVER_NODES=$5

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

echo "label the admin master node (for admin purposes only)..."
for i in $(seq 0 $(expr $N_SERVER_NODES - 1))
do
  kubectl label nodes node$i.$EXPERIMENT_NAME.$EMULAB_GROUP_NAME.emulab.net kubernetes.io/role=admin --overwrite
done

echo "label the worker nodes (where the experiment happens)..."
N_NODES=$(kubectl get nodes -o go-template="{{len .items}}")

for i in $(seq $(expr $N_SERVER_NODES) $(expr $N_SERVER_NODES))
do
  kubectl label nodes node$i.$EXPERIMENT_NAME.$EMULAB_GROUP_NAME.emulab.net kubernetes.io/role=server --overwrite
done

for i in $(seq $(expr $N_SERVER_NODES + 1) $(expr $N_NODES - 1))
do
  kubectl label nodes node$i.$EXPERIMENT_NAME.$EMULAB_GROUP_NAME.emulab.net kubernetes.io/role=client --overwrite
done
