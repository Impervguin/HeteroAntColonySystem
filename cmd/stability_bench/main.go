package main

import (
	"HeteroAntColonySystem/internal/core/colony"
	"HeteroAntColonySystem/internal/strategies/apply"
	"HeteroAntColonySystem/internal/strategies/crossover"
	"HeteroAntColonySystem/internal/strategies/mutation"
	"HeteroAntColonySystem/internal/strategies/optimisation"
	"HeteroAntColonySystem/internal/strategies/path"
	"HeteroAntColonySystem/internal/strategies/selection"
	"HeteroAntColonySystem/pkg/algo/aco"
	"HeteroAntColonySystem/pkg/algo/greedy"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/tsplib"
	"HeteroAntColonySystem/pkg/tsplib/adapters"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

func main() {
	// Определение флагов
	outputFile := flag.String("o", "results.csv", "Output CSV file")
	runsPerConfig := flag.Int("runs", 10, "Number of runs per configuration")

	// Флаг для генераций (можно указывать несколько раз)
	var alphas paramFlags
	var betas paramFlags
	var tspFiles fileFlags
	flag.Var(&alphas, "a", "Alpha parameter (can be specified multiple times)")
	flag.Var(&betas, "b", "Beta parameter (can be specified multiple times)")
	flag.Var(&tspFiles, "f", "TSP file (can be specified multiple times)")

	flag.Parse()

	if len(alphas) == 0 {
		fmt.Println("Error: at least one generation count required (-g flag)")
		flag.Usage()
		os.Exit(1)
	}

	if len(betas) == 0 {
		fmt.Println("Error: at least one generation count required (-g flag)")
		flag.Usage()
		os.Exit(1)
	}

	if len(tspFiles) == 0 {
		fmt.Println("Error: at least one TSP file required")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Output file: %s\n", *outputFile)
	fmt.Printf("Runs per config: %d\n", *runsPerConfig)
	fmt.Printf("Alpha parameters: %v\n", alphas)
	fmt.Printf("Beta parameters: %v\n", betas)
	fmt.Printf("TSP files: %v\n", tspFiles)

	writeHeaderIfNeeded(*outputFile)

	for _, file := range tspFiles {
		fmt.Printf("[STABILITY] Processing %s\n", file)

		g := parseTSPFile(file)
		greed := greedy.NewGreedyAlgorithm(g)
		greed.Run()
		fmt.Printf("[STABILITY] File %s - greed score: %f\n", file, greed.Score())

		colonySize := uint(g.Len())
		initialPheromone := float64(colonySize) / greed.Score()
		parentCount := uint(0.4 * float64(colonySize))
		genCount := 500

		for _, alpha := range alphas {
			for _, beta := range betas {
				fmt.Printf("[STABILITY] HACO\n")
				for run := 1; run <= *runsPerConfig; run++ {
					fmt.Printf("[STABILITY] %s - HACO  %f/%f - Run %d/%d\n", file, alpha, beta, run, *runsPerConfig)

					runtime.GC()
					var memStats runtime.MemStats
					runtime.ReadMemStats(&memStats)
					memBefore := memStats.Alloc

					start := time.Now()
					haco := createHACO(uint(genCount), colonySize, initialPheromone, parentCount, alpha, beta)
					haco.Prepare(g)
					haco.Run()
					end := time.Now()

					runtime.ReadMemStats(&memStats)
					memoryKB := (memStats.Alloc - memBefore) / 1024

					writeResult(*outputFile, []string{
						file, fmt.Sprintf("%f", alpha), fmt.Sprintf("%f", beta), fmt.Sprintf("%d", run),
						"haco", fmt.Sprintf("%f", haco.Score()),
						fmt.Sprintf("%d", end.UnixMilli()-start.UnixMilli()),
						fmt.Sprintf("%d", memoryKB),
					})
				}

				fmt.Printf("[STABILITY] ACO\n")
				for run := 1; run <= *runsPerConfig; run++ {
					fmt.Printf("[STABILITY] %s - ACO  %f/%f - Run %d/%d\n", file, alpha, beta, run, *runsPerConfig)

					runtime.GC()
					var memStats runtime.MemStats
					runtime.ReadMemStats(&memStats)
					memBefore := memStats.Alloc

					start := time.Now()
					aco := CreateACO(g, uint(genCount), colonySize, initialPheromone, parentCount, alpha, beta)
					aco.Run()
					end := time.Now()

					runtime.ReadMemStats(&memStats)
					memoryKB := (memStats.Alloc - memBefore) / 1024

					writeResult(*outputFile, []string{
						file, fmt.Sprintf("%f", alpha), fmt.Sprintf("%f", beta), fmt.Sprintf("%d", run),
						"aco", fmt.Sprintf("%f", aco.BestScore()),
						fmt.Sprintf("%d", end.UnixMilli()-start.UnixMilli()),
						fmt.Sprintf("%d", memoryKB),
					})
				}
			}
		}
	}
}

type paramFlags []float64

func (g *paramFlags) String() string {
	return fmt.Sprintf("%v", *g)
}

func (g *paramFlags) Set(value string) error {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("invalid generation count: %s", value)
	}
	*g = append(*g, val)
	return nil
}

type fileFlags []string

func (f *fileFlags) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *fileFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func createHACO(genCount uint, colonySize uint, initialPheromone float64, parentCount uint, alpha, beta float64) *colony.HeteroAntColony {
	haco, _ := colony.NewHeteroAntColony(
		colony.WithDefaultAlpha(alpha),
		colony.WithDefaultBeta(beta),
		colony.WithEvaporationRate(0.5),
		colony.WithInitialPheromone(initialPheromone),
		colony.WithPheromoneMultiplier(1),
		colony.WithColonySize(colonySize),
		colony.WithGenerationCount(genCount),
		colony.WithGenerationPeriod(5),
		colony.WithParentCount(parentCount),
		colony.WithPathChoiceStrategy(path.NewPahtClassicStrategy()),
		colony.WithPheromoneApplyingStrategy(apply.NewApplyClassicStrategy()),
		colony.WithLocalOptimisationStrategy(optimisation.NewNoOpLocalOptimisation()),
		colony.WithCrossoverStrategy(crossover.NewBLXCrossoverStrategy(0.5)),
		colony.WithMutationStrategy(mutation.NewGaussMutationStrategy(0.05, 0)),
		colony.WithParentSelectionStrategy(selection.NewRouletteSelectionStrategy()),
	)
	return haco
}

func CreateACO(gr *graph.Graph, genCount uint, colonySize uint, initialPheromone float64, parentCount uint, alpha, beta float64) *aco.AntColony {
	col, _ := aco.NewAntColony(
		gr,
		aco.WithColonySize(colonySize),
		aco.WithGenerationCount(genCount),
		aco.WithAlpha(alpha),
		aco.WithBeta(beta),
		aco.WithPheromoneMultiplier(1),
		aco.WithEvaporationRate(0.5),
		aco.WithInitialPheromone(initialPheromone),
	)
	return col
}

func parseTSPFile(file string) *graph.Graph {
	f, _ := os.Open(file)
	defer f.Close()
	parser := tsplib.NewTSPLIBParser(adapters.GetRegistry())
	g, err := parser.Parse(f)
	if err != nil {
		panic(err)
	}
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
		writer.Write([]string{"file", "alpha", "beta", "run", "name", "score", "duration_ms", "memory_kb"})
		writer.Flush()
	}
}
