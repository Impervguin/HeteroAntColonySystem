package ant

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
