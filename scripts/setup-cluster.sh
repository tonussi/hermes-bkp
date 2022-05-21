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

echo "label server nodes..."
for i in $(seq 0 $(expr $N_SERVER_NODES - 1))
do
  kubectl label node node$i.$EXPERIMENT_NAME.$EMULAB_GROUP_NAME.emulab.net key=server --overwrite
  # kubectl taint nodes node$i.$EXPERIMENT_NAME.$EMULAB_GROUP_NAME.emulab.net node-role.kubernetes.io/master-
done

echo "label client nodes..."
N_NODES=$(kubectl get nodes -o go-template="{{len .items}}")

for i in $(seq $(expr $N_SERVER_NODES) $(expr $N_NODES))
do
  kubectl label nodes node$i.$EXPERIMENT_NAME.$EMULAB_GROUP_NAME.emulab.net key=client --overwrite
  # kubectl taint nodes node$i.$EXPERIMENT_NAME.$EMULAB_GROUP_NAME.emulab.net node-role.kubernetes.io/master-
done

# kubectl label nodes node0.hermes-lucas.lptonussi.emulab.net key=server --overwrite
# kubectl label nodes node1.hermes-lucas.lptonussi.emulab.net key=client --overwrite
# kubectl label nodes node2.hermes-lucas.lptonussi.emulab.net key=client --overwrite
