package choose

import (
	"HeteroAntColonySystem/internal/core"
	"HeteroAntColonySystem/pkg/graph"
	"math"
	"math/rand"
	"time"
)

type ClassicChoosePath struct {
	r *rand.Rand
}

var _ core.ChoosePathStrategy = &ClassicChoosePath{}

func NewClassicChoosePath() *ClassicChoosePath {
	return &ClassicChoosePath{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *ClassicChoosePath) ChooseNext(
	state core.AntInWorkView,
	ant *core.HeteroAnt,
) (*graph.Vertex, bool) {
	gr := state.Graph()
	current := state.Current()
	pm := state.PheromoneMap()

	type candidate struct {
		vertex *graph.Vertex
		prob   float64
	}

	candidates := make([]candidate, 0)

	gr.ForEachSource(current, func(e *graph.Edge) bool {
		if state.Visited(e.Target()) {
			return false
		}

		heuristic := math.Pow(1.0/e.Weight(), ant.Beta())
		pher := math.Pow(pm.Get(e), ant.Alpha())

		candidates = append(candidates, candidate{
			vertex: e.Target(),
			prob:   heuristic * pher,
		})
		return false
	})

	if len(candidates) == 0 {
		return nil, true
	}

	sum := 0.0
	for _, c := range candidates {
		sum += c.prob
	}

	r := c.r.Float64() * sum

	acc := 0.0
	for _, c := range candidates {
		acc += c.prob
		if r <= acc {
			return c.vertex, false
		}
	}

	return candidates[len(candidates)-1].vertex, false
}
