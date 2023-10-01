package astTraversal

import "errors"

var (
	ErrInvalidNodeType = errors.New("invalid node type")
	ErrInvalidIndex    = errors.New("invalid return index")
	ErrBuiltInFunction = errors.New("builtin function")
)
