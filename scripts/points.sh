#!/usr/bin/env sh

bash ./scripts/experiment.sh kubernetes 1 2 0 raw/e3/128-w
sleep 10

bash ./scripts/experiment.sh kubernetes 1 4 0 raw/e3/128-w
sleep 10

bash ./scripts/experiment.sh kubernetes 1 8 0 raw/e3/128-w
sleep 10

bash ./scripts/experiment.sh kubernetes 2 4 0 raw/e3/128-w
sleep 10

bash ./scripts/experiment.sh kubernetes 1 16 0 raw/e3/128-w
sleep 10

bash ./scripts/experiment.sh kubernetes 2 8 0 raw/e3/128-w
sleep 10

bash ./scripts/experiment.sh kubernetes 1 32 0 raw/e3/128-w
sleep 10

bash ./scripts/experiment.sh kubernetes 2 16 0 raw/e3/128-w
sleep 10

bash ./scripts/experiment.sh kubernetes 1 64 0 raw/e3/128-w
sleep 10

bash ./scripts/experiment.sh kubernetes 2 32 0 raw/e3/128-w
sleep 10

bash ./scripts/experiment.sh kubernetes 2 36 0 raw/e3/128-w
sleep 10

bash ./scripts/experiment.sh kubernetes 3 24 0 raw/e3/128-w
sleep 10
