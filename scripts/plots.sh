#!/usr/bin/env sh

EXPERIMENT=$1

for f in $EXPERIMENT/latency/*.log
do
	echo "plotting $f"
	python plotters/plot-latency.py "$f"
done

for f in $EXPERIMENT/throughput/*.log
do
	echo "plotting $f"
	python plotters/plot-throughput.py "$f"
done
