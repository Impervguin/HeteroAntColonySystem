package main

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/core/colony"
	"HeteroAntColonySystem/internal/observers"
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
	var genCounts generationFlags
	var tspFiles fileFlags
	flag.Var(&genCounts, "g", "Generation count (can be specified multiple times)")
	flag.Var(&tspFiles, "f", "TSP file (can be specified multiple times)")

	flag.Parse()

	if len(genCounts) == 0 {
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
	fmt.Printf("Generation counts: %v\n", genCounts)
	fmt.Printf("TSP files: %v\n", tspFiles)

	writeHeaderIfNeeded(*outputFile)

	for _, file := range tspFiles {
		fmt.Printf("[CONVERGENCE] Processing %s\n", file)

		g := parseTSPFile(file)
		greed := greedy.NewGreedyAlgorithm(g)
		greed.Run()
		fmt.Printf("[CONVERGENCE] File %s - greed score: %f\n", file, greed.Score())

		colonySize := uint(g.Len())
		initialPheromone := float64(colonySize) / greed.Score()
		parentCount := uint(0.4 * float64(colonySize))

		fmt.Printf("[CONVERGENCE] Haco no local optimisation\n")
		for _, genCount := range genCounts {
			for run := 1; run <= *runsPerConfig; run++ {
				fmt.Printf("[CONVERGENCE] %s - HACO no local/%d - Run %d/%d\n", file, genCount, run, *runsPerConfig)

				runtime.GC()
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)
				memBefore := memStats.Alloc

				start := time.Now()
				haco, itb, auc := createHACO(genCount, colonySize, initialPheromone, parentCount, optimisation.NewNoOpLocalOptimisation())
				haco.Prepare(g)
				haco.Run()
				end := time.Now()

				runtime.ReadMemStats(&memStats)
				memoryKB := (memStats.Alloc - memBefore) / 1024

				writeResult(*outputFile, []string{
					file, fmt.Sprintf("%d", genCount), fmt.Sprintf("%d", run),
					"haco_no_local", fmt.Sprintf("%f", haco.Score()),
					fmt.Sprintf("%f", auc.AreaUnderCurve()),
					fmt.Sprintf("%d", itb.IterationsToBest()),
					fmt.Sprintf("%d", end.UnixMilli()-start.UnixMilli()),
					fmt.Sprintf("%d", memoryKB),
				})
			}
		}

		fmt.Printf("[CONVERGENCE] Haco with local optimisation\n")
		for _, genCount := range genCounts {
			for run := 1; run <= *runsPerConfig; run++ {
				fmt.Printf("[CONVERGENCE] %s - HACO with local/%d - Run %d/%d\n", file, genCount, run, *runsPerConfig)

				runtime.GC()
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)
				memBefore := memStats.Alloc

				start := time.Now()
				haco, itb, auc := createHACO(genCount, colonySize, initialPheromone, parentCount, optimisation.NewTwoOptLocalOptimisation())
				haco.Prepare(g)
				haco.Run()
				end := time.Now()

				runtime.ReadMemStats(&memStats)
				memoryKB := (memStats.Alloc - memBefore) / 1024

				writeResult(*outputFile, []string{
					file, fmt.Sprintf("%d", genCount), fmt.Sprintf("%d", run),
					"haco_with_local", fmt.Sprintf("%f", haco.Score()),
					fmt.Sprintf("%f", auc.AreaUnderCurve()),
					fmt.Sprintf("%d", itb.IterationsToBest()),
					fmt.Sprintf("%d", end.UnixMilli()-start.UnixMilli()),
					fmt.Sprintf("%d", memoryKB),
				})
			}
		}

		fmt.Printf("[CONVERGENCE] ACO\n")
		for _, genCount := range genCounts {
			for run := 1; run <= *runsPerConfig; run++ {
				fmt.Printf("[CONVERGENCE] %s - ACO/%d - Run %d/%d\n", file, genCount, run, *runsPerConfig)

				runtime.GC()
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)
				memBefore := memStats.Alloc

				start := time.Now()
				aco, itb, auc := CreateACO(g, genCount, colonySize, initialPheromone, parentCount)
				aco.Run()
				end := time.Now()

				runtime.ReadMemStats(&memStats)
				memoryKB := (memStats.Alloc - memBefore) / 1024

				writeResult(*outputFile, []string{
					file, fmt.Sprintf("%d", genCount), fmt.Sprintf("%d", run),
					"aco", fmt.Sprintf("%f", aco.BestScore()),
					fmt.Sprintf("%f", auc.AreaUnderCurve()),
					fmt.Sprintf("%d", itb.IterationsToBest()),
					fmt.Sprintf("%d", end.UnixMilli()-start.UnixMilli()),
					fmt.Sprintf("%d", memoryKB),
				})
			}
		}
	}
}

type generationFlags []uint

func (g *generationFlags) String() string {
	return fmt.Sprintf("%v", *g)
}

func (g *generationFlags) Set(value string) error {
	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid generation count: %s", value)
	}
	*g = append(*g, uint(val))
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

func createHACO(genCount uint, colonySize uint, initialPheromone float64, parentCount uint, localOpt ant.LocalOptimisationStrategy) (*colony.HeteroAntColony, *observers.IterationsToBestObserver, *observers.AreaUnderCurveObserver) {
	itb := observers.NewIterationsToBestObserver()
	auc := observers.NewAreaUnderCurveObserver()
	haco, _ := colony.NewHeteroAntColony(
		colony.WithDefaultAlpha(1),
		colony.WithDefaultBeta(5),
		colony.WithEvaporationRate(0.5),
		colony.WithInitialPheromone(initialPheromone),
		colony.WithPheromoneMultiplier(1),
		colony.WithColonySize(colonySize),
		colony.WithGenerationCount(genCount),
		colony.WithGenerationPeriod(5),
		colony.WithParentCount(parentCount),
		colony.WithPathChoiceStrategy(path.NewPahtClassicStrategy()),
		colony.WithPheromoneApplyingStrategy(apply.NewApplyClassicStrategy()),
		colony.WithLocalOptimisationStrategy(localOpt),
		colony.WithCrossoverStrategy(crossover.NewBLXCrossoverStrategy(0.5)),
		colony.WithMutationStrategy(mutation.NewGaussMutationStrategy(0.05, 0)),
		colony.WithParentSelectionStrategy(selection.NewRouletteSelectionStrategy()),
		colony.WithColonyObserver(itb),
		colony.WithColonyObserver(auc),
	)
	return haco, itb, auc
}

func CreateACO(gr *graph.Graph, genCount uint, colonySize uint, initialPheromone float64, parentCount uint) (*aco.AntColony, *aco.IterationsToBestObserver, *aco.AreaUnderCurveObserver) {
	itb := aco.NewIterationsToBestObserver()
	auc := aco.NewAreaUnderCurveObserver()
	col, _ := aco.NewAntColony(
		gr,
		aco.WithColonySize(colonySize),
		aco.WithGenerationCount(genCount),
		aco.WithAlpha(1),
		aco.WithBeta(5),
		aco.WithPheromoneMultiplier(1),
		aco.WithEvaporationRate(0.5),
		aco.WithInitialPheromone(initialPheromone),
		aco.WithObservers(itb, auc),
	)
	return col, itb, auc
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
		writer.Write([]string{"file", "gensize", "run", "name", "score", "auc", "itb", "duration_ms", "memory_kb"})
		writer.Flush()
	}
}
