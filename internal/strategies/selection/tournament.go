package selection

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/core/colony"
	"math/rand/v2"
)

type TournamentSelectionStrategy struct {
	k uint
}

var _ colony.ParentSelectionStrategy = (*TournamentSelectionStrategy)(nil)

func NewTournamentSelectionStrategy(k uint) *TournamentSelectionStrategy {
	return &TournamentSelectionStrategy{
		k: k,
	}
}

func (s *TournamentSelectionStrategy) SelectParents(ants []ant.AntView, n uint) []ant.AntView {
	res := make([]ant.AntView, 0, n)

	samples := make([]ant.AntView, s.k)
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
