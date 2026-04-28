package observers

import (
	"HeteroAntColonySystem/internal/core/colony"
	"HeteroAntColonySystem/pkg/graph"
)

type BestPathObserver struct {
	bestPaths  map[uint][]*graph.Vertex
	bestScores map[uint]float64
}

func NewBestPathObserver() *BestPathObserver {
	return &BestPathObserver{
		bestPaths:  make(map[uint][]*graph.Vertex),
		bestScores: make(map[uint]float64),
	}
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
	o.bestScores[dto.Generation] = dto.Best.Score()
}

func (o *BestPathObserver) Path(gen uint) ([]*graph.Vertex, float64) {
	return o.bestPaths[gen], o.bestScores[gen]
}
