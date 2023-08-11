package astTraversal

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func setupTestCallExpressionTraverser() (*ast.FuncDecl, *BaseTraverser, error) {
	traverser, err := createTraverserFromTestFile("callExpression.go")
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

func TestCallExpressionTraverser_FuncName(t *testing.T) {
	function, traverser, err := setupTestCallExpressionTraverser()
	assert.NoError(t, err)

	callExpression, err := traverser.CallExpression(function.Body.List[0].(*ast.ExprStmt).X)
	assert.NoError(t, err)

	assert.Equal(t, "Println", callExpression.FuncName())
}

func TestCallExpressionTraverser_FuncResult(t *testing.T) {
	function, traverser, err := setupTestCallExpressionTraverser()
	assert.NoError(t, err)

	callExpression, err := traverser.CallExpression(function.Body.List[0].(*ast.ExprStmt).X)
	assert.NoError(t, err)

	funcResult := callExpression.FuncResult()

	assert.Equal(t, "Println", funcResult.Type)
	assert.Equal(t, "fmt", funcResult.Package.Path())
}

func TestCallExpressionTraverser_IsExternal(t *testing.T) {
	function, traverser, err := setupTestCallExpressionTraverser()
	assert.NoError(t, err)

	callExpression, err := traverser.CallExpression(function.Body.List[0].(*ast.ExprStmt).X)
	assert.NoError(t, err)

	assert.True(t, callExpression.IsExternal())

	callExpression, err = traverser.CallExpression(function.Body.List[3].(*ast.ExprStmt).X)
	assert.NoError(t, err)

	assert.False(t, callExpression.IsExternal())
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

func TestCallExpressionTraverser_ReturnResult(t *testing.T) {
	function, traverser, err := setupTestCallExpressionTraverser()
	assert.NoError(t, err)

	callExpression, err := traverser.CallExpression(function.Body.List[2].(*ast.ExprStmt).X)
	assert.NoError(t, err)

	result, err := callExpression.ReturnResult(0)
	assert.NoError(t, err)
	assert.Equal(t, "bool", result.Type)

	traverser.Reset()

	result, err = callExpression.ReturnResult(1)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidIndex)
	assert.Equal(t, Result{}, result)
}

func TestCallExpressionTraverser_ArgResult(t *testing.T) {
	function, traverser, err := setupTestCallExpressionTraverser()
	assert.NoError(t, err)

	callExpression, err := traverser.CallExpression(function.Body.List[3].(*ast.ExprStmt).X)
	assert.NoError(t, err)

	result, err := callExpression.ArgResult(0)
	assert.NoError(t, err)
	assert.Equal(t, "int", result.Type)

	traverser.Reset()

	result, err = callExpression.ArgResult(1)
	assert.NoError(t, err)
	assert.Equal(t, "string", result.Type)

	traverser.Reset()

	result, err = callExpression.ArgResult(2)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidIndex)
	assert.Equal(t, Result{}, result)
}

func TestCallExpressionTraverser_ArgIndex(t *testing.T) {
	function, traverser, err := setupTestCallExpressionTraverser()
	assert.NoError(t, err)

	callExpression, err := traverser.CallExpression(function.Body.List[2].(*ast.ExprStmt).X)
	assert.NoError(t, err)

	index, ok := callExpression.ArgIndex("str")
	assert.True(t, ok)
	assert.Equal(t, 0, index)

	index, ok = callExpression.ArgIndex("bool")
	assert.False(t, ok)
	assert.Equal(t, 0, index)
}

func TestCallExpressionTraverser_Args(t *testing.T) {
	function, traverser, err := setupTestCallExpressionTraverser()
	assert.NoError(t, err)

	callExpression, err := traverser.CallExpression(function.Body.List[3].(*ast.ExprStmt).X)
	assert.NoError(t, err)

	results := callExpression.Args()
	assert.Len(t, results, 2)
}
