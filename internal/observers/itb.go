package observers

import (
	"HeteroAntColonySystem/internal/core/colony"
	"math"
)

type IterationsToBestObserver struct {
	bestScore float64
	iteration uint
}

func NewIterationsToBestObserver() *IterationsToBestObserver {
	return &IterationsToBestObserver{
		bestScore: math.Inf(1),
		iteration: 0,
	}
}

var _ colony.ColonyObserver = &IterationsToBestObserver{}

func (o *IterationsToBestObserver) Observe(dto *colony.ColonyObserverDTO) {
	if dto.Best != nil && dto.Best.Score() < o.bestScore {
		o.bestScore = dto.Best.Score()
		o.iteration = dto.Generation
	}
}

func (o *IterationsToBestObserver) IterationsToBest() uint {
	return o.iteration
}
