package mutation

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/strategies/apply"
	"HeteroAntColonySystem/internal/strategies/path"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockAntViewUniform for mutation testing
type MockAntViewUniform struct {
	alphaVal          float64
	betaVal           float64
	pheromoneVal      float64
	pathStrategy      ant.PathChoiceStrategy
	pheromoneStrategy ant.PheromoneApplyingStrategy
}

func (m *MockAntViewUniform) Alpha() float64 {
	return m.alphaVal
}
func (m *MockAntViewUniform) Beta() float64 {
	return m.betaVal
}
func (m *MockAntViewUniform) PheromoneMultiplier() float64 {
	return m.pheromoneVal
}
func (m *MockAntViewUniform) PathStrategy() ant.PathChoiceStrategy {
	return m.pathStrategy
}
func (m *MockAntViewUniform) PheromoneApplyStrategy() ant.PheromoneApplyingStrategy {
	return m.pheromoneStrategy
}

// Unused methods
func (m *MockAntViewUniform) Graph() *graph.Graph                   { return nil }
func (m *MockAntViewUniform) PheromoneMap() *pheromone.PheromoneMap { return nil }
func (m *MockAntViewUniform) Current() *graph.Vertex                { return nil }
func (m *MockAntViewUniform) Visited(v *graph.Vertex) bool          { return false }
func (m *MockAntViewUniform) Path() []*graph.Vertex                 { return nil }
func (m *MockAntViewUniform) Score() float64                        { return 0 }
func (m *MockAntViewUniform) SumScore() float64                     { return 0 }

func TestUniformMutationStrategy_Mutate(t *testing.T) {
	l, r := -1.0, 2.0
	strategy := NewUniformMutationStrategy(l, r)

	parent := &MockAntViewUniform{
		alphaVal:          5.0,
		betaVal:           3.0,
		pheromoneVal:      1.5,
		pathStrategy:      path.NewPahtClassicStrategy(),
		pheromoneStrategy: apply.NewApplyClassicStrategy(),
	}

	// Run many samples to verify range
	minAlpha := math.Inf(1)
	maxAlpha := math.Inf(-1)
	minBeta := math.Inf(1)
	maxBeta := math.Inf(-1)
	samples := 10000

	for i := 0; i < samples; i++ {
		child := strategy.Mutate(parent)
		assert.NotNil(t, child)

		a := child.Alpha()
		b := child.Beta()
		if a < minAlpha {
			minAlpha = a
		}
		if a > maxAlpha {
			maxAlpha = a
		}
		if b < minBeta {
			minBeta = b
		}
		if b > maxBeta {
			maxBeta = b
		}

		// Check inheritance
		assert.InDelta(t, parent.PheromoneMultiplier(), child.PheromoneMultiplier(), 1e-9, "pheromone multiplier")
		assert.Equal(t, parent.PathStrategy(), child.PathStrategy())
		assert.Equal(t, parent.PheromoneApplyStrategy(), child.PheromoneApplyStrategy())
	}

	// Expected raw range before clamping: [alpha + l, alpha + r)
	// alpha=5, l=-1, r=2 => [4, 7)
	expectedMinAlpha := parent.Alpha() + l
	expectedMaxAlpha := parent.Alpha() + r
	// Since l negative but alpha+l positive, no clamping expected.
	assert.InDelta(t, expectedMinAlpha, minAlpha, 0.1, "min alpha out of range")
	assert.InDelta(t, expectedMaxAlpha, maxAlpha, 0.1, "max alpha out of range")

	// Beta similarly
	expectedMinBeta := parent.Beta() + l
	expectedMaxBeta := parent.Beta() + r
	assert.InDelta(t, expectedMinBeta, minBeta, 0.1, "min beta out of range")
	assert.InDelta(t, expectedMaxBeta, maxBeta, 0.1, "max beta out of range")
}

// Test clamping to zero when negative
func TestUniformMutationStrategy_ClampZero(t *testing.T) {
	l, r := -10.0, -5.0 // both negative, will produce negative values
	strategy := NewUniformMutationStrategy(l, r)

	parent := &MockAntViewUniform{
		alphaVal:          5.0,
		betaVal:           3.0,
		pheromoneVal:      1.5,
		pathStrategy:      path.NewPahtClassicStrategy(),
		pheromoneStrategy: apply.NewApplyClassicStrategy(),
	}

	// Run many samples; after clamping should be zero
	samples := 1000
	for i := 0; i < samples; i++ {
		child := strategy.Mutate(parent)
		assert.NotNil(t, child)
		// Since alpha + rand*(r-l) + l where l negative and r negative, the max is alpha + r (still negative if alpha small)
		// With alpha=2, r=-5 => max = -3 <0, so after clamping should be 0.
		assert.Equal(t, 0.0, child.Alpha(), "alpha should be clamped to zero")
		assert.Equal(t, 0.0, child.Beta(), "beta should be clamped to zero")
		// Other parameters inherited
		assert.InDelta(t, parent.PheromoneMultiplier(), child.PheromoneMultiplier(), 1e-9, "pheromone multiplier")
		assert.Equal(t, parent.PathStrategy(), child.PathStrategy())
		assert.Equal(t, parent.PheromoneApplyStrategy(), child.PheromoneApplyStrategy())
	}
}
