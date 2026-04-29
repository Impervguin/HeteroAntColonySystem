package dto

import (
	"HeteroAntColonySystem/api/common/dto"
	"HeteroAntColonySystem/pkg/graph"
)

type graphJson = dto.GraphJson

func mapGraph(g *graph.Graph) *graphJson {
	return dto.MapGraph(g)
}
