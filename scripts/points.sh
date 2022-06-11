#!/usr/bin/env sh

for i in $(seq 100)
do
    bash ./scripts/experiment.sh kubernetes $i 2 50 e3/1e6-8-rw
done
