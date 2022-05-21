#!/usr/bin/env sh
DOCKERHUB_USER_NAME=$1

docker build -t $DOCKERHUB_USER_NAME/public:hermes -f docker/hermes.dockerfile .
# docker build -t $DOCKERHUB_USER_NAME/public:tcp-kv-client -f docker/tcp-kv-client.dockerfile .
# docker build -t $DOCKERHUB_USER_NAME/public:tcp-kv-hashicorp-raft -f docker/tcp-kv-hashicorp-raft.dockerfile .
# docker build -t $DOCKERHUB_USER_NAME/public:tcp-kv-server -f docker/tcp-kv-server.dockerfile .

# docker push $DOCKERHUB_USER_NAME/public:hermes
# docker push $DOCKERHUB_USER_NAME/public:tcp-kv-client
# docker push $DOCKERHUB_USER_NAME/public:tcp-kv-hashicorp-raft
# docker push $DOCKERHUB_USER_NAME/public:tcp-kv-server

# docker run --network=host --name hermes $DOCKERHUB_USER_NAME/public:hermes
# docker run --network=host --name tcp-kv-client $DOCKERHUB_USER_NAME/public:tcp-kv-client
# docker run --network=host --name tcp-kv-hashicorp-raft $DOCKERHUB_USER_NAME/public:tcp-kv-hashicorp-raft
# docker run --network=host --name tcp-kv-server $DOCKERHUB_USER_NAME/public:tcp-kv-server
