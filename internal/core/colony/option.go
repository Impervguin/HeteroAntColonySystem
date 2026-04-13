package colony

import (
	"HeteroAntColonySystem/internal/core/ant"
)

// Package config provides configuration options for the Hetero Ant Colony Optimization algorithm.
// It includes functional options for configuring colony parameters.

type HeteroAntColonyOption func(colony *ColonyConfig)

// ColonyConfig holds the configuration parameters for a heterogeneous ant colony.
// This is used with functional options to configure the colony.
type ColonyConfig struct {
	PathChoice     ant.PathChoiceStrategy
	PheromoneApply ant.PheromoneApplyingStrategy
	ParentSelect   ParentSelectionStrategy
	Crossover      CrossoverStrategy
	Mutation       MutationStrategy

	DefaultAlpha        float64
	DefaultBeta         float64
	PheromoneMultiplier float64
	EvaporationRate     float64
	InitialPheromone    float64

	GenerationPeriod uint
	ParentCount      uint

	GenerationCount uint
	ColonySize      uint
}

// Default configuration values.
const (
	DefaultAlpha            = 1.0
	DefaultBeta             = 1.0
	DefaultPheromoneMult    = 1.0
	DefaultEvaporation      = 0.2
	DefaultPheromone        = 1.0
	DefaultGenerations      = 100
	DefaultColonySize       = 100
	DefaultGenerationPeriod = 5
	DefaultParentCount      = 20
)

// WithDefaultAlpha sets the default alpha parameter for ants.
func WithDefaultAlpha(alpha float64) HeteroAntColonyOption {
	return func(colony *ColonyConfig) {
		colony.DefaultAlpha = alpha
	}
}

// WithDefaultBeta sets the default beta parameter for ants.
func WithDefaultBeta(beta float64) HeteroAntColonyOption {
	return func(colony *ColonyConfig) {
		colony.DefaultBeta = beta
	}
}

// WithPheromoneMultiplier sets the pheromone multiplier for ants.
func WithPheromoneMultiplier(multiplier float64) HeteroAntColonyOption {
	return func(colony *ColonyConfig) {
		colony.PheromoneMultiplier = multiplier
	}
}

// WithEvaporationRate sets the pheromone evaporation rate.
func WithEvaporationRate(rate float64) HeteroAntColonyOption {
	return func(colony *ColonyConfig) {
		colony.EvaporationRate = rate
	}
}

// WithInitialPheromone sets the initial pheromone level on all edges.
func WithInitialPheromone(initial float64) HeteroAntColonyOption {
	return func(colony *ColonyConfig) {
		colony.InitialPheromone = initial
	}
}

// WithGenerationCount sets the number of generations to run.
func WithGenerationCount(count uint) HeteroAntColonyOption {
	return func(colony *ColonyConfig) {
		colony.GenerationCount = count
	}
}

// WithColonySize sets the number of ants in the colony.
func WithColonySize(size uint) HeteroAntColonyOption {
	return func(colony *ColonyConfig) {
		colony.ColonySize = size
	}
}

func WithGenerationPeriod(period uint) HeteroAntColonyOption {
	return func(colony *ColonyConfig) {
		colony.GenerationPeriod = period
	}
}

func WithParentCount(count uint) HeteroAntColonyOption {
	return func(colony *ColonyConfig) {
		colony.ParentCount = count
	}
}

// WithPathChoiceStrategy sets the path selection strategy for ants.
func WithPathChoiceStrategy(choice ant.PathChoiceStrategy) HeteroAntColonyOption {
	return func(colony *ColonyConfig) {
		colony.PathChoice = choice
	}
}

// WithPheromoneApplyingStrategy sets the pheromone update strategy.
func WithPheromoneApplyingStrategy(apply ant.PheromoneApplyingStrategy) HeteroAntColonyOption {
	return func(colony *ColonyConfig) {
		colony.PheromoneApply = apply
	}
}

// WithParentSelectionStrategy sets the parent selection strategy.
func WithParentSelectionStrategy(sel ParentSelectionStrategy) HeteroAntColonyOption {
	return func(colony *ColonyConfig) {
		colony.ParentSelect = sel
	}
}

// WithCrossoverStrategy sets the crossover strategy.
func WithCrossoverStrategy(crossover CrossoverStrategy) HeteroAntColonyOption {
	return func(colony *ColonyConfig) {
		colony.Crossover = crossover
	}
}

// WithMutationStrategy sets the mutation strategy.
func WithMutationStrategy(mutation MutationStrategy) HeteroAntColonyOption {
	return func(colony *ColonyConfig) {
		colony.Mutation = mutation
	}
}
