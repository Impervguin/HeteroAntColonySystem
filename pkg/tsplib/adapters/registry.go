package adapters

import (
	"sync"

	"HeteroAntColonySystem/pkg/tsplib"
)

var (
	once     sync.Once
	instance *AdapterRegistry
)

type AdapterRegistry struct {
	adapters map[string]tsplib.TSPLIBAdapter
}

func GetRegistry() *AdapterRegistry {
	once.Do(func() {
		instance = &AdapterRegistry{
			adapters: make(map[string]tsplib.TSPLIBAdapter),
		}
	})
	return instance
}

func (r *AdapterRegistry) RegisterAdapter(adapter tsplib.TSPLIBAdapter) {
	r.adapters[adapter.Name()] = adapter
}

func (r *AdapterRegistry) Get(weightType string, weightFormat string) tsplib.TSPLIBAdapter {
	if weightType == "" {
		weightType = tsplib.WeightTypeEUC2D
	}
	if weightFormat == "" {
		weightFormat = tsplib.WeightFormatFUNCTION
	}

	for _, adapter := range r.adapters {
		if adapter.CanHandle(weightType, weightFormat) {
			return adapter
		}
	}
	return nil
}

func (r *AdapterRegistry) ListAdapters() []string {
	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	return names
}
