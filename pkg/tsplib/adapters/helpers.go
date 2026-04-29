package adapters

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/tsplib"
)

// NodeParser handles common node parsing logic
type NodeParser struct {
	ExpectedDimension int
	MinFields         int
	ParseField        func(fields []string, index int) (float64, error)
}

// Node represents a parsed node with ID and coordinates
type Node struct {
	ID     int
	Coords []float64
}

// ParseNodes reads node coordinates from a section
func ParseNodes(r io.Reader, expectedDim int, minFields int) ([]Node, error) {
	scanner := bufio.NewScanner(r)

	if !scanner.Scan() {
		return nil, fmt.Errorf("%w: empty data section", tsplib.ErrInvalidData)
	}

	header := strings.TrimSpace(scanner.Text())
	if header != tsplib.SectionNodeCoord && header != tsplib.SectionDisplayData {
		return nil, fmt.Errorf("%w: expected node section, got %s", tsplib.ErrInvalidFormat, header)
	}

	nodes := make([]Node, 0, expectedDim)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line == tsplib.SectionEOF {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < minFields {
			continue
		}

		id, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("%w: invalid node ID %s", tsplib.ErrInvalidData, fields[0])
		}

		coords := make([]float64, 0, len(fields)-1)
		for i := 1; i < len(fields); i++ {
			val, err := strconv.ParseFloat(fields[i], 64)
			if err != nil {
				return nil, fmt.Errorf("%w: invalid coordinate for node %d: %s",
					tsplib.ErrInvalidData, id, fields[i])
			}
			coords = append(coords, val)
		}

		nodes = append(nodes, Node{ID: id, Coords: coords})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	if len(nodes) != expectedDim {
		return nil, fmt.Errorf("%w: expected %d nodes, got %d",
			tsplib.ErrInvalidData, expectedDim, len(nodes))
	}

	return nodes, nil
}

type MetadataFabric func(n Node) any

// CreateVertices creates graph vertices from nodes
func CreateVertices(nodes []Node, mf MetadataFabric) ([]*graph.Vertex, map[int]*graph.Vertex) {
	vertices := make([]*graph.Vertex, 0, len(nodes))
	vertexMap := make(map[int]*graph.Vertex, len(nodes))
	vertexNames := GenerateVertexNameSequence(len(nodes))
	for i, n := range nodes {
		name := vertexNames[i]
		v := graph.NewVertex(name)
		v.SetMetadata(mf(n))
		vertices = append(vertices, v)
		vertexMap[n.ID] = v
	}

	return vertices, vertexMap
}

// EdgeCalculator defines a function that calculates edge weight between two nodes
type EdgeCalculator func(n1, n2 Node) float64

// BuildCompleteGraph creates all edges for a complete graph
func BuildCompleteGraph(
	nodes []Node,
	vertexMap map[int]*graph.Vertex,
	calculator EdgeCalculator,
) []*graph.Edge {

	edges := make([]*graph.Edge, 0, len(nodes)*(len(nodes)-1))

	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes); j++ {
			if i == j {
				continue
			}

			weight := calculator(nodes[i], nodes[j])
			edge := graph.NewEdge(vertexMap[nodes[i].ID], vertexMap[nodes[j].ID], weight)
			edges = append(edges, edge)
		}
	}

	return edges
}

func GenerateVertexNameSequence(count int) []string {
	if count <= 0 {
		return []string{}
	}

	result := make([]string, count)

	for i := 0; i < count; i++ {
		result[i] = vertexName(i + 1)
	}

	return result
}

func vertexName(n int) string {
	result := ""

	for n > 0 {
		n--
		result = string(rune('A'+n%26)) + result
		n /= 26
	}

	return result
}

// MustGetCoord returns coordinate at index, panics if out of range
func MustGetCoord(n Node, idx int) float64 {
	if idx >= len(n.Coords) {
		panic(fmt.Sprintf("node %d missing coordinate at index %d", n.ID, idx))
	}
	return n.Coords[idx]
}

// GetCoord safely returns coordinate at index, with default value
func GetCoord(n Node, idx int, defaultValue float64) float64 {
	if idx >= len(n.Coords) {
		return defaultValue
	}
	return n.Coords[idx]
}

// EuclideanDistance calculates 2D Euclidean distance
func EuclideanDistance(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2
	return math.Sqrt(dx*dx + dy*dy)
}

// EuclideanDistance3D calculates 3D Euclidean distance
func EuclideanDistance3D(x1, y1, z1, x2, y2, z2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2
	dz := z1 - z2
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// ManhattanDistance calculates 2D Manhattan distance
func ManhattanDistance(x1, y1, x2, y2 float64) float64 {
	return math.Abs(x1-x2) + math.Abs(y1-y2)
}

// ManhattanDistance3D calculates 3D Manhattan distance
func ManhattanDistance3D(x1, y1, z1, x2, y2, z2 float64) float64 {
	return math.Abs(x1-x2) + math.Abs(y1-y2) + math.Abs(z1-z2)
}

// MaxDistance calculates 2D Chebyshev (max) distance
func MaxDistance(x1, y1, x2, y2 float64) float64 {
	dx := math.Abs(x1 - x2)
	dy := math.Abs(y1 - y2)
	if dx > dy {
		return dx
	}
	return dy
}

// MaxDistance3D calculates 3D Chebyshev (max) distance
func MaxDistance3D(x1, y1, z1, x2, y2, z2 float64) float64 {
	dx := math.Abs(x1 - x2)
	dy := math.Abs(y1 - y2)
	dz := math.Abs(z1 - z2)

	max := dx
	if dy > max {
		max = dy
	}
	if dz > max {
		max = dz
	}
	return max
}

// CeilDistance calculates CEIL_2D distance
func CeilDistance(x1, y1, x2, y2 float64) float64 {
	return math.Ceil(EuclideanDistance(x1, y1, x2, y2))
}

// RoundDistance implements ATT rounding
func RoundDistance(dx, dy float64) float64 {
	r := math.Sqrt((dx*dx + dy*dy) / 10.0)
	t := math.Round(r)
	if t < r {
		return t + 1.0
	}
	return t
}
