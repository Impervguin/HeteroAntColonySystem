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

type LocalOptimisationStrategy interface {
	// Optimise optimises the ants path
	// Should change the contents of the ants path by swapping vertices
	// But not it's size or vertieces itself
	Optimise([]*graph.Vertex, *graph.Graph)
}