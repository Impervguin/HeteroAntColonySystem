package core

import (
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
)

// HeteroAnt represents a single ant in the Ant Colony System with
// heterogeneous parameters (alpha and beta) that control the influence
// of pheromone trails and heuristic distance during path selection.
type HeteroAnt struct {
	// alpha — influence of pheromone on choosing an edge
	alpha float64

	// beta — influence of heuristic (distance) on choosing an edge
	beta float64

	// strategy defines how the ant selects the next vertex
	strategy ChoosePathStrategy

	// inWork stores the current state of the ant during traversal
	inWork *AntState

	// result stores the completed tour and its score after traversal
	result *AntResult
}

// NewHeteroAnt creates a new ant with the given alpha, beta, and path selection strategy.
func NewHeteroAnt(alpha, beta float64, strategy ChoosePathStrategy) *HeteroAnt {
	return &HeteroAnt{
		alpha:    alpha,
		beta:     beta,
		strategy: strategy,
	}
}

// Alpha returns the ant's pheromone influence coefficient.
func (a *HeteroAnt) Alpha() float64 {
	return a.alpha
}

// Beta returns the ant's heuristic influence coefficient.
func (a *HeteroAnt) Beta() float64 {
	return a.beta
}

// AntState stores the current state of an ant while constructing a tour.
type AntState struct {
	graph *graph.Graph
	pm    *pheromone.PheromoneMap

	// current vertex where the ant is located
	current *graph.Vertex

	// route contains vertices visited in order so far
	route []*graph.Vertex

	// visited keeps track of vertices already visited
	visited map[*graph.Vertex]struct{}
}

// AntResult stores the completed tour of an ant along with its total score.
type AntResult struct {
	// tour is the sequence of visited vertices forming the tour
	tour []*graph.Vertex

	// visited keeps track of vertices visited during the tour
	visited map[*graph.Vertex]struct{}

	// score is the total weight (length) of the tour
	score float64
}

// StartAnt initializes the ant's state for a new traversal starting at the given vertex.
func (a *HeteroAnt) StartAnt(gr *graph.Graph, pheromoneMap *pheromone.PheromoneMap, initial *graph.Vertex) {
	a.inWork = &AntState{
		graph:   gr,
		current: initial,
		pm:      pheromoneMap,
		route:   []*graph.Vertex{initial},
		visited: map[*graph.Vertex]struct{}{initial: {}},
	}
}

// step performs a single step of the ant's traversal, returning true if the tour is complete.
func (a *HeteroAnt) step() bool {
	next, done := a.strategy.ChooseNext(a.inWork, a)
	if done {
		return true
	}
	a.inWork.route = append(a.inWork.route, next)
	a.inWork.current = next
	a.inWork.visited[next] = struct{}{}
	return false
}

// Run executes the ant's traversal until a complete tour is constructed and
// calculates its total score.
func (a *HeteroAnt) Run() {
	for !a.step() {
	}
	res := AntResult{
		tour:    a.inWork.route,
		visited: a.inWork.visited,
	}
	score := 0.
	for i := 0; i < len(res.tour)-1; i++ {
		e, _ := a.inWork.graph.Edge(res.tour[i], res.tour[i+1])
		score += e.Weight()
	}
	e, _ := a.inWork.graph.Edge(res.tour[len(res.tour)-1], res.tour[0])
	score += e.Weight()

	res.score = score
	a.result = &res
	a.inWork = nil
}

// Score returns the total score (length) of the ant's completed tour.
func (a *HeteroAnt) Score() float64 {
	if a.result == nil {
		return 0.
	}
	return a.result.score
}

// Visited returns true if the given vertex has already been visited by the ant.
func (a *AntState) Visited(v *graph.Vertex) bool {
	_, ok := a.visited[v]
	return ok
}

// Current returns the vertex where the ant currently resides.
func (a *AntState) Current() *graph.Vertex {
	return a.current
}

// Graph returns the graph used by the ant during traversal.
func (a *AntState) Graph() *graph.Graph {
	return a.graph
}

// PheromoneMap returns the pheromone map used for probabilistic decisions.
func (a *AntState) PheromoneMap() *pheromone.PheromoneMap {
	return a.pm
}
