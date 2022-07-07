#!/usr/bin/env sh

kubectl label nodes gke-hermes-default-pool-bdd746dd-0jbk kubernetes.io/role=server
kubectl label nodes gke-hermes-default-pool-bdd746dd-7dbl kubernetes.io/role=server
kubectl label nodes gke-hermes-default-pool-bdd746dd-dzc5 kubernetes.io/role=server
kubectl label nodes gke-hermes-default-pool-bdd746dd-hhlw kubernetes.io/role=client
kubectl label nodes gke-hermes-default-pool-bdd746dd-n64v kubernetes.io/role=client
kubectl label nodes gke-hermes-default-pool-bdd746dd-zq59 kubernetes.io/role=client

kubectl label nodes gke-hermes-default-pool-bdd746dd-0jbk role=server
kubectl label nodes gke-hermes-default-pool-bdd746dd-7dbl role=server
kubectl label nodes gke-hermes-default-pool-bdd746dd-dzc5 role=server
kubectl label nodes gke-hermes-default-pool-bdd746dd-hhlw role=client
kubectl label nodes gke-hermes-default-pool-bdd746dd-n64v role=client
kubectl label nodes gke-hermes-default-pool-bdd746dd-zq59 role=client

kubectl label nodes gke-hermes-hermes-pool-341ebb48-fxj1 role=server
kubectl label nodes gke-hermes-hermes-pool-341ebb48-r1vr role=server
kubectl label nodes gke-hermes-hermes-pool-341ebb48-v1p5 role=server
kubectl label nodes gke-hermes-hermes-pool-341ebb48-xntb role=client
kubectl label nodes gke-hermes-hermes-pool-341ebb48-zkh7 role=client
