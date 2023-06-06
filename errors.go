package gengo

import "errors"

var (
	ErrConfigNotFound     = errors.New("config not found")
	ErrConfigPortRequired = errors.New("config port is required")
)
