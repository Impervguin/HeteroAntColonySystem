package ant

import (
	"HeteroAntColonySystem/internal/core/errors"
	"HeteroAntColonySystem/internal/core/strategy"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
)

// HeteroAnt represents an ant in the Ant Colony Optimization algorithm.
// It contains both configuration parameters (alpha, beta, pheromoneMultiplier)
// and runtime state (current position, path, visited vertices).
//
// The ant uses strategy pattern for path selection and pheromone application,
// allowing for different ACO variants to be implemented.
type HeteroAnt struct {
	// Configuration parameters
	alpha               float64
	beta                float64
	pheromoneMultiplier float64

	// Strategies
	chooseNext     strategy.PathChoiceStrategy
	pheromoneApply strategy.PheromoneApplyingStrategy

	// Runtime state (initialized in Prepare)
	g       *graph.Graph
	pm      *pheromone.PheromoneMap
	current *graph.Vertex
	path    []*graph.Vertex
	visited map[*graph.Vertex]struct{}
	done    bool
	score   float64
}

// NewHeteroAnt creates a new heterogeneous ant with the given parameters and strategies.
func NewHeteroAnt(alpha, beta, pherMultiplier float64, chooseNext strategy.PathChoiceStrategy, pherAppl strategy.PheromoneApplyingStrategy) *HeteroAnt {
	return &HeteroAnt{
		alpha:               alpha,
		beta:                beta,
		pheromoneMultiplier: pherMultiplier,
		chooseNext:          chooseNext,
		pheromoneApply:      pherAppl,
	}
}

// Alpha returns the alpha parameter for this ant.
func (a *HeteroAnt) Alpha() float64 {
	return a.alpha
}

// Beta returns the beta parameter for this ant.
func (a *HeteroAnt) Beta() float64 {
	return a.beta
}

// PheromoneMultiplier returns the pheromone multiplier for this ant.
func (a *HeteroAnt) PheromoneMultiplier() float64 {
	return a.pheromoneMultiplier
}

// PathStrategy returns the path choice strategy used by this ant.
func (a *HeteroAnt) PathStrategy() strategy.PathChoiceStrategy {
	return a.chooseNext
}

// PheromoneApplyStrategy returns the pheromone applying strategy used by this ant.
func (a *HeteroAnt) PheromoneApplyStrategy() strategy.PheromoneApplyingStrategy {
	return a.pheromoneApply
}

// Prepare initializes the ant for a new tour on the given graph.
// It resets the ant's state and selects a random starting vertex.
func (a *HeteroAnt) Prepare(g *graph.Graph, pm *pheromone.PheromoneMap) {
	a.g = g
	a.pm = pm
	a.current = g.RandomVertex()
	a.path = make([]*graph.Vertex, 0, g.Len())
	a.visited = make(map[*graph.Vertex]struct{}, g.Len())
	a.done = false
	a.score = -1

	// Mark starting vertex as visited
	a.path = append(a.path, a.current)
	a.visited[a.current] = struct{}{}
}

// Run executes the ant's path construction algorithm.
// The ant continues to select next vertices until all have been visited.
func (a *HeteroAnt) Run() error {
	if a.g == nil {
		return errors.ErrAntNotPrepared
	}
	for !a.step() {
		// Continue until step returns true (no more vertices)
	}
	a.done = true
	a.calculateScore()
	return nil
}

// step performs one iteration of path construction.
// Returns true if the path is complete (no more unvisited vertices).
func (a *HeteroAnt) step() bool {
	next := a.chooseNext.ChooseNext(a)
	if next == nil {
		return true // Path complete
	}
	a.current = next
	a.visited[a.current] = struct{}{}
	a.path = append(a.path, a.current)
	return false
}

// calculateScore computes the total weight of the ant's path.
func (a *HeteroAnt) calculateScore() {
	a.score = 0
	if len(a.path) <= 1 {
		return
	}

	// Sum edges in the path
	for i := 0; i < len(a.path)-1; i++ {
		v1, v2 := a.path[i], a.path[i+1]
		e, _ := a.g.Edge(v1, v2)
		a.score += e.Weight()
	}

	// Add return edge to complete the tour
	wrapE, _ := a.g.Edge(a.path[len(a.path)-1], a.path[0])
	a.score += wrapE.Weight()
}

// ApplyPheromone updates the pheromone levels based on this ant's path.
// The amount of pheromone deposited is inversely proportional to the path score.
func (a *HeteroAnt) ApplyPheromone() error {
	if a.g == nil || !a.done {
		return errors.ErrAntNotDone
	}

	a.pheromoneApply.ApplyPheromone(a)
	return nil
}

// Score returns the total weight of the ant's path.
// Returns -1 if the ant hasn't completed its path.
func (a *HeteroAnt) Score() float64 {
	if a.g == nil || !a.done {
		return -1
	}

	if a.score == -1 {
		a.calculateScore()
	}

	return a.score
}

// Path returns the vertices in the ant's path.
// Returns nil if the ant hasn't completed its path.
func (a *HeteroAnt) Path() []*graph.Vertex {
	if a.g == nil || !a.done {
		return nil
	}
	return a.path
}

// Visited checks if the given vertex has been visited by this ant.
func (a *HeteroAnt) Visited(v *graph.Vertex) bool {
	if a.g == nil {
		return false
	}
	_, ok := a.visited[v]
	return ok
}

// Graph returns the graph the ant is currently operating on.
func (a *HeteroAnt) Graph() *graph.Graph {
	return a.g
}

// PheromoneMap returns the pheromone map the ant is using.
func (a *HeteroAnt) PheromoneMap() *pheromone.PheromoneMap {
	return a.pm
}

// Current returns the current vertex the ant is at.
// Returns nil if the ant hasn't been prepared or has completed its path.
func (a *HeteroAnt) Current() *graph.Vertex {
	if a.g == nil || a.done {
		return nil
	}
	return a.current
}

func (a *HeteroAnt) FullCopy() *HeteroAnt {
	path := make([]*graph.Vertex, 0, len(a.path))
	for _, v := range a.path {
		path = append(path, v)
	}
	visited := make(map[*graph.Vertex]struct{}, len(a.visited))
	for v := range a.visited {
		visited[v] = struct{}{}
	}
	return &HeteroAnt{
		alpha:               a.alpha,
		beta:                a.beta,
		pheromoneMultiplier: a.pheromoneMultiplier,
		chooseNext:          a.chooseNext,
		pheromoneApply:      a.pheromoneApply,
		g:                   a.g,
		pm:                  a.pm,
		current:             a.current,
		path:                path,
		visited:             visited,
		done:                a.done,
		score:               a.score,
	}
}
