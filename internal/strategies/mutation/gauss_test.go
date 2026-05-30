package mutation

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/strategies/apply"
	"HeteroAntColonySystem/internal/strategies/path"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockAntViewGauss for mutation testing
type MockAntViewGauss struct {
	alphaVal          float64
	betaVal           float64
	pheromoneVal      float64
	pathStrategy      ant.PathChoiceStrategy
	pheromoneStrategy ant.PheromoneApplyingStrategy
}

func (m *MockAntViewGauss) Alpha() float64 {
	return m.alphaVal
}
func (m *MockAntViewGauss) Beta() float64 {
	return m.betaVal
}
func (m *MockAntViewGauss) PheromoneMultiplier() float64 {
	return m.pheromoneVal
}
func (m *MockAntViewGauss) PathStrategy() ant.PathChoiceStrategy {
	return m.pathStrategy
}
func (m *MockAntViewGauss) PheromoneApplyStrategy() ant.PheromoneApplyingStrategy {
	return m.pheromoneStrategy
}

// Unused methods
func (m *MockAntViewGauss) Graph() *graph.Graph                   { return nil }
func (m *MockAntViewGauss) PheromoneMap() *pheromone.PheromoneMap { return nil }
func (m *MockAntViewGauss) Current() *graph.Vertex                { return nil }
func (m *MockAntViewGauss) Visited(v *graph.Vertex) bool          { return false }
func (m *MockAntViewGauss) Path() []*graph.Vertex                 { return nil }
func (m *MockAntViewGauss) Score() float64                        { return 0 }
func (m *MockAntViewGauss) SumScore() float64                     { return 0 }

func TestGaussMutationStrategy_Mutate(t *testing.T) {
	sigma := 2.0
	mu := 1.0
	strategy := NewGaussMutationStrategy(sigma, mu)

	parent := &MockAntViewUniform{
		alphaVal:          5.0,
		betaVal:           3.0,
		pheromoneVal:      1.5,
		pathStrategy:      path.NewPahtClassicStrategy(),
		pheromoneStrategy: apply.NewApplyClassicStrategy(),
	}
	// Collect deltas
	var alphaDeltas, betaDeltas []float64
	samples := 50000
	for i := 0; i < samples; i++ {
		child := strategy.Mutate(parent)
		assert.NotNil(t, child)

		alphaDeltas = append(alphaDeltas, child.Alpha()-parent.Alpha())
		betaDeltas = append(betaDeltas, child.Beta()-parent.Beta())

		// Inheritance
		assert.InDelta(t, parent.PheromoneMultiplier(), child.PheromoneMultiplier(), 1e-9, "pheromone multiplier")
		assert.Equal(t, parent.PathStrategy(), child.PathStrategy())
		assert.Equal(t, parent.PheromoneApplyStrategy(), child.PheromoneApplyStrategy())
	}

	// Compute sample mean and variance for alpha
	meanAlpha := 0.0
	for _, v := range alphaDeltas {
		meanAlpha += v
	}
	meanAlpha /= float64(len(alphaDeltas))

	var varianceAlpha float64
	for _, v := range alphaDeltas {
		diff := v - meanAlpha
		varianceAlpha += diff * diff
	}
	varianceAlpha /= float64(len(alphaDeltas))

	// Same for beta
	meanBeta := 0.0
	for _, v := range betaDeltas {
		meanBeta += v
	}
	meanBeta /= float64(len(betaDeltas))

	var varianceBeta float64
	for _, v := range betaDeltas {
		diff := v - meanBeta
		varianceBeta += diff * diff
	}
	varianceBeta /= float64(len(betaDeltas))

	// Assert mean close to mu, variance close to sigma^2
	toleranceMean := 0.1 // allowable error in mean
	toleranceVar := 0.2  // allowable error in variance
	assert.InDelta(t, mu, meanAlpha, toleranceMean, "alpha mean")
	assert.InDelta(t, sigma*sigma, varianceAlpha, toleranceVar, "alpha variance")
	assert.InDelta(t, mu, meanBeta, toleranceMean, "beta mean")
	assert.InDelta(t, sigma*sigma, varianceBeta, toleranceVar, "beta variance")
}

// Test zero sigma (should return parent values exactly)
func TestGaussMutationStrategy_ZeroSigma(t *testing.T) {
	sigma := 0.0
	mu := 5.0
	strategy := NewGaussMutationStrategy(sigma, mu)

	parent := &MockAntViewUniform{
		alphaVal:          5.0,
		betaVal:           3.0,
		pheromoneVal:      1.5,
		pathStrategy:      path.NewPahtClassicStrategy(),
		pheromoneStrategy: apply.NewApplyClassicStrategy(),
	}

	// With sigma=0, normalRand returns mu always? Actually rand.NormFloat64()*0 + mu = mu.
	// So alpha' = alpha + mu, beta' = beta + mu.
	child := strategy.Mutate(parent)
	assert.NotNil(t, child)
	assert.Equal(t, parent.Alpha()+mu, child.Alpha(), "alpha")
	assert.Equal(t, parent.Beta()+mu, child.Beta(), "beta")
	// Inheritance
	assert.InDelta(t, parent.PheromoneMultiplier(), child.PheromoneMultiplier(), 1e-9, "pheromone multiplier")
	assert.Equal(t, parent.PathStrategy(), child.PathStrategy())
	assert.Equal(t, parent.PheromoneApplyStrategy(), child.PheromoneApplyStrategy())
}
