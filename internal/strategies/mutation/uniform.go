package mutation

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/core/colony"
	"math/rand/v2"
)

type UniformMutationStrategy struct {
	l, r float64
}

var _ colony.MutationStrategy = (*UniformMutationStrategy)(nil)

func NewUniformMutationStrategy(l, r float64) *UniformMutationStrategy {
	return &UniformMutationStrategy{
		l: l,
		r: r,
	}
}

func (s *UniformMutationStrategy) Mutate(a ant.AntView) *ant.HeteroAnt {
	alpha := s.l + rand.Float64()*(s.r-s.l)
	beta := s.l + rand.Float64()*(s.r-s.l)

	return ant.NewHeteroAnt(
		a.Alpha()+alpha,
		a.Beta()+beta,
		a.PheromoneMultiplier(),
		a.PathStrategy(),
		a.PheromoneApplyStrategy(),
	)
}
