package strategy

import (
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
)

// AntView provides a read-only interface for ant strategies to access
// necessary information without exposing internal implementation details

type AntView interface {
	Graph() *graph.Graph
	PheromoneMap() *pheromone.PheromoneMap
	Current() *graph.Vertex
	Visited(vertex *graph.Vertex) bool

	Path() []*graph.Vertex
	Score() float64

	Alpha() float64
	Beta() float64
	PheromoneMultiplier() float64

	PathStrategy() PathChoiceStrategy
	PheromoneApplyStrategy() PheromoneApplyingStrategy
}

type PathChoiceStrategy interface {
	// ChooseNext Chooses next vertex for ants path
	// must return nil if path done
	ChooseNext(ant AntView) *graph.Vertex
}

type PheromoneApplyingStrategy interface {
	ApplyPheromone(ant AntView)
}

type ParentSelectionStrategy interface {
	SelectParents(ants []AntView, n uint) []AntView
}

type CrossoverStrategy interface {
	Crossover(p1, p2 AntView) AntView
}

type MutationStrategy interface {
	Mutate(ant AntView) AntView
}
