#!/usr/bin/env sh

echo "running experiment 1 clients and $(expr 2 \/ 1) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 1 $(expr 2 \/ 1) 50 e3/1000-rw
sleep 10

echo "running experiment 1 clients and $(expr 4 \/ 1) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 1 $(expr 4 \/ 1) 50 e3/1000-rw
sleep 10

echo "running experiment 2 clients and $(expr 4 \/ 2) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 2 $(expr 4 \/ 2) 50 e3/1000-rw
sleep 10

echo "running experiment 1 clients and $(expr 6 \/ 1) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 1 $(expr 6 \/ 1) 50 e3/1000-rw
sleep 10

echo "running experiment 1 clients and $(expr 8 \/ 1) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 1 $(expr 8 \/ 1) 50 e3/1000-rw
sleep 10

echo "running experiment 2 clients and $(expr 8 \/ 2) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 2 $(expr 8 \/ 2) 50 e3/1000-rw
sleep 10

echo "running experiment 1 clients and $(expr 12 \/ 1) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 1 $(expr 12 \/ 1) 50 e3/1000-rw
sleep 10

echo "running experiment 2 clients and $(expr 12 \/ 2) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 2 $(expr 12 \/ 2) 50 e3/1000-rw
sleep 10

echo "running experiment 1 clients and $(expr 16 \/ 1) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 1 $(expr 16 \/ 1) 50 e3/1000-rw
sleep 10

echo "running experiment 2 clients and $(expr 16 \/ 2) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 2 $(expr 16 \/ 2) 50 e3/1000-rw
sleep 10

echo "running experiment 1 clients and $(expr 20 \/ 1) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 1 $(expr 20 \/ 1) 50 e3/1000-rw
sleep 10

echo "running experiment 2 clients and $(expr 20 \/ 2) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 2 $(expr 20 \/ 2) 50 e3/1000-rw
sleep 10

echo "running experiment 1 clients and $(expr 24 \/ 1) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 1 $(expr 24 \/ 1) 50 e3/1000-rw
sleep 10

echo "running experiment 2 clients and $(expr 24 \/ 2) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 2 $(expr 24 \/ 2) 50 e3/1000-rw
sleep 10

echo "running experiment 3 clients and $(expr 30 \/ 3) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 3 $(expr 30 \/ 3) 50 e3/1000-rw
sleep 10

echo "running experiment 1 clients and $(expr 32 \/ 1) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 1 $(expr 32 \/ 1) 50 e3/1000-rw
sleep 10

echo "running experiment 2 clients and $(expr 32 \/ 2) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 2 $(expr 32 \/ 2) 50 e3/1000-rw
sleep 10

echo "running experiment 1 clients and $(expr 36 \/ 1) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 1 $(expr 36 \/ 1) 50 e3/1000-rw
sleep 10

echo "running experiment 3 clients and $(expr 36 \/ 3) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 3 $(expr 36 \/ 3) 50 e3/1000-rw
sleep 10

echo "running experiment 3 clients and $(expr 42 \/ 3) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 3 $(expr 42 \/ 3) 50 e3/1000-rw
sleep 10

echo "running experiment 1 clients and $(expr 48 \/ 1) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 1 $(expr 48 \/ 1) 50 e3/1000-rw
sleep 10

echo "running experiment 2 clients and $(expr 48 \/ 2) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 2 $(expr 48 \/ 2) 50 e3/1000-rw
sleep 10

echo "running experiment 3 clients and $(expr 48 \/ 3) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 3 $(expr 48 \/ 3) 50 e3/1000-rw
sleep 10

echo "running experiment 1 clients and $(expr 64 \/ 1) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 1 $(expr 64 \/ 1) 50 e3/1000-rw
sleep 10

echo "running experiment 2 clients and $(expr 64 \/ 2) threads with a read rate of 50%"
sh ./scripts/experiment.sh kubernetes 2 $(expr 64 \/ 2) 50 e3/1000-rw
sleep 10
