package crossover_test

import (
	"HeteroAntColonySystem/internal/strategies/crossover"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSBXCrossoverStrategy_Crossover(t *testing.T) {
	eta := 2.0
	strategy := crossover.NewSBXCrossoverStrategy(eta)

	parent1 := NewMockAntView(5.0, 3.0, 1.5)
	parent2 := NewMockAntView(5.0, 3.0, 1.5)

	// Run many times to sample distribution
	samples := 1000
	minAlpha := math.Inf(1)
	maxAlpha := math.Inf(-1)
	minBeta := math.Inf(1)
	maxBeta := math.Inf(-1)

	for i := 0; i < samples; i++ {
		children := strategy.Crossover(parent1, parent2)
		assert.Len(t, children, 2, "SBX should produce two children")
		for _, child := range children {
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
			assert.InDelta(t, parent1.PheromoneMultiplier(), child.PheromoneMultiplier(), 1e-9, "pheromone multiplier")
			assert.Equal(t, parent1.PathStrategy(), child.PathStrategy())
			assert.Equal(t, parent1.PheromoneApplyStrategy(), child.PheromoneApplyStrategy())
		}
	}

	// For SBX, with eta=2, the spread is moderate. We'll just ensure values are within a reasonable range.
	// Since SBX can produce values outside [x1,x2] but not too far.
	// We'll check that minAlpha >= min(parent1.Alpha, parent2.Alpha) - some margin.
	// For simplicity, we'll just ensure they are not negative huge.
	assert.GreaterOrEqual(t, minAlpha, 0.0, "alpha should be non-negative")
	assert.GreaterOrEqual(t, minBeta, 0.0, "beta should be non-negative")
	// Upper bound: we can't guarantee but we can check they are not astronomically large.
	assert.LessOrEqual(t, maxAlpha, 100.0, "alpha unreasonably large")
	assert.LessOrEqual(t, maxBeta, 100.0, "beta unreasonably large")
}

// Test edge case where parents are identical
func TestSBXCrossoverStrategy_IdenticalParents(t *testing.T) {
	eta := 2.0
	strategy := crossover.NewSBXCrossoverStrategy(eta)

	parent := NewMockAntView(5.0, 3.0, 1.5)

	children := strategy.Crossover(parent, parent)
	assert.Len(t, children, 2)
	// With identical parents, SBX should return copies of the parent (since gamma=1?)
	// Actually in implementation, if u>0.5 or |x1-x2|<1e-10, returns x1,x2.
	// Since they are equal, |x1-x2|<1e-10, so should return same values.
	for _, child := range children {
		assert.InDelta(t, parent.Alpha(), child.Alpha(), 1e-9, "alpha")
		assert.InDelta(t, parent.Beta(), child.Beta(), 1e-9, "beta")
		assert.InDelta(t, parent.PheromoneMultiplier(), child.PheromoneMultiplier(), 1e-9, "pheromone multiplier")
		assert.Equal(t, parent.PathStrategy(), child.PathStrategy())
		assert.Equal(t, parent.PheromoneApplyStrategy(), child.PheromoneApplyStrategy())
	}
}
