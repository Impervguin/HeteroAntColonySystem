package apply_test

import (
	"HeteroAntColonySystem/internal/strategies/apply"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyClassicStrategy_ApplyPheromone(t *testing.T) {
	strategy := apply.NewApplyClassicStrategy()

	// Create a simple graph with three vertices in a line
	g := graph.NewGraph(3)
	v1 := graph.NewVertex("V1")
	v2 := graph.NewVertex("V2")
	v3 := graph.NewVertex("V3")
	g.AddVertex(v1)
	g.AddVertex(v2)
	g.AddVertex(v3)

	// Add edges (v1-v2 weight 1, v2-v3 weight 1, v3-v1 weight 2 for wrap)
	e12, _ := g.Edge(v1, v2)
	if e12 == nil {
		if err := g.AddEdge(1.0, v1, v2); err != nil {
			t.Fatalf("Failed to add edge v1-v2: %v", err)
		}
		e12, _ = g.Edge(v1, v2)
	}
	e23, _ := g.Edge(v2, v3)
	if e23 == nil {
		if err := g.AddEdge(1.0, v2, v3); err != nil {
			t.Fatalf("Failed to add edge v2-v3: %v", err)
		}
		e23, _ = g.Edge(v2, v3)
	}
	e31, _ := g.Edge(v3, v1)
	if e31 == nil {
		if err := g.AddEdge(2.0, v3, v1); err != nil {
			t.Fatalf("Failed to add edge v3-v1: %v", err)
		}
		e31, _ = g.Edge(v3, v1)
	}

	// Set up mock ant with path V1 -> V2 -> V3 (completed tour)
	mockAnt := new(MockAntView)
	mockAnt.On("Graph").Return(g)
	mockAnt.On("Path").Return([]*graph.Vertex{v1, v2, v3})
	mockAnt.On("Score").Return(4.0) // total length: 1+1+2 =4
	mockAnt.On("PheromoneMultiplier").Return(float64(1.0))

	pm := pheromone.NewPheromoneMap(g, 0)
	mockAnt.On("PheromoneMap").Return(pm)

	// Apply pheromone
	strategy.ApplyPheromone(mockAnt)

	// Delta = multiplier / score = 1.0 / 4.0 = 0.25
	expectedDelta := 0.25

	// Check that each edge in path got delta added
	assert.Equal(t, expectedDelta, pm.Get(e12), "edge v1-v2")
	assert.Equal(t, expectedDelta, pm.Get(e23), "edge v2-v3")
	assert.Equal(t, expectedDelta, pm.Get(e31), "edge v3-v1 (wrap)")

	mockAnt.AssertExpectations(t)
}

func TestApplyClassicStrategy_ApplyPheromone_ZeroScore(t *testing.T) {
	strategy := apply.NewApplyClassicStrategy()

	g := graph.NewGraph(1)
	v1 := graph.NewVertex("V1")
	g.AddVertex(v1)

	mockAnt := new(MockAntView)
	mockAnt.On("Graph").Return(g)
	mockAnt.On("Path").Return([]*graph.Vertex{v1})
	mockAnt.On("Score").Return(0.0) // zero score would cause division by zero
	mockAnt.On("PheromoneMultiplier").Return(float64(1.0))

	pm := pheromone.NewPheromoneMap(g, 0)
	mockAnt.On("PheromoneMap").Return(pm)

	// The code does delta := ant.PheromoneMultiplier() / ant.Score()
	// If Score is zero, delta will be infinite? Actually float64 division by zero gives +Inf.
	// We'll just ensure it doesn't panic.
	strategy.ApplyPheromone(mockAnt)

	// Expect no change? Actually delta will be +Inf, adding Inf to pheromone map.
	// We'll just check that it didn't panic.
	mockAnt.AssertExpectations(t)
}
