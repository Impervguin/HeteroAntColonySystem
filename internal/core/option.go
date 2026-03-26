package core

type HeteroAntColonyOption func(colony *HeteroAntColony)

var (
	DefaultDefaultAlpha        float64 = 1
	DefaultDefaultBeta         float64 = 1
	DefaultPheromoneMultiplier float64 = 1
	DefaultEvaporationRate     float64 = 0.2
	DefaultInitialPheromone    float64 = 1.0
	DefaultGenerationCount     uint    = 100
	DefaultColonySize          uint    = 100
)

func WithDefaultAlpha(alpha float64) HeteroAntColonyOption {
	return func(colony *HeteroAntColony) {
		colony.defaultAlpha = alpha
	}
}

func WithDefaultBeta(beta float64) HeteroAntColonyOption {
	return func(colony *HeteroAntColony) {
		colony.defaultBeta = beta
	}
}

func WithPheromoneMultiplier(multiplier float64) HeteroAntColonyOption {
	return func(colony *HeteroAntColony) {
		colony.pheromoneMultiplier = multiplier
	}
}

func WithEvaporationRate(rate float64) HeteroAntColonyOption {
	return func(colony *HeteroAntColony) {
		colony.evaporationRate = rate
	}
}

func WithInitialPheromone(initial float64) HeteroAntColonyOption {
	return func(colony *HeteroAntColony) {
		colony.initialPheromone = initial
	}
}

func WithGenerationCount(count uint) HeteroAntColonyOption {
	return func(colony *HeteroAntColony) {
		colony.generationCount = count
	}
}

func WithColonySize(size uint) HeteroAntColonyOption {
	return func(colony *HeteroAntColony) {
		colony.colonySize = size
	}
}

func WithPathChoiceStrategy(choice PathChoiceStrategy) HeteroAntColonyOption {
	return func(colony *HeteroAntColony) {
		colony.pathChoice = choice
	}
}
func WithPheromoneApplyingStrategy(apply PheromoneApplyingStrategy) HeteroAntColonyOption {
	return func(colony *HeteroAntColony) {
		colony.pheromoneApply = apply
	}
}
