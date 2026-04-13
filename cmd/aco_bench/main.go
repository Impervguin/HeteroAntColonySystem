package main

import (
	"HeteroAntColonySystem/pkg/algo/aco"
	"HeteroAntColonySystem/pkg/tsplib"
	"HeteroAntColonySystem/pkg/tsplib/adapters"
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

type FileLock struct {
	f *os.File
}

func (l *FileLock) Lock() error {
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_EX)
}

func (l *FileLock) Unlock() error {
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}

func NewFileLock(filename string) (*FileLock, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &FileLock{f: f}, nil
}

func writeResultWithLock(csvPath string, record []string) error {
	lock, err := NewFileLock(csvPath)
	if err != nil {
		return fmt.Errorf("failed to create file lock: %v", err)
	}
	defer lock.f.Close()
	defer lock.Unlock()

	if err := lock.Lock(); err != nil {
		return fmt.Errorf("failed to lock file: %v", err)
	}

	file, err := os.OpenFile(csvPath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record: %v", err)
	}

	return nil
}

func writeHeaderIfNeeded(csvPath string) error {
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		lock, err := NewFileLock(csvPath)
		if err != nil {
			return fmt.Errorf("failed to create file lock: %v", err)
		}
		defer lock.f.Close()
		defer lock.Unlock()

		if err := lock.Lock(); err != nil {
			return fmt.Errorf("failed to lock file: %v", err)
		}

		file, err := os.Create(csvPath)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		header := []string{"file", "algorithm", "run", "score", "duration_ms", "memory_kb"}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("failed to write header: %v", err)
		}
		writer.Flush()
	}
	return nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main_aco.go <output.csv> <runs_per_file> <tsp_file1> <tsp_file2> ...")
		os.Exit(1)
	}

	outputFile := os.Args[1]
	runsPerFile, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(fmt.Sprintf("Invalid runs_per_file: %v", err))
	}
	tspFiles := os.Args[3:]

	if err := writeHeaderIfNeeded(outputFile); err != nil {
		panic(err)
	}

	for _, file := range tspFiles {
		fmt.Printf("[ACO] Processing %s\n", file)

		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", file, err)
			continue
		}

		parser := tsplib.NewTSPLIBParser(adapters.GetRegistry())
		g, err := parser.Parse(f)
		f.Close()

		if err != nil {
			fmt.Printf("Error parsing file %s: %v\n", file, err)
			continue
		}

		for run := 1; run <= runsPerFile; run++ {
			fmt.Printf("[ACO] %s - Run %d/%d\n", file, run, runsPerFile)

			var memStats runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&memStats)
			memBefore := memStats.Alloc

			start := time.Now()

			antColony, err := aco.NewAntColony(g,
				aco.WithAlpha(2),
				aco.WithBeta(1.8),
				aco.WithEvaporationRate(0.5),
				aco.WithInitialPheromone(1),
				aco.WithPheromoneMultiplier(4),
				aco.WithColonySize(500),
				aco.WithGenerationCount(500),
			)

			if err != nil {
				fmt.Printf("Error creating ACO for %s: %v\n", file, err)
				continue
			}

			antColony.Run()

			runtime.ReadMemStats(&memStats)
			memAfter := memStats.Alloc
			memoryKB := (memAfter - memBefore) / 1024

			duration := time.Since(start)

			record := []string{
				file,
				"aco",
				fmt.Sprintf("%d", run),
				fmt.Sprintf("%f", antColony.BestScore()),
				fmt.Sprintf("%d", duration.Milliseconds()),
				fmt.Sprintf("%d", memoryKB),
			}

			if err := writeResultWithLock(outputFile, record); err != nil {
				fmt.Printf("Error writing record: %v\n", err)
			}
		}
	}

	fmt.Printf("[ACO] Benchmark completed. Results saved to %s\n", outputFile)
}
