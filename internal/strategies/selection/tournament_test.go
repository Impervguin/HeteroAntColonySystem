package selection_test

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/strategies/selection"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTournamentSelectionStrategy_SelectParents(t *testing.T) {
	// Use a fixed seed for deterministic randomness
	strategy := selection.NewTournamentSelectionStrategy(3) // k=3

	// Create ants with known scores
	ants := []ant.AntView{
		&MockAntView{sumScoreVal: 10.0, name: "A"},
		&MockAntView{sumScoreVal: 5.0, name: "B"}, // best (lowest)
		&MockAntView{sumScoreVal: 20.0, name: "C"},
		&MockAntView{sumScoreVal: 15.0, name: "D"},
	}

	// Run many iterations to estimate selection probability
	counts := map[string]int{}
	iterations := 30000
	n := uint(5) // select 5 parents each iteration

	for i := 0; i < iterations; i++ {
		selected := strategy.SelectParents(ants, n)
		for _, ant := range selected {
			// ant is MockAntView
			mockAnt := ant.(*MockAntView)
			counts[mockAnt.name]++
		}
	}

	// Ant B (best) should have highest count
	bestCount := counts["B"]
	assert.GreaterOrEqual(t, bestCount, counts["A"], "best ant should be selected more often than A")
	assert.GreaterOrEqual(t, bestCount, counts["C"], "best ant should be selected more often than C")
	assert.GreaterOrEqual(t, bestCount, counts["D"], "best ant should be selected more often than D")

	// Also check that total selections equal iterations * n
	total := 0
	for _, v := range counts {
		total += v
	}
	assert.Equal(t, iterations*int(n), total, "total selections mismatch")
}

// Test k=1 (uniform random)
func TestTournamentSelectionStrategy_K1(t *testing.T) {
	strategy := selection.NewTournamentSelectionStrategy(1)

	ants := []ant.AntView{
		&MockAntView{sumScoreVal: 1.0, name: "X"},
		&MockAntView{sumScoreVal: 2.0, name: "Y"},
		&MockAntView{sumScoreVal: 3.0, name: "Z"},
	}

	counts := map[string]int{}
	iterations := 9000
	n := uint(3)

	for i := 0; i < iterations; i++ {
		selected := strategy.SelectParents(ants, n)
		for _, ant := range selected {
			mockAnt := ant.(*MockAntView)
			counts[mockAnt.name]++
		}
	}

	// With k=1, each selection is uniform random over ants.
	// Expected count per ant = iterations * n / 3
	expected := float64(iterations*int(n)) / 3.0
	tolerance := expected * 0.05 // 5% tolerance
	for _, name := range []string{"X", "Y", "Z"} {
		assert.InDelta(t, float64(counts[name]), expected, tolerance, "selection count for %s", name)
	}
}

// Test empty ants
func TestTournamentSelectionStrategy_Empty(t *testing.T) {
	strategy := selection.NewTournamentSelectionStrategy(3)
	selected := strategy.SelectParents([]ant.AntView{}, 5)
	assert.Empty(t, selected)
}

// Test zero n
func TestTournamentSelectionStrategy_ZeroN(t *testing.T) {
	strategy := selection.NewTournamentSelectionStrategy(3)
	ants := []ant.AntView{&MockAntView{sumScoreVal: 1.0}}
	selected := strategy.SelectParents(ants, 0)
	assert.Empty(t, selected)
}
