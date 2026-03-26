package apply

import "HeteroAntColonySystem/internal/core"

type ApplyClassicStrategy struct{}

var _ core.PheromoneApplyingStrategy = &ApplyClassicStrategy{}

func NewApplyClassicStrategy() *ApplyClassicStrategy {
	return &ApplyClassicStrategy{}
}

func (*ApplyClassicStrategy) ApplyPheromone(ant core.AntView) {
	g := ant.Graph()
	pm := ant.PheromoneMap()
	path := ant.Path()

	delta := ant.PheromoneMultiplier() / ant.Score()

	for i := 0; i < len(path)-1; i++ {
		s, t := path[i], path[i+1]
		e, _ := g.Edge(s, t)
		pm.Add(e, delta)
	}
	wrapE, _ := g.Edge(path[len(path)-1], path[0])
	pm.Add(wrapE, delta)
}
