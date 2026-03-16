package tsplib

import "fmt"

var (
	ErrInvalidFormat   = fmt.Errorf("invalid TSPLIB format")
	ErrUnsupportedType = fmt.Errorf("unsupported TSPLIB type")
	ErrInvalidData     = fmt.Errorf("invalid data")
	ErrAdapterNotFound = fmt.Errorf("adapter not found")
	ErrSectionNotFound = fmt.Errorf("data section not found")
)
