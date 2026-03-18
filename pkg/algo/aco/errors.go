package aco

import "errors"

var (
	ErrColonySizeNotSet          = errors.New("colony size not set")
	ErrGenerationCountNotSet     = errors.New("generation count not set")
	ErrAlphaNotSet               = errors.New("alpha not set")
	ErrBetaNotSet                = errors.New("beta not set")
	ErrPheromoneMultiplierNotSet = errors.New("pheromone multiplier not set")
	ErrEvaporationRateNotSet     = errors.New("evaporation rate not set")
	ErrInitialPheromoneNotSet    = errors.New("initial pheromone not set")
	ErrInitialVertexNotSet       = errors.New("initial vertex not set")
)
