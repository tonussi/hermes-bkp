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

# kubectl label nodes node0.simulation-1.lptonussi.emulab.net key=server --overwrite
# kubectl label nodes node1.simulation-1.lptonussi.emulab.net key=client --overwrite
# kubectl label nodes node2.simulation-1.lptonussi.emulab.net key=client --overwrite
#  minikube start --nodes 3 -p 10.10.1.2 10.10.1.1 10.10.1.3
# kubectl label nodes 10.10.1.2 key=server --overwrite
# kubectl label nodes 10.10.1.1-m02 key=client --overwrite
# kubectl label nodes 10.10.1.3-m03 key=client --overwrite
kubectl label nodes multinode-demo key=server --overwrite
kubectl label nodes multinode-demo-m02 key=client --overwrite
kubectl label nodes multinode-demo-m03 key=client --overwrite
kubectl label nodes multinode-demo kubernetes.io/role=landowner --overwrite
kubectl label nodes multinode-demo-m02 kubernetes.io/role=hardworker --overwrite
kubectl label nodes multinode-demo-m03 kubernetes.io/role=hardworker --overwrite
kubectl get nodes
