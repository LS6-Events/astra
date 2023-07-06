package gengo

import "errors"

// This file contains all the errors that can be returned by the generator
var (
	ErrConfigNotFound     = errors.New("config not found")
	ErrConfigPortRequired = errors.New("config port is required")

	ErrInputModeNotFound = errors.New("input mode not found")

	ErrOutputModeNotFound     = errors.New("output mode not found")
	ErrOutputFilePathRequired = errors.New("output file path is required")
)
