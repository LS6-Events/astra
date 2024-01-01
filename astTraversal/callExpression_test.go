package astTraversal

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func setupTestCallExpressionTraverser() (*ast.FuncDecl, *BaseTraverser, error) {
	traverser, err := CreateTraverserFromTestFile("callExpression.go")
	if err != nil {
		return nil, nil, err
	}

	var function *ast.FuncDecl
	ast.Inspect(traverser.ActiveFile().AST, func(node ast.Node) bool {
		if node == nil {
			return false
		}

		funcDecl, ok := node.(*ast.FuncDecl)
		if ok && funcDecl.Name.Name == "callExpr" {
			function = funcDecl
			return false
		}

		return true
	})

	return function, traverser, nil
}

func TestTraverser_CallExpression(t *testing.T) {
	var node ast.Node
	node = &ast.CallExpr{}

	_, traverser, err := setupTestCallExpressionTraverser()
	assert.NoError(t, err)

	callExpression, err := traverser.CallExpression(node)
	assert.NoError(t, err)
	assert.NotNil(t, callExpression)

	node = &ast.Ident{}
	callExpression, err = traverser.CallExpression(node)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidNodeType)
	assert.Nil(t, callExpression)
}

func TestCallExpressionTraverser_Function(t *testing.T) {
	function, traverser, err := setupTestCallExpressionTraverser()
	assert.NoError(t, err)

	callExpression, err := traverser.CallExpression(function.Body.List[0].(*ast.ExprStmt).X)
	assert.NoError(t, err)

	functionTraverser, err := callExpression.Function()
	assert.NoError(t, err)

	assert.NotNil(t, functionTraverser.Node)
}

func TestCallExpressionTraverser_Args(t *testing.T) {
	function, traverser, err := setupTestCallExpressionTraverser()
	assert.NoError(t, err)

	callExpression, err := traverser.CallExpression(function.Body.List[3].(*ast.ExprStmt).X)
	assert.NoError(t, err)

	results := callExpression.Args()
	assert.Len(t, results, 2)
}
