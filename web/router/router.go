package router

import (
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"

	"HeteroAntColonySystem/api/utils/ginerr"
	"HeteroAntColonySystem/web/router/dto"
	"HeteroAntColonySystem/web/templates"
	"HeteroAntColonySystem/web/templates/components"
)

type Router struct {
	apiServer string
	apiBase   string
}

func Setup(r *gin.Engine, apiServer, apiBase string) {
	router := &Router{
		apiServer: apiServer,
		apiBase:   apiBase,
	}
	gr := r.Group("/view")
	gr.GET("/", router.Get)

	r.Static("/static", "./web/static")

	r.GET(apiBase+"/*action", func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host, req.Host = apiServer, apiServer
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	r.POST(apiBase+"/*action", func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host, req.Host = apiServer, apiServer
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	rd := r.Group("/render")
	rd.POST("/graph-stats", router.RenderGraphStats)
	rd.POST("/runtime-stats", router.RenderRuntimeStats)
}

var (
	defaultPageData = templates.PageData{
		APIBase: "/api/v1",

		// defaults
		DefaultAlpha:        1,
		DefaultBeta:         1,
		PheromoneMultiplier: 1,
		EvaporationRate:     0.5,
		InitialPheromone:    1,

		GenerationCount:  300,
		ColonySize:       100,
		GenerationPeriod: 5,
		ParentCount:      25,

		// стратегии (простые дефолты)
		SelectionType: "roulette",
		TournamentK:   3,

		CrossoverType:  "blx",
		CrossoverEta:   1.0,
		CrossoverGamma: 0.5,

		MutationType: "gauss",
		MutationMin:  -0.1,
		MutationMax:  0.1,
		MutationMean: 0,
		MutationStd:  0.05,

		LocalOptimisation: "noop",
	}
)

func (r *Router) Get(c *gin.Context) {
	files, err := GetFiles(r.apiServer, r.apiBase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginerr.ErrJSONBody(err))
		return
	}

	pd := defaultPageData
	pd.Files = files.Files
	templates.Page(pd).Render(c.Request.Context(), c.Writer)
}

func (r *Router) RenderGraphStats(c *gin.Context) {
	req, err := dto.DeserializeGraphStatsRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		return
	}

	components.GraphStats(req).Render(c.Request.Context(), c.Writer)
}

func (r *Router) RenderRuntimeStats(c *gin.Context) {
	req, err := dto.DeserializeRuntimeStatsRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		return
	}

	components.RuntimeStats(req).Render(c.Request.Context(), c.Writer)
}
