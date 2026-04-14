#!/bin/bash

# Общий файл для результатов
OUTPUT_FILE="results.csv"
RUNS_PER_FILE=10

# Список файлов для тестирования
TSP_FILES="tsp/ulysses22.tsp tsp/eil51.tsp tsp/eil76.tsp tsp/berlin52.tsp tsp/rd100.tsp tsp/tsp225.tsp"

# Удаляем старый файл результатов если существует
rm -f $OUTPUT_FILE

go run cmd/haco_bench/main.go $OUTPUT_FILE $RUNS_PER_FILE $TSP_FILES
go run cmd/aco_bench/main.go $OUTPUT_FILE $RUNS_PER_FILE $TSP_FILES

echo "All benchmarks completed!"s
echo "Results saved to $OUTPUT_FILE"

# Показываем статистику
echo -e "\nResults summary:"
echo "Total lines: $(wc -l < $OUTPUT_FILE)"
echo -e "\nFirst 10 lines:"
head -10 $OUTPUT_FILE