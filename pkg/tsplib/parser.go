package tsplib

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"HeteroAntColonySystem/pkg/graph"
)

type Problem struct {
	Name             string
	Type             string
	Comment          string
	Dimension        int
	EdgeWeightType   string
	EdgeWeightFormat string
}

func NewProblem() *Problem {
	return &Problem{
		Type:             TypeTSP,
		EdgeWeightType:   WeightTypeEUC2D,
		EdgeWeightFormat: WeightFormatFUNCTION,
	}
}

type TSPLIBParser struct {
	registry AdapterRegistry
}

func NewTSPLIBParser(registry AdapterRegistry) *TSPLIBParser {
	return &TSPLIBParser{
		registry: registry,
	}
}

func (p *TSPLIBParser) Parse(r io.Reader) (*graph.Graph, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	problem, dataSection, err := p.parseMetadata(string(content))
	if err != nil {
		return nil, err
	}

	adapter := p.registry.Get(problem.EdgeWeightType, problem.EdgeWeightFormat)
	if adapter == nil {
		return nil, fmt.Errorf("%w: for type %s format %s",
			ErrAdapterNotFound, problem.EdgeWeightType, problem.EdgeWeightFormat)
	}

	dataReader := strings.NewReader(dataSection)
	vertices, edges, err := adapter.Parse(dataReader, problem)
	if err != nil {
		return nil, fmt.Errorf("adapter %s failed: %w", adapter.Name(), err)
	}

	if len(vertices) != problem.Dimension {
		return nil, fmt.Errorf("%w: expected %d vertices, got %d",
			ErrInvalidData, problem.Dimension, len(vertices))
	}

	g := p.buildGraph(problem, vertices, edges)
	return g, nil
}

func (p *TSPLIBParser) buildGraph(problem *Problem, vertices []*graph.Vertex, edges []*graph.Edge) *graph.Graph {
	g := graph.NewGraph(uint(problem.Dimension))

	for _, v := range vertices {
		g.AddVertex(v)
	}

	for _, e := range edges {
		g.AddEdge(e.Weight(), e.Source(), e.Target())
	}

	adapter := p.registry.Get(problem.EdgeWeightType, problem.EdgeWeightFormat)
	g.SetMetadataType(adapter.MetadataType())

	return g
}

func (p *TSPLIBParser) parseMetadata(content string) (*Problem, string, error) {
	problem := NewProblem()

	scanner := bufio.NewScanner(strings.NewReader(content))
	var dataSection strings.Builder
	inDataSection := false
	dataSectionStarted := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if !dataSectionStarted {
			if line == SectionNodeCoord ||
				line == SectionEdgeWeight ||
				line == SectionDisplayData {
				dataSectionStarted = true
				inDataSection = true
				dataSection.WriteString(line + "\n")
				continue
			}
		}

		if line == SectionEOF {
			break
		}

		if inDataSection {
			dataSection.WriteString(line + "\n")
		} else {
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "COMMENT:") {
				problem.Comment += strings.TrimPrefix(line, "COMMENT:") + "\n"
				continue
			}

			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "NAME":
				problem.Name = value
			case "TYPE":
				if value != "" {
					problem.Type = value
				}
			case "DIMENSION":
				dim, err := strconv.Atoi(value)
				if err != nil {
					return nil, "", fmt.Errorf("%w: invalid dimension", ErrInvalidFormat)
				}
				problem.Dimension = dim
			case "EDGE_WEIGHT_TYPE":
				if value != "" {
					problem.EdgeWeightType = value
				}
			case "EDGE_WEIGHT_FORMAT":
				if value != "" {
					problem.EdgeWeightFormat = value
				}
			}
		}
	}

	if problem.Dimension == 0 {
		return nil, "", fmt.Errorf("%w: missing DIMENSION", ErrInvalidFormat)
	}

	if !dataSectionStarted {
		return nil, "", fmt.Errorf("%w: no data section found", ErrSectionNotFound)
	}

	return problem, dataSection.String(), scanner.Err()
}
