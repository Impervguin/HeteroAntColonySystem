package crossover

import (
	"HeteroAntColonySystem/internal/core"
	"HeteroAntColonySystem/pkg/graph"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockChoose struct {
}

var _ core.ChoosePathStrategy = &mockChoose{}

func (m *mockChoose) ChooseNext(state core.AntInWorkView, ant *core.HeteroAnt) (*graph.Vertex, bool) {
	return nil, true
}

func TestAriphmeticCrossover(t *testing.T) {
	ant := core.NewHeteroAnt(1, 1, &mockChoose{})
	other := core.NewHeteroAnt(2, 2, &mockChoose{})

	crossover := NewAriphmeticCrossover(&mockChoose{})
	newAnt := crossover.Crossover(ant, other)

	require.Equal(t, newAnt.Alpha(), 1.5)
	require.Equal(t, newAnt.Beta(), 1.5)
}

func TestAriphmeticCrossoverZero(t *testing.T) {
	ant := core.NewHeteroAnt(1, 1, &mockChoose{})
	other := core.NewHeteroAnt(0, 0, &mockChoose{})

	crossover := NewAriphmeticCrossover(&mockChoose{})
	newAnt := crossover.Crossover(ant, other)

	require.Equal(t, newAnt.Alpha(), 0.5)
	require.Equal(t, newAnt.Beta(), 0.5)
}

func TestAriphmeticCrossoverNegative(t *testing.T) {
	ant := core.NewHeteroAnt(1, 1, &mockChoose{})
	other := core.NewHeteroAnt(-1, -1, &mockChoose{})

	crossover := NewAriphmeticCrossover(&mockChoose{})
	newAnt := crossover.Crossover(ant, other)

	require.Equal(t, newAnt.Alpha(), 0.)
	require.Equal(t, newAnt.Beta(), 0.)
}

func TestAriphmeticCrossoverSame(t *testing.T) {
	ant := core.NewHeteroAnt(1, 1, &mockChoose{})
	other := core.NewHeteroAnt(1, 1, &mockChoose{})

	crossover := NewAriphmeticCrossover(&mockChoose{})
	newAnt := crossover.Crossover(ant, other)

	require.Equal(t, newAnt.Alpha(), 1.)
	require.Equal(t, newAnt.Beta(), 1.)
}
