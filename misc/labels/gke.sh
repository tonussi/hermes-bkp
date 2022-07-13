#!/usr/bin/env sh

kubectl label nodes gke-hermes-default-pool-ba8b6fd5-6ms5 kubernetes.io/role=server
kubectl label nodes gke-hermes-default-pool-ba8b6fd5-6vbt kubernetes.io/role=server
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

kubectl label nodes gke-hermes-default-pool-ba8b6fd5-6ms5 kubernetes.io/role=server --overwrite
kubectl label nodes gke-hermes-default-pool-ba8b6fd5-6vbt kubernetes.io/role=server --overwrite
kubectl label nodes gke-hermes-default-pool-ba8b6fd5-hjgm kubernetes.io/role=server --overwrite
kubectl label nodes gke-hermes-default-pool-ba8b6fd5-tzfx kubernetes.io/role=client --overwrite
kubectl label nodes gke-hermes-default-pool-ba8b6fd5-wrbq kubernetes.io/role=client --overwrite

kubectl label nodes gke-hermes-default-pool-ba8b6fd5-6ms5 role=server --overwrite
kubectl label nodes gke-hermes-default-pool-ba8b6fd5-6vbt role=server --overwrite
kubectl label nodes gke-hermes-default-pool-ba8b6fd5-hjgm role=server --overwrite
kubectl label nodes gke-hermes-default-pool-ba8b6fd5-tzfx role=client --overwrite
kubectl label nodes gke-hermes-default-pool-ba8b6fd5-wrbq role=client --overwrite
