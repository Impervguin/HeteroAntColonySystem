package dto

import (
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/tsplib/adapters"
	"errors"

	"github.com/google/uuid"
)

type GraphJson struct {
	Nodes        []GraphNode `json:"nodes" binding:"required"`
	Edges        []GraphEdge `json:"edges" binding:"required"`
	MetadataType string      `json:"metadata_type" binding:"required"`
}

type GraphNode struct {
	ID       string       `json:"id" binding:"required"`
	Name     string       `json:"name" binding:"required"`
	Metadata NodeMetadata `json:"metadata" binding:"required"`
}

type NodeMetadata struct {
	X *float64 `json:"x,omitempty"`
	Y *float64 `json:"y,omitempty"`
	Z *float64 `json:"z,omitempty"`

	Lat *float64 `json:"lat,omitempty"`
	Lon *float64 `json:"lon,omitempty"`
}

type GraphEdge struct {
	Source string  `json:"source" binding:"required"`
	Target string  `json:"target" binding:"required"`
	Weight float64 `json:"weight"`
}

func MapGraph(g *graph.Graph) *GraphJson {
	nodes := make([]GraphNode, 0, g.Len())
	edges := make([]GraphEdge, 0, g.EdgeLen())

	g.ForEachVertex(func(v *graph.Vertex) bool {
		nodes = append(nodes, GraphNode{
			ID:       v.ID().String(),
			Name:     v.Name(),
			Metadata: MapToNodeMetadata(v.Metadata()),
		})
		return false
	})

	g.ForEachEdge(func(e *graph.Edge) bool {
		edges = append(edges, GraphEdge{
			Source: e.Source().ID().String(),
			Target: e.Target().ID().String(),
			Weight: e.Weight(),
		})
		return false
	})

	return &GraphJson{
		Nodes:        nodes,
		Edges:        edges,
		MetadataType: MapMetadataType(g.MetadataType()),
	}
}

func (json *GraphJson) Parse() (*graph.Graph, error) {
	g := graph.NewGraph(uint(len(json.Nodes)))

	metaType, err := ParseMetadataType(json.MetadataType)
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
		metadata, err := ParseMetadata(metaType, node.Metadata)
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

func MapMetadataType(metadataType any) string {
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

func MapToNodeMetadata(metadata any) NodeMetadata {
	if metadata == nil {
		return NodeMetadata{}
	}
	switch m := metadata.(type) {
	case *adapters.Manhattan2DMetadata:
		return NodeMetadata{X: &m.X, Y: &m.Y}
	case *adapters.Euclidean2DMetadata:
		return NodeMetadata{X: &m.X, Y: &m.Y}
	case *adapters.GEOMetadata:
		return NodeMetadata{Lat: &m.Lat, Lon: &m.Lon}
	case *adapters.Euclidean3DMetadata:
		return NodeMetadata{X: &m.X, Y: &m.Y, Z: &m.Z}
	default:
		return NodeMetadata{}
	}
}

func ParseMetadataType(metadataType string) (any, error) {
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

func ParseMetadata(metadataType any, metadata NodeMetadata) (any, error) {
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
