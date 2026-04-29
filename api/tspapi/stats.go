package tspapi

import (
	"HeteroAntColonySystem/api/tspapi/dto"
	"HeteroAntColonySystem/pkg/graph"
	"math"
)

func factorial(n int) float64 {
	if n == 0 {
		return 1
	}
	factorial := 1.
	for i := 2; i <= n; i++ {
		factorial *= float64(i)
	}
	return factorial
}
func CalculateGraphStats(g *graph.Graph) *dto.GraphStatsResponse {
	avgEdgeWeight := 0.0
	maxEdgeWeight := 0.0
	minEdgeWeight := math.MaxFloat64

	g.ForEachEdge(func(e *graph.Edge) bool {
		avgEdgeWeight += e.Weight()
		if maxEdgeWeight < e.Weight() {
			maxEdgeWeight = e.Weight()
		}
		if minEdgeWeight > e.Weight() {
			minEdgeWeight = e.Weight()
		}
		return false
	})
	avgEdgeWeight /= float64(g.EdgeLen())
	expectedPathLength := avgEdgeWeight * float64(g.Len())
	possibleSolutions := factorial(g.Len()-1) / 2

	recommendedPheromoneMultiplier := avgEdgeWeight * (math.Sqrt(float64(g.Len())) / 2)
	recommendedEvaporationRate := math.Min(0.5, math.Max(0.05, float64(g.Len())/10.))
	return &dto.GraphStatsResponse{
		NodesCount:                     uint(g.Len()),
		EdgesCount:                     uint64(g.EdgeLen()),
		PossibleSolutions:              possibleSolutions,
		AvgEdgeWeight:                  avgEdgeWeight,
		MaxEdgeWeight:                  maxEdgeWeight,
		MinEdgeWeight:                  minEdgeWeight,
		ExpectedPathLength:             expectedPathLength,
		RecommendedPheromoneMultiplier: recommendedPheromoneMultiplier,
		RecommendedEvaporationRate:     recommendedEvaporationRate,
	}
}
