package main

import (
	"HeteroAntColonySystem/web/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	router.Setup(r, "localhost:8080", "/api/v1")

	r.Run(":3000")
}
