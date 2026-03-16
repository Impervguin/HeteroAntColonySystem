package core

import (
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
)

type AntInWorkView interface {
	Graph() *graph.Graph
	PheromoneMap() *pheromone.PheromoneMap
	Current() *graph.Vertex
	Visited(v *graph.Vertex) bool
}

// ChoosePathStrategy defines the rule used by an ant to select the next vertex
// during tour construction. It typically uses pheromone levels, heuristic
// information, and the ant's parameters (alpha and beta).
// The returned bool indicates whether the tour is finished.
type ChoosePathStrategy interface {
	ChooseNext(state AntInWorkView, ant *HeteroAnt) (*graph.Vertex, bool)
}

// SelectionStrategy defines the method used in the genetic algorithm
// to select a subset of ants from the candidate population based on
// their fitness (tour quality).
type SelectionStrategy interface {
	Select(candidates []*HeteroAnt, count int) []*HeteroAnt
}

// CrossoverStrategy defines how two parent ants are combined to produce
// a new ant with mixed parameters (e.g., alpha and beta). This is part
// of the genetic algorithm used to evolve heterogeneous ant parameters.
type CrossoverStrategy interface {
	Crossover(ant *HeteroAnt, other *HeteroAnt) *HeteroAnt
}

// MutationStrategy defines how an ant's parameters are randomly modified
// to introduce diversity into the population during the genetic algorithm process.
type MutationStrategy interface {
	Mutate(ant *HeteroAnt) *HeteroAnt
}
