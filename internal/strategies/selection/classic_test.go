package selection

import (
	"HeteroAntColonySystem/internal/core"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockChoose struct {
	path []*graph.Vertex
}

var _ core.ChoosePathStrategy = &mockChoose{}

func (m *mockChoose) ChooseNext(state core.AntInWorkView, ant *core.HeteroAnt) (*graph.Vertex, bool) {
	for _, v := range m.path {
		if !state.Visited(v) {
			return v, false
		}
	}
	return nil, true
}

func TestClassicSelectionSelectsCountBest(t *testing.T) {
	g := graph.NewGraph(4)
	v1 := graph.NewVertex("A")
	v2 := graph.NewVertex("B")
	v3 := graph.NewVertex("C")
	v4 := graph.NewVertex("D")
	g.AddVertex(v1)
	g.AddVertex(v2)
	g.AddVertex(v3)
	g.AddVertex(v4)
	g.AddEdge(2, v1, v2)
	g.AddEdge(2, v2, v1)
	g.AddEdge(5, v1, v3)
	g.AddEdge(5, v3, v1)
	g.AddEdge(3, v1, v4)
	g.AddEdge(3, v4, v1)
	g.AddEdge(7, v2, v3)
	g.AddEdge(7, v3, v2)
	g.AddEdge(8, v2, v4)
	g.AddEdge(8, v4, v2)
	g.AddEdge(3, v3, v4)
	g.AddEdge(3, v4, v3)

	ants := []*core.HeteroAnt{
		core.NewHeteroAnt(1, 1, &mockChoose{
			path: []*graph.Vertex{v1, v2, v3, v4}, // 2 + 7 + 3 + 3 = 15
		}),
		core.NewHeteroAnt(1, 1, &mockChoose{
			path: []*graph.Vertex{v1, v3, v2, v4}, // 5 + 7 + 8 + 3 = 23
		}),
		core.NewHeteroAnt(1, 1, &mockChoose{
			path: []*graph.Vertex{v1, v4, v3, v2}, // 3 + 3 + 7 + 2 = 15
		}),
		core.NewHeteroAnt(1, 1, &mockChoose{
			path: []*graph.Vertex{v1, v2, v4, v3}, // 2 + 8 + 3 + 5 = 18
		}),
		core.NewHeteroAnt(1, 1, &mockChoose{
			path: []*graph.Vertex{v1, v3, v4, v2}, // 5 + 3 + 3 + 2 = 13
		}),
		core.NewHeteroAnt(1, 1, &mockChoose{
			path: []*graph.Vertex{v1, v4, v2, v3}, // 3 + 8 + 7 + 5 = 23
		}),
	}

	expectedScores := []float64{
		13, 15, 18,
	}

	pm := pheromone.NewPheromoneMap(g, 1)

	for _, ant := range ants {
		ant.StartAnt(g, pm, v1)
		ant.Run()
	}

	selection := NewClassicSelection()
	res := selection.Select(ants, 3)

	require.Len(t, res, 3)
	for _, e := range res {
		require.Contains(t, expectedScores, e.Score())
	}
}

func TestClassicSelectionReturnsEmptySliceWhenCountIsZero(t *testing.T) {
	g := graph.NewGraph(3)
	v1 := graph.NewVertex("A")
	v2 := graph.NewVertex("B")
	v3 := graph.NewVertex("C")
	g.AddVertex(v1)
	g.AddVertex(v2)
	g.AddVertex(v3)
	g.AddEdge(1, v1, v2)
	g.AddEdge(1, v2, v1)
	g.AddEdge(1, v2, v3)
	g.AddEdge(1, v3, v2)
	g.AddEdge(1, v3, v1)
	g.AddEdge(1, v1, v3)

	pm := pheromone.NewPheromoneMap(g, 1)

	ant := core.NewHeteroAnt(1, 1, &mockChoose{})

	ant.StartAnt(g, pm, v1)
	ant.Run()

	selection := NewClassicSelection()
	res := selection.Select([]*core.HeteroAnt{}, 0)

	require.Len(t, res, 0)
}
