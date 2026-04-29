package dto

import "github.com/gin-gonic/gin"

type ListTSPResponse struct {
	Files []string `json:"files"`
}

func SerializeListTSPResponse(_ *gin.Context, files []string) *ListTSPResponse {
	return &ListTSPResponse{
		Files: files,
	}
}
