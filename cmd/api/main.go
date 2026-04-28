package main

import (
	"HeteroAntColonySystem/api/hacoapi"
	"HeteroAntColonySystem/api/tspapi"
	"HeteroAntColonySystem/pkg/tsplib"
	"HeteroAntColonySystem/pkg/tsplib/adapters"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Server
	engine := gin.Default()
	engine.Use(gin.Recovery())

	// Infrastructure
	parser := tsplib.NewTSPLIBParser(adapters.GetRegistry())
	if parser == nil {
		panic("Failed to create TSPLIB parser")
	}

	fs := os.DirFS("tsp")

	// API
	apiGroup := engine.Group("/api/v1")

	tspRouter := tspapi.NewRouter(apiGroup, parser, fs)
	_ = tspRouter

	hacoRouter := hacoapi.NewRouter(apiGroup)
	_ = hacoRouter

	// Start server
	engine.Run(":8080")
}
