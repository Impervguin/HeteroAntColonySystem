package crossover_test

import (
	"HeteroAntColonySystem/internal/strategies/crossover"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBLXCrossoverStrategy_Crossover(t *testing.T) {
	gamma := 0.5
	strategy := crossover.NewBLXCrossoverStrategy(gamma)

	parent1 := NewMockAntView(2.0, 3.0, 1.5)
	parent2 := NewMockAntView(4.0, 6.0, 2.5)

	// Run many times to sample distribution
	samples := 1000
	minAlpha := math.Inf(1)
	maxAlpha := math.Inf(-1)
	minBeta := math.Inf(1)
	maxBeta := math.Inf(-1)

	for i := 0; i < samples; i++ {
		children := strategy.Crossover(parent1, parent2)
		assert.Len(t, children, 1)
		child := children[0]
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

	// Compute expected blend range for alpha
	// blend(x1,x2) = lower + rand*(upper-lower)
	// lower = min - gamma*(max-min)
	// upper = max + gamma*(max-min)
	// For alpha: min=2, max=4, interval=2
	// lower = 2 - 0.5*2 = 1
	// upper = 4 + 0.5*2 = 5
	expectedAlphaMin := 1.0
	expectedAlphaMax := 5.0
	// For beta: min=3, max=6, interval=3
	// lower = 3 - 0.5*3 = 1.5
	// upper = 6 + 0.5*3 = 7.5
	expectedBetaMin := 1.5
	expectedBetaMax := 7.5

	assert.InDelta(t, expectedAlphaMin, minAlpha, 0.1, "min alpha out of expected range")
	assert.InDelta(t, expectedAlphaMax, maxAlpha, 0.1, "max alpha out of expected range")
	assert.InDelta(t, expectedBetaMin, minBeta, 0.1, "min beta out of expected range")
	assert.InDelta(t, expectedBetaMax, maxBeta, 0.1, "max beta out of expected range")
}
