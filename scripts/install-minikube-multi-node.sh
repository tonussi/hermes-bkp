#!/usr/bin/env sh

set -exu

sudo sh -c "echo 'deb [arch=amd64] https://download.virtualbox.org/virtualbox/debian bullseye contrib' >> /etc/apt/sources.list.d/virtualbox.list"

sudo apt -yq update
sudo apt -yq install git docker docker.io virtualbox curl wget apt-transport-https ca-certificates virtualbox virtualbox-ext-pack

sudo curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo apt -yq update
sudo apt -yq install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl

# Make root mounted as rshared to fix kube-dns issues.
sudo mount --make-rshared /
# Download kubectl, which is a requirement for using minikube.
curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v1.9.0/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/

# Download minikube.
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
sudo minikube start --vm-driver=none --kubernetes-version=v1.10.0 --force
# Fix the kubectl context, as it's often stale.
minikube update-context
# Wait for Kubernetes to be up and ready.
JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'; until kubectl get nodes -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 1; done
# kube-addon-manager is responsible for managing other kubernetes components, such as kube-dns, dashboard, storage-provisioner..
JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'; until kubectl -n kube-system get pods -lcomponent=kube-addon-manager -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 1;echo "waiting for kube-addon-manager to be available"; kubectl get pods --all-namespaces; done
# Wait for kube-dns to be ready.
JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'; until kubectl -n kube-system get pods -lk8s-app=kube-dns -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 1;echo "waiting for kube-dns to be available"; kubectl get pods --all-namespaces; done

# Install stafy

echo "let minikube handle $HOME/.minikube; if you need space you will have to remove $HOME/.minikube"
wget https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo cp minikube-linux-amd64 /usr/local/bin/minikube
sudo chmod 755 /usr/local/bin/minikube
minikube version
curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl
kubectl version -o json
sudo chmod +x /users/$USER/.minikube
sudo groupadd dockersudo usermod -aG docker $USER
docker run hello-world
sudo chown -R $USER $HOME/.minikube; chmod -R u+wrx $HOME/.minikube
echo "https://minikube.sigs.k8s.io/docs/tutorials/multi_node/"
# use this if you dont have space in disk
# rm -rf $HOME/.minikube
# check $HOME
# df
# minikube start --nodes 3 -p multinode-demo
