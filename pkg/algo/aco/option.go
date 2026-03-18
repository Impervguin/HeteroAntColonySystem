package aco

type ColonyOption func(*AntColony)

const (
	DefaultColonySize          = 100
	DefaultGenerationSize      = 100
	DefaultAlpha               = 1.0
	DefaultBeta                = 1.0
	DefaultPheromoneMultiplier = 1.0
	DefaultEvaporationRate     = 0.2
	DefaultInitialPheromone    = 1.0
)

func WithColonySize(size uint) ColonyOption {
	return func(c *AntColony) {
		c.colonySize = size
	}
}

func WithGenerationCount(count uint) ColonyOption {
	return func(c *AntColony) {
		c.generationCount = count
	}
}

func WithAlpha(alpha float64) ColonyOption {
	return func(c *AntColony) {
		c.alpha = alpha
	}
}

func WithBeta(beta float64) ColonyOption {
	return func(c *AntColony) {
		c.beta = beta
	}
}

func WithPheromoneMultiplier(multiplier float64) ColonyOption {
	return func(c *AntColony) {
		c.pheromoneMultiplier = multiplier
	}
}

func WithInitialPheromone(initial float64) ColonyOption {
	return func(c *AntColony) {
		c.initialPheromone = initial
	}
}

func WithEvaporationRate(rate float64) ColonyOption {
	return func(c *AntColony) {
		c.evaporationRate = rate
	}
}
