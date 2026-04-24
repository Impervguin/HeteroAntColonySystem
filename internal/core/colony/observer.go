package colony

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
)

type ColonyObserverDTO struct {
	C     *HeteroAntColony
	Ants  []ant.AntView
	Pm    *pheromone.PheromoneMap
	G     *graph.Graph
	Best  *ant.HeteroAnt
	Score float64

	Generation uint
}

type ColonyObserver interface {
	Observe(dto *ColonyObserverDTO)
}
