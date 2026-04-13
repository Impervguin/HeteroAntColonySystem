package selection

import (
	"HeteroAntColonySystem/internal/core/strategy"
	"math/rand/v2"
)

type TournamentSelectionStrategy struct {
	k uint
}

var _ strategy.ParentSelectionStrategy = (*TournamentSelectionStrategy)(nil)

func NewTournamentSelectionStrategy(k uint) *TournamentSelectionStrategy {
	return &TournamentSelectionStrategy{
		k: k,
	}
}

func (s *TournamentSelectionStrategy) SelectParents(ants []strategy.AntView, n uint) []strategy.AntView {
	res := make([]strategy.AntView, 0, n)

	samples := make([]strategy.AntView, s.k)
	for i := 0; uint(i) < n; i++ {
		for j := 0; uint(j) < s.k; j++ {
			r := rand.IntN(len(ants))
			samples[j] = ants[r]
		}

		best := samples[0]
		for _, ant := range samples {
			if ant.Score() > best.Score() {
				best = ant
			}
		}
		res = append(res, best)
	}

	return res
}
