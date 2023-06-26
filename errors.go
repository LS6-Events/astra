package gengo

import "errors"

var (
	ErrConfigNotFound     = errors.New("config not found")
	ErrConfigPortRequired = errors.New("config port is required")

	ErrFailedToGetCaller = errors.New("failed to get caller")

	ErrInputModeNotFound = errors.New("input mode not found")

	ErrOutputModeNotFound     = errors.New("output mode not found")
	ErrOutputFilePathRequired = errors.New("output file path is required")
)
