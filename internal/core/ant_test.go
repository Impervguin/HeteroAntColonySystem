package core_test

import (
	"HeteroAntColonySystem/internal/core"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
	"testing"

	"github.com/stretchr/testify/require"
)

type testChooseStrategy struct{}

func (t *testChooseStrategy) ChooseNext(state core.AntInWorkView, ant *core.HeteroAnt) (*graph.Vertex, bool) {
	var next *graph.Vertex

	state.Graph().ForEachSource(state.Current(), func(e *graph.Edge) bool {
		if !state.Visited(e.Target()) {
			next = e.Target()
			return true
		}
		return false
	})

	if next == nil {
		return nil, true
	}

	return next, false
}

var _ core.ChoosePathStrategy = &testChooseStrategy{}

func TestAntBuildsRoute(t *testing.T) {
	g := graph.NewGraph(3)

	v1 := graph.NewVertex("A")
	v2 := graph.NewVertex("B")
	v3 := graph.NewVertex("C")

	g.AddVertex(v1)
	g.AddVertex(v2)
	g.AddVertex(v3)

	g.AddEdge(1, v1, v2)
	g.AddEdge(1, v2, v3)
	g.AddEdge(1, v3, v1)

	g.AddEdge(1, v1, v3)
	g.AddEdge(1, v3, v2)
	g.AddEdge(1, v2, v1)

	pm := pheromone.NewPheromoneMap(g, 1)

	ant := core.NewHeteroAnt(1, 1, &testChooseStrategy{})

	ant.StartAnt(g, pm, v1)
	ant.Run()

	if ant.Score() == 0 {
		t.Fatalf("expected score > 0")
	}

	tour := ant.Tour()
	require.Len(t, tour, 3)

	require.Contains(t, tour, v1)
	require.Contains(t, tour, v2)
	require.Contains(t, tour, v3)
}

func TestAntStopsWhenNoMoves(t *testing.T) {
	g := graph.NewGraph(1)

	v1 := graph.NewVertex("A")

	g.AddVertex(v1)

	pm := pheromone.NewPheromoneMap(g, 1)

	ant := core.NewHeteroAnt(1, 1, &testChooseStrategy{})

	ant.StartAnt(g, pm, v1)
	ant.Run()

	if ant.Score() != 0 {
		t.Fatalf("expected score 0 for single vertex")
	}
}
