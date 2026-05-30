package crossover_test

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"

	"github.com/stretchr/testify/mock"
)

// MockAntView is a minimal mock of ant.AntView for crossover testing
type MockAntView struct {
	mock.Mock
	alphaVal          float64
	betaVal           float64
	pheromoneVal      float64
	pathStrategy      ant.PathChoiceStrategy
	pheromoneStrategy ant.PheromoneApplyingStrategy
}

func NewMockAntView(alphaVal, betaVal, pheromoneVal float64) *MockAntView {
	return &MockAntView{
		alphaVal:          alphaVal,
		betaVal:           betaVal,
		pheromoneVal:      pheromoneVal,
		pathStrategy:      &MockPathStrategy{},
		pheromoneStrategy: &MockPheromoneStrategy{},
	}
}

func (m *MockAntView) Alpha() float64 {
	return m.alphaVal
}
func (m *MockAntView) Beta() float64 {
	return m.betaVal
}
func (m *MockAntView) PheromoneMultiplier() float64 {
	return m.pheromoneVal
}
func (m *MockAntView) PathStrategy() ant.PathChoiceStrategy {
	return m.pathStrategy
}
func (m *MockAntView) PheromoneApplyStrategy() ant.PheromoneApplyingStrategy {
	return m.pheromoneStrategy
}

// Unused methods for crossover - return dummy values
func (m *MockAntView) Graph() *graph.Graph                   { return nil }
func (m *MockAntView) PheromoneMap() *pheromone.PheromoneMap { return nil }
func (m *MockAntView) Current() *graph.Vertex                { return nil }
func (m *MockAntView) Visited(v *graph.Vertex) bool          { return false }
func (m *MockAntView) Path() []*graph.Vertex                 { return nil }
func (m *MockAntView) Score() float64                        { return 0 }
func (m *MockAntView) SumScore() float64                     { return 0 }

type MockPheromoneStrategy struct {
	mock.Mock
}

var _ ant.PheromoneApplyingStrategy = &MockPheromoneStrategy{}

func (m *MockPheromoneStrategy) ApplyPheromone(ant ant.AntView) {
	m.Called(ant)
}

type MockPathStrategy struct {
	mock.Mock
}

var _ ant.PathChoiceStrategy = &MockPathStrategy{}

func (m *MockPathStrategy) ChooseNext(ant ant.AntView) *graph.Vertex {
	args := m.Called(ant)
	return args.Get(0).(*graph.Vertex)
}
