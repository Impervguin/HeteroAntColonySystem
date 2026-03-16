package mutation

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

func TestUniformMutationInRange(t *testing.T) {
	ant := core.NewHeteroAnt(1, 1, &mockChoose{})

	mutation := NewUniformMutation(-0.5, 1, &mockChoose{})
	for i := 0; i < 1000; i++ {
		newAnt := mutation.Mutate(ant)
		if newAnt.Alpha()-ant.Alpha() < -0.5 || newAnt.Alpha()-ant.Alpha() > 1 {
			t.Fatalf("alpha out of range")
		}
		if newAnt.Beta()-ant.Beta() < -0.5 || newAnt.Beta()-ant.Beta() > 1 {
			t.Fatalf("beta out of range")
		}
	}
}

func TestUniformMutationRandom(t *testing.T) {
	ant := core.NewHeteroAnt(1, 1, &mockChoose{})

	mutation := NewUniformMutation(-0.5, 1, &mockChoose{})
	alphas := make(map[float64]struct{})
	betas := make(map[float64]struct{})
	for i := 0; i < 1000; i++ {
		newAnt := mutation.Mutate(ant)
		alphas[newAnt.Alpha()] = struct{}{}
		betas[newAnt.Beta()] = struct{}{}
	}
	require.Len(t, alphas, 1000)
	require.Len(t, betas, 1000)
}

func TestUniformMutationZeroRange(t *testing.T) {
	ant := core.NewHeteroAnt(1, 1, &mockChoose{})

	mutation := NewUniformMutation(0, 0, &mockChoose{})
	alphas := make(map[float64]struct{})
	betas := make(map[float64]struct{})
	for i := 0; i < 1000; i++ {
		newAnt := mutation.Mutate(ant)
		alphas[newAnt.Alpha()] = struct{}{}
		betas[newAnt.Beta()] = struct{}{}
	}
	require.Len(t, alphas, 1)
	require.Len(t, betas, 1)
}
