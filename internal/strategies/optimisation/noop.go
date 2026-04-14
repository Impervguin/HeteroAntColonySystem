package optimisation

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/pkg/graph"
)

type NoOpLocalOptimisation struct {
}

func NewNoOpLocalOptimisation() *NoOpLocalOptimisation {
	return &NoOpLocalOptimisation{}
}

var _ ant.LocalOptimisationStrategy = &NoOpLocalOptimisation{}

func (s *NoOpLocalOptimisation) Optimise(_ []*graph.Vertex, _ *graph.Graph) {
}
