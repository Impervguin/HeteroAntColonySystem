package choose_test

import (
	"HeteroAntColonySystem/internal/core"
	"HeteroAntColonySystem/internal/strategies/choose"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
	"testing"
)

type mockState struct {
	graph   *graph.Graph
	pm      *pheromone.PheromoneMap
	current *graph.Vertex
	visited map[*graph.Vertex]struct{}
}

func (m *mockState) Graph() *graph.Graph {
	return m.graph
}

func (m *mockState) PheromoneMap() *pheromone.PheromoneMap {
	return m.pm
}

func (m *mockState) Current() *graph.Vertex {
	return m.current
}

func (m *mockState) Visited(v *graph.Vertex) bool {
	_, ok := m.visited[v]
	return ok
}

func TestClassicChoosePathCanSelectAllVertices(t *testing.T) {

	g := graph.NewGraph(3)

	v1 := graph.NewVertex("A")
	v2 := graph.NewVertex("B")
	v3 := graph.NewVertex("C")
	v4 := graph.NewVertex("D")
	v5 := graph.NewVertex("E")
	v6 := graph.NewVertex("F")

	g.AddVertex(v1)
	g.AddVertex(v2)
	g.AddVertex(v3)
	g.AddVertex(v4)
	g.AddVertex(v5)
	g.AddVertex(v6)

	g.AddEdge(1, v1, v2)
	g.AddEdge(1, v1, v3)
	g.AddEdge(1, v1, v4)
	g.AddEdge(1, v1, v5)
	g.AddEdge(1, v1, v6)

	pm := pheromone.NewPheromoneMap(g, 1)

	state := &mockState{
		graph:   g,
		pm:      pm,
		current: v1,
		visited: map[*graph.Vertex]struct{}{
			v1: {},
		},
	}

	strategy := choose.NewClassicChoosePath()
	ant := core.NewHeteroAnt(1, 1, strategy)

	found := map[*graph.Vertex]bool{}

	for i := 0; i < 1000; i++ {

		v, done := strategy.ChooseNext(state, ant)

		if done {
			t.Fatalf("unexpected done=true")
		}

		found[v] = true
	}

	if !found[v2] || !found[v3] || !found[v4] || !found[v5] || !found[v6] {
		t.Fatalf("not all vertices were selected: %+v", found)
	}
}

func TestClassicChoosePathReturnsDoneWhenFinished(t *testing.T) {

	g := graph.NewGraph(2)

	v1 := graph.NewVertex("A")
	v2 := graph.NewVertex("B")

	g.AddVertex(v1)
	g.AddVertex(v2)

	g.AddEdge(1, v1, v2)

	pm := pheromone.NewPheromoneMap(g, 1)

	state := &mockState{
		graph:   g,
		pm:      pm,
		current: v1,
		visited: map[*graph.Vertex]struct{}{
			v1: {},
			v2: {},
		},
	}

	strategy := choose.NewClassicChoosePath()
	ant := core.NewHeteroAnt(1, 1, strategy)

	v, done := strategy.ChooseNext(state, ant)

	if !done {
		t.Fatalf("expected done=true")
	}

	if v != nil {
		t.Fatalf("expected nil vertex when finished")
	}
}
