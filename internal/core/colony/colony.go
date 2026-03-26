package colony

import (
	"HeteroAntColonySystem/internal/core/ant"
	"HeteroAntColonySystem/internal/core/config"
	"HeteroAntColonySystem/internal/core/errors"
	"HeteroAntColonySystem/internal/core/strategy"
	"HeteroAntColonySystem/pkg/graph"
	"HeteroAntColonySystem/pkg/pheromone"
)

// HeteroAntColony represents a colony of heterogeneous ants in the ACO algorithm.
// It manages the evolution process across multiple generations, coordinating
// ant movement, pheromone updates, and solution optimization.
//
// The colony uses functional options pattern for configuration, allowing flexible
// setup of algorithm parameters and strategies.
type HeteroAntColony struct {
	// Configuration (immutable after creation)
	pathChoice     strategy.PathChoiceStrategy
	pheromoneApply strategy.PheromoneApplyingStrategy

	defaultAlpha        float64
	defaultBeta         float64
	pheromoneMultiplier float64
	evaporationRate     float64
	initialPheromone    float64

	generationCount uint
	colonySize      uint

	// Runtime state
	g    *graph.Graph
	pm   *pheromone.PheromoneMap
	ants []*ant.HeteroAnt

	// Best solution tracking
	best  *ant.HeteroAnt
	score float64
}

// NewHeteroAntColony creates a new heterogeneous ant colony with the given options.
// Returns an error if required strategies are not set or if parameters are invalid.
func NewHeteroAntColony(opts ...config.HeteroAntColonyOption) (*HeteroAntColony, error) {
	// Create config with default values
	cfg := &config.ColonyConfig{
		PathChoice:          nil,
		PheromoneApply:      nil,
		DefaultAlpha:        config.DefaultAlpha,
		DefaultBeta:         config.DefaultBeta,
		PheromoneMultiplier: config.DefaultPheromoneMult,
		EvaporationRate:     config.DefaultEvaporation,
		InitialPheromone:    config.DefaultPheromone,
		GenerationCount:     config.DefaultGenerations,
		ColonySize:          config.DefaultColonySize,
	}

	// Apply user-provided options
	for _, opt := range opts {
		opt(cfg)
	}

	// Validate required fields
	if cfg.PathChoice == nil {
		return nil, errors.ErrPathChoiceStrategyNotSet
	}
	if cfg.PheromoneApply == nil {
		return nil, errors.ErrPheromoneApplyStrategyNotSet
	}

	// Validate parameters
	if cfg.ColonySize == 0 {
		return nil, errors.ErrColonySizeInvalid
	}
	if cfg.GenerationCount == 0 {
		return nil, errors.ErrGenerationCountInvalid
	}
	if cfg.PheromoneMultiplier <= 0 {
		return nil, errors.ErrPheromoneMultiplierInvalid
	}
	if cfg.EvaporationRate <= 0 {
		return nil, errors.ErrEvaporationRateInvalid
	}
	if cfg.InitialPheromone <= 0 {
		return nil, errors.ErrInitialPheromoneInvalid
	}

	// Create colony with config
	return &HeteroAntColony{
		pathChoice:          cfg.PathChoice,
		pheromoneApply:      cfg.PheromoneApply,
		defaultAlpha:        cfg.DefaultAlpha,
		defaultBeta:         cfg.DefaultBeta,
		pheromoneMultiplier: cfg.PheromoneMultiplier,
		evaporationRate:     cfg.EvaporationRate,
		initialPheromone:    cfg.InitialPheromone,
		generationCount:     cfg.GenerationCount,
		colonySize:          cfg.ColonySize,
	}, nil
}

// Prepare initializes the colony for execution on the given graph.
// This must be called before Run.
func (c *HeteroAntColony) Prepare(g *graph.Graph) error {
	c.g = g
	c.pm = pheromone.NewPheromoneMap(g, c.initialPheromone)
	c.ants = make([]*ant.HeteroAnt, 0, c.colonySize)
	c.best = nil
	c.score = 0
	return nil
}

// Run executes the ant colony optimization algorithm.
// It runs for generationCount iterations, each time:
// 1. Having all ants construct paths
// 2. Finding the best solution in this generation
// 3. Evaporating pheromones
// 4. Applying pheromone updates based on ant paths
func (c *HeteroAntColony) Run() error {
	// Initialize ants for first generation
	for i := uint(0); i < c.colonySize; i++ {
		c.ants = append(c.ants, ant.NewHeteroAnt(
			c.defaultAlpha,
			c.defaultBeta,
			c.pheromoneMultiplier,
			c.pathChoice,
			c.pheromoneApply,
		))
	}

	// Main optimization loop
	for gen := uint(0); gen < c.generationCount; gen++ {
		// Phase 1: All ants construct paths
		for _, a := range c.ants {
			a.Prepare(c.g, c.pm)
			if err := a.Run(); err != nil {
				return err
			}
		}

		// Phase 2: Find best ant in this generation
		bestInGen := c.ants[0]
		for _, a := range c.ants {
			if a.Score() < bestInGen.Score() {
				bestInGen = a
			}
		}

		// Update global best if needed
		if c.best == nil || bestInGen.Score() < c.best.Score() {
			c.best = bestInGen.FullCopy()
		}

		// Phase 3: Evaporate pheromones
		c.evaporatePheromones()

		// Phase 4: Apply pheromone updates
		for _, a := range c.ants {
			if err := a.ApplyPheromone(); err != nil {
				return err
			}
		}

		// Phase 5: Prepare next generation (reuse ant objects with same config)
		c.prepareNextGeneration()
	}

	// Final score
	if c.best != nil {
		c.score = c.best.Score()
	}

	return nil
}

// evaporatePheromones applies the evaporation rate to all edges in the graph.
func (c *HeteroAntColony) evaporatePheromones() {
	factor := 1 - c.evaporationRate
	c.g.ForEachEdge(func(e *graph.Edge) bool {
		c.pm.Update(e, c.pm.Get(e)*factor)
		return false
	})
}

// prepareNextGeneration resets ant paths while keeping their configurations.
// This allows ants to learn from previous generations' experiences.
func (c *HeteroAntColony) prepareNextGeneration() {
	// Ants retain their alpha, beta, pheromoneMultiplier, and strategies
	// They will be re-prepared with new paths in the next iteration
	// This avoids recreating ant objects unnecessarily
}

// Score returns the score of the best solution found.
func (c *HeteroAntColony) Score() float64 {
	return c.score
}

// BestPath returns the vertices in the best solution found.
func (c *HeteroAntColony) BestPath() []*graph.Vertex {
	if c.best == nil {
		return nil
	}
	return c.best.Path()
}
