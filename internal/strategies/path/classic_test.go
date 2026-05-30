package path_test

import (
	"HeteroAntColonySystem/internal/strategies/path"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathClassicStrategy_ChooseNext(t *testing.T) {
	strategy := path.NewPahtClassicStrategy()

	// Create a simple graph with two vertices and one edge
	g := graph.NewGraph(2)
	v1 := graph.NewVertex("V1")
	v2 := graph.NewVertex("V2")
	g.AddVertex(v1)
	g.AddVertex(v2)
	// Add edge with weight 1.0
	if err := g.AddEdge(1.0, v1, v2); err != nil {
		t.Fatalf("Failed to add edge: %v", err)
	}
	if err := g.AddEdge(1.0, v2, v1); err != nil {
		t.Fatalf("Failed to add reverse edge: %v", err)
	}

	// Create pheromone map with initial 0, then set edge (v1,v2) to 1.0
	pm := pheromone.NewPheromoneMap(g, 0.0)
	e12, _ := g.Edge(v1, v2)
	pm.Add(e12, 1.0)

	stub := NewAntViewStub(g, pm, v1, []*graph.Vertex{v1}, map[*graph.Vertex]struct{}{v1: {}}, 1.0, 1.0, 1.0, strategy, 0.0, 0.0)

	next := strategy.ChooseNext(stub)

	// Expect v2 as the only unvisited neighbor
	assert.Equal(t, v2, next)
}

func TestPathClassicStrategy_ChooseNext_NoUnvisited(t *testing.T) {
	strategy := path.NewPahtClassicStrategy()

	g := graph.NewGraph(1)
	v1 := graph.NewVertex("V1")
	g.AddVertex(v1)

	pm := pheromone.NewPheromoneMap(g, 0.0)

	stub := NewAntViewStub(g, pm, v1, []*graph.Vertex{v1}, map[*graph.Vertex]struct{}{v1: {}}, 1.0, 1.0, 1.0, strategy, 0.0, 0.0)

	next := strategy.ChooseNext(stub)
	assert.Nil(t, next)
}

func TestPathClassicStrategy_ChooseNext_Probability(t *testing.T) {
	strategy := path.NewPahtClassicStrategy()

	// Graph with three vertices: current v1, neighbors v2 and v3
	g := graph.NewGraph(3)
	v1 := graph.NewVertex("V1")
	v2 := graph.NewVertex("V2")
	v3 := graph.NewVertex("V3")
	g.AddVertex(v1)
	g.AddVertex(v2)
	g.AddVertex(v3)

	// Edges with weights: v1-v2 weight 1, v1-v3 weight 2
	e12, _ := g.Edge(v1, v2)
	if e12 == nil {
		if err := g.AddEdge(1.0, v1, v2); err != nil {
			t.Fatalf("Failed to add edge v1-v2: %v", err)
		}
		e12, _ = g.Edge(v1, v2)
	}
	e13, _ := g.Edge(v1, v3)
	if e13 == nil {
		if err := g.AddEdge(2.0, v1, v3); err != nil {
			t.Fatalf("Failed to add edge v1-v3: %v", err)
		}
		e13, _ = g.Edge(v1, v3)
	}
	// Add reverse edges
	g.AddEdge(1.0, v2, v1)
	g.AddEdge(2.0, v3, v1)

	pm := pheromone.NewPheromoneMap(g, 0.0)
	// Set pheromone levels: both edges have pheromone 1.0
	pm.Add(e12, 1.0)
	pm.Add(e13, 1.0)

	stub := NewAntViewStub(g, pm, v1, []*graph.Vertex{v1}, map[*graph.Vertex]struct{}{v1: {}}, 1.0, 1.0, 1.0, strategy, 0.0, 0.0)

	// Run many times to check distribution
	counts := map[*graph.Vertex]int{v2: 0, v3: 0}
	samples := 10000
	for i := 0; i < samples; i++ {
		next := strategy.ChooseNext(stub)
		if next == v2 {
			counts[v2]++
		} else if next == v3 {
			counts[v3]++
		} else {
			t.Errorf("Unexpected next vertex: %v", next)
		}
	}
	// Expected probability proportional to (1/weight^beta) * pheromone^alpha
	// weight v1-v2 = 1 => 1/1^1 = 1
	// weight v1-v3 = 2 => 1/2^1 = 0.5
	// pheromone both 1 => 1^1 =1
	// So probabilities: 1 and 0.5 => normalized: v2: 1/(1+0.5)=2/3, v3: 0.5/(1+0.5)=1/3
	expectedV2 := float64(samples) * 2.0 / 3.0
	expectedV3 := float64(samples) * 1.0 / 3.0
	tolerance := float64(samples) * 0.05 // 5% tolerance
	if diff := math.Abs(float64(counts[v2]) - expectedV2); diff > tolerance {
		t.Errorf("v2 count out of range: got %d, expected %.2f, diff %.2f", counts[v2], expectedV2, diff)
	}
	if diff := math.Abs(float64(counts[v3]) - expectedV3); diff > tolerance {
		t.Errorf("v3 count out of range: got %d, expected %.2f, diff %.2f", counts[v3], expectedV3, diff)
	}
}
