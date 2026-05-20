#!/bin/bash

OUTPUT_FILE="stability.csv"
RUNS_PER_CONF=20

TSP_FILES="-f tsp/eil76.tsp"
ALPHAS="-a 0.2 -a 0.4 -a 0.6 -a 0.8 -a 1 -a 1.2 -a 1.4 -a 1.6 -a 1.8 -a 2"
BETAS="-b 0.5 -b 1 -b 1.5 -b 2 -b 2.5 -b 3 -b 3.5 -b 4 -b 4.5 -b 5"


rm -f $OUTPUT_FILE

go run cmd/stability_bench/main.go -runs $RUNS_PER_CONF -o $OUTPUT_FILE  $ALPHAS $BETAS $TSP_FILES

echo "All benchmarks completed!"s
echo "Results saved to $OUTPUT_FILE"

echo -e "\nResults summary:"
echo "Total lines: $(wc -l < $OUTPUT_FILE)"
echo -e "\nFirst 10 lines:"
head -10 $OUTPUT_FILE