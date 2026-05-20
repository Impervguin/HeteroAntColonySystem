#!/bin/bash

OUTPUT_FILE="haco.csv"
RUNS_PER_FILE=20

# TSP_FILES="tsp/ulysses22.tsp tsp/swiss42.tsp tsp/eil51.tsp tsp/berlin52.tsp tsp/eil76.tsp"
TSP_FILES="tsp/st70.tsp"

rm -f $OUTPUT_FILE

go run cmd/haco_bench/main.go $OUTPUT_FILE $RUNS_PER_FILE $TSP_FILES

echo "All benchmarks completed!"s
echo "Results saved to $OUTPUT_FILE"

echo -e "\nResults summary:"
echo "Total lines: $(wc -l < $OUTPUT_FILE)"
echo -e "\nFirst 10 lines:"
head -10 $OUTPUT_FILE