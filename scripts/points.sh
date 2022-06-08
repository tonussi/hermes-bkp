#!/usr/bin/env sh

sleep 5
bash ./scripts/experiment.sh kubernetes 1 2 50 e3/1000-rw
sleep 5
bash ./scripts/experiment.sh kubernetes 1 4 50 e3/1000-rw

sleep 5
bash ./scripts/experiment.sh kubernetes 2 2 50 e3/1000-rw
sleep 5
bash ./scripts/experiment.sh kubernetes 2 $(expr 2 \* 2) 50 e3/1000-rw

sleep 5
bash ./scripts/experiment.sh kubernetes 3 5 50 e3/1000-rw
sleep 5
bash ./scripts/experiment.sh kubernetes 3 $(expr 2 \* 3) 50 e3/1000-rw

sleep 5
bash ./scripts/experiment.sh kubernetes 4 2 50 e3/1000-rw
sleep 5
bash ./scripts/experiment.sh kubernetes 4 $(expr 2 \* 4) 50 e3/1000-rw

sleep 5
bash ./scripts/experiment.sh kubernetes 5 2 50 e3/1000-rw
sleep 5
bash ./scripts/experiment.sh kubernetes 5 $(expr 2 \* 5) 50 e3/1000-rw

sleep 5
bash ./scripts/experiment.sh kubernetes 6 2 50 e3/1000-rw
sleep 5
bash ./scripts/experiment.sh kubernetes 6 $(expr 2 \* 6) 50 e3/1000-rw

sleep 5
bash ./scripts/experiment.sh kubernetes 7 2 50 e3/1000-rw
sleep 5
bash ./scripts/experiment.sh kubernetes 7 $(expr 2 \* 7) 50 e3/1000-rw

sleep 5
bash ./scripts/experiment.sh kubernetes 8 2 50 e3/1000-rw
sleep 5
bash ./scripts/experiment.sh kubernetes 8 $(expr 2 \* 8) 50 e3/1000-rw

sleep 5
bash ./scripts/experiment.sh kubernetes 9 2 50 e3/1000-rw
sleep 5
bash ./scripts/experiment.sh kubernetes 9 $(expr 2 \* 9) 50 e3/1000-rw

sleep 5
bash ./scripts/experiment.sh kubernetes 10 2 50 e3/1000-rw
sleep 5
bash ./scripts/experiment.sh kubernetes 10 $(expr 2 \* 10) 50 e3/1000-rw

