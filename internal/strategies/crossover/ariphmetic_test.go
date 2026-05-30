package crossover_test

import (
	"HeteroAntColonySystem/internal/strategies/crossover"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAriphmeticCrossoverStrategy_Crossover(t *testing.T) {
	strategy := crossover.NewAriphmeticCrossoverStrategy()
	parent1 := NewMockAntView(2.0, 3.0, 1.5)
	parent2 := NewMockAntView(4.0, 6.0, 2.5)

	children := strategy.Crossover(parent1, parent2)
	assert.Len(t, children, 1, "expected one child")

	child := children[0]
	assert.NotNil(t, child)

	// Expected alpha = (2+4)/2 = 3
	assert.InDelta(t, 3.0, child.Alpha(), 1e-9, "alpha")
	// Expected beta = (3+6)/2 = 4.5
	assert.InDelta(t, 4.5, child.Beta(), 1e-9, "beta")
	// Pheromone multiplier should be taken from parent1 (as per implementation)
	assert.InDelta(t, 1.5, child.PheromoneMultiplier(), 1e-9, "pheromone multiplier")
	// Path strategy should be from parent1
	assert.Equal(t, parent1.PathStrategy(), child.PathStrategy())
	// Pheromone apply strategy should be from parent1
	assert.Equal(t, parent1.PheromoneApplyStrategy(), child.PheromoneApplyStrategy())
}
