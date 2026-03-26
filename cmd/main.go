package main

import (
	"HeteroAntColonySystem/internal/core/colony"
	"HeteroAntColonySystem/internal/core/config"
	"HeteroAntColonySystem/internal/strategies/apply"
	"HeteroAntColonySystem/internal/strategies/path"
	"HeteroAntColonySystem/pkg/algo/aco"
	"HeteroAntColonySystem/pkg/tsplib"
	"HeteroAntColonySystem/pkg/tsplib/adapters"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("tsp/ulysses22.tsp")
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
		config.WithDefaultAlpha(2),
		config.WithDefaultBeta(1.8),
		config.WithEvaporationRate(0.2),
		config.WithInitialPheromone(1),
		config.WithPheromoneMultiplier(1),
		config.WithColonySize(500),
		config.WithGenerationCount(500),
		config.WithPathChoiceStrategy(path.NewPahtClassicStrategy()),
		config.WithPheromoneApplyingStrategy(apply.NewApplyClassicStrategy()),
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
		aco.WithAlpha(2),
		aco.WithBeta(1.8),
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
