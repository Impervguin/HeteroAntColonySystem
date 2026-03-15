package graph

import (
	"sync"
)

type Graph struct {
	vertices map[*Vertex]struct{}
	edges    []*Edge
	// sourceMap maps a vertex to all edges that have it as source
	sourceMap map[*Vertex][]*Edge
	// targetMap maps a vertex to all edges that have it as target
	targetMap map[*Vertex][]*Edge
	dist      map[*Vertex]map[*Vertex]*Edge

	mut sync.RWMutex
}

// vertexCount is the number of vertices planed in the graph
func NewGraph(vertCount uint) *Graph {
	g := &Graph{
		vertices:  make(map[*Vertex]struct{}, vertCount),
		edges:     make([]*Edge, 0, vertCount*vertCount),
		sourceMap: make(map[*Vertex][]*Edge, vertCount),
		targetMap: make(map[*Vertex][]*Edge, vertCount),
		dist:      make(map[*Vertex]map[*Vertex]*Edge, vertCount),
	}
	return g
}

// AddVertex adds a vertex to the graph
func (g *Graph) AddVertex(v *Vertex) error {
	g.mut.Lock()
	defer g.mut.Unlock()

	g.vertices[v] = struct{}{}
	g.sourceMap[v] = make([]*Edge, 0)
	g.targetMap[v] = make([]*Edge, 0)
	return nil
}

// AddEdge adds an edge to the graph
func (g *Graph) AddEdge(weight float64, source *Vertex, target *Vertex) error {
	g.mut.Lock()
	defer g.mut.Unlock()

	if _, ok := g.vertices[source]; !ok {
		return ErrVertexNotFound
	}

	if _, ok := g.vertices[target]; !ok {
		return ErrVertexNotFound
	}

	edge := NewEdge(source, target, weight)
	g.edges = append(g.edges, edge)
	g.sourceMap[source] = append(g.sourceMap[source], edge)
	g.targetMap[target] = append(g.targetMap[target], edge)
	if _, ok := g.dist[source]; !ok {
		g.dist[source] = make(map[*Vertex]*Edge)
	}
	g.dist[source][target] = edge
	return nil
}

// ForEachVertex iterates over all vertices of the graph
// f should not call any methods on the graph
func (g *Graph) ForEachVertex(f func(v *Vertex)) {
	g.mut.RLock()
	defer g.mut.RUnlock()

	for v := range g.vertices {
		f(v)
	}
}

// Vertices returns the vertices of the graph
// as a channel
func (g *Graph) VerticesChan() <-chan *Vertex {
	ch := make(chan *Vertex)

	go func() {
		g.mut.RLock()
		defer g.mut.RUnlock()

		for v := range g.vertices {
			ch <- v
		}
		close(ch)
	}()

	return ch
}

// ForEachEdge iterates over all edges of the graph
// f should not call any methods on the graph
func (g *Graph) ForEachEdge(f func(e *Edge)) {
	g.mut.RLock()
	defer g.mut.RUnlock()

	for _, e := range g.edges {
		f(e)
	}
}

// Edges returns the edges of the graph as a channel
func (g *Graph) EdgesChan() <-chan *Edge {
	ch := make(chan *Edge)

	go func() {
		g.mut.RLock()
		defer g.mut.RUnlock()

		for _, e := range g.edges {
			ch <- e
		}
		close(ch)
	}()

	return ch
}

// ForEachSource iterates over all edges that have the given vertex as source
// f should not call any methods on the graph
func (g *Graph) ForEachSource(v *Vertex, f func(v *Edge)) {
	g.mut.RLock()
	defer g.mut.RUnlock()

	edges, ok := g.sourceMap[v]
	if !ok {
		return
	}

	for _, e := range edges {
		f(e)
	}
}

// SourceChan returns the edges that have the given vertex as source
func (g *Graph) SourceChan(v *Vertex) <-chan *Edge {
	ch := make(chan *Edge)

	go func() {
		g.mut.RLock()
		defer g.mut.RUnlock()

		for _, e := range g.sourceMap[v] {
			ch <- e
		}
		close(ch)
	}()

	return ch
}

// ForEachTarget iterates over all edges that have the given vertex as target
// f should not call any methods on the graph
func (g *Graph) ForEachTarget(v *Vertex, f func(v *Edge)) {
	g.mut.RLock()
	defer g.mut.RUnlock()

	edges, ok := g.targetMap[v]
	if !ok {
		return
	}

	for _, e := range edges {
		f(e)
	}
}

// TargetChan returns the edges that have the given vertex as target
func (g *Graph) TargetChan(v *Vertex) <-chan *Edge {
	ch := make(chan *Edge)

	go func() {
		g.mut.RLock()
		defer g.mut.RUnlock()

		for _, e := range g.targetMap[v] {
			ch <- e
		}
		close(ch)
	}()

	return ch
}

func (g *Graph) Edge(source *Vertex, target *Vertex) (*Edge, bool) {
	g.mut.RLock()
	defer g.mut.RUnlock()

	edge, ok := g.dist[source][target]
	return edge, ok
}
