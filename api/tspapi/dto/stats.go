package dto

import (
	"HeteroAntColonySystem/pkg/graph"

	"github.com/gin-gonic/gin"
)

type GraphStatsRequest struct {
	Graph *graphJson `json:"graph" binding:"required"`
}

type GraphStatsResponse struct {
	NodesCount                     uint    `json:"nodes_count"`
	EdgesCount                     uint64  `json:"edges_count"`
	PossibleSolutions              float64 `json:"possible_solutions"`
	AvgEdgeWeight                  float64 `json:"avg_edge_weight"`
	MaxEdgeWeight                  float64 `json:"max_edge_weight"`
	MinEdgeWeight                  float64 `json:"min_edge_weight"`
	ExpectedPathLength             float64 `json:"expected_path_length"`
	RecommendedPheromoneMultiplier float64 `json:"recommended_pheromone_multiplier"`
	RecommendedEvaporationRate     float64 `json:"recommended_evaporation_rate"`
}

func DeserializeGetGraphStatsRequest(c *gin.Context) (*graph.Graph, error) {
	var req GraphStatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	g, err := req.Graph.Parse()
	if err != nil {
		return nil, err
	}

	return g, nil

}

func SerializeGetGraphStatsResponse(_ *gin.Context, stats *GraphStatsResponse) any {
	return stats
}
