package observers

import (
	"HeteroAntColonySystem/internal/core/colony"
	"HeteroAntColonySystem/pkg/graph"
)

type BestPathObserver struct {
	bestPaths map[uint][]*graph.Vertex
}

var _ colony.ColonyObserver = (*BestPathObserver)(nil)

func (o *BestPathObserver) Observe(dto *colony.ColonyObserverDTO) {
	if dto.Best == nil {
		return
	}
	bestPath := dto.Best.Path()
	cpy := make([]*graph.Vertex, 0, len(bestPath))
	for _, v := range bestPath {
		cpy = append(cpy, v)
	}
	o.bestPaths[dto.Generation] = cpy
}

func (o *BestPathObserver) Path(gen uint) []*graph.Vertex {
	return o.bestPaths[gen]
}
