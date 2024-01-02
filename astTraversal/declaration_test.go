package astTraversal

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func setupTestDeclarationTraverser() (*BaseTraverser, error) {
	return CreateTraverserFromTestFile("declaration.go")
}

func TestTraverser_Declaration(t *testing.T) {
	var node ast.Node = &ast.GenDecl{}

	traverser, err := setupTestDeclarationTraverser()
	assert.NoError(t, err)

	declaration, err := traverser.Declaration(node, "")
	assert.NoError(t, err)
	assert.NotNil(t, declaration)
}
