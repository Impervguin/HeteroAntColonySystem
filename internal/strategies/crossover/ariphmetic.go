package crossover

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/core/colony"
)

type AriphmeticCrossoverStrategy struct{}

func NewAriphmeticCrossoverStrategy() *AriphmeticCrossoverStrategy {
	return &AriphmeticCrossoverStrategy{}
}

var _ colony.CrossoverStrategy = (*AriphmeticCrossoverStrategy)(nil)

func (s *AriphmeticCrossoverStrategy) Crossover(p1, p2 ant.AntView) *ant.HeteroAnt {
	return ant.NewHeteroAnt(
		(p1.Alpha()+p2.Alpha())/2,
		(p1.Beta()+p2.Beta())/2,
		(p1.PheromoneMultiplier()+p2.PheromoneMultiplier())/2,
		p1.PathStrategy(),
		p1.PheromoneApplyStrategy(),
	)
}
