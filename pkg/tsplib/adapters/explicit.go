package adapters

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/tsplib"
)

type ExplicitAdapter struct{}

func init() {
	GetRegistry().RegisterAdapter(&ExplicitAdapter{})
}

func (a *ExplicitAdapter) Name() string {
	return "EXPLICIT_FULL_MATRIX"
}

func (a *ExplicitAdapter) CanHandle(weightType string, weightFormat string) bool {
	return weightType == tsplib.WeightTypeEXPLICIT && weightFormat == tsplib.WeightFormatFULL_MATRIX
}

type ExplicitMetadata struct{}

func (a *ExplicitAdapter) MetadataType() any {
	return &ExplicitMetadata{}
}

func (a *ExplicitAdapter) Parse(r io.Reader, problem *tsplib.Problem) ([]*graph.Vertex, []*graph.Edge, error) {
	scanner := bufio.NewScanner(r)

	if !scanner.Scan() {
		return nil, nil, fmt.Errorf("%w: empty data section", tsplib.ErrInvalidData)
	}

	header := strings.TrimSpace(scanner.Text())
	if header != tsplib.SectionEdgeWeight {
		return nil, nil, fmt.Errorf("%w: expected EDGE_WEIGHT_SECTION, got %s", tsplib.ErrInvalidFormat, header)
	}

	var weights []float64
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line == tsplib.SectionEOF {
			continue
		}
		fields := strings.Fields(line)
		for _, f := range fields {
			w, err := strconv.ParseFloat(f, 64)
			if err != nil {
				return nil, nil, fmt.Errorf("%w: invalid weight %s", tsplib.ErrInvalidData, f)
			}
			weights = append(weights, w)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("scanner error: %w", err)
	}

	dim := problem.Dimension
	if len(weights) != dim*dim {
		return nil, nil, fmt.Errorf("%w: expected %d weights, got %d", tsplib.ErrInvalidData, dim*dim, len(weights))
	}

	vertices := make([]*graph.Vertex, dim)
	vertexMap := make(map[int]*graph.Vertex, dim)
	names := GenerateVertexNameSequence(dim)
	for i := 0; i < dim; i++ {
		v := graph.NewVertex(names[i])
		v.SetMetadata(&ExplicitMetadata{})
		vertices[i] = v
		vertexMap[i+1] = v // TSPLIB IDs are 1-based
	}

	edges := make([]*graph.Edge, 0, dim*(dim-1))
	idx := 0
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			if i == j {
				idx++
				continue
			}
			w := weights[idx]
			edge := graph.NewEdge(vertexMap[i+1], vertexMap[j+1], w)
			edges = append(edges, edge)
			idx++
		}
	}

	return vertices, edges, nil
}
