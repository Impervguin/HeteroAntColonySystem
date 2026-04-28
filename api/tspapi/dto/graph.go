package dto

import (
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/tsplib/adapters"
)

type graphJson struct {
	Nodes        []graphNode `json:"nodes"`
	Edges        []graphEdge `json:"edges"`
	MetadataType string      `json:"metadata_type"`
}

type graphNode struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Metadata any    `json:"metadata"`
}

type graphEdge struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Weight float64 `json:"weight"`
}

func mapGraph(g *graph.Graph) *graphJson {
	nodes := make([]graphNode, 0, g.Len())
	edges := make([]graphEdge, 0, g.EdgeLen())

	g.ForEachVertex(func(v *graph.Vertex) bool {
		nodes = append(nodes, graphNode{
			ID:       v.ID().String(),
			Name:     v.Name(),
			Metadata: mapMetadata(v.Metadata()),
		})
		return false
	})

	g.ForEachEdge(func(e *graph.Edge) bool {
		edges = append(edges, graphEdge{
			Source: e.Source().ID().String(),
			Target: e.Target().ID().String(),
			Weight: e.Weight(),
		})
		return false
	})

	return &graphJson{
		Nodes:        nodes,
		Edges:        edges,
		MetadataType: mapMetadataType(g.MetadataType()),
	}
}

func mapMetadata(metadata any) map[string]any {
	if metadata == nil {
		return nil
	}
	switch m := metadata.(type) {
	case *adapters.Manhattan2DMetadata:
		return map[string]any{
			"x": m.X,
			"y": m.Y,
		}
	case *adapters.Euclidean2DMetadata:
		return map[string]any{
			"x": m.X,
			"y": m.Y,
		}
	case *adapters.GEOMetadata:
		return map[string]any{
			"lat": m.Lat,
			"lon": m.Lon,
		}
	case *adapters.Euclidean3DMetadata:
		return map[string]any{
			"x": m.X,
			"y": m.Y,
			"z": m.Z,
		}
	default:
		return nil
	}
}

func mapMetadataType(metadataType any) string {
	switch metadataType.(type) {
	case *adapters.Manhattan2DMetadata:
		return "manhattan_2d"
	case *adapters.Euclidean2DMetadata:
		return "euclidean_2d"
	case *adapters.GEOMetadata:
		return "geo"
	case *adapters.Euclidean3DMetadata:
		return "euclidean_3d"
	default:
		return ""
	}
}
