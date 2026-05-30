package apply_test

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"

	"github.com/stretchr/testify/mock"
)

// MockAntView is a mock implementation of ant.AntView for testing
type MockAntView struct {
	mock.Mock
}

func (m *MockAntView) Graph() *graph.Graph {
	args := m.Called()
	return args.Get(0).(*graph.Graph)
}
func (m *MockAntView) PheromoneMap() *pheromone.PheromoneMap {
	args := m.Called()
	return args.Get(0).(*pheromone.PheromoneMap)
}
func (m *MockAntView) Current() *graph.Vertex {
	args := m.Called()
	return args.Get(0).(*graph.Vertex)
}
func (m *MockAntView) Visited(v *graph.Vertex) bool {
	args := m.Called(v)
	return args.Bool(0)
}
func (m *MockAntView) Path() []*graph.Vertex {
	args := m.Called()
	return args.Get(0).([]*graph.Vertex)
}
func (m *MockAntView) Score() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}
func (m *MockAntView) SumScore() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}
func (m *MockAntView) Alpha() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}
func (m *MockAntView) Beta() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}
func (m *MockAntView) PheromoneMultiplier() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}
func (m *MockAntView) PathStrategy() ant.PathChoiceStrategy {
	args := m.Called()
	return args.Get(0).(ant.PathChoiceStrategy)
}
func (m *MockAntView) PheromoneApplyStrategy() ant.PheromoneApplyingStrategy {
	args := m.Called()
	return args.Get(0).(ant.PheromoneApplyingStrategy)
}
