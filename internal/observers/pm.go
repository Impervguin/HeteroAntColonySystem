package observers

import (
	"HeteroAntColonySystem/internal/core/colony"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
)

type PheromoneMapObserver struct {
	pheromoneMaps map[uint]*pheromone.PheromoneMap
}

var _ colony.ColonyObserver = (*PheromoneMapObserver)(nil)

func (o *PheromoneMapObserver) Observe(dto *colony.ColonyObserverDTO) {
	pmCpy := pheromone.NewPheromoneMap(dto.G, 0)
	dto.Pm.ForEachEdgeRead(func(e *graph.Edge, pheromone float64) {
		pmCpy.Update(e, pheromone)
	})
	o.pheromoneMaps[dto.Generation] = pmCpy
}

func (o *PheromoneMapObserver) Map(gen uint) *pheromone.PheromoneMap {
	return o.pheromoneMaps[gen]
}
