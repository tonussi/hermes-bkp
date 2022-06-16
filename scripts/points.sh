#!/usr/bin/env sh

# for i in $(seq 2 2 14)
# do
#     if [ `expr $i % 2` -eq 0 ]
#     then
#         echo $i
#         bash ./scripts/experiment.sh kubernetes 1 $i 50 raw/e3/128-3-rw
#     fi
# done

# for i in $(seq 2 2 14)
# do
#     if [ `expr $i % 2` -eq 0 ]
#     then
#         echo $i
#         bash ./scripts/experiment.sh kubernetes 2 $i 50 raw/e3/128-3-rw
#         sleep 10
#     fi
# done

# for i in $(seq 2 2 14)
# do
#     if [ `expr $i % 2` -eq 0 ]
#     then
#         echo $i
#         bash ./scripts/experiment.sh kubernetes 3 $i 50 raw/e3/128-3-rw
#         sleep 10
#     fi
# done

# for i in $(seq 2 2 14)
# do
#     if [ `expr $i % 2` -eq 0 ]
#     then
#         echo $i
#         bash ./scripts/experiment.sh kubernetes 4 $i 50 raw/e3/128-3-rw
#         sleep 10
#     fi
# done

bash ./scripts/experiment.sh kubernetes 1 2 50 raw/e3/128-3-rw
sleep 10

bash ./scripts/experiment.sh kubernetes 1 4 50 raw/e3/128-3-rw
sleep 10

bash ./scripts/experiment.sh kubernetes 1 8 50 raw/e3/128-3-rw
sleep 10

bash ./scripts/experiment.sh kubernetes 2 4 50 raw/e3/128-3-rw
sleep 10

bash ./scripts/experiment.sh kubernetes 1 16 50 raw/e3/128-3-rw
sleep 10

bash ./scripts/experiment.sh kubernetes 2 8 50 raw/e3/128-3-rw
sleep 10

bash ./scripts/experiment.sh kubernetes 1 32 50 raw/e3/128-3-rw
sleep 10

bash ./scripts/experiment.sh kubernetes 2 16 50 raw/e3/128-3-rw
sleep 10

bash ./scripts/experiment.sh kubernetes 1 64 50 raw/e3/128-3-rw
sleep 10

bash ./scripts/experiment.sh kubernetes 2 32 50 raw/e3/128-3-rw
sleep 10

bash ./scripts/experiment.sh kubernetes 2 36 50 raw/e3/128-3-rw
sleep 10

bash ./scripts/experiment.sh kubernetes 3 24 50 raw/e3/128-3-rw
sleep 10
