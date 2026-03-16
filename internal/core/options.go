package core

type ColonyOption func(*HeteroAntColony)

const (
	DefaultColonySize          = 100
	DefaultGenerationSize      = 100
	DefaultParentCount         = 2
	DefaultAlpha               = 1.0
	DefaultBeta                = 1.0
	DefaultPheromoneMultiplier = 1.0
	DefaultGenerationPeriod    = 5
	DefaultEvaporationRate     = 0.2
	DefaultInitialPheromone    = 1.0
)

func WithColonySize(size uint) ColonyOption {
	return func(c *HeteroAntColony) {
		c.colonySize = size
	}
}

func WithGenerationCount(count uint) ColonyOption {
	return func(c *HeteroAntColony) {
		c.generationCount = count
	}
}

func WithParentCount(count uint) ColonyOption {
	return func(c *HeteroAntColony) {
		c.parentCount = count
	}
}

func WithDefaultAlpha(alpha float64) ColonyOption {
	return func(c *HeteroAntColony) {
		c.defaultAlpha = alpha
	}
}

func WithDefaultBeta(beta float64) ColonyOption {
	return func(c *HeteroAntColony) {
		c.defaultBeta = beta
	}
}

func WithPheromoneMultiplier(multiplier float64) ColonyOption {
	return func(c *HeteroAntColony) {
		c.pheromoneMultiplier = multiplier
	}
}

func WithInitialPheromone(initial float64) ColonyOption {
	return func(c *HeteroAntColony) {
		c.initialPheromone = initial
	}
}

func WithGenerationPeriod(period uint) ColonyOption {
	return func(c *HeteroAntColony) {
		c.generationPeriod = period
	}
}

func WithChooseStrategy(strategy ChoosePathStrategy) ColonyOption {
	return func(c *HeteroAntColony) {
		c.choose = strategy
	}
}

func WithSelectionStrategy(strategy SelectionStrategy) ColonyOption {
	return func(c *HeteroAntColony) {
		c.selection = strategy
	}
}

func WithCrossoverStrategy(strategy CrossoverStrategy) ColonyOption {
	return func(c *HeteroAntColony) {
		c.crossover = strategy
	}
}

func WithMutationStrategy(strategy MutationStrategy) ColonyOption {
	return func(c *HeteroAntColony) {
		c.mutation = strategy
	}
}
