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
	"HeteroAntColonySystem/pkg/algo/greedy"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/tsplib"
	"HeteroAntColonySystem/pkg/tsplib/adapters"
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

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
	runsPerConfig, _ := strconv.Atoi(os.Args[2])
	tspFiles := os.Args[3:]

	writeHeaderIfNeeded(outputFile)
	configs := buildConfigs()

	for _, file := range tspFiles {
		fmt.Printf("[HACO] Processing %s\n", file)

		g := parseTSPFile(file)
		greed := greedy.NewGreedyAlgorithm(g)
		greed.Run()
		fmt.Printf("[HACO] File %s - greed score: %f\n", file, greed.Score())

		colonySize := uint(g.Len())
		initialPheromone := float64(colonySize) / greed.Score()
		parentCount := uint(0.4 * float64(colonySize))

		for configIdx, config := range configs {
			fmt.Printf("[HACO] %s - Config %d/%d: %s\n", file, configIdx+1, len(configs), config.Name)

			for run := 1; run <= runsPerConfig; run++ {
				fmt.Printf("[HACO] %s - %s - Run %d/%d\n", file, config.Name, run, runsPerConfig)

				runtime.GC()
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)
				memBefore := memStats.Alloc

				start := time.Now()
				haco := createHACO(config, colonySize, initialPheromone, parentCount)
				haco.Prepare(g)
				haco.Run()

				runtime.ReadMemStats(&memStats)
				memoryKB := (memStats.Alloc - memBefore) / 1024

				writeResult(outputFile, []string{
					file, fmt.Sprintf("%d", run), config.Name,
					fmt.Sprintf("%f", haco.Score()),
					fmt.Sprintf("%d", time.Since(start).Milliseconds()),
					fmt.Sprintf("%d", memoryKB),
				})
			}
		}
	}
	fmt.Printf("[HACO] Benchmark completed. Results saved to %s\n", outputFile)
}

func buildConfigs() []ExperimentConfig {
	mutations := []colony.MutationStrategy{
		mutation.NewUniformMutationStrategy(-0.1, 0.1),
		mutation.NewGaussMutationStrategy(0.05, 0),
	}
	selections := []colony.ParentSelectionStrategy{
		selection.NewBestSelectionStrategy(),
		selection.NewTournamentSelectionStrategy(5),
		selection.NewRouletteSelectionStrategy(),
	}
	crossovers := []colony.CrossoverStrategy{
		crossover.NewAriphmeticCrossoverStrategy(),
		crossover.NewBLXCrossoverStrategy(0.5),
		crossover.NewSBXCrossoverStrategy(2),
	}

	var configs []ExperimentConfig
	for _, mut := range mutations {
		mutName := getMutationName(mut)
		for _, sel := range selections {
			selName := getSelectionName(sel)
			for _, cro := range crossovers {
				configs = append(configs, ExperimentConfig{
					Name:              fmt.Sprintf("%s_%s_%s", mutName, selName, getCrossoverName(cro)),
					LocalOptimization: optimisation.NewNoOpLocalOptimisation(),
					SelectionStrategy: sel,
					MutationStrategy:  mut,
					CrossoverStrategy: cro,
				})
			}
		}
	}
	return configs
}

func createHACO(config ExperimentConfig, colonySize uint, initialPheromone float64, parentCount uint) *colony.HeteroAntColony {
	haco, _ := colony.NewHeteroAntColony(
		colony.WithDefaultAlpha(1),
		colony.WithDefaultBeta(3),
		colony.WithEvaporationRate(0.5),
		colony.WithInitialPheromone(initialPheromone),
		colony.WithPheromoneMultiplier(1),
		colony.WithColonySize(colonySize),
		colony.WithGenerationCount(300),
		colony.WithGenerationPeriod(5),
		colony.WithParentCount(parentCount),
		colony.WithPathChoiceStrategy(path.NewPahtClassicStrategy()),
		colony.WithPheromoneApplyingStrategy(apply.NewApplyClassicStrategy()),
		colony.WithLocalOptimisationStrategy(config.LocalOptimization),
		colony.WithCrossoverStrategy(config.CrossoverStrategy),
		colony.WithMutationStrategy(config.MutationStrategy),
		colony.WithParentSelectionStrategy(config.SelectionStrategy),
	)
	return haco
}

func parseTSPFile(file string) *graph.Graph {
	f, _ := os.Open(file)
	defer f.Close()
	parser := tsplib.NewTSPLIBParser(adapters.GetRegistry())
	g, _ := parser.Parse(f)
	return g
}

func writeResult(csvPath string, record []string) error {
	file, err := os.OpenFile(csvPath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	return writer.Write(record)
}

func writeHeaderIfNeeded(csvPath string) {
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		file, _ := os.Create(csvPath)
		defer file.Close()
		writer := csv.NewWriter(file)
		writer.Write([]string{"file", "run", "name", "score", "duration_ms", "memory_kb"})
		writer.Flush()
	}
}

func getMutationName(mut colony.MutationStrategy) string {
	switch mut.(type) {
	case *mutation.UniformMutationStrategy:
		return "uniform"
	default:
		return "gauss"
	}
}

func getSelectionName(sel colony.ParentSelectionStrategy) string {
	switch sel.(type) {
	case *selection.BestSelectionStrategy:
		return "best"
	case *selection.TournamentSelectionStrategy:
		return "tournament"
	default:
		return "roulette"
	}
}

func getCrossoverName(cro colony.CrossoverStrategy) string {
	switch cro.(type) {
	case *crossover.AriphmeticCrossoverStrategy:
		return "ariphmetic"
	case *crossover.BLXCrossoverStrategy:
		return "blx"
	default:
		return "sbx"
	}
}
