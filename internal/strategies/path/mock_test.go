package path_test

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"

	"github.com/stretchr/testify/mock"
)

type MockPheromoneStrategy struct {
	mock.Mock
}

var _ ant.PheromoneApplyingStrategy = &MockPheromoneStrategy{}

func (m *MockPheromoneStrategy) ApplyPheromone(ant ant.AntView) {
	m.Called(ant)
}

// AntViewStub is a stub implementation of ant.AntView for testing
type AntViewStub struct {
	g                *graph.Graph
	pm               *pheromone.PheromoneMap
	current          *graph.Vertex
	path             []*graph.Vertex
	visitedMap       map[*graph.Vertex]struct{}
	alphaVal         float64
	betaVal          float64
	pheromoneMultVal float64
	pathChoice       ant.PathChoiceStrategy
	pheromoneApply   ant.PheromoneApplyingStrategy
	scoreVal         float64
	sumScoreVal      float64
}

func NewAntViewStub(g *graph.Graph, pm *pheromone.PheromoneMap, current *graph.Vertex, path []*graph.Vertex, visited map[*graph.Vertex]struct{}, alpha, beta, pheromoneMult float64, pathChoice ant.PathChoiceStrategy, score, sumScore float64) *AntViewStub {
	return &AntViewStub{
		g:                g,
		pm:               pm,
		current:          current,
		path:             path,
		visitedMap:       visited,
		alphaVal:         alpha,
		betaVal:          beta,
		pheromoneMultVal: pheromoneMult,
		pathChoice:       pathChoice,
		pheromoneApply:   &MockPheromoneStrategy{},
		scoreVal:         score,
		sumScoreVal:      sumScore,
	}
}

func (s *AntViewStub) Graph() *graph.Graph {
	return s.g
}
func (s *AntViewStub) PheromoneMap() *pheromone.PheromoneMap {
	return s.pm
}
func (s *AntViewStub) Current() *graph.Vertex {
	return s.current
}
func (s *AntViewStub) Visited(v *graph.Vertex) bool {
	_, ok := s.visitedMap[v]
	return ok
}
func (s *AntViewStub) Path() []*graph.Vertex {
	return s.path
}
func (s *AntViewStub) Score() float64 {
	return s.scoreVal
}
func (s *AntViewStub) SumScore() float64 {
	return s.sumScoreVal
}
func (s *AntViewStub) Alpha() float64 {
	return s.alphaVal
}
func (s *AntViewStub) Beta() float64 {
	return s.betaVal
}
func (s *AntViewStub) PheromoneMultiplier() float64 {
	return s.pheromoneMultVal
}
func (s *AntViewStub) PathStrategy() ant.PathChoiceStrategy {
	return s.pathChoice
}
func (s *AntViewStub) PheromoneApplyStrategy() ant.PheromoneApplyingStrategy {
	return s.pheromoneApply
}
