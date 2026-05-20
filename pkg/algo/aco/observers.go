package aco

import "math"

type AreaUnderCurveObserver struct {
	auc       float64
	lastScore float64
}

func NewAreaUnderCurveObserver() *AreaUnderCurveObserver {
	return &AreaUnderCurveObserver{
		lastScore: -1,
		auc:       0,
	}
}

var _ AntColonyObserver = (*AreaUnderCurveObserver)(nil)

func (o *AreaUnderCurveObserver) Observe(c *AntColony) {
	curBestScore := math.Inf(1)
	for _, ant := range c.generation {
		if ant.score < curBestScore {
			curBestScore = ant.score
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

type IterationsToBestObserver struct {
	bestScore float64
	iteration uint
}

func NewIterationsToBestObserver() *IterationsToBestObserver {
	return &IterationsToBestObserver{
		bestScore: math.Inf(1),
		iteration: 0,
	}
}

var _ AntColonyObserver = (*IterationsToBestObserver)(nil)

func (o *IterationsToBestObserver) Observe(c *AntColony) {
	if c.bestTour != nil && c.bestScore < o.bestScore {
		o.bestScore = c.bestScore
		o.iteration = c.currentGeneration
	}
}

func (o *IterationsToBestObserver) IterationsToBest() uint {
	return o.iteration
}
