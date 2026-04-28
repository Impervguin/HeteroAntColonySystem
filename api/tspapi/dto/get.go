package dto

import (
	"HeteroAntColonySystem/pkg/graph"

	"github.com/gin-gonic/gin"
)

type GetTSPRequest struct {
	File string `uri:"file" binding:"required"`
}

func DeserializeGetTSPRequest(c *gin.Context) (*GetTSPRequest, error) {
	var req GetTSPRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

type GetTSPResponse struct {
	Graph *graphJson `json:"graph"`
}

func SerializeGetTSPResponse(_ *gin.Context, g *graph.Graph) *GetTSPResponse {
	return &GetTSPResponse{
		Graph: mapGraph(g),
	}
}
