package aco

import (
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
)

type ChoosePathStrategy interface {
	ChooseNext(ant *Ant) (*graph.Vertex, bool)
}

type Ant struct {
	g       *graph.Graph
	pm      *pheromone.PheromoneMap
	initial *graph.Vertex
	choice  ChoosePathStrategy

	current *graph.Vertex
	route   []*graph.Vertex
	visited map[*graph.Vertex]struct{}
	score   float64
	done    bool
}

func NewAnt(g *graph.Graph, pm *pheromone.PheromoneMap, initial *graph.Vertex, choice ChoosePathStrategy) *Ant {
	return &Ant{
		g:       g,
		pm:      pm,
		initial: initial,
		current: initial,
		route:   []*graph.Vertex{initial},
		visited: map[*graph.Vertex]struct{}{initial: {}},
		score:   0.,
		choice:  choice,
		done:    false,
	}
}

func (a *Ant) Run() {
	if a.done {
		return
	}
	for !a.step() {
	}

	score := 0.
	if len(a.route) > 1 {
		for i := 0; i < len(a.route)-1; i++ {
			e, _ := a.g.Edge(a.route[i], a.route[i+1])
			score += e.Weight()
		}
		e, _ := a.g.Edge(a.route[len(a.route)-1], a.route[0])
		score += e.Weight()
	}

	a.score = score
	a.done = true
}

func (a *Ant) Tour() []*graph.Vertex {
	return a.route
}

func (a *Ant) Score() float64 {
	return a.score
}

func (a *Ant) step() bool {
	next, done := a.choice.ChooseNext(a)
	if done {
		return true
	}
	a.route = append(a.route, next)
	a.current = next
	a.visited[next] = struct{}{}
	return false
}
