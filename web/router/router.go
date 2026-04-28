package router

import (
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"

	"HeteroAntColonySystem/web/templates"
)

func Setup(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		templates.Page(
			templates.PageData{
				APIBase: "/api/v1",

				// defaults
				DefaultAlpha:        1,
				DefaultBeta:         1,
				PheromoneMultiplier: 2500,
				EvaporationRate:     0.01,
				InitialPheromone:    2500,

				GenerationCount:  100,
				ColonySize:       10,
				GenerationPeriod: 10,
				ParentCount:      10,

				// стратегии (простые дефолты)
				SelectionType: "best",
				TournamentK:   3,

				CrossoverType: "arithmetic",

				MutationType: "uniform",
				MutationMin:  -0.2,
				MutationMax:  0.2,
				MutationMean: 0,
				MutationStd:  0.3,

				LocalOptimisation: "noop",
			},
		).Render(c.Request.Context(), c.Writer)
	})

	r.Static("/static", "./web/static")

	r.GET("/api/v1/*action", func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host, req.Host = "localhost:8080", "localhost:8080"
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/api/v1/*action", func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host, req.Host = "localhost:8080", "localhost:8080"
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}
