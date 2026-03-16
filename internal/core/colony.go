package core

import (
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
	"math/rand/v2"
)

type HeteroAntColony struct {
	g *graph.Graph

	// colonySize is the number of ants in the colony
	colonySize uint

	// generationCount is the number of generations to run
	generationCount uint

	// generationPeriod is the number of generations between
	// crossover and mutation
	generationPeriod uint

	// parentCount is the number of parents to select for crossover
	// must be less than colonySize
	parentCount uint

	// ant colony algorithm parameters
	defaultAlpha        float64
	defaultBeta         float64
	evaporationRate     float64
	pheromoneMultiplier float64
	initialPheromone    float64

	// strategy for ants
	choose ChoosePathStrategy

	// strategies for producing offsprings
	selection SelectionStrategy
	crossover CrossoverStrategy
	mutation  MutationStrategy

	// in work state of the ant colony
	state *AntColonyState

	// result of the ant colony work
	result *AntColonyResult
}

type AntColonyState struct {
	g             *graph.Graph
	pm            *pheromone.PheromoneMap
	initialVertex *graph.Vertex

	currentGeneration uint
	generation        []*HeteroAnt

	bestPerGeneration []*HeteroAnt
	bestTour          []*graph.Vertex
	bestScore         float64
}

type AntColonyResult struct {
	bestPerGeneration []*HeteroAnt

	bestTour  []*graph.Vertex
	bestScore float64
}

func NewHeteroAntColony(g *graph.Graph, opts ...ColonyOption) (*HeteroAntColony, error) {
	c := &HeteroAntColony{
		g:                   g,
		colonySize:          DefaultColonySize,
		generationCount:     DefaultGenerationSize,
		parentCount:         DefaultParentCount,
		defaultAlpha:        DefaultAlpha,
		defaultBeta:         DefaultBeta,
		evaporationRate:     DefaultEvaporationRate,
		pheromoneMultiplier: DefaultPheromone,
		generationPeriod:    DefaultGenerationPeriod,
		choose:              nil,
		selection:           nil,
		crossover:           nil,
		mutation:            nil,
		state:               nil,
		result:              nil,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.choose == nil {
		return nil, ErrChooseStrategyNotSet
	}

	if c.selection == nil {
		return nil, ErrSelectionStrategyNotSet
	}

	if c.crossover == nil {
		return nil, ErrCrossoverStrategyNotSet
	}

	if c.mutation == nil {
		return nil, ErrMutationStrategyNotSet
	}

	return c, nil
}

func (c *HeteroAntColony) Run() {
	c.state = &AntColonyState{
		g:                 c.g,
		pm:                pheromone.NewPheromoneMap(c.g, c.initialPheromone),
		currentGeneration: 0,
		generation:        make([]*HeteroAnt, 0, c.colonySize),
		bestPerGeneration: make([]*HeteroAnt, 0, c.generationCount),
		bestTour:          nil,
		bestScore:         0.,
	}
	state := c.state

	// Set any as initial vertex
	state.g.ForEachVertex(func(v *graph.Vertex) bool {
		state.initialVertex = v
		return true
	})

	// Init first generation
	for i := 0; i < int(c.colonySize); i++ {
		ant := NewHeteroAnt(c.defaultAlpha, c.defaultBeta, c.choose)
		state.generation = append(c.state.generation, ant)
	}

	state.currentGeneration++

	for state.currentGeneration < c.generationCount {
		// Ant search iteration
		c.iteration()

		// Update best ants
		c.updateBest()

		// Next generation
		state.currentGeneration++
		if state.currentGeneration%c.generationPeriod == 0 {
			// Creating new ants via genetic algorithm
			c.evolve()
		} else {
			// Copying ants from previous generation
			c.stagnate()
		}
	}

	// Final result
	c.result = &AntColonyResult{
		bestPerGeneration: state.bestPerGeneration,
		bestTour:          state.bestTour,
		bestScore:         state.bestScore,
	}

	c.state = nil
}

func (c *HeteroAntColony) iteration() {
	state := c.state
	if state == nil {
		return
	}

	// 1. Route building
	for _, ant := range state.generation {
		ant.StartAnt(state.g, state.pm, state.initialVertex)
		ant.Run()
	}

	// 2. Pheromone evaporation
	state.g.ForEachEdge(func(e *graph.Edge) bool {
		state.pm.Update(e, state.pm.Get(e)*(1.0-c.evaporationRate))
		return false
	})

	// 3. Pheromone update
	for _, ant := range state.generation {

		score := ant.Score()
		if score == 0 {
			continue
		}

		delta := c.pheromoneMultiplier / score
		tour := ant.result.tour

		for i := 0; i < len(tour)-1; i++ {

			e, ok := state.g.Edge(tour[i], tour[i+1])
			if !ok {
				continue
			}

			state.pm.Add(e, delta)
		}

		// wrap around
		e, ok := state.g.Edge(tour[len(tour)-1], tour[0])
		if ok {
			state.pm.Add(e, delta)
		}
	}
}

func (c *HeteroAntColony) findBestGenerationAnt() *HeteroAnt {
	state := c.state
	if state == nil {
		return nil
	}

	best := state.generation[0]
	for _, ant := range state.generation {
		if ant.Score() < best.Score() {
			best = ant
		}
	}
	return best
}

func (c *HeteroAntColony) updateBest() {
	state := c.state
	if state == nil {
		return
	}

	// Finding best ant
	ant := c.findBestGenerationAnt()
	state.bestPerGeneration = append(state.bestPerGeneration, ant)

	// Best tour
	if ant.Score() < state.bestScore || state.bestTour == nil {
		state.bestTour = ant.Tour()
		state.bestScore = ant.Score()
	}
}

func (c *HeteroAntColony) evolve() {
	state := c.state
	if state == nil {
		return
	}
	newGeneration := make([]*HeteroAnt, 0, c.colonySize)

	// 1. Selection
	parents := c.selection.Select(state.generation, c.parentCount)

	for i := uint(0); i < c.colonySize; i++ {
		// 2. Crossover
		var parent1 uint = rand.UintN(c.parentCount)
		var parent2 uint = rand.UintN(c.parentCount)
		for ; parent1 == parent2; parent2 = rand.UintN(c.parentCount) {
		}

		child := c.crossover.Crossover(parents[parent1], parents[parent2])

		// 3. Mutation
		child = c.mutation.Mutate(child)

		// 4. Add to generation
		newGeneration = append(newGeneration, child)
	}

	state.generation = newGeneration
}

func (c *HeteroAntColony) stagnate() {
	state := c.state
	if state == nil {
		return
	}

	newGeneration := make([]*HeteroAnt, 0, c.colonySize)

	for _, ant := range state.generation {
		newGeneration = append(newGeneration, NewHeteroAnt(ant.Alpha(), ant.Beta(), c.choose))
	}
	state.generation = newGeneration
}
