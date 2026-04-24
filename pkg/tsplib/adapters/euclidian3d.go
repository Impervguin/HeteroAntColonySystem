// pkg/tsplib/adapters/euclidean3d.go
package adapters

import (
	"io"
	"math"

	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/tsplib"
)

type Euclidean3DAdapter struct{}

func init() {
	GetRegistry().RegisterAdapter(&Euclidean3DAdapter{})
}

func (a *Euclidean3DAdapter) Name() string {
	return "EUC_3D"
}

func (a *Euclidean3DAdapter) CanHandle(weightType string, weightFormat string) bool {
	return weightType == tsplib.WeightTypeEUC3D && weightFormat == tsplib.WeightFormatFUNCTION
}

type Euclidean3DMetadata struct {
	X, Y, Z float64
}

func (a *Euclidean3DAdapter) MetadataType() any {
	return &Euclidean3DMetadata{}
}

func (a *Euclidean3DAdapter) Parse(r io.Reader, problem *tsplib.Problem) ([]*graph.Vertex, []*graph.Edge, error) {
	nodes, err := ParseNodes(r, problem.Dimension, 4)
	if err != nil {
		return nil, nil, err
	}

	vertices, vertexMap := CreateVertices(nodes, func(n Node) any {
		return &Euclidean3DMetadata{
			X: MustGetCoord(n, 0),
			Y: MustGetCoord(n, 1),
			Z: MustGetCoord(n, 2),
		}
	})

	edges := BuildCompleteGraph(nodes, vertexMap, func(n1, n2 Node) float64 {
		x1 := MustGetCoord(n1, 0)
		y1 := MustGetCoord(n1, 1)
		z1 := MustGetCoord(n1, 2)
		x2 := MustGetCoord(n2, 0)
		y2 := MustGetCoord(n2, 1)
		z2 := MustGetCoord(n2, 2)

		return math.Round(EuclideanDistance3D(x1, y1, z1, x2, y2, z2))
	})

	return vertices, edges, nil
}
