package aco

import (
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
)

type AntColony struct {
	gr     *graph.Graph
	pm     *pheromone.PheromoneMap
	choice ChoosePathStrategy

	colonySize      uint
	generationCount uint

	alpha               float64
	beta                float64
	pheromoneMultiplier float64
	evaporationRate     float64
	initialPheromone    float64

	currentGeneration uint
	generation        []*Ant

	bestPerGeneration []*Ant
	bestTour          []*graph.Vertex
	bestScore         float64
}

func NewAntColony(gr *graph.Graph, opts ...ColonyOption) (*AntColony, error) {
	c := &AntColony{
		gr: gr,

		colonySize:          DefaultColonySize,
		generationCount:     DefaultGenerationSize,
		alpha:               DefaultAlpha,
		beta:                DefaultBeta,
		pheromoneMultiplier: DefaultPheromoneMultiplier,
		evaporationRate:     DefaultEvaporationRate,
		initialPheromone:    DefaultInitialPheromone,
	}
	for _, opt := range opts {
		opt(c)
	}

	if c.colonySize <= 0 {
		return nil, ErrColonySizeNotSet
	}
	if c.generationCount <= 0 {
		return nil, ErrGenerationCountNotSet
	}
	if c.alpha <= 0 {
		return nil, ErrAlphaNotSet
	}
	if c.beta <= 0 {
		return nil, ErrBetaNotSet
	}
	if c.pheromoneMultiplier <= 0 {
		return nil, ErrPheromoneMultiplierNotSet
	}
	if c.evaporationRate <= 0 {
		return nil, ErrEvaporationRateNotSet
	}
	if c.initialPheromone <= 0 {
		return nil, ErrInitialPheromoneNotSet
	}

	c.choice = NewClassicStrategy(c.alpha, c.beta)
	c.bestPerGeneration = make([]*Ant, 0, c.generationCount)

	return c, nil
}

func (c *AntColony) Run() {

	c.pm = pheromone.NewPheromoneMap(c.gr, c.initialPheromone)
	c.generation = make([]*Ant, 0, c.colonySize)
	for i := 0; i < int(c.colonySize); i++ {
		c.generation = append(c.generation, NewAnt(c.gr, c.pm, c.gr.RandomVertex(), c.choice))
	}

	c.currentGeneration = 0

	for c.currentGeneration < c.generationCount {
		// Ant search iteration
		c.iteration()

		// Update best ants
		c.updateBest()

		// Copying ants from previous generation
		c.prepareGeneration()

		c.currentGeneration++
	}
}

func (c *AntColony) iteration() {
	// 1. Route building
	for _, ant := range c.generation {
		ant.Run()
	}

	// 2. Pheromone evaporation
	c.gr.ForEachEdge(func(e *graph.Edge) bool {
		c.pm.Update(e, c.pm.Get(e)*(1.0-c.evaporationRate))
		return false
	})

	// 3. Pheromone update
	for _, ant := range c.generation {

		score := ant.score
		if score == 0 {
			continue
		}

		delta := c.pheromoneMultiplier / score
		tour := ant.route

		for i := 0; i < len(tour)-1; i++ {

			e, ok := c.gr.Edge(tour[i], tour[i+1])
			if !ok {
				continue
			}

			c.pm.Add(e, delta)
		}

		// wrap around
		e, ok := c.gr.Edge(tour[len(tour)-1], tour[0])
		if ok {
			c.pm.Add(e, delta)
		}
	}
}

func (c *AntColony) updateBest() {
	bestInGeneration := c.generation[0]
	for _, ant := range c.generation {
		if ant.score < bestInGeneration.score {
			bestInGeneration = ant
		}
	}

	if bestInGeneration.score < c.bestScore || c.bestTour == nil {
		c.bestTour = bestInGeneration.route
		c.bestScore = bestInGeneration.score
	}
}

func (c *AntColony) prepareGeneration() {
	newGeneration := make([]*Ant, 0, c.colonySize)

	for i := 0; i < int(c.colonySize); i++ {
		newGeneration = append(newGeneration, NewAnt(c.gr, c.pm, c.gr.RandomVertex(), c.choice))
	}

	c.generation = newGeneration
}

func (c *AntColony) BestTour() []*graph.Vertex {
	return c.bestTour
}

func (c *AntColony) BestScore() float64 {
	return c.bestScore
}

func (c *AntColony) BestPerGeneration() []*Ant {
	return c.bestPerGeneration
}
