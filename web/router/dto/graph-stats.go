package dto

import (
	tspdto "HeteroAntColonySystem/api/tspapi/dto"

	"github.com/gin-gonic/gin"
)

type GraphStatsRequest = tspdto.GraphStatsResponse

func DeserializeGraphStatsRequest(c *gin.Context) (*GraphStatsRequest, error) {
	var req GraphStatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	return &req, nil
}
