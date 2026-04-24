package tsplib

import (
	"HeteroAntColonySystem/pkg/graph"
	"io"
)

// TSPLIBAdapter interface for parsing different TSPLIB formats
type TSPLIBAdapter interface {
	// Parse reads the data section and returns vertices and edges
	Parse(r io.Reader, problem *Problem) ([]*graph.Vertex, []*graph.Edge, error)

	// CanHandle returns true if adapter can handle this combination
	CanHandle(weightType string, weightFormat string) bool

	// Name returns adapter name
	Name() string

	// MetadataType returns the type of metadata stored in the vertices
	MetadataType() any
}

type AdapterRegistry interface {
	Get(weightType string, weightFormat string) TSPLIBAdapter
}
