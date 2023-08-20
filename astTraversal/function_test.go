package astTraversal

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestBaseTraverser_Function(t *testing.T) {
	traverser, err := createTraverserFromTestFile("declaration.go")
	assert.NoError(t, err)

	t.Run("should return a function traverser for a function literal", func(t *testing.T) {
		funcLit := traverser.ActiveFile().AST.Decls[5].(*ast.GenDecl).Specs[0].(*ast.ValueSpec).Values[0].(*ast.FuncLit)
		function, err := traverser.Function(funcLit)
		assert.NoError(t, err)
		assert.NotNil(t, function)
		assert.Equal(t, function.Node, funcLit)
	})

	t.Run("should return a function traverser for a function declaration", func(t *testing.T) {
		funcDecl := traverser.ActiveFile().AST.Decls[3].(*ast.FuncDecl)
		function, err := traverser.Function(funcDecl)
		assert.NoError(t, err)
		assert.NotNil(t, function)
		assert.IsType(t, function.Node, &ast.FuncLit{}) // We convert to func lit internally
	})

	t.Run("should return an error for a non-function node", func(t *testing.T) {
		_, err := traverser.Function(traverser.ActiveFile().AST.Decls[1])
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidNodeType)
	})
}

func TestFunctionTraverser_Arguments(t *testing.T) {
	traverser, err := createTraverserFromTestFile("declaration.go")
	assert.NoError(t, err)

	t.Run("should return the arguments of a function", func(t *testing.T) {
		function, err := traverser.Function(traverser.ActiveFile().AST.Decls[4])
		assert.NoError(t, err)
		assert.NotNil(t, function)

		arguments := function.Arguments()
		assert.Len(t, arguments, 2)
	})

	t.Run("should return an empty array for a function with no arguments", func(t *testing.T) {
		function, err := traverser.Function(traverser.ActiveFile().AST.Decls[3])
		assert.NoError(t, err)
		assert.NotNil(t, function)

		arguments := function.Arguments()
		assert.Len(t, arguments, 0)
	})
}

func TestFunctionTraverser_Results(t *testing.T) {
	traverser, err := createTraverserFromTestFile("declaration.go")
	assert.NoError(t, err)

	t.Run("should return the results of a function", func(t *testing.T) {
		function, err := traverser.Function(traverser.ActiveFile().AST.Decls[4])
		assert.NoError(t, err)
		assert.NotNil(t, function)

		results := function.Results()
		assert.Len(t, results, 2)
	})

	t.Run("should return an empty array for a function with void return", func(t *testing.T) {
		function, err := traverser.Function(traverser.ActiveFile().AST.Decls[3])
		assert.NoError(t, err)
		assert.NotNil(t, function)

		results := function.Results()
		assert.Len(t, results, 0)
	})
}

func TestFunctionTraverser_FindArgumentNameByType(t *testing.T) {
	traverser, err := createTraverserFromTestFile("declaration.go")
	assert.NoError(t, err)

	t.Run("should return the name of an argument by type", func(t *testing.T) {
		function, err := traverser.Function(traverser.ActiveFile().AST.Decls[4])
		assert.NoError(t, err)
		assert.NotNil(t, function)

		name := function.FindArgumentNameByType("string", "", false)
		assert.Equal(t, "param1", name)
	})

	t.Run("should return the empty if it can't find it", func(t *testing.T) {
		function, err := traverser.Function(traverser.ActiveFile().AST.Decls[4])
		assert.NoError(t, err)
		assert.NotNil(t, function)

		name := function.FindArgumentNameByType("float64", "", false)
		assert.Equal(t, "", name)
	})
}
