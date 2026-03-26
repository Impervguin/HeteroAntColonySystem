package core

import "errors"

var (
	ErrAntNotPrepared = errors.New("ant was not prepared")
	ErrAntNotDone     = errors.New("ant was not done")
)

var (
	ErrGenerationCountInvalid     = errors.New("generation count must be greater than 0")
	ErrColonySizeInvalid          = errors.New("colony size must be greater than 0")
	ErrPheromoneMultiplierInvalid = errors.New("pheromone multiplier must be greater than 0")
	ErrEvaporationRateInvalid     = errors.New("evaporation rate must be greater than 0")
	ErrInitialPheromoneInvalid    = errors.New("initial pheromone must be greater than 0")

	ErrPathChoiceStrategyNotSet     = errors.New("path choice strategy not set")
	ErrPheromoneApplyStrategyNotSet = errors.New("pheromone apply strategy not set")
)
