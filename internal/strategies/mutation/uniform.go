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
	alpha := a.Alpha() + s.l + rand.Float64()*(s.r-s.l)
	beta := a.Beta() + s.l + rand.Float64()*(s.r-s.l)
	if alpha < 0 {
		alpha = 0
	}
	if beta < 0 {
		beta = 0
	}

	return ant.NewHeteroAnt(
		alpha,
		beta,
		a.PheromoneMultiplier(),
		a.PathStrategy(),
		a.PheromoneApplyStrategy(),
	)
}
