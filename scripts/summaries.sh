#!/usr/bin/env sh

for f in logs/e3/*
do
	echo "plotting $f"
	python plotters/plot-results.py "$f/throughput" "$f/latency"
done
