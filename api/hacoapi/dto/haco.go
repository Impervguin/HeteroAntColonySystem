package dto

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/core/colony"
	"HeteroAntColonySystem/internal/strategies/crossover"
	"HeteroAntColonySystem/internal/strategies/mutation"
	"HeteroAntColonySystem/internal/strategies/optimisation"
	"HeteroAntColonySystem/internal/strategies/selection"
	"HeteroAntColonySystem/pkg/graph"
	"errors"

	"github.com/gin-gonic/gin"
)

type HacoRunRequest struct {
	Graph graphJson `json:"graph" binding:"required"`

	DefaultAlpha        float64 `json:"default_alpha" binding:"required"`
	DefaultBeta         float64 `json:"default_beta" binding:"required"`
	PheromoneMultiplier float64 `json:"pheromone_multiplier" binding:"required"`
	EvaporationRate     float64 `json:"evaporation_rate" binding:"required"`
	InitialPheromone    float64 `json:"initial_pheromone" binding:"required"`

	GenerationCount  uint `json:"generation_count" binding:"required"`
	ColonySize       uint `json:"colony_size" binding:"required"`
	GenerationPeriod uint `json:"generation_period" binding:"required"`
	ParentCount      uint `json:"parent_count" binding:"required"`

	Selection         HacoSelectionStrategy         `json:"selection" binding:"required"`
	Crossover         HacoCrossoverStrategy         `json:"crossover" binding:"required"`
	Mutation          HacoMutationStrategy          `json:"mutation" binding:"required"`
	LocalOptimisation HacoLocalOptimisationStrategy `json:"local_optimisation" binding:"required"`
}

func DeserializeHacoRunRequest(c *gin.Context) (*HacoRunRequest, error) {
	var req HacoRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	if err := req.Selection.Verify(); err != nil {
		return nil, err
	}
	if err := req.Crossover.Verify(); err != nil {
		return nil, err
	}
	if err := req.Mutation.Verify(); err != nil {
		return nil, err
	}
	if err := req.LocalOptimisation.Verify(); err != nil {
		return nil, err
	}

	return &req, nil
}

type HacoSelectionStrategy struct {
	T string `json:"type" binding:"required,oneof=best tournament"`

	// Tournament
	K *uint `json:"k"`
}

func (s *HacoSelectionStrategy) Verify() error {
	switch s.T {
	case "best":
		return nil
	case "tournament":
		if s.K == nil {
			return errors.New("tournament selection strategy requires k")
		}
	default:
		return errors.New("unknown selection strategy")
	}
	return nil
}

func (s *HacoSelectionStrategy) Get() colony.ParentSelectionStrategy {
	switch s.T {
	case "best":
		return selection.NewBestSelectionStrategy()
	case "tournament":
		return selection.NewTournamentSelectionStrategy(*s.K)
	default:
		panic("unknown selection strategy")
	}
}

type HacoCrossoverStrategy struct {
	T string `json:"type" binding:"required,oneof=arithmetic"`
}

func (s *HacoCrossoverStrategy) Verify() error {
	switch s.T {
	case "arithmetic":
		return nil
	default:
		return errors.New("unknown crossover strategy")
	}
}

// Get returns the crossover strategy
func (s *HacoCrossoverStrategy) Get() colony.CrossoverStrategy {
	switch s.T {
	case "arithmetic":
		return crossover.NewAriphmeticCrossoverStrategy()
	default:
		panic("unknown crossover strategy")
	}
}

type HacoMutationStrategy struct {
	T string `json:"type" binding:"required,oneof=uniform gauss"`

	// Uniform
	Min *float64 `json:"min"`
	Max *float64 `json:"max"`

	// Gauss
	Mean *float64 `json:"mean"`
	Std  *float64 `json:"std"`
}

func (s *HacoMutationStrategy) Verify() error {
	switch s.T {
	case "uniform":
		if s.Min == nil || s.Max == nil {
			return errors.New("uniform mutation strategy requires min and max")
		}
	case "gauss":
		if s.Mean == nil || s.Std == nil {
			return errors.New("gauss mutation strategy requires mean and std")
		}
	default:
		return errors.New("unknown mutation strategy")
	}
	return nil
}

func (s *HacoMutationStrategy) Get() colony.MutationStrategy {
	switch s.T {
	case "uniform":
		return mutation.NewUniformMutationStrategy(*s.Min, *s.Max)
	case "gauss":
		return mutation.NewGaussMutationStrategy(*s.Mean, *s.Std)
	default:
		panic("unknown mutation strategy")
	}
}

type HacoLocalOptimisationStrategy struct {
	T string `json:"type" binding:"required,oneof=noop 2opt"`
}

func (s *HacoLocalOptimisationStrategy) Verify() error {
	switch s.T {
	case "noop":
		return nil
	case "2opt":
		return nil
	default:
		return errors.New("unknown local optimisation strategy")
	}
}

func (s *HacoLocalOptimisationStrategy) Get() ant.LocalOptimisationStrategy {
	switch s.T {
	case "noop":
		return optimisation.NewNoOpLocalOptimisation()
	case "2opt":
		return optimisation.NewTwoOptLocalOptimisation()
	default:
		panic("unknown local optimisation strategy")
	}
}

type HacoRunResponse struct {
	BestScore float64  `json:"best_score"`
	BestPath  []string `json:"best_path"`
}

func SerializeHacoRunResponse(_ *gin.Context, bestPath []*graph.Vertex, bestScore float64) any {
	path := make([]string, 0, len(bestPath))
	for _, v := range bestPath {
		path = append(path, v.ID().String())
	}
	return &HacoRunResponse{
		BestScore: bestScore,
		BestPath:  path,
	}
}
