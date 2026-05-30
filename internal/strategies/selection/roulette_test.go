package selection_test

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/strategies/selection"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouletteSelectionStrategy_SelectParents(t *testing.T) {
	strategy := selection.NewRouletteSelectionStrategy()

	// Create ants with different sum scores
	ants := []ant.AntView{
		&MockAntView{sumScoreVal: 5.0, name: "A"}, // best (lowest)
		&MockAntView{sumScoreVal: 10.0, name: "B"},
		&MockAntView{sumScoreVal: 20.0, name: "C"},
	}

	// Compute expected fitness
	// fitness = 1/(sumScore - minScore + 1)
	// minScore = 5
	// fitnessA = 1/(5-5+1) = 1/1 = 1
	// fitnessB = 1/(10-5+1) = 1/6 ≈ 0.1667
	// fitnessC = 1/(20-5+1) = 1/16 = 0.0625
	// total fitness = 1 + 0.1667 + 0.0625 = 1.2292
	// probabilities: A: 1/1.2292 ≈ 0.813, B: 0.1667/1.2292≈0.136, C:0.0625/1.2292≈0.051

	iterations := 50000
	n := uint(3) // select 3 parents each iteration

	counts := map[string]int{}

	for i := 0; i < iterations; i++ {
		selected := strategy.SelectParents(ants, n)
		for _, ant := range selected {
			mockAnt := ant.(*MockAntView)
			counts[mockAnt.name]++
		}
	}

	// Compute expected counts
	totalSelects := iterations * int(n)
	expectedA := float64(totalSelects) * 1.0 / (1.0 + 1.0/6.0 + 1.0/16.0)
	expectedB := float64(totalSelects) * (1.0 / 6.0) / (1.0 + 1.0/6.0 + 1.0/16.0)
	expectedC := float64(totalSelects) * (1.0 / 16.0) / (1.0 + 1.0/6.0 + 1.0/16.0)
	tolerance := float64(totalSelects) * 0.02 // 2% tolerance

	assert.InDelta(t, float64(counts["A"]), expectedA, tolerance, "ant A count")
	assert.InDelta(t, float64(counts["B"]), expectedB, tolerance, "ant B count")
	assert.InDelta(t, float64(counts["C"]), expectedC, tolerance, "ant C count")
}

// Test empty ants
func TestRouletteSelectionStrategy_Empty(t *testing.T) {
	strategy := selection.NewRouletteSelectionStrategy()
	selected := strategy.SelectParents([]ant.AntView{}, 5)
	assert.Empty(t, selected)
}

// Test all same scores (should be uniform)
func TestRouletteSelectionStrategy_EqualScores(t *testing.T) {
	strategy := selection.NewRouletteSelectionStrategy()

	ants := []ant.AntView{
		&MockAntView{sumScoreVal: 10.0, name: "X"},
		&MockAntView{sumScoreVal: 10.0, name: "Y"},
		&MockAntView{sumScoreVal: 10.0, name: "Z"},
	}

	counts := map[string]int{}
	iterations := 30000
	n := uint(3)

	for i := 0; i < iterations; i++ {
		selected := strategy.SelectParents(ants, n)
		for _, ant := range selected {
			mockAnt := ant.(*MockAntView)
			counts[mockAnt.name]++
		}
	}

	expected := float64(iterations*int(n)) / 3.0
	tolerance := expected * 0.03
	for _, name := range []string{"X", "Y", "Z"} {
		assert.InDelta(t, float64(counts[name]), expected, tolerance, "equal score selection for %s", name)
	}
}

// Test zero n
func TestRouletteSelectionStrategy_ZeroN(t *testing.T) {
	strategy := selection.NewRouletteSelectionStrategy()
	ants := []ant.AntView{&MockAntView{sumScoreVal: 1.0}}
	selected := strategy.SelectParents(ants, 0)
	assert.Empty(t, selected)
}
