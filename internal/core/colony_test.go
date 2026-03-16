package core

import (
	"HeteroAntColonySystem/pkg/graph"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockChoose struct{}

func (m *mockChoose) ChooseNext(state AntInWorkView, ant *HeteroAnt) (*graph.Vertex, bool) {
	var next *graph.Vertex = nil
	state.Graph().ForEachVertex(func(v *graph.Vertex) bool {
		if !state.Visited(v) {
			next = v
			return true
		}
		return false
	})
	if next != nil {
		return next, false
	}
	return nil, true
}

type mockSelection struct{}

func (m *mockSelection) Select(candidates []*HeteroAnt, count uint) []*HeteroAnt {
	return candidates[:count]
}

type mockCrossover struct{}

func (m *mockCrossover) Crossover(a, b *HeteroAnt) *HeteroAnt {
	return NewHeteroAnt(a.Alpha(), b.Beta(), a.strategy)
}

type mockMutation struct{}

func (m *mockMutation) Mutate(a *HeteroAnt) *HeteroAnt {
	return a
}

func buildTriangleGraph() *graph.Graph {
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
	return g
}

func TestNewHeteroAntColonySuccess(t *testing.T) {
	g := buildTriangleGraph()
	c, err := NewHeteroAntColony(g,
		WithChooseStrategy(&mockChoose{}),
		WithSelectionStrategy(&mockSelection{}),
		WithCrossoverStrategy(&mockCrossover{}),
		WithMutationStrategy(&mockMutation{}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.g != g {
		t.Fatalf("graph not set correctly")
	}
}

func TestRunInitializesState(t *testing.T) {
	g := buildTriangleGraph()
	c, _ := NewHeteroAntColony(g,
		WithChooseStrategy(&mockChoose{}),
		WithSelectionStrategy(&mockSelection{}),
		WithCrossoverStrategy(&mockCrossover{}),
		WithMutationStrategy(&mockMutation{}),
	)
	c.Run()
	if c.result == nil {
		t.Fatalf("result should not be nil after Run")
	}
	if len(c.result.bestTour) != 3 {
		t.Fatalf("expected best tour length 3, got %d", len(c.result.bestTour))
	}
}

func TestEvolveNewGeneration(t *testing.T) {
	g := buildTriangleGraph()
	c, _ := NewHeteroAntColony(g,
		WithChooseStrategy(&mockChoose{}),
		WithSelectionStrategy(&mockSelection{}),
		WithCrossoverStrategy(&mockCrossover{}),
		WithMutationStrategy(&mockMutation{}),
	)
	c.colonySize = 5
	c.parentCount = 2

	state := &AntColonyState{
		generation: make([]*HeteroAnt, 5),
	}
	for i := 0; i < 5; i++ {
		state.generation[i] = NewHeteroAnt(rand.Float64(), rand.Float64(), &mockChoose{})
	}
	c.state = state
	oldGeneration := c.state.generation

	c.evolve()

	require.Len(t, c.state.generation, 5)
	for i := 0; i < 5; i++ {
		require.NotContains(t, c.state.generation, oldGeneration[i])
	}
}

func TestStagnateCopiesAntParameters(t *testing.T) {
	g := buildTriangleGraph()
	c, _ := NewHeteroAntColony(g,
		WithChooseStrategy(&mockChoose{}),
		WithSelectionStrategy(&mockSelection{}),
		WithCrossoverStrategy(&mockCrossover{}),
		WithMutationStrategy(&mockMutation{}),
	)
	c.colonySize = 3

	state := &AntColonyState{
		generation: []*HeteroAnt{
			NewHeteroAnt(1.2, 0.8, &mockChoose{}),
			NewHeteroAnt(1.5, 1.0, &mockChoose{}),
			NewHeteroAnt(2.0, 0.5, &mockChoose{}),
		},
	}
	c.state = state
	c.stagnate()

	for i, ant := range c.state.generation {
		if ant.Alpha() != state.generation[i].Alpha() || ant.Beta() != state.generation[i].Beta() {
			t.Fatalf("stagnated ant parameters mismatch")
		}
	}
}

func TestColony_RunSimpleGraph(t *testing.T) {
	g := buildTriangleGraph()
	c, _ := NewHeteroAntColony(g,
		WithChooseStrategy(&mockChoose{}),
		WithSelectionStrategy(&mockSelection{}),
		WithCrossoverStrategy(&mockCrossover{}),
		WithMutationStrategy(&mockMutation{}),
		WithParentCount(2),
		WithColonySize(5),
		WithGenerationCount(10),
		WithGenerationPeriod(5),
	)
	c.Run()
	if c.result == nil || c.result.bestTour == nil || c.result.bestScore <= 0 {
		t.Fatalf("expected valid result with positive score")
	}
}
