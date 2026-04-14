package optimisation

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/pkg/graph"
)

type TwoOptLocalOptimisation struct {
}

func NewTwoOptLocalOptimisation() *TwoOptLocalOptimisation {
	return &TwoOptLocalOptimisation{}
}

var _ ant.LocalOptimisationStrategy = &TwoOptLocalOptimisation{}

func (s *TwoOptLocalOptimisation) Optimise(path []*graph.Vertex, g *graph.Graph) {
	if len(path) < 4 {
		return
	}

	n := len(path)
	improved := true
	for improved {
		improved = false
		for i := 0; i < n-2; i++ {
			for j := i + 2; j < n-1; j++ {
				e1, ok1 := g.Edge(path[i], path[i+1])
				e2, ok2 := g.Edge(path[j], path[j+1])
				if !ok1 || !ok2 {
					continue
				}
				if s.shouldSwap(e1, e2, g) {
					s.swapEdges(path, i, j)
					improved = true
					break
				}
			}

			// Enclosure
			if improved {
				break
			}
			
			if i < n-2 {
				e1, ok1 := g.Edge(path[i], path[i+1])
				e2, ok2 := g.Edge(path[n-1], path[0])

				if ok1 && ok2 && s.shouldSwap(e1, e2, g) {
					s.swapEdges(path, i, n-1)
					improved = true
				}
			}
		}
	}
}

func (s *TwoOptLocalOptimisation) shouldSwap(e1, e2 *graph.Edge, g *graph.Graph) bool {
	v1, v2 := e1.Source(), e1.Target()
	u1, u2 := e2.Source(), e2.Target()
	curWeight := e1.Weight() + e2.Weight()

	newEdge1, ok1 := g.Edge(v1, u1)
	newEdge2, ok2 := g.Edge(v2, u2)
	if !ok1 || !ok2 {
		return false
	}
	newWeight := newEdge1.Weight() + newEdge2.Weight()

	return newWeight < curWeight
}

func (s *TwoOptLocalOptimisation) swapEdges(path []*graph.Vertex, i, j int) {
	s.reverse(path, i+1, j)
}

// reverse переворачивает сегмент пути между индексами start и end включительно
func (s *TwoOptLocalOptimisation) reverse(path []*graph.Vertex, start, end int) {
	for start < end {
		path[start], path[end] = path[end], path[start]
		start++
		end--
	}
}
