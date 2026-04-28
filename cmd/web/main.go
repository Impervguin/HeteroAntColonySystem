package main

import (
	"HeteroAntColonySystem/web/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	router.Setup(r)

	r.Run(":3000")
}
