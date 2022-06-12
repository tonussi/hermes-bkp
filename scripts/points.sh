#!/usr/bin/env sh

for i in $(seq 4)
do
    bash ./scripts/experiment.sh kubernetes 1 $i 50 e3/1e5-11-rw-microsecs
done

for i in $(seq 4)
do
    bash ./scripts/experiment.sh kubernetes 2 $i 50 e3/1e5-11-rw-microsecs
done

for i in $(seq 4)
do
    bash ./scripts/experiment.sh kubernetes 3 $i 50 e3/1e5-11-rw-microsecs
done
