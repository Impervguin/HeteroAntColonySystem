package selection

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/core/colony"
	"slices"
)

type BestSelectionStrategy struct{}

func NewBestSelectionStrategy() *BestSelectionStrategy {
	return &BestSelectionStrategy{}
}

var _ colony.ParentSelectionStrategy = (*BestSelectionStrategy)(nil)

func (s *BestSelectionStrategy) SelectParents(ants []ant.AntView, n uint) []ant.AntView {
	tmp := make([]ant.AntView, 0, n)
	for _, ant := range ants {
		tmp = append(tmp, ant)
	}

	slices.SortFunc(tmp, func(a, b ant.AntView) int {
		return int(a.Score() - b.Score())
	})

	return tmp[:n]
}
