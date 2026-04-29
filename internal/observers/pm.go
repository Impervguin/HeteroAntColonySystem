package observers

import (
	"HeteroAntColonySystem/internal/core/colony"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
)

type PheromoneMapObserver struct {
	pheromoneMaps map[uint]*pheromone.PheromoneMap
}

func NewPheromoneMapObserver(gens uint, g *graph.Graph) *PheromoneMapObserver {
	o := &PheromoneMapObserver{
		pheromoneMaps: make(map[uint]*pheromone.PheromoneMap, gens),
	}

	// pre-allocate pheromone maps
	for i := uint(0); i < gens; i++ {
		o.pheromoneMaps[i] = pheromone.NewPheromoneMap(g, 0)
	}
	return o
}

var _ colony.ColonyObserver = (*PheromoneMapObserver)(nil)

func (o *PheromoneMapObserver) Observe(dto *colony.ColonyObserverDTO) {
	pm := o.pheromoneMaps[dto.Generation]
	dto.Pm.ForEachEdgeRead(func(e *graph.Edge, pheromone float64) {
		pm.Update(e, pheromone)
	})
}

func (o *PheromoneMapObserver) Map(gen uint) *pheromone.PheromoneMap {
	return o.pheromoneMaps[gen]
}
