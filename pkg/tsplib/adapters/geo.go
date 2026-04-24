// pkg/tsplib/adapters/geo.go
package adapters

import (
	"io"
	"math"

	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/tsplib"
)

type GEOAdapter struct{}

func init() {
	GetRegistry().RegisterAdapter(&GEOAdapter{})
}

func (a *GEOAdapter) Name() string {
	return "GEO"
}

func (a *GEOAdapter) CanHandle(weightType string, weightFormat string) bool {
	return weightType == tsplib.WeightTypeGEO && weightFormat == tsplib.WeightFormatFUNCTION
}

type GEOMetadata struct {
	Lat, Lon float64
}

func (a *GEOAdapter) MetadataType() any {
	return &GEOMetadata{}
}

func (a *GEOAdapter) Parse(r io.Reader, problem *tsplib.Problem) ([]*graph.Vertex, []*graph.Edge, error) {
	nodes, err := ParseNodes(r, problem.Dimension, 3)
	if err != nil {
		return nil, nil, err
	}

	vertices, vertexMap := CreateVertices(nodes, func(n Node) any {
		return &GEOMetadata{
			Lat: MustGetCoord(n, 0),
			Lon: MustGetCoord(n, 1),
		}
	})

	edges := BuildCompleteGraph(nodes, vertexMap, func(n1, n2 Node) float64 {
		lat1 := MustGetCoord(n1, 0)
		lon1 := MustGetCoord(n1, 1)
		lat2 := MustGetCoord(n2, 0)
		lon2 := MustGetCoord(n2, 1)

		return a.geoDistance(lat1, lon1, lat2, lon2)
	})

	return vertices, edges, nil
}

func (a *GEOAdapter) geoDistance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := a.degToRad(lat1)
	lon1Rad := a.degToRad(lon1)
	lat2Rad := a.degToRad(lat2)
	lon2Rad := a.degToRad(lon2)

	q1 := math.Cos(lon1Rad - lon2Rad)
	q2 := math.Cos(lat1Rad - lat2Rad)
	q3 := math.Cos(lat1Rad + lat2Rad)

	const RRR = 6378.388

	dist := RRR*math.Acos(0.5*((1.0+q1)*q2-(1.0-q1)*q3)) + 1.0

	return dist
}

func (a *GEOAdapter) degToRad(deg float64) float64 {
	degInt := math.Floor(deg)
	minutes := deg - degInt

	rad := math.Pi * (degInt + 5.0*minutes/3.0) / 180.0

	return rad
}
