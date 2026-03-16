package selection

import (
	"HeteroAntColonySystem/internal/core"
	"container/heap"
)

// ClassicSelection takes count best ants from the candidates
// as fitness function it uses inverse of path length
type ClassicSelection struct {
}

func NewClassicSelection() *ClassicSelection {
	return &ClassicSelection{}
}

var _ core.SelectionStrategy = &ClassicSelection{}

func fitness(a *core.HeteroAnt) float64 {
	return 1.0 / a.Score()
}

type antHeap []*core.HeteroAnt

func (h antHeap) Len() int {
	return len(h)
}

// The more fitness the earlier in the heap
func (h antHeap) Less(i, j int) bool {
	return fitness(h[i]) > fitness(h[j])
}

func (h antHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *antHeap) Push(x interface{}) {
	*h = append(*h, x.(*core.HeteroAnt))
}

func (h *antHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (s *ClassicSelection) Select(candidates []*core.HeteroAnt, count uint) []*core.HeteroAnt {
	h := antHeap{}
	heap.Init(&h)
	for _, ant := range candidates {
		heap.Push(&h, ant)
	}

	res := make([]*core.HeteroAnt, 0, len(candidates))
	for i := uint(0); i < count; i++ {
		res = append(res, heap.Pop(&h).(*core.HeteroAnt))
	}
	return res
}
