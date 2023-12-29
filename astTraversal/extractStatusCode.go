package astTraversal

import (
	"errors"
	"go/ast"
	"strconv"
)

// ExtractStatusCode extracts the status code from a handler, assuming it's the first argument.
func (t *BaseTraverser) ExtractStatusCode(status ast.Node) (int, error) {
	exprNode, ok := status.(ast.Expr)
	if !ok {
		return 0, ErrInvalidNodeType
	}

	expr := t.Expression(exprNode)

	constant, err := expr.Value()
	if err != nil {
		return 0, err
	}

	statusCode, err := strconv.Atoi(constant)
	if err != nil {
		return 0, errors.Join(errors.New("status code is not a number"), err)
	}

	return statusCode, nil
}
