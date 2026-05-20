package crossover

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/core/colony"
	"math"
	"math/rand/v2"
)

type SBXCrossoverStrategy struct {
	eta float64
}

func NewSBXCrossoverStrategy(eta float64) *SBXCrossoverStrategy {
	return &SBXCrossoverStrategy{
		eta: eta,
	}
}

var _ colony.CrossoverStrategy = (*SBXCrossoverStrategy)(nil)

func (s *SBXCrossoverStrategy) Crossover(p1, p2 ant.AntView) []*ant.HeteroAnt {
	alpha1, alpha2 := s.sbxCross(p1.Alpha(), p2.Alpha())
	beta1, beta2 := s.sbxCross(p1.Beta(), p2.Beta())

	return []*ant.HeteroAnt{
		ant.NewHeteroAnt(
			alpha1,
			beta1,
			p1.PheromoneMultiplier(),
			p1.PathStrategy(),
			p1.PheromoneApplyStrategy(),
		),
		ant.NewHeteroAnt(
			alpha2,
			beta2,
			p1.PheromoneMultiplier(),
			p1.PathStrategy(),
			p1.PheromoneApplyStrategy(),
		),
	}
}

func (s *SBXCrossoverStrategy) sbxCross(x1, x2 float64) (float64, float64) {
	if rand.Float64() > 0.5 {
		return x1, x2
	}

	if math.Abs(x1-x2) < 1e-10 {
		return x1, x2
	}

	if x1 > x2 {
		x1, x2 = x2, x1
	}

	var gamma float64
	u := rand.Float64()

	if u <= 0.5 {
		gamma = math.Pow(2*u, 1/(s.eta+1))
	} else {
		gamma = math.Pow(1/(2*(1-u)), 1/(s.eta+1))
	}

	c1 := 0.5 * ((1+gamma)*x1 + (1-gamma)*x2)
	c2 := 0.5 * ((1-gamma)*x1 + (1+gamma)*x2)

	return c1, c2
}
