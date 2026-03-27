package selection

import (
	"HeteroAntColonySystem/internal/core/strategy"
	"slices"
)

type BestSelectionStrategy struct{}

func NewBestSelectionStrategy() *BestSelectionStrategy {
	return &BestSelectionStrategy{}
}

func (s *BestSelectionStrategy) SelectParents(ants []strategy.AntView, n uint) []strategy.AntView {
	tmp := make([]strategy.AntView, 0, n)
	for _, ant := range ants {
		tmp = append(tmp, ant)
	}

	slices.SortFunc(tmp, func(a, b strategy.AntView) int {
		return int(a.Score() - b.Score())
	})

	return tmp[:n]
}
