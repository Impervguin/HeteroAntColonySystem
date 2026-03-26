package core

import (
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
)

type HeteroAntColony struct {
	pathChoice     PathChoiceStrategy
	pheromoneApply PheromoneApplyingStrategy

	defaultAlpha        float64
	defaultBeta         float64
	pheromoneMultiplier float64
	evaporationRate     float64
	initialPheromone    float64

	generationCount uint
	colonySize      uint

	*HeteroColonyWork
}

type HeteroColonyWork struct {
	g    *graph.Graph
	pm   *pheromone.PheromoneMap
	ants []*HeteroAnt

	bestPerGeneration []*HeteroAnt
	best              *HeteroAnt
	score             float64
}

func NewHeteroAntColony(opts ...HeteroAntColonyOption) (*HeteroAntColony, error) {
	c := &HeteroAntColony{
		pathChoice:     nil,
		pheromoneApply: nil,

		defaultAlpha:        DefaultDefaultAlpha,
		defaultBeta:         DefaultDefaultBeta,
		pheromoneMultiplier: DefaultPheromoneMultiplier,
		evaporationRate:     DefaultEvaporationRate,
		initialPheromone:    DefaultInitialPheromone,

		generationCount: DefaultGenerationCount,
		colonySize:      DefaultColonySize,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.pathChoice == nil {
		return nil, ErrPathChoiceStrategyNotSet
	}
	if c.pheromoneApply == nil {
		return nil, ErrPheromoneApplyStrategyNotSet
	}
	if c.colonySize <= 0 {
		return nil, ErrColonySizeInvalid
	}
	if c.generationCount <= 0 {
		return nil, ErrGenerationCountInvalid
	}
	if c.pheromoneMultiplier <= 0 {
		return nil, ErrPheromoneMultiplierInvalid
	}
	if c.evaporationRate <= 0 {
		return nil, ErrEvaporationRateInvalid
	}
	if c.initialPheromone <= 0 {
		return nil, ErrInitialPheromoneInvalid
	}

	return c, nil
}

func (c *HeteroAntColony) Prepare(g *graph.Graph) error {
	c.HeteroColonyWork = &HeteroColonyWork{
		g:    g,
		pm:   pheromone.NewPheromoneMap(g, c.initialPheromone),
		ants: make([]*HeteroAnt, 0, c.colonySize),

		bestPerGeneration: make([]*HeteroAnt, 0, c.generationCount),
		best:              nil,
		score:             0,
	}

	return nil
}

func (c *HeteroAntColony) Run() error {
	// first generation
	for i := 0; uint(i) < c.colonySize; i++ {
		c.ants = append(c.ants, NewHeteroAnt(c.defaultAlpha, c.defaultBeta, c.pheromoneMultiplier, c.pathChoice, c.pheromoneApply))
	}

	for i := 0; uint(i) < c.generationCount; i++ {
		// calculate paths
		for _, ant := range c.ants {
			ant.Prepare(c.g, c.pm)
			err := ant.Run()
			if err != nil {
				return err
			}
		}
		// update best
		best := c.ants[0]
		for _, ant := range c.ants {
			if ant.Score() < best.Score() {
				best = ant
			}
		}
		c.bestPerGeneration = append(c.bestPerGeneration, best)
		if c.best == nil || best.Score() < c.best.Score() {
			c.best = best
		}

		// evaporate pheromon
		c.g.ForEachEdge(func(e *graph.Edge) bool {
			c.pm.Update(e, c.pm.Get(e)*(1-c.evaporationRate))
			return false
		})

		// apply pheromone
		for _, ant := range c.ants {
			err := ant.ApplyPheromone()
			if err != nil {
				return err
			}
		}

		// update ants
		nextGen := make([]*HeteroAnt, 0, len(c.ants))
		for _, ant := range c.ants {
			nextGen = append(nextGen, NewHeteroAnt(ant.Alpha(), ant.Beta(), ant.PheromoneMultiplier(), ant.PathStrategy(), ant.PheromoneApplyStrategy()))
		}
		c.ants = nextGen
	}

	c.score = c.best.Score()
	return nil
}

func (c *HeteroAntColony) Score() float64 {
	if c.HeteroColonyWork == nil {
		return 0
	}
	return c.score
}

func (c *HeteroAntColony) BestPath() []*graph.Vertex {
	return c.best.Path()
}
