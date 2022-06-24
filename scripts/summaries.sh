#!/usr/bin/env sh

EXPERIMENT=$1

python plotters/plot-results.py "$EXPERIMENT/throughput" "$EXPERIMENT/latency"
