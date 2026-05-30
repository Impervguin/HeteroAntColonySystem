package optimisation

import (
	"HeteroAntColonySystem/pkg/graph"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoOpLocalOptimisation_Optimise(t *testing.T) {
	strategy := NewNoOpLocalOptimisation()

	// Create a simple graph with three vertices in a triangle
	g := graph.NewGraph(3)
	vA := graph.NewVertex("A")
	vB := graph.NewVertex("B")
	vC := graph.NewVertex("C")
	g.AddVertex(vA)
	g.AddVertex(vB)
	g.AddVertex(vC)
	// Add edges with weight 1.0
	g.AddEdge(1.0, vA, vB)
	g.AddEdge(1.0, vB, vA)
	g.AddEdge(1.0, vB, vC)
	g.AddEdge(1.0, vC, vB)
	g.AddEdge(1.0, vC, vA)
	g.AddEdge(1.0, vA, vC)

	// Create a path that is not optimal (A->C->B)
	path := []*graph.Vertex{vA, vC, vB}
	original := make([]*graph.Vertex, len(path))
	copy(original, path)

	// Apply optimisation
	strategy.Optimise(path, g)

	// Path should remain unchanged
	assert.Equal(t, original, path, "path should not be changed by NoOp")
}

// Test with nil path and nil graph (should not panic)
func TestNoOpLocalOptimisation_NilInputs(t *testing.T) {
	strategy := NewNoOpLocalOptimisation()
	strategy.Optimise(nil, nil)
	// No assertion needed, just ensure no panic
}