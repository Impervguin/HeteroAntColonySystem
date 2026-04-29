package dto

import "github.com/gin-gonic/gin"

type RuntimeStatsRequest struct {
	Score float64 `json:"score"`
	// Milliseconds
	TotalTime float64 `json:"total_time"`
	AvgTime   float64 `json:"avg_time"`
	MaxTime   float64 `json:"max_time"`
	MinTime   float64 `json:"min_time"`

	MaxMemory uint64 `json:"max_memory"`
	MinMemory uint64 `json:"min_memory"`
	AvgMemory uint64 `json:"avg_memory"`

	SeenPaths uint `json:"seen_paths"`
}

func DeserializeRuntimeStatsRequest(c *gin.Context) (*RuntimeStatsRequest, error) {
	var req RuntimeStatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	return &req, nil
}
