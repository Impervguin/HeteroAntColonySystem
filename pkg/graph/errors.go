package graph

import "errors"

var (
	ErrVertexNotFound  = errors.New("vertex not found")
	ErrVertexAlreadyIn = errors.New("vertex already added")
)
