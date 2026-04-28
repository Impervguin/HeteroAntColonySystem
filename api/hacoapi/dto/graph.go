package dto

import (
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/tsplib/adapters"
	"errors"

	"github.com/google/uuid"
)

type graphJson struct {
	Nodes        []graphNode `json:"nodes" binding:"required"`
	Edges        []graphEdge `json:"edges" binding:"required"`
	MetadataType string      `json:"metadata_type" binding:"required"`
}

type graphNode struct {
	ID       string       `json:"id" binding:"required"`
	Name     string       `json:"name" binding:"required"`
	Metadata nodeMetadata `json:"metadata" binding:"required"`
}

type nodeMetadata struct {
	X *float64 `json:"x"`
	Y *float64 `json:"y"`
	Z *float64 `json:"z"`

	Lat *float64 `json:"lat"`
	Lon *float64 `json:"lon"`
}

type graphEdge struct {
	Source string  `json:"source" binding:"required"`
	Target string  `json:"target" binding:"required"`
	Weight float64 `json:"weight"`
}

func (json *graphJson) Parse() (*graph.Graph, error) {
	g := graph.NewGraph(uint(len(json.Nodes)))

	metaType, err := parseMetadataType(json.MetadataType)
	if err != nil {
		return nil, err
	}

	g.SetMetadataType(metaType)

	vmap := make(map[uuid.UUID]*graph.Vertex, len(json.Nodes))
	for _, node := range json.Nodes {
		id, err := uuid.Parse(node.ID)
		if err != nil {
			return nil, err
		}
		metadata, err := parseMetadata(metaType, node.Metadata)
		if err != nil {
			return nil, err
		}

		v := graph.NewVertex(node.Name, graph.WithID(id), graph.WithMetadata(metadata))
		g.AddVertex(v)
		vmap[id] = v
	}

	for _, edge := range json.Edges {
		sid, err := uuid.Parse(edge.Source)
		if err != nil {
			return nil, err
		}
		tid, err := uuid.Parse(edge.Target)
		if err != nil {
			return nil, err
		}

		if _, oks := vmap[sid]; !oks {
			return nil, errors.New("invalid edge source")
		}

		if _, okt := vmap[tid]; !okt {
			return nil, errors.New("invalid edge target")
		}

		g.AddEdge(edge.Weight, vmap[sid], vmap[tid])
	}

	return g, nil
}

func parseMetadataType(metadataType string) (any, error) {
	switch metadataType {
	case "manhattan_2d":
		return &adapters.Manhattan2DMetadata{}, nil
	case "euclidean_2d":
		return &adapters.Euclidean2DMetadata{}, nil
	case "geo":
		return &adapters.GEOMetadata{}, nil
	case "euclidean_3d":
		return &adapters.Euclidean3DMetadata{}, nil
	default:
		return nil, errors.New("unknown metadata type")
	}
}

func parseMetadata(metadataType any, metadata nodeMetadata) (any, error) {
	switch metadataType.(type) {
	case *adapters.Manhattan2DMetadata:
		if metadata.X == nil || metadata.Y == nil {
			return nil, errors.New("missing x or y coordinate for manhattan 2d metadata")
		}
		return &adapters.Manhattan2DMetadata{
			X: *metadata.X,
			Y: *metadata.Y,
		}, nil
	case *adapters.Euclidean2DMetadata:
		if metadata.X == nil || metadata.Y == nil {
			return nil, errors.New("missing x or y coordinate for euclidean 2d metadata")
		}
		return &adapters.Euclidean2DMetadata{
			X: *metadata.X,
			Y: *metadata.Y,
		}, nil
	case *adapters.GEOMetadata:
		if metadata.Lat == nil || metadata.Lon == nil {
			return nil, errors.New("missing lat or lon coordinate for geo metadata")
		}
		return &adapters.GEOMetadata{
			Lat: *metadata.Lat,
			Lon: *metadata.Lon,
		}, nil
	case *adapters.Euclidean3DMetadata:
		if metadata.X == nil || metadata.Y == nil || metadata.Z == nil {
			return nil, errors.New("missing x, y or z coordinate for euclidean 3d metadata")
		}
		return &adapters.Euclidean3DMetadata{
			X: *metadata.X,
			Y: *metadata.Y,
			Z: *metadata.Z,
		}, nil
	default:
		return nil, errors.New("unknown metadata type")
	}
}
