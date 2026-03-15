package pheromone

import (
	"HeteroAntColonySystem/pkg/graph"
	"sync"
)

type PheromoneMap struct {
	gr *graph.Graph
	pm map[*graph.Edge]float64

	mut sync.RWMutex
}

func NewPheromoneMap(gr *graph.Graph, initial float64) *PheromoneMap {
	pm := &PheromoneMap{
		gr:  gr,
		pm:  make(map[*graph.Edge]float64, 0),
		mut: sync.RWMutex{},
	}

	gr.ForEachEdge(func(e *graph.Edge) {
		pm.pm[e] = initial
	})

	return pm
}

func (pm *PheromoneMap) Update(e *graph.Edge, value float64) {
	pm.mut.Lock()
	defer pm.mut.Unlock()

	pm.pm[e] = value
}

func (pg *PheromoneMap) Add(e *graph.Edge, delta float64) {
	pg.mut.Lock()
	defer pg.mut.Unlock()

	pg.pm[e] += delta
}

func (pm *PheromoneMap) Get(e *graph.Edge) float64 {
	pm.mut.RLock()
	defer pm.mut.RUnlock()

	return pm.pm[e]
}

func (pm *PheromoneMap) ForEachEdge(f func(e *graph.Edge, pheromone float64) float64) {
	pm.mut.Lock()
	defer pm.mut.Unlock()

	for e, p := range pm.pm {
		pm.pm[e] = f(e, p)
	}
}

func (pm *PheromoneMap) ForEachEdgeRead(f func(e *graph.Edge, pheromone float64)) {
	pm.mut.Lock()
	defer pm.mut.Unlock()

	for e, p := range pm.pm {
		f(e, p)
	}
}
