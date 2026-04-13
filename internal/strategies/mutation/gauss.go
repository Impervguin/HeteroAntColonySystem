package mutation

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/core/colony"
	"math/rand/v2"
)

type GaussMutationStrategy struct {
	sigma float64
	mu    float64
}

var _ colony.MutationStrategy = (*GaussMutationStrategy)(nil)

func normalRand(sigma, mu float64) float64 {
	return rand.NormFloat64()*sigma + mu
}

func NewGaussMutationStrategy(sigma, mu float64) *GaussMutationStrategy {
	return &GaussMutationStrategy{
		sigma: sigma,
		mu:    mu,
	}
}

func (s *GaussMutationStrategy) Mutate(a ant.AntView) *ant.HeteroAnt {
	alpha := normalRand(s.sigma, s.mu)
	beta := normalRand(s.sigma, s.mu)

	return ant.NewHeteroAnt(
		a.Alpha()+alpha,
		a.Beta()+beta,
		a.PheromoneMultiplier(),
		a.PathStrategy(),
		a.PheromoneApplyStrategy(),
	)
}
