package aco

import (
	"HeteroAntColonySystem/pkg/graph"
	"math"
	"math/rand/v2"
)

type ClassicStrategy struct {
	alpha float64
	beta  float64
}

var _ ChoosePathStrategy = &ClassicStrategy{}

func NewClassicStrategy(alpha, beta float64) *ClassicStrategy {
	return &ClassicStrategy{
		alpha: alpha,
		beta:  beta,
	}
}

func (c *ClassicStrategy) ChooseNext(ant *Ant) (*graph.Vertex, bool) {
	gr := ant.g
	current := ant.current
	pm := ant.pm

	type candidate struct {
		vertex *graph.Vertex
		prob   float64
	}

	candidates := make([]candidate, 0)

	gr.ForEachSource(current, func(e *graph.Edge) bool {
		if _, ok := ant.visited[e.Target()]; ok {
			return false
		}

		heuristic := math.Pow(1.0/e.Weight(), c.beta)
		pher := math.Pow(pm.Get(e), c.alpha)

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

	r := rand.Float64() * sum

	acc := 0.0
	for _, c := range candidates {
		acc += c.prob
		if r <= acc {
			return c.vertex, false
		}
	}
	return candidates[len(candidates)-1].vertex, false
}
