package colony

import "HeteroAntColonySystem/internal/core/ant"

type ParentSelectionStrategy interface {
	SelectParents(ants []ant.AntView, n uint) []ant.AntView
}

type CrossoverStrategy interface {
	Crossover(p1, p2 ant.AntView) *ant.HeteroAnt
}

type MutationStrategy interface {
	Mutate(ant ant.AntView) *ant.HeteroAnt
}
