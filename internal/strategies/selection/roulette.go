package selection

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/core/colony"
	"math/rand/v2"
)

type RouletteSelectionStrategy struct{}

var _ colony.ParentSelectionStrategy = (*RouletteSelectionStrategy)(nil)

func NewRouletteSelectionStrategy() *RouletteSelectionStrategy {
	return &RouletteSelectionStrategy{}
}

func (s *RouletteSelectionStrategy) SelectParents(ants []ant.AntView, n uint) []ant.AntView {
	res := make([]ant.AntView, 0, n)

	if len(ants) == 0 {
		return res
	}

	totalFitness := s.calculateTotalFitness(ants)

	for i := 0; uint(i) < n; i++ {
		selected := s.selectOne(ants, totalFitness)
		res = append(res, selected)
	}

	return res
}

func (s *RouletteSelectionStrategy) calculateTotalFitness(ants []ant.AntView) float64 {
	var total float64
	minScore := s.findMinScore(ants)

	for _, ant := range ants {
		fitness := 1.0 / (ant.SumScore() - minScore + 1.0)
		total += fitness
	}

	return total
}

func (s *RouletteSelectionStrategy) findMinScore(ants []ant.AntView) float64 {
	if len(ants) == 0 {
		return 0
	}

	minScore := ants[0].SumScore()
	for _, ant := range ants[1:] {
		if ant.SumScore() < minScore {
			minScore = ant.SumScore()
		}
	}
	return minScore
}

func (s *RouletteSelectionStrategy) selectOne(ants []ant.AntView, totalFitness float64) ant.AntView {
	rouletteValue := rand.Float64() * totalFitness

	var accumulated float64
	minScore := s.findMinScore(ants)

	for _, ant := range ants {
		fitness := 1.0 / (ant.SumScore() - minScore + 1.0)
		accumulated += fitness

		if accumulated >= rouletteValue {
			return ant
		}
	}

	return ants[len(ants)-1]
}
