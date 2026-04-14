package optimisation_test

import (
	"HeteroAntColonySystem/internal/strategies/optimisation"
	"HeteroAntColonySystem/pkg/graph"
	"fmt"
	"math"
	"math/rand/v2"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TwoOptLocalOptimisationSuite struct {
	suite.Suite

	opt *optimisation.TwoOptLocalOptimisation
}

func (s *TwoOptLocalOptimisationSuite) SetupTest() {
	s.opt = optimisation.NewTwoOptLocalOptimisation()
}

func (s *TwoOptLocalOptimisationSuite) TestTwoOptLocalOptimisation_Square() {
	g := graph.NewGraph(4)

	// Add vertices
	v1 := graph.NewVertex("A(0,0)")
	v2 := graph.NewVertex("B(10,0)")
	v3 := graph.NewVertex("C(10,10)")
	v4 := graph.NewVertex("D(0,10)")
	g.AddVertex(v1)
	g.AddVertex(v2)
	g.AddVertex(v3)
	g.AddVertex(v4)

	// Add edges
	err := euclideanDistancesEdges(g)
	s.NoError(err)

	route := []*graph.Vertex{
		v1,
		v3,
		v2,
		v4,
	}

	expected := []*graph.Vertex{
		v1,
		v2,
		v3,
		v4,
	}

	s.opt.Optimise(route, g)
	s.Equal(expected, route)
}

// Сhecks intersection of two edges with 2 vertices between them
func (s *TwoOptLocalOptimisationSuite) TestTwoOptLocalOptimisation_Cities() {
	g := graph.NewGraph(6)

	cities := make([]*graph.Vertex, 6)
	cities[0] = graph.NewVertex("A(0,0)")
	cities[1] = graph.NewVertex("B(5,5)")
	cities[2] = graph.NewVertex("C(10,5)")
	cities[3] = graph.NewVertex("D(15,0)")
	cities[4] = graph.NewVertex("E(10,-5)")
	cities[5] = graph.NewVertex("F(5,-5)")

	for _, v := range cities {
		g.AddVertex(v)
	}

	err := euclideanDistancesEdges(g)
	s.NoError(err)

	route := []*graph.Vertex{
		cities[0],
		cities[4],
		cities[3],
		cities[2],
		cities[1],
		cities[5],
	}

	expected := []*graph.Vertex{
		cities[0],
		cities[1],
		cities[2],
		cities[3],
		cities[4],
		cities[5],
	}

	s.opt.Optimise(route, g)
	s.Equal(expected, route)
}

// Сhecks intersection of two edges when on of edges in [n-1, 0]
func (s *TwoOptLocalOptimisationSuite) TestTwoOptLocalOptimisation_Enclosure() {
	g := graph.NewGraph(6)

	cities := make([]*graph.Vertex, 6)
	cities[0] = graph.NewVertex("A(0,0)")
	cities[1] = graph.NewVertex("B(5,5)")
	cities[2] = graph.NewVertex("C(10,5)")
	cities[3] = graph.NewVertex("D(15,0)")
	cities[4] = graph.NewVertex("E(10,-5)")
	cities[5] = graph.NewVertex("F(5,-5)")

	for _, v := range cities {
		g.AddVertex(v)
	}

	err := euclideanDistancesEdges(g)
	s.NoError(err)

	route := []*graph.Vertex{
		cities[0],
		cities[1],
		cities[2],
		cities[5],
		cities[4],
		cities[3],
	}

	expected := []*graph.Vertex{
		cities[0],
		cities[1],
		cities[2],
		cities[3],
		cities[4],
		cities[5],
	}

	s.opt.Optimise(route, g)
	s.Equal(expected, route)
}

// Test that 2-opt never worsens the route over 100 random iterations
func (s *TwoOptLocalOptimisationSuite) TestTwoOptNeverWorsens_RandomGraphs() {
	sizes := []int{5, 8, 10, 12}

	for _, n := range sizes {
		failCount := 0
		improvements := 0

		for iter := 0; iter < 100; iter++ {
			g := s.createRandomCompleteGraph(n)
			path := s.createRandomHamiltonianCycle(g)

			originalLength := s.calculatePathLength(path, g)
			s.opt.Optimise(path, g)
			newLength := s.calculatePathLength(path, g)

			if newLength > originalLength+1e-9 {
				failCount++
				s.T().Errorf("n=%d iter=%d: worsened %.4f -> %.4f",
					n, iter, originalLength, newLength)
			} else if newLength < originalLength-1e-9 {
				improvements++
			}
		}

		s.T().Logf("n=%d: %d improvements, %d worsenings out of 100", n, improvements, failCount)
		s.Equal(0, failCount, "2-opt should never worsen the route")
	}
}

// Helper functions

func (s *TwoOptLocalOptimisationSuite) createRandomCompleteGraph(n int) *graph.Graph {
	g := graph.NewGraph(uint(n))
	vertices := make([]*graph.Vertex, n)

	for i := 0; i < n; i++ {
		x, y := rand.IntN(100), rand.IntN(100)
		vertices[i] = graph.NewVertex(fmt.Sprintf("V%d(%d,%d)", i, x, y))
		g.AddVertex(vertices[i])
	}

	err := euclideanDistancesEdges(g)
	s.NoError(err)

	return g
}

func (s *TwoOptLocalOptimisationSuite) createRandomHamiltonianCycle(g *graph.Graph) []*graph.Vertex {
	vertices := make([]*graph.Vertex, 0, g.Len())
	g.ForEachVertex(func(v *graph.Vertex) bool {
		vertices = append(vertices, v)
		return false
	})
	rand.Shuffle(len(vertices), func(i, j int) { vertices[i], vertices[j] = vertices[j], vertices[i] })
	return vertices
}

func (s *TwoOptLocalOptimisationSuite) calculatePathLength(path []*graph.Vertex, g *graph.Graph) float64 {
	if len(path) == 0 {
		return 0
	}
	total := 0.0
	for i := 0; i < len(path)-1; i++ {
		if edge, ok := g.Edge(path[i], path[i+1]); ok {
			total += edge.Weight()
		}
	}
	if edge, ok := g.Edge(path[len(path)-1], path[0]); ok {
		total += edge.Weight()
	}
	return total
}

func TestTwoOptLocalOptimisationSuite(t *testing.T) {
	suite.Run(t, new(TwoOptLocalOptimisationSuite))
}

// Names should be in format "{name}({x},{y})"
// Example: "A(0,0)", "B(10,5)", "City1(3,7)"
func euclideanDistancesEdges(g *graph.Graph) error {
	type vertex struct {
		name string
		x, y int
		ptr  *graph.Vertex
	}

	vertices := make([]vertex, 0, g.Len())

	re := regexp.MustCompile(`^(.+?)\((-?\d+),(-?\d+)\)$`)

	g.ForEachVertex(func(v *graph.Vertex) bool {
		name := v.Name()
		matches := re.FindStringSubmatch(name)

		if len(matches) != 4 {
			return false
		}

		vertexName := matches[1]
		x, err1 := strconv.Atoi(matches[2])
		y, err2 := strconv.Atoi(matches[3])

		if err1 != nil || err2 != nil {
			return false
		}

		vertices = append(vertices, vertex{
			name: vertexName,
			x:    x,
			y:    y,
			ptr:  v,
		})

		return false
	})

	if len(vertices) == 0 {
		return fmt.Errorf("no valid vertices found")
	}

	for i := 0; i < len(vertices); i++ {
		for j := i + 1; j < len(vertices); j++ {
			dx := float64(vertices[i].x - vertices[j].x)
			dy := float64(vertices[i].y - vertices[j].y)
			distance := math.Sqrt(dx*dx + dy*dy)

			if err := g.AddEdge(distance, vertices[i].ptr, vertices[j].ptr); err != nil {
				return err
			}
			if err := g.AddEdge(distance, vertices[j].ptr, vertices[i].ptr); err != nil {
				return err
			}
		}
	}

	return nil
}
