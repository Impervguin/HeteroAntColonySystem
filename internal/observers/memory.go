package observers

import (
	"HeteroAntColonySystem/internal/core/colony"
	"runtime"
)

type RunMemory struct {
	Run uint
	Memory
}

type Memory struct {
	Heap uint64
	Sys  uint64
}

type MemoryObserver struct {
	start Memory
	end   Memory
	stats []RunMemory
}

func NewMemoryObserver(gen uint) *MemoryObserver {
	return &MemoryObserver{
		stats: make([]RunMemory, 0, gen),
	}
}

var _ colony.ColonyObserver = (*MemoryObserver)(nil)

func (o *MemoryObserver) Start() {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	o.start = Memory{
		Heap: m.HeapAlloc,
		Sys:  m.Sys,
	}
}

func (o *MemoryObserver) Observe(dto *colony.ColonyObserverDTO) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	o.stats = append(o.stats, RunMemory{
		Run: dto.Generation,
		Memory: Memory{
			Heap: m.HeapAlloc,
			Sys:  m.Sys,
		},
	})
}

func (o *MemoryObserver) End() {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	o.end = Memory{
		Heap: m.HeapAlloc,
		Sys:  m.Sys,
	}
}

type MemoryData struct {
	Stats []RunMemory
	Start Memory
	End   Memory
}

func (o *MemoryObserver) Data() MemoryData {
	return MemoryData{
		Stats: o.stats,
		Start: o.start,
		End:   o.end,
	}
}
