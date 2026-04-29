package dto

import (
	"HeteroAntColonySystem/internal/observers"
	"HeteroAntColonySystem/pkg/graph"

	"github.com/gin-gonic/gin"
)

type HacoRunDetailsResponse struct {
	BestScore         float64      `json:"best_score"`
	BestPath          []string     `json:"best_path"`
	AvgCoeffs         []avgCoeffs  `json:"avg_coeffs"`
	FinalPheromoneMap pheromoneMap `json:"final_pheromone_map"`
	BestPaths         []path       `json:"best_paths"`
	Memory            memData      `json:"memory"`
	Time              timeData     `json:"time"`
}

type avgCoeffs struct {
	Alpha float64 `json:"alpha"`
	Beta  float64 `json:"beta"`
	Run   uint    `json:"run"`
}

type pheromoneMap struct {
	Items []pheromoneItem `json:"items"`
}

type pheromoneItem struct {
	SourceID  string  `json:"source_id"`
	TargetID  string  `json:"target_id"`
	Pheromone float64 `json:"pheromone"`
}

type path struct {
	Run   uint     `json:"run"`
	Path  []string `json:"path"`
	Score float64  `json:"score"`
}

type memData struct {
	Stats []memStat `json:"stats"`
	Start memStat   `json:"start"`
	End   memStat   `json:"end"`
}

type timeData struct {
	Stats []timeStat `json:"runs"`
	Start timeStat   `json:"start"`
	End   timeStat   `json:"end"`
}

type memStat struct {
	Run  uint   `json:"run"`
	Heap uint64 `json:"heap"`
	Sys  uint64 `json:"sys"`
}

type timeStat struct {
	Run uint `json:"run"`

	// Milliseconds
	Moment *float64 `json:"moment,omitempty"`
	Time   float64  `json:"time"`
}

func SerializeHacoRunDetailsResponse(_ *gin.Context,
	bestPath []*graph.Vertex,
	bestScore float64,
	coeffObserver *observers.AntParamsObserver,
	pmObserver *observers.PheromoneMapObserver,
	pathObserver *observers.BestPathObserver,
	memoryObserver *observers.MemoryObserver,
	timeObserver *observers.TimeObserver,
	gens uint,
) any {
	bestPaths := make([]path, 0, gens)
	for i := uint(0); i < gens; i++ {
		p, s := pathObserver.Path(i)
		arr := make([]string, 0, len(p))
		for _, v := range p {
			arr = append(arr, v.ID().String())
		}
		bestPaths = append(bestPaths, path{
			Run:   i,
			Path:  arr,
			Score: s,
		})
	}

	coeffs := make([]avgCoeffs, 0, gens)
	for i := uint(0); i < gens; i++ {
		alpha, beta := coeffObserver.Params(i)
		coeffs = append(coeffs, avgCoeffs{
			Alpha: alpha,
			Beta:  beta,
			Run:   i,
		})
	}

	items := make([]pheromoneItem, 0, gens)
	pm := pmObserver.Map(gens - 1)
	pm.ForEachEdgeRead(func(e *graph.Edge, pheromone float64) {
		items = append(items, pheromoneItem{
			SourceID:  e.Source().ID().String(),
			TargetID:  e.Target().ID().String(),
			Pheromone: pheromone,
		})
	})

	p := make([]string, 0, len(bestPath))
	for _, v := range bestPath {
		p = append(p, v.ID().String())
	}

	mData := memoryObserver.Data()
	tData := timeObserver.Data()
	memStats := make([]memStat, 0, len(mData.Stats))
	for _, s := range mData.Stats {
		memStats = append(memStats, memStat{
			Run:  s.Run,
			Heap: s.Memory.Heap,
			Sys:  s.Memory.Sys,
		})
	}

	timeStats := make([]timeStat, 0, len(tData.Runs))
	for _, s := range tData.Runs {
		timeStats = append(timeStats, timeStat{
			Run:  s.Run,
			Time: float64(s.Time.Milliseconds()) + (float64(s.Time.Microseconds()) / 1000),
		})
	}

	startT := float64(tData.StartTime.UnixMilli()) + (float64(tData.StartTime.UnixMicro()) / 1000)
	endT := float64(tData.EndTime.UnixMilli()) + (float64(tData.EndTime.UnixMicro()) / 1000)

	return &HacoRunDetailsResponse{
		BestScore: bestScore,
		BestPath:  p,
		BestPaths: bestPaths,
		AvgCoeffs: coeffs,
		FinalPheromoneMap: pheromoneMap{
			Items: items,
		},
		Memory: memData{
			Stats: memStats,
			Start: memStat{
				Heap: mData.Start.Heap,
				Sys:  mData.Start.Sys,
			},
			End: memStat{
				Heap: mData.End.Heap,
				Sys:  mData.End.Sys,
			},
		},
		Time: timeData{
			Stats: timeStats,
			Start: timeStat{
				Moment: &startT,
			},
			End: timeStat{
				Moment: &endT,
			},
		},
	}
}
