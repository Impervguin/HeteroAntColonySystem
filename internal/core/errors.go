package core

import "errors"

var (
	ErrChooseStrategyNotSet    = errors.New("choose strategy not set")
	ErrSelectionStrategyNotSet = errors.New("selection strategy not set")
	ErrCrossoverStrategyNotSet = errors.New("crossover strategy not set")
	ErrMutationStrategyNotSet  = errors.New("mutation strategy not set")
)
