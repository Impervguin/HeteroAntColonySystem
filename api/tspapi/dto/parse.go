package dto

import (
	"HeteroAntColonySystem/pkg/graph"
	"io"

	"github.com/gin-gonic/gin"
)

type ParseTSPRequest struct {
	File io.Reader
	Name string
	Size int64
}

const (
	ParseTSPRequestFile = "file"
)

func DeserializeParseTSPRequest(c *gin.Context) (*ParseTSPRequest, error) {
	h, err := c.FormFile(ParseTSPRequestFile)
	if err != nil {
		return nil, err
	}
	f, err := h.Open()
	if err != nil {
		return nil, err
	}

	return &ParseTSPRequest{
		File: f,
		Name: h.Filename,
		Size: h.Size,
	}, nil
}

type ParseTSPResponse struct {
	Graph *graphJson `json:"graph"`
}

func SerializeParseTSPResponse(_ *gin.Context, g *graph.Graph) *ParseTSPResponse {
	return &ParseTSPResponse{
		Graph: mapGraph(g),
	}
}
