package core

import (
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
	"fmt"
)

type HeteroAnt struct {
	alpha               float64
	beta                float64
	pheromoneMultiplier float64

	chooseNext     PathChoiceStrategy
	pheromoneApply PheromoneApplyingStrategy
	*HeteroAntWork
}

type HeteroAntWork struct {
	g  *graph.Graph
	pm *pheromone.PheromoneMap

	current *graph.Vertex
	path    []*graph.Vertex
	visited map[*graph.Vertex]struct{}
	done    bool

	score float64
}

func NewHeteroAnt(alpha, beta, pherMultiplier float64, chooseNext PathChoiceStrategy, pherAppl PheromoneApplyingStrategy) *HeteroAnt {
	return &HeteroAnt{
		alpha:               alpha,
		beta:                beta,
		pheromoneMultiplier: pherMultiplier,
		chooseNext:          chooseNext,
		pheromoneApply:      pherAppl,
		HeteroAntWork:       nil,
	}
}

func (a *HeteroAnt) PathStrategy() PathChoiceStrategy {
	return a.chooseNext
}

func (a *HeteroAnt) PheromoneApplyStrategy() PheromoneApplyingStrategy {
	return a.pheromoneApply
}

func (a *HeteroAnt) Alpha() float64 {
	return a.alpha
}

func (a *HeteroAnt) Beta() float64 {
	return a.beta
}

func (a *HeteroAnt) PheromoneMultiplier() float64 {
	return a.pheromoneMultiplier
}

func (a *HeteroAnt) Prepare(g *graph.Graph, pm *pheromone.PheromoneMap) {
	a.HeteroAntWork = &HeteroAntWork{
		g:  g,
		pm: pm,

		current: g.RandomVertex(),
		path:    make([]*graph.Vertex, 0, g.Len()),
		visited: make(map[*graph.Vertex]struct{}, g.Len()),
		done:    false,
		score:   -1,
	}

	a.path = append(a.path, a.current)
	a.visited[a.current] = struct{}{}
}

func (a *HeteroAnt) Run() error {
	if a.HeteroAntWork == nil {
		return ErrAntNotPrepared
	}
	for !a.step() {
	}
	a.done = true
	a.calculateScore()
	return nil
}

func (a *HeteroAnt) ApplyPheromone() error {
	if a.HeteroAntWork == nil || !a.done {
		return ErrAntNotDone
	}

	a.pheromoneApply.ApplyPheromone(a)
	return nil
}

func (a *HeteroAnt) Score() float64 {
	if a.HeteroAntWork == nil || !a.done {
		return -1
	}

	if a.score == -1 {
		a.calculateScore()
	}

	return a.score
}

func (a *HeteroAnt) Path() []*graph.Vertex {
	if a.HeteroAntWork == nil || !a.done {
		return nil
	}
	return a.path
}

func (a *HeteroAnt) Visited(v *graph.Vertex) bool {
	if a.HeteroAntWork == nil {
		fmt.Println("nil")
		return false
	}
	_, ok := a.visited[v]
	return ok
}

func (a *HeteroAnt) Graph() *graph.Graph {
	if a.HeteroAntWork == nil {
		return nil
	}
	return a.g
}

func (a *HeteroAnt) PheromoneMap() *pheromone.PheromoneMap {
	if a.HeteroAntWork == nil {
		return nil
	}
	return a.HeteroAntWork.pm
}

func (a *HeteroAnt) Current() *graph.Vertex {
	if a.HeteroAntWork == nil || a.done {
		return nil
	}
	return a.current
}

func (a *HeteroAnt) step() bool {
	next := a.chooseNext.ChooseNext(a)
	// fmt.Println(next)
	if next == nil {
		return true
	}
	a.current = next
	a.visited[a.current] = struct{}{}
	a.path = append(a.path, a.current)
	return false
}

func (a *HeteroAnt) calculateScore() {
	a.score = 0
	if len(a.path) <= 1 {
		return
	}

	for i := 0; i < len(a.path)-1; i++ {
		v1, v2 := a.path[i], a.path[i+1]
		e, _ := a.g.Edge(v1, v2)
		a.score += e.Weight()
	}
	wrapE, _ := a.g.Edge(a.path[len(a.path)-1], a.path[0])
	a.score += wrapE.Weight()
}
