#!/usr/bin/env sh

for f in logs/e1/1e6-rw/latency/*.log
do
	echo "plotting $f"
	python plotters/plot-latency.py "$f"
done

for f in logs/e1/1e6-rw/throughput/*.log
do
	echo "plotting $f"
	python plotters/plot-throughput.py "$f"
done
