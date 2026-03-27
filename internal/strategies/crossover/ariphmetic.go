package crossover

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/core/strategy"
)

type AriphmeticCrossoverStrategy struct{}

func NewAriphmeticCrossoverStrategy() *AriphmeticCrossoverStrategy {
	return &AriphmeticCrossoverStrategy{}
}

func (s *AriphmeticCrossoverStrategy) Crossover(p1, p2 strategy.AntView) strategy.AntView {
	return ant.NewHeteroAnt(
		(p1.Alpha()+p2.Alpha())/2,
		(p1.Beta()+p2.Beta())/2,
		(p1.PheromoneMultiplier()+p2.PheromoneMultiplier())/2,
		p1.PathStrategy(),
		p1.PheromoneApplyStrategy(),
	)
}
