package ant

import "HeteroAntColonySystem/pkg/graph"

type PathChoiceStrategy interface {
	// ChooseNext Chooses next vertex for ants path
	// must return nil if path done
	ChooseNext(ant AntView) *graph.Vertex
}

type PheromoneApplyingStrategy interface {
	ApplyPheromone(ant AntView)
}
