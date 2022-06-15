#!/usr/bin/env sh

# for i in $(seq 2 2 14)
# do
#     if [ `expr $i % 2` -eq 0 ]
#     then
#         echo $i
#         bash ./scripts/experiment.sh kubernetes 1 $i 50 e3/128-rw
#     fi
# done

# for i in $(seq 2 2 14)
# do
#     if [ `expr $i % 2` -eq 0 ]
#     then
#         echo $i
#         bash ./scripts/experiment.sh kubernetes 2 $i 50 e3/128-rw
#     fi
# done

# for i in $(seq 2 2 14)
# do
#     if [ `expr $i % 2` -eq 0 ]
#     then
#         echo $i
#         bash ./scripts/experiment.sh kubernetes 3 $i 50 e3/128-rw
#     fi
# done

# for i in $(seq 2 2 14)
# do
#     if [ `expr $i % 2` -eq 0 ]
#     then
#         echo $i
#         bash ./scripts/experiment.sh kubernetes 4 $i 50 e3/128-rw
#     fi
# done

bash ./scripts/experiment.sh kubernetes 1 $(expr 2 \/ 1) 50 e3/128-rw
bash ./scripts/experiment.sh kubernetes 1 $(expr 4 \/ 1) 50 e3/128-rw
bash ./scripts/experiment.sh kubernetes 1 $(expr 8 \/ 1) 50 e3/128-rw
bash ./scripts/experiment.sh kubernetes 2 $(expr 8 \/ 2) 50 e3/128-rw
bash ./scripts/experiment.sh kubernetes 1 $(expr 16 \/ 1) 50 e3/128-rw
bash ./scripts/experiment.sh kubernetes 2 $(expr 16 \/ 2) 50 e3/128-rw
bash ./scripts/experiment.sh kubernetes 1 $(expr 32 \/ 1) 50 e3/128-rw
bash ./scripts/experiment.sh kubernetes 2 $(expr 32 \/ 2) 50 e3/128-rw
bash ./scripts/experiment.sh kubernetes 1 $(expr 64 \/ 1) 50 e3/128-rw
bash ./scripts/experiment.sh kubernetes 2 $(expr 64 \/ 2) 50 e3/128-rw

# bash ./scripts/experiment.sh kubernetes 2 $(expr 72 \/ 2) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 3 $(expr 72 \/ 3) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 2 $(expr 76 \/ 2) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 2 $(expr 84 \/ 2) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 2 $(expr 96 \/ 2) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 3 $(expr 96 \/ 3) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 2 $(expr 104 \/ 2) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 3 $(expr 108 \/ 3) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 3 $(expr 126 \/ 3) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 1 $(expr 128 \/ 1) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 2 $(expr 128 \/ 2) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 3 $(expr 132 \/ 3) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 3 $(expr 192 \/ 3) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 4 $(expr 208 \/ 4) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 2 $(expr 256 \/ 2) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 4 $(expr 256 \/ 4) 50 e3/128-rw
# bash ./scripts/experiment.sh kubernetes 4 $(expr 288 \/ 4) 50 e3/128-rw
