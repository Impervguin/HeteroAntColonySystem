package greedy

import (
	"HeteroAntColonySystem/pkg/graph"
	"slices"
)

type GreedyAlgorithm struct {
	g       *graph.Graph
	initial *graph.Vertex

	current *graph.Vertex
	route   []*graph.Vertex
	score   float64
}

func NewGreedyAlgorithm(g *graph.Graph) *GreedyAlgorithm {
	var v *graph.Vertex
	g.ForEachVertex(func(v1 *graph.Vertex) bool {
		v = v1
		return true
	})
	return &GreedyAlgorithm{
		g:       g,
		initial: v,
	}
}

func (a *GreedyAlgorithm) Run() {
	a.current = a.initial
	a.route = []*graph.Vertex{a.initial}
	a.score = 0.

	for {
		next, done := a.chooseNext()
		if done {
			break
		}
		a.route = append(a.route, next)
		e, _ := a.g.Edge(a.current, next)
		a.current = next
		a.score += e.Weight()
	}

	laste, _ := a.g.Edge(a.route[len(a.route)-1], a.route[0])
	a.score += laste.Weight()
}

func (a *GreedyAlgorithm) Tour() []*graph.Vertex {
	return a.route
}

func (a *GreedyAlgorithm) Score() float64 {
	return a.score
}

func (a *GreedyAlgorithm) chooseNext() (*graph.Vertex, bool) {
	gr := a.g

	var enext *graph.Edge
	gr.ForEachSource(a.current, func(e *graph.Edge) bool {
		if slices.Contains(a.route, e.Target()) {
			return false
		}

		if enext == nil {
			enext = e
			return false
		}

		if e.Weight() < enext.Weight() {
			enext = e
		}

		return false
	})
	if enext != nil {
		return enext.Target(), false
	}
	return nil, true
}
