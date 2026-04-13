package main

import (
	"HeteroAntColonySystem/internal/core/colony"
	"HeteroAntColonySystem/internal/core/config"
	"HeteroAntColonySystem/internal/strategies/apply"
	"HeteroAntColonySystem/internal/strategies/crossover"
	"HeteroAntColonySystem/internal/strategies/mutation"
	"HeteroAntColonySystem/internal/strategies/path"
	"HeteroAntColonySystem/internal/strategies/selection"
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
	// Создаем блокировку для файла
	lock, err := NewFileLock(csvPath)
	if err != nil {
		return fmt.Errorf("failed to create file lock: %v", err)
	}
	defer lock.f.Close()
	defer lock.Unlock()

	// Блокируем файл для эксклюзивного доступа
	if err := lock.Lock(); err != nil {
		return fmt.Errorf("failed to lock file: %v", err)
	}

	// Открываем файл для добавления
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
	// Проверяем существует ли файл
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		// Файл не существует, создаем и пишем заголовок
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
		fmt.Println("Usage: go run main_haco.go <output.csv> <runs_per_file> <tsp_file1> <tsp_file2> ...")
		os.Exit(1)
	}

	outputFile := os.Args[1]
	runsPerFile, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(fmt.Sprintf("Invalid runs_per_file: %v", err))
	}
	tspFiles := os.Args[3:]

	// Записываем заголовок если нужно
	if err := writeHeaderIfNeeded(outputFile); err != nil {
		panic(err)
	}

	// Process each file sequentially
	for _, file := range tspFiles {
		fmt.Printf("[HACO] Processing %s\n", file)

		// Parse TSP file once per file
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
			fmt.Printf("[HACO] %s - Run %d/%d\n", file, run, runsPerFile)

			// Memory measurement before
			var memStats runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&memStats)
			memBefore := memStats.Alloc

			start := time.Now()

			// Configure HACO
			haco, err := colony.NewHeteroAntColony(
				config.WithDefaultAlpha(2),
				config.WithDefaultBeta(1.8),
				config.WithEvaporationRate(0.5),
				config.WithInitialPheromone(1),
				config.WithPheromoneMultiplier(4),
				config.WithColonySize(400),
				config.WithGenerationCount(400),
				config.WithGenerationPeriod(10),
				config.WithParentCount(20),
				config.WithPathChoiceStrategy(path.NewPahtClassicStrategy()),
				config.WithPheromoneApplyingStrategy(apply.NewApplyClassicStrategy()),
				config.WithCrossoverStrategy(crossover.NewAriphmeticCrossoverStrategy()),
				config.WithMutationStrategy(mutation.NewUniformMutationStrategy(-0.2, 0.2)),
				config.WithParentSelectionStrategy(selection.NewBestSelectionStrategy()),
			)

			if err != nil {
				fmt.Printf("Error creating HACO for %s: %v\n", file, err)
				continue
			}

			haco.Prepare(g)
			haco.Run()

			// Memory measurement after
			runtime.ReadMemStats(&memStats)
			memAfter := memStats.Alloc
			memoryUsed := memAfter - memBefore
			memoryKB := memoryUsed / 1024

			duration := time.Since(start)

			// Write result to CSV with file locking
			record := []string{
				file,
				"haco",
				fmt.Sprintf("%d", run),
				fmt.Sprintf("%f", haco.Score()),
				fmt.Sprintf("%d", duration.Milliseconds()),
				fmt.Sprintf("%d", memoryKB),
			}

			if err := writeResultWithLock(outputFile, record); err != nil {
				fmt.Printf("Error writing record: %v\n", err)
			}
		}
	}

	fmt.Printf("[HACO] Benchmark completed. Results saved to %s\n", outputFile)
}
