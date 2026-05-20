package crossover

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/core/colony"
	"math"
	"math/rand/v2"
)

type BLXCrossoverStrategy struct {
	gamma float64 // параметр альфа (обычно 0.3-0.5)
}

func NewBLXCrossoverStrategy(gamma float64) *BLXCrossoverStrategy {
	return &BLXCrossoverStrategy{
		gamma: gamma,
	}
}

var _ colony.CrossoverStrategy = (*BLXCrossoverStrategy)(nil)

func (s *BLXCrossoverStrategy) Crossover(p1, p2 ant.AntView) []*ant.HeteroAnt {
	// Для каждой координаты вычисляем границы блендинга
	alpha := s.blend(p1.Alpha(), p2.Alpha())
	beta := s.blend(p1.Beta(), p2.Beta())

	return []*ant.HeteroAnt{
		ant.NewHeteroAnt(
			alpha,
			beta,
			p1.PheromoneMultiplier(),
			p1.PathStrategy(),
			p1.PheromoneApplyStrategy(),
		),
	}
}

func (s *BLXCrossoverStrategy) blend(x1, x2 float64) float64 {
	minVal := math.Min(x1, x2)
	maxVal := math.Max(x1, x2)

	interval := maxVal - minVal
	lower := minVal - s.gamma*interval
	upper := maxVal + s.gamma*interval

	return lower + rand.Float64()*(upper-lower)
}
