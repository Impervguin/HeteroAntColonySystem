package graph

import "github.com/google/uuid"

type Vertex struct {
	id       uuid.UUID
	name     string
	metadata any
}

func NewVertex(name string, opts ...VertexOption) *Vertex {
	v := &Vertex{
		id:   uuid.New(),
		name: name,
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

func (v *Vertex) ID() uuid.UUID {
	return v.id
}

func (v *Vertex) Name() string {
	return v.name
}

func (v *Vertex) UpdateName(name string) {
	v.name = name
}

func (v *Vertex) Metadata() any {
	return v.metadata
}

func (v *Vertex) SetMetadata(metadata any) {
	v.metadata = metadata
}

type VertexOption func(v *Vertex)

func WithMetadata(metadata any) VertexOption {
	return func(v *Vertex) {
		v.metadata = metadata
	}
}

func WithID(id uuid.UUID) VertexOption {
	return func(v *Vertex) {
		v.id = id
	}
}
