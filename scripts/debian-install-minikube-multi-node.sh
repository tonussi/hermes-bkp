#!/usr/bin/env sh

# Install deps
sudo sh -c "echo 'deb [arch=amd64] https://download.virtualbox.org/virtualbox/debian bullseye contrib' > /etc/apt/sources.list.d/virtualbox.list"
sudo apt -yq update
sudo apt -yq install git docker docker.io virtualbox curl wget apt-transport-https ca-certificates virtualbox
sudo usermod -aG docker $USER && newgrp docker

# Make root mounted as rshared to fix kube-dns issues.
# For multiple physical nodes
# sudo mount --make-rshared /

# Download kubectl, which is a requirement for using minikube.
curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
chmod +x kubectl && sudo mv kubectl /usr/local/bin/

# Download minikube.
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/

# Config minikube
sudo chmod 755 /usr/local/bin/minikube
minikube version

# After minikube start
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
