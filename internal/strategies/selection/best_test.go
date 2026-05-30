package selection_test

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/strategies/selection"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBestSelectionStrategy_SelectParents(t *testing.T) {
	strategy := selection.NewBestSelectionStrategy()

	// Create ants with different sum scores
	ants := []ant.AntView{
		&MockAntView{sumScoreVal: 10.0},
		&MockAntView{sumScoreVal: 5.0},
		&MockAntView{sumScoreVal: 20.0},
		&MockAntView{sumScoreVal: 15.0},
	}

	// Select top 2 (lowest sum score)
	selected := strategy.SelectParents(ants, 2)
	assert.Len(t, selected, 2)

	// Expect the ants with scores 5.0 and 10.0
	// Since SortFunc sorts ascending by a.SumScore() - b.SumScore()
	// So order should be 5,10,15,20
	assert.InDelta(t, 5.0, selected[0].SumScore(), 1e-9, "first selected")
	assert.InDelta(t, 10.0, selected[1].SumScore(), 1e-9, "second selected")
}

// Test selecting more than available
func TestBestSelectionStrategy_SelectMoreThanAvailable(t *testing.T) {
	strategy := selection.NewBestSelectionStrategy()

	ants := []ant.AntView{
		&MockAntView{sumScoreVal: 1.0},
		&MockAntView{sumScoreVal: 2.0},
	}

	selected := strategy.SelectParents(ants, 5) // request 5 but only 2 available
	assert.Len(t, selected, 2, "should return all available when n > len(ants)")
	assert.InDelta(t, 1.0, selected[0].SumScore(), 1e-9)
	assert.InDelta(t, 2.0, selected[1].SumScore(), 1e-9)
}

// Test empty input
func TestBestSelectionStrategy_EmptyInput(t *testing.T) {
	strategy := selection.NewBestSelectionStrategy()
	selected := strategy.SelectParents([]ant.AntView{}, 3)
	assert.Empty(t, selected)
}
