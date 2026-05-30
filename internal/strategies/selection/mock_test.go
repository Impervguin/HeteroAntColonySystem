package selection_test

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"

	"github.com/stretchr/testify/mock"
)

// MockAntView for selection testing - only SumScore is used
type MockAntView struct {
	mock.Mock
	sumScoreVal float64
	name        string
}

func (m *MockAntView) SumScore() float64 {
	return m.sumScoreVal
}

// Unused methods - return dummy values
func (m *MockAntView) Graph() *graph.Graph                                   { return nil }
func (m *MockAntView) PheromoneMap() *pheromone.PheromoneMap                 { return nil }
func (m *MockAntView) Current() *graph.Vertex                                { return nil }
func (m *MockAntView) Visited(v *graph.Vertex) bool                          { return false }
func (m *MockAntView) Path() []*graph.Vertex                                 { return nil }
func (m *MockAntView) Score() float64                                        { return 0 }
func (m *MockAntView) Alpha() float64                                        { return 0 }
func (m *MockAntView) Beta() float64                                         { return 0 }
func (m *MockAntView) PheromoneMultiplier() float64                          { return 0 }
func (m *MockAntView) PathStrategy() ant.PathChoiceStrategy                  { return nil }
func (m *MockAntView) PheromoneApplyStrategy() ant.PheromoneApplyingStrategy { return nil }
