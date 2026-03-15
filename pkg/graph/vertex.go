package graph

import "github.com/google/uuid"

type Vertex struct {
	id   uuid.UUID
	name string
}

func NewVertex(name string) *Vertex {
	return &Vertex{
		id:   uuid.New(),
		name: name,
	}
}

//go:inline
func (v *Vertex) ID() uuid.UUID {
	return v.id
}

//go:inline
func (v *Vertex) Name() string {
	return v.name
}

//go:inline
func (v *Vertex) UpdateName(name string) {
	v.name = name
}
