#!/bin/bash

# Общий файл для результатов
OUTPUT_FILE="results.csv"
RUNS_PER_FILE=10

# Список файлов для тестирования
TSP_FILES="tsp/ulysses16.tsp tsp/ulysses22.tsp tsp/eil51.tsp tsp/eil76.tsp tsp/berlin52.tsp"

# Удаляем старый файл результатов если существует
rm -f $OUTPUT_FILE

# Запускаем три алгоритма параллельно с записью в один файл
go run cmd/haco_bench/main.go $OUTPUT_FILE $RUNS_PER_FILE $TSP_FILES &
go run cmd/aco_bench/main.go $OUTPUT_FILE $RUNS_PER_FILE $TSP_FILES &
# go run cmd/greedy_bench/main.go $OUTPUT_FILE $RUNS_PER_FILE $TSP_FILES &

# Ждем завершения всех процессов
wait

echo "All benchmarks completed!"
echo "Results saved to $OUTPUT_FILE"

# Показываем статистику
echo -e "\nResults summary:"
echo "Total lines: $(wc -l < $OUTPUT_FILE)"
echo -e "\nFirst 10 lines:"
head -10 $OUTPUT_FILE