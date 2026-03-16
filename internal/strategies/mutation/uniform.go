package mutation

import (
	"HeteroAntColonySystem/internal/core"
	"math/rand"
	"time"
)

type UniformMutation struct {
	r *rand.Rand
	a float64
	b float64

	choose core.ChoosePathStrategy
}

func NewUniformMutation(a, b float64, choose core.ChoosePathStrategy) *UniformMutation {
	if a > b {
		a, b = b, a
	}
	return &UniformMutation{
		a:      a,
		b:      b,
		r:      rand.New(rand.NewSource(time.Now().UnixNano())),
		choose: choose,
	}
}

var _ core.MutationStrategy = &UniformMutation{}

func (m *UniformMutation) Mutate(ant *core.HeteroAnt) *core.HeteroAnt {
	adelta := m.r.Float64()*(m.b-m.a) + m.a
	bdelta := m.r.Float64()*(m.b-m.a) + m.a
	return core.NewHeteroAnt(ant.Alpha()+adelta, ant.Beta()+bdelta, m.choose)
}
