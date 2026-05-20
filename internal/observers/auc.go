package observers

import (
	"HeteroAntColonySystem/internal/core/colony"
	"math"
)

type AreaUnderCurveObserver struct {
	lastScore float64
	auc       float64
}

func NewAreaUnderCurveObserver() *AreaUnderCurveObserver {
	return &AreaUnderCurveObserver{
		lastScore: -1,
		auc:       0,
	}
}

var _ colony.ColonyObserver = (*AreaUnderCurveObserver)(nil)

func (o *AreaUnderCurveObserver) Observe(dto *colony.ColonyObserverDTO) {
	curBestScore := math.Inf(1)
	for _, ant := range dto.Ants {
		if ant.Score() < curBestScore {
			curBestScore = ant.Score()
		}
	}

	if o.lastScore != -1 {
		// Trapezoidal rule
		o.auc += (curBestScore + o.lastScore) / 2
	}
	o.lastScore = curBestScore
}

func (o *AreaUnderCurveObserver) AreaUnderCurve() float64 {
	return o.auc
}
