package observers

import (
	"HeteroAntColonySystem/internal/core/colony"
	"HeteroAntColonySystem/pkg/graph"
)

type BestPathObserver struct {
	bestPaths  map[uint][]*graph.Vertex
	bestScores map[uint]float64
}

func NewBestPathObserver(gen, vertexCount uint) *BestPathObserver {
	o := &BestPathObserver{
		bestPaths:  make(map[uint][]*graph.Vertex, gen),
		bestScores: make(map[uint]float64, gen),
	}

	// pre-allocate best paths
	for i := uint(0); i < gen; i++ {
		o.bestPaths[i] = make([]*graph.Vertex, 0, vertexCount)
	}

	return o
}

var _ colony.ColonyObserver = (*BestPathObserver)(nil)

func (o *BestPathObserver) Observe(dto *colony.ColonyObserverDTO) {
	if dto.Best == nil {
		return
	}
	// find best in current generation
	bestAnt := dto.Ants[0]
	for _, ant := range dto.Ants {
		if ant.Score() < bestAnt.Score() {
			bestAnt = ant
		}
	}

	bestPath := bestAnt.Path()
	for _, v := range bestPath {
		o.bestPaths[dto.Generation] = append(o.bestPaths[dto.Generation], v)
	}
	o.bestScores[dto.Generation] = bestAnt.Score()
}

func (o *BestPathObserver) Path(gen uint) ([]*graph.Vertex, float64) {
	return o.bestPaths[gen], o.bestScores[gen]
}
