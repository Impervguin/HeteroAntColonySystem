package graph

import "github.com/google/uuid"

type Vertex struct {
	id       uuid.UUID
	name     string
	metadata any
}

func NewVertex(name string) *Vertex {
	return &Vertex{
		id:   uuid.New(),
		name: name,
	}
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
