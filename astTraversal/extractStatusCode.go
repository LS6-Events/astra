package astTraversal

import (
	"go/ast"
	"strconv"
)

// ExtractStatusCode extracts the status code from a handler, assuming it's the first argument
func (t *BaseTraverser) ExtractStatusCode(status ast.Node) (int, error) {
	expr := t.Expression(status)

	constant, err := expr.Value()
	if err != nil {
		return 0, err
	}

	statusCode, err := strconv.Atoi(constant)
	if err != nil {
		return 0, err
	}

	return statusCode, nil
}
