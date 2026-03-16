package graph

import (
	"testing"
)

func TestVertexCreation(t *testing.T) {
	v := NewVertex("A")

	if v.Name() != "A" {
		t.Fatalf("expected name A, got %s", v.Name())
	}

	if v.ID() == [16]byte{} {
		t.Fatal("expected non-empty UUID")
	}
}

func TestVertexUpdateName(t *testing.T) {
	v := NewVertex("A")
	v.UpdateName("B")

	if v.Name() != "B" {
		t.Fatalf("expected name B, got %s", v.Name())
	}
}

func TestEdgeCreation(t *testing.T) {
	v1 := NewVertex("A")
	v2 := NewVertex("B")

	e := NewEdge(v1, v2, 10)

	if e.Source() != v1 {
		t.Fatal("wrong source vertex")
	}

	if e.Target() != v2 {
		t.Fatal("wrong target vertex")
	}

	if e.Weight() != 10 {
		t.Fatalf("expected weight 10 got %f", e.Weight())
	}
}

func TestEdgeUpdateWeight(t *testing.T) {
	v1 := NewVertex("A")
	v2 := NewVertex("B")

	e := NewEdge(v1, v2, 10)
	e.UpdateWeight(20)

	if e.Weight() != 20 {
		t.Fatalf("expected weight 20 got %f", e.Weight())
	}
}

func TestEdgeReverse(t *testing.T) {
	v1 := NewVertex("A")
	v2 := NewVertex("B")

	e := NewEdge(v1, v2, 5)
	r := e.Reverse()

	if r.Source() != v2 || r.Target() != v1 {
		t.Fatal("reverse edge incorrect")
	}

	if r.Weight() != 5 {
		t.Fatal("reverse edge weight mismatch")
	}
}

func TestGraphAddVertex(t *testing.T) {
	g := NewGraph(2)

	v := NewVertex("A")

	err := g.AddVertex(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	count := 0
	g.ForEachVertex(func(v *Vertex) bool {
		count++
		return false
	})

	if count != 1 {
		t.Fatalf("expected 1 vertex got %d", count)
	}
}

func TestGraphAddEdge(t *testing.T) {
	g := NewGraph(2)

	v1 := NewVertex("A")
	v2 := NewVertex("B")

	g.AddVertex(v1)
	g.AddVertex(v2)

	err := g.AddEdge(10, v1, v2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	count := 0
	g.ForEachEdge(func(e *Edge) bool {
		count++
		return false
	})

	if count != 1 {
		t.Fatalf("expected 1 edge got %d", count)
	}
}

func TestGraphAddEdgeVertexNotFound(t *testing.T) {
	g := NewGraph(2)

	v1 := NewVertex("A")
	v2 := NewVertex("B")

	g.AddVertex(v1)

	err := g.AddEdge(10, v1, v2)

	if err != ErrVertexNotFound {
		t.Fatalf("expected ErrVertexNotFound got %v", err)
	}
}

func TestGraphEdgeLookup(t *testing.T) {
	g := NewGraph(2)

	v1 := NewVertex("A")
	v2 := NewVertex("B")

	g.AddVertex(v1)
	g.AddVertex(v2)

	g.AddEdge(7, v1, v2)

	e, ok := g.Edge(v1, v2)

	if !ok {
		t.Fatal("edge should exist")
	}

	if e.Weight() != 7 {
		t.Fatalf("expected weight 7 got %f", e.Weight())
	}
}

func TestGraphForEachSource(t *testing.T) {
	g := NewGraph(3)

	v1 := NewVertex("A")
	v2 := NewVertex("B")
	v3 := NewVertex("C")

	g.AddVertex(v1)
	g.AddVertex(v2)
	g.AddVertex(v3)

	g.AddEdge(1, v1, v2)
	g.AddEdge(2, v1, v3)

	count := 0

	g.ForEachSource(v1, func(e *Edge) bool {
		count++
		return false
	})

	if count != 2 {
		t.Fatalf("expected 2 edges got %d", count)
	}
}

func TestGraphForEachTarget(t *testing.T) {
	g := NewGraph(3)

	v1 := NewVertex("A")
	v2 := NewVertex("B")
	v3 := NewVertex("C")

	g.AddVertex(v1)
	g.AddVertex(v2)
	g.AddVertex(v3)

	g.AddEdge(1, v1, v3)
	g.AddEdge(2, v2, v3)

	count := 0

	g.ForEachTarget(v3, func(e *Edge) bool {
		count++
		return false
	})

	if count != 2 {
		t.Fatalf("expected 2 edges got %d", count)
	}
}

func TestVerticesChan(t *testing.T) {
	g := NewGraph(2)

	v1 := NewVertex("A")
	v2 := NewVertex("B")

	g.AddVertex(v1)
	g.AddVertex(v2)

	count := 0
	for range g.VerticesChan() {
		count++
	}

	if count != 2 {
		t.Fatalf("expected 2 vertices got %d", count)
	}
}

func TestEdgesChan(t *testing.T) {
	g := NewGraph(2)

	v1 := NewVertex("A")
	v2 := NewVertex("B")

	g.AddVertex(v1)
	g.AddVertex(v2)

	g.AddEdge(1, v1, v2)

	count := 0
	for range g.EdgesChan() {
		count++
	}

	if count != 1 {
		t.Fatalf("expected 1 edge got %d", count)
	}
}

func TestSourceChan(t *testing.T) {
	g := NewGraph(3)

	v1 := NewVertex("A")
	v2 := NewVertex("B")
	v3 := NewVertex("C")

	g.AddVertex(v1)
	g.AddVertex(v2)
	g.AddVertex(v3)

	g.AddEdge(1, v1, v2)
	g.AddEdge(2, v1, v3)

	count := 0
	for range g.SourceChan(v1) {
		count++
	}

	if count != 2 {
		t.Fatalf("expected 2 edges got %d", count)
	}
}

func TestTargetChan(t *testing.T) {
	g := NewGraph(3)

	v1 := NewVertex("A")
	v2 := NewVertex("B")
	v3 := NewVertex("C")

	g.AddVertex(v1)
	g.AddVertex(v2)
	g.AddVertex(v3)

	g.AddEdge(1, v1, v3)
	g.AddEdge(2, v2, v3)

	count := 0
	for range g.TargetChan(v3) {
		count++
	}

	if count != 2 {
		t.Fatalf("expected 2 edges got %d", count)
	}
}
