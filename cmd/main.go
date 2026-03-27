package main

import (
	"HeteroAntColonySystem/internal/core/colony"
	"HeteroAntColonySystem/internal/core/config"
	"HeteroAntColonySystem/internal/strategies/apply"
	"HeteroAntColonySystem/internal/strategies/crossover"
	"HeteroAntColonySystem/internal/strategies/mutation"
	"HeteroAntColonySystem/internal/strategies/path"
	"HeteroAntColonySystem/internal/strategies/selection"
	"HeteroAntColonySystem/pkg/algo/aco"
	"HeteroAntColonySystem/pkg/tsplib"
	"HeteroAntColonySystem/pkg/tsplib/adapters"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("tsp/tsp225.tsp")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	parser := tsplib.NewTSPLIBParser(adapters.GetRegistry())
	g, err := parser.Parse(f)
	if err != nil {
		panic(err)
	}

	haco, err := colony.NewHeteroAntColony(
		config.WithDefaultAlpha(1),
		config.WithDefaultBeta(1),
		config.WithEvaporationRate(0.2),
		config.WithInitialPheromone(1),
		config.WithPheromoneMultiplier(1),
		config.WithColonySize(500),
		config.WithGenerationCount(500),
		config.WithGenerationPeriod(10),
		config.WithParentCount(20),
		config.WithPathChoiceStrategy(path.NewPahtClassicStrategy()),
		config.WithPheromoneApplyingStrategy(apply.NewApplyClassicStrategy()),
		config.WithCrossoverStrategy(crossover.NewAriphmeticCrossoverStrategy()),
		config.WithMutationStrategy(mutation.NewUniformMutationStrategy(-0.2, 0.2)),
		config.WithParentSelectionStrategy(selection.NewBestSelectionStrategy()),
	)

	if err != nil {
		panic(err)
	}

	haco.Prepare(g)
	haco.Run()

	fmt.Println("HACO best path:", haco.BestPath())
	fmt.Println("HACO best score:", haco.Score())

	// ACO
	aco, err := aco.NewAntColony(g,
		aco.WithAlpha(1),
		aco.WithBeta(1),
		aco.WithEvaporationRate(0.2),
		aco.WithInitialPheromone(1),
		aco.WithPheromoneMultiplier(1),
		aco.WithColonySize(500),
		aco.WithGenerationCount(500),
	)

	if err != nil {
		panic(err)
	}

	aco.Run()
	fmt.Println("ACO best solution:", aco.BestTour())
	fmt.Println("ACO score:", aco.BestScore())
}
