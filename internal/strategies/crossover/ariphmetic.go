package crossover

import "HeteroAntColonySystem/internal/core"

// AriphmeticCrossover creates a new ant with mixed parameters
// by taking the average of the parameters of the parents
type AriphmeticCrossover struct {
	choose core.ChoosePathStrategy
}

func NewAriphmeticCrossover(choose core.ChoosePathStrategy) *AriphmeticCrossover {
	return &AriphmeticCrossover{
		choose: choose,
	}
}

var _ core.CrossoverStrategy = &AriphmeticCrossover{}

func (c *AriphmeticCrossover) Crossover(a, b *core.HeteroAnt) *core.HeteroAnt {
	return core.NewHeteroAnt((a.Alpha()+b.Alpha())/2, (a.Beta()+b.Beta())/2, c.choose)
}
