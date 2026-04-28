package dto

import (
	"HeteroAntColonySystem/internal/observers"
	"HeteroAntColonySystem/pkg/graph"

	"github.com/gin-gonic/gin"
)

type HacoRunDetailsResponse struct {
	BestScore         float64      `json:"best_score"`
	BestPath          []string     `json:"best_path"`
	AvgCoeffs         []avgCoeffs  `json:"avg_coeffs"`
	FinalPheromoneMap pheromoneMap `json:"final_pheromone_map"`
	BestPaths         []path       `json:"best_paths"`
}

type avgCoeffs struct {
	Alpha float64 `json:"alpha"`
	Beta  float64 `json:"beta"`
	Run   uint    `json:"run"`
}

type pheromoneMap struct {
	Items []pheromoneItem `json:"items"`
}

type pheromoneItem struct {
	SourceID  string  `json:"source_id"`
	TargetID  string  `json:"target_id"`
	Pheromone float64 `json:"weight"`
}

type path struct {
	Run   uint     `json:"run"`
	Path  []string `json:"path"`
	Score float64  `json:"score"`
}

func SerializeHacoRunDetailsResponse(_ *gin.Context,
	bestPath []*graph.Vertex,
	bestScore float64,
	coeffObserver *observers.AntParamsObserver,
	pmObserver *observers.PheromoneMapObserver,
	pathObserver *observers.BestPathObserver,
	gens uint,
) any {
	bestPaths := make([]path, 0, gens)
	for i := uint(0); i < gens; i++ {
		p, s := pathObserver.Path(i)
		arr := make([]string, 0, len(p))
		for _, v := range p {
			arr = append(arr, v.ID().String())
		}
		bestPaths = append(bestPaths, path{
			Run:   i,
			Path:  arr,
			Score: s,
		})
	}

	coeffs := make([]avgCoeffs, 0, gens)
	for i := uint(0); i < gens; i++ {
		alpha, beta := coeffObserver.Params(i)
		coeffs = append(coeffs, avgCoeffs{
			Alpha: alpha,
			Beta:  beta,
			Run:   i,
		})
	}

	items := make([]pheromoneItem, 0, gens)
	pm := pmObserver.Map(gens - 1)
	pm.ForEachEdgeRead(func(e *graph.Edge, pheromone float64) {
		items = append(items, pheromoneItem{
			SourceID:  e.Source().ID().String(),
			TargetID:  e.Target().ID().String(),
			Pheromone: pheromone,
		})
	})

	p := make([]string, 0, len(bestPath))
	for _, v := range bestPath {
		p = append(p, v.ID().String())
	}

	return &HacoRunDetailsResponse{
		BestScore: bestScore,
		BestPath:  p,
		BestPaths: bestPaths,
		AvgCoeffs: coeffs,
		FinalPheromoneMap: pheromoneMap{
			Items: items,
		},
	}
}
