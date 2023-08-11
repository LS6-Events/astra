package astTraversal

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func setupTestDeclarationTraverser() (*BaseTraverser, error) {
	return createTraverserFromTestFile("declaration.go")
}

func TestTraverser_Declaration(t *testing.T) {
	var node ast.Node
	node = &ast.GenDecl{}

	traverser, err := setupTestDeclarationTraverser()
	assert.NoError(t, err)

	declaration, err := traverser.Declaration(node, "")
	assert.NoError(t, err)
	assert.NotNil(t, declaration)
}

func TestDeclarationTraverser_Result(t *testing.T) {
	traverser, err := setupTestDeclarationTraverser()
	assert.NoError(t, err)

	t.Run("should return the result of a single assignment statement", func(t *testing.T) {
		decl, err := traverser.Declaration(traverser.ActiveFile().AST.Decls[3].(*ast.FuncDecl).Body.List[0].(*ast.AssignStmt), "assignStmt")
		assert.NoError(t, err)
		assert.NotNil(t, decl)

		result, err := decl.Result("assignStmt")
		assert.NoError(t, err)
		assert.Equal(t, "string", result.Type)
	})

	t.Run("should return the result of a multiple assignment statement", func(t *testing.T) {
		decl, err := traverser.Declaration(traverser.ActiveFile().AST.Decls[3].(*ast.FuncDecl).Body.List[1].(*ast.AssignStmt), "assignStmt")
		assert.NoError(t, err)
		assert.NotNil(t, decl)

		result, err := decl.Result("var1")
		assert.NoError(t, err)
		assert.Equal(t, "string", result.Type)

		result, err = decl.Result("var2")
		assert.NoError(t, err)
		assert.Equal(t, "int", result.Type)
	})

	t.Run("should return the result of a root level declaration string", func(t *testing.T) {
		decl, err := traverser.Declaration(traverser.ActiveFile().AST.Decls[1], "MyVar1")
		assert.NoError(t, err)

		result, err := decl.Result("MyVar1")
		assert.NoError(t, err)
		assert.Equal(t, "string", result.Type)
	})

	t.Run("should return the result of a root level declaration float", func(t *testing.T) {
		decl, err := traverser.Declaration(traverser.ActiveFile().AST.Decls[1], "MyVar2")
		assert.NoError(t, err)

		result, err := decl.Result("MyVar2")
		assert.NoError(t, err)
		assert.Equal(t, "float", result.Type)
	})

	t.Run("should return the result of a root level constant declaration string", func(t *testing.T) {
		decl, err := traverser.Declaration(traverser.ActiveFile().AST.Decls[2], "MyConst1")
		assert.NoError(t, err)

		result, err := decl.Result("MyConst1")
		assert.NoError(t, err)
		assert.Equal(t, "string", result.Type)
		assert.Equal(t, "MyConst1", result.ConstantValue)
	})

	t.Run("should return the result of a root level constant declaration int", func(t *testing.T) {
		decl, err := traverser.Declaration(traverser.ActiveFile().AST.Decls[2], "MyConst2")
		assert.NoError(t, err)

		result, err := decl.Result("MyConst2")
		assert.NoError(t, err)
		assert.Equal(t, "int", result.Type)
		assert.Equal(t, "1234", result.ConstantValue)
	})
}
