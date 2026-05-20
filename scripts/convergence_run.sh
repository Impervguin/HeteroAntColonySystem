#!/bin/bash

OUTPUT_FILE="convergence_big.csv"
RUNS_PER_CONF=20

# TSP_FILES="-f tsp/ulysses22.tsp -f tsp/swiss42.tsp -f tsp/eil51.tsp -f tsp/berlin52.tsp -f tsp/eil76.tsp"
# TSP_FILES="-f tsp/swiss42.tsp -f tsp/eil51.tsp -f tsp/berlin52.tsp -f tsp/eil76.tsp -f eil101.tsp -f rd100.tsp"
# TSP_FILES="-f tsp/eil101.tsp -f tsp/rd100.tsp"
TSP_FILES="-f tsp/tsp225.tsp -f tsp/rd400.tsp"
GEN_SIZES="-g 1000"

rm -f $OUTPUT_FILE

go run cmd/convergence_bench/main.go -runs $RUNS_PER_CONF -o $OUTPUT_FILE  $GEN_SIZES $TSP_FILES

echo "All benchmarks completed!"s
echo "Results saved to $OUTPUT_FILE"

echo -e "\nResults summary:"
echo "Total lines: $(wc -l < $OUTPUT_FILE)"
echo -e "\nFirst 10 lines:"
head -10 $OUTPUT_FILE