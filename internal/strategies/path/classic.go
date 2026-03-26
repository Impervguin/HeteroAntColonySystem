package path

import (
	"HeteroAntColonySystem/internal/core/strategy"
	"HeteroAntColonySystem/pkg/graph"
	"math"
	"math/rand/v2"
)

// PathClassicStrategy implements the classic ACO path selection algorithm
// using probability based on pheromone levels and edge weights
// Implements strategy.PathChoiceStrategy interface


type PathClassicStrategy struct{}

var _ strategy.PathChoiceStrategy = &PathClassicStrategy{}

func NewPahtClassicStrategy() *PathClassicStrategy {
	return &PathClassicStrategy{}
}

func (s *PathClassicStrategy) ChooseNext(ant strategy.AntView) *graph.Vertex {
	g := ant.Graph()
	pm := ant.PheromoneMap()
	current := ant.Current()

	type NextCandidate struct {
		v           *graph.Vertex
		probability float64
	}
	candidates := make([]NextCandidate, 0)
	var accumulatedProbability float64

	g.ForEachSource(current, func(edge *graph.Edge) bool {
		t := edge.Target()
		if ant.Visited(t) {
			return false
		}
		w := math.Pow(1./edge.Weight(), ant.Beta())
		p := math.Pow(pm.Get(edge), ant.Alpha())

		probability := w * p
		accumulatedProbability += probability
		candidates = append(candidates, NextCandidate{
			v:           t,
			probability: probability,
		})
		return false
	})

	if len(candidates) == 0 {
		return nil
	}

	r := rand.Float64() * accumulatedProbability

	cumulative := 0.
	for _, candidate := range candidates {
		cumulative += candidate.probability
		if r <= cumulative {
			return candidate.v
		}
	}

	return candidates[len(candidates)-1].v
}
