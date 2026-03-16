package adapters

import (
	"io"
	"math"

	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/tsplib"
)

type Euclidean2DAdapter struct{}

func init() {
	GetRegistry().RegisterAdapter(&Euclidean2DAdapter{})
}

func (a *Euclidean2DAdapter) Name() string {
	return "EUC_2D"
}

func (a *Euclidean2DAdapter) CanHandle(weightType string, weightFormat string) bool {
	return weightType == tsplib.WeightTypeEUC2D && weightFormat == tsplib.WeightFormatFUNCTION
}

func (a *Euclidean2DAdapter) Parse(r io.Reader, problem *tsplib.Problem) ([]*graph.Vertex, []*graph.Edge, error) {
	nodes, err := ParseNodes(r, problem.Dimension, 3)
	if err != nil {
		return nil, nil, err
	}

	vertices, vertexMap := CreateVertices(nodes)

	edges := BuildCompleteGraph(nodes, vertexMap, func(n1, n2 Node) float64 {
		x1 := MustGetCoord(n1, 0)
		y1 := MustGetCoord(n1, 1)
		x2 := MustGetCoord(n2, 0)
		y2 := MustGetCoord(n2, 1)

		return math.Round(EuclideanDistance(x1, y1, x2, y2))
	})

	return vertices, edges, nil
}
