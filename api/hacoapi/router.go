package hacoapi

import (
	"HeteroAntColonySystem/api/hacoapi/dto"
	"HeteroAntColonySystem/api/utils/ginerr"
	"HeteroAntColonySystem/internal/core/colony"
	"HeteroAntColonySystem/internal/observers"
	"HeteroAntColonySystem/internal/strategies/apply"
	"HeteroAntColonySystem/internal/strategies/path"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
}

func NewRouter(r *gin.RouterGroup) *Router {
	hacp := &Router{}
	gr := r.Group("/haco")

	gr.POST("/run", hacp.Run)
	gr.POST("/run/details", hacp.RunDetails)

	return hacp
}

func (r *Router) Run(c *gin.Context) {
	req, err := dto.DeserializeHacoRunRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		return
	}

	colony, err := r.colonyFromRequest(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		return
	}

	g, err := req.Graph.Parse()
	if err != nil {
		c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		return
	}

	err = colony.Prepare(g)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginerr.ErrJSONBody(err))
		return
	}

	err = colony.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginerr.ErrJSONBody(err))
		return
	}

	c.JSON(http.StatusOK, dto.SerializeHacoRunResponse(c, colony.BestPath(), colony.Score()))
}

func (r *Router) RunDetails(c *gin.Context) {

	req, err := dto.DeserializeHacoRunRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		return
	}

	pathObserver := observers.NewBestPathObserver()
	coeffObserver := observers.NewAntParamsObserver()
	pmObserver := observers.NewPheromoneMapObserver()

	colony, err := r.colonyFromRequest(
		req,
		colony.WithColonyObserver(pathObserver),
		colony.WithColonyObserver(coeffObserver),
		colony.WithColonyObserver(pmObserver),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		return
	}

	g, err := req.Graph.Parse()
	if err != nil {
		c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		return
	}

	err = colony.Prepare(g)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginerr.ErrJSONBody(err))
		return
	}

	err = colony.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginerr.ErrJSONBody(err))
		return
	}

	c.JSON(http.StatusOK, dto.SerializeHacoRunDetailsResponse(c, colony.BestPath(), colony.Score(), coeffObserver, pmObserver, pathObserver, req.GenerationCount))
}

func (r *Router) colonyFromRequest(req *dto.HacoRunRequest, additional ...colony.HeteroAntColonyOption) (*colony.HeteroAntColony, error) {
	selectionStrategy := req.Selection.Get()
	crossoverStrategy := req.Crossover.Get()
	mutationStrategy := req.Mutation.Get()
	localOptimisationStrategy := req.LocalOptimisation.Get()
	applyStrategy := apply.NewApplyClassicStrategy()
	pathStrategy := path.NewPahtClassicStrategy()

	options := []colony.HeteroAntColonyOption{
		colony.WithDefaultAlpha(req.DefaultAlpha),
		colony.WithDefaultBeta(req.DefaultBeta),
		colony.WithPheromoneMultiplier(req.PheromoneMultiplier),
		colony.WithEvaporationRate(req.EvaporationRate),
		colony.WithInitialPheromone(req.InitialPheromone),
		colony.WithGenerationCount(req.GenerationCount),
		colony.WithColonySize(req.ColonySize),
		colony.WithGenerationPeriod(req.GenerationPeriod),
		colony.WithParentCount(req.ParentCount),
		colony.WithPathChoiceStrategy(pathStrategy),
		colony.WithPheromoneApplyingStrategy(applyStrategy),
		colony.WithLocalOptimisationStrategy(localOptimisationStrategy),
		colony.WithParentSelectionStrategy(selectionStrategy),
		colony.WithCrossoverStrategy(crossoverStrategy),
		colony.WithMutationStrategy(mutationStrategy),
	}

	return colony.NewHeteroAntColony(append(options, additional...)...)
}
