package observers

import (
	"HeteroAntColonySystem/internal/core/colony"
	"time"
)

type RunTime struct {
	Run    uint
	Moment time.Time
	Time   time.Duration
}

type TimeObserver struct {
	startTime time.Time
	endTime   time.Time
	runs      []RunTime
}

func NewTimeObserver(gens uint) *TimeObserver {
	return &TimeObserver{
		runs: make([]RunTime, 0, gens),
	}
}

var _ colony.ColonyObserver = (*TimeObserver)(nil)

func (o *TimeObserver) Start() {
	o.startTime = time.Now()
}

func (o *TimeObserver) End() {
	o.endTime = time.Now()
}

func (o *TimeObserver) Observe(dto *colony.ColonyObserverDTO) {
	if len(o.runs) == 0 {
		o.runs = append(o.runs, RunTime{
			Run:    dto.Generation,
			Moment: time.Now(),
			Time:   time.Since(o.startTime),
		})
		return
	}
	lastRun := o.runs[len(o.runs)-1]
	o.runs = append(o.runs, RunTime{
		Run:    dto.Generation,
		Moment: time.Now(),
		Time:   time.Since(lastRun.Moment),
	})
}

type TimeData struct {
	Runs      []RunTime
	StartTime time.Time
	EndTime   time.Time
}

func (o *TimeObserver) Data() TimeData {
	return TimeData{
		Runs:      o.runs,
		StartTime: o.startTime,
		EndTime:   o.endTime,
	}
}
