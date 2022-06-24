#!/usr/bin/env sh

for i in $(seq 2 2 64)
do
    if [ `expr $i % 2` -eq 0 ]
    then
        echo $i
        bash ./scripts/experiment.sh kubernetes 1 $i 50 $1
        sleep 10
    fi
done

for i in $(seq 2 2 14)
do
    if [ `expr $i % 2` -eq 0 ]
    then
        echo $i
        bash ./scripts/experiment.sh kubernetes 2 $i 50 $1
        sleep 10
    fi
done

for i in $(seq 2 2 14)
do
    if [ `expr $i % 2` -eq 0 ]
    then
        echo $i
        bash ./scripts/experiment.sh kubernetes 3 $i 50 $1
        sleep 10
    fi
done

for i in $(seq 2 2 14)
do
    if [ `expr $i % 2` -eq 0 ]
    then
        echo $i
        bash ./scripts/experiment.sh kubernetes 4 $i 50 $1
        sleep 10
    fi
done
