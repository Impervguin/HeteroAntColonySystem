package graph

type Edge struct {
	source *Vertex
	target *Vertex
	weight float64
}

func NewEdge(source *Vertex, target *Vertex, weight float64) *Edge {
	return &Edge{
		source: source,
		target: target,
		weight: weight,
	}
}

func (e *Edge) Source() *Vertex {
	return e.source
}

func (e *Edge) Target() *Vertex {
	return e.target
}

func (e *Edge) Weight() float64 {
	return e.weight
}

func (e *Edge) UpdateWeight(weight float64) {
	e.weight = weight
}

func (e *Edge) Reverse() *Edge {
	return NewEdge(e.target, e.source, e.weight)
}
