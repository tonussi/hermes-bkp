#!/usr/bin/env sh
EXPERIMENT_NAME=$1
ROLE_LABEL=$2
shift 2

echo "label $ROLE_LABEL nodels..."
for i in $@
do
  kubectl label nodes node$i.$EXPERIMENT_NAME.scalablesmr.emulab.net role=$ROLE_LABEL
done