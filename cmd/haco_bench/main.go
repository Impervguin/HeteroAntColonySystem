package main

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/core/colony"
	"HeteroAntColonySystem/internal/strategies/apply"
	"HeteroAntColonySystem/internal/strategies/crossover"
	"HeteroAntColonySystem/internal/strategies/mutation"
	"HeteroAntColonySystem/internal/strategies/optimisation"
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
		header := []string{
			"file", "run", "score", "duration_ms", "memory_kb",
			"local_optimization", "selection_strategy", "mutation_strategy",
			"mutation_rate", "crossover_strategy",
		}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("failed to write header: %v", err)
		}
		writer.Flush()
	}
	return nil
}

// Configuration struct for experiment
type ExperimentConfig struct {
	Name              string
	LocalOptimization ant.LocalOptimisationStrategy
	SelectionStrategy colony.ParentSelectionStrategy
	MutationStrategy  colony.MutationStrategy
	CrossoverStrategy colony.CrossoverStrategy
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main_haco.go <output.csv> <runs_per_config> <tsp_file1> <tsp_file2> ...")
		os.Exit(1)
	}

	outputFile := os.Args[1]
	runsPerConfig, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(fmt.Sprintf("Invalid runs_per_config: %v", err))
	}
	tspFiles := os.Args[3:]

	if err := writeHeaderIfNeeded(outputFile); err != nil {
		panic(err)
	}

	// Define experiment configurations
	configs := []ExperimentConfig{
		{
			Name:              "best_uniform_noopt",
			LocalOptimization: optimisation.NewNoOpLocalOptimisation(),
			SelectionStrategy: selection.NewBestSelectionStrategy(),
			MutationStrategy:  mutation.NewUniformMutationStrategy(-0.2, 0.2),
			CrossoverStrategy: crossover.NewAriphmeticCrossoverStrategy(),
		},
		{
			Name:              "best_uniform_2opt",
			LocalOptimization: optimisation.NewTwoOptLocalOptimisation(),
			SelectionStrategy: selection.NewBestSelectionStrategy(),
			MutationStrategy:  mutation.NewUniformMutationStrategy(-0.2, 0.2),
			CrossoverStrategy: crossover.NewAriphmeticCrossoverStrategy(),
		},
		{
			Name:              "best_gauss_noopt",
			LocalOptimization: optimisation.NewNoOpLocalOptimisation(),
			SelectionStrategy: selection.NewBestSelectionStrategy(),
			MutationStrategy:  mutation.NewGaussMutationStrategy(0.2, 0),
			CrossoverStrategy: crossover.NewAriphmeticCrossoverStrategy(),
		},
		{
			Name:              "best_gauss_2opt",
			LocalOptimization: optimisation.NewTwoOptLocalOptimisation(),
			SelectionStrategy: selection.NewBestSelectionStrategy(),
			MutationStrategy:  mutation.NewGaussMutationStrategy(0.2, 0),
			CrossoverStrategy: crossover.NewAriphmeticCrossoverStrategy(),
		},
		{
			Name:              "tournament_uniform_noopt",
			LocalOptimization: optimisation.NewNoOpLocalOptimisation(),
			SelectionStrategy: selection.NewTournamentSelectionStrategy(3),
			MutationStrategy:  mutation.NewUniformMutationStrategy(-0.2, 0.2),
			CrossoverStrategy: crossover.NewAriphmeticCrossoverStrategy(),
		},
		{
			Name:              "tournament_uniform_2opt",
			LocalOptimization: optimisation.NewTwoOptLocalOptimisation(),
			SelectionStrategy: selection.NewTournamentSelectionStrategy(3),
			MutationStrategy:  mutation.NewUniformMutationStrategy(-0.2, 0.2),
			CrossoverStrategy: crossover.NewAriphmeticCrossoverStrategy(),
		},
		{
			Name:              "tournament_gauss_noopt",
			LocalOptimization: optimisation.NewNoOpLocalOptimisation(),
			SelectionStrategy: selection.NewTournamentSelectionStrategy(3),
			MutationStrategy:  mutation.NewGaussMutationStrategy(0.2, 0),
			CrossoverStrategy: crossover.NewAriphmeticCrossoverStrategy(),
		},
		{
			Name:              "tournament_gauss_2opt",
			LocalOptimization: optimisation.NewTwoOptLocalOptimisation(),
			SelectionStrategy: selection.NewTournamentSelectionStrategy(3),
			MutationStrategy:  mutation.NewGaussMutationStrategy(0.2, 0),
			CrossoverStrategy: crossover.NewAriphmeticCrossoverStrategy(),
		},
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

		// Run each configuration
		for configIdx, config := range configs {
			fmt.Printf("[HACO] %s - Config %d/%d: %s\n", file, configIdx+1, len(configs), config.Name)

			for run := 1; run <= runsPerConfig; run++ {
				fmt.Printf("[HACO] %s - %s - Run %d/%d\n", file, config.Name, run, runsPerConfig)

				// Memory measurement before
				var memStats runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&memStats)
				memBefore := memStats.Alloc

				start := time.Now()

				// Build colony options
				options := []colony.HeteroAntColonyOption{
					colony.WithDefaultAlpha(1),
					colony.WithDefaultBeta(1),
					colony.WithEvaporationRate(0.2),
					colony.WithInitialPheromone(1),
					colony.WithPheromoneMultiplier(2),
					colony.WithColonySize(300),
					colony.WithGenerationCount(300),
					colony.WithGenerationPeriod(10),
					colony.WithParentCount(20),
					colony.WithPathChoiceStrategy(path.NewPahtClassicStrategy()),
					colony.WithPheromoneApplyingStrategy(apply.NewApplyClassicStrategy()),
					colony.WithLocalOptimisationStrategy(config.LocalOptimization),
					colony.WithCrossoverStrategy(config.CrossoverStrategy),
					colony.WithMutationStrategy(config.MutationStrategy),
					colony.WithParentSelectionStrategy(config.SelectionStrategy),
				}

				haco, err := colony.NewHeteroAntColony(options...)
				if err != nil {
					fmt.Printf("Error creating HACO for %s config %s: %v\n", file, config.Name, err)
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

				// Write result to CSV
				record := []string{
					file,
					config.Name,
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
	}

	fmt.Printf("[HACO] Benchmark completed. Results saved to %s\n", outputFile)
}
