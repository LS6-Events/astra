package astTraversal

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func setupTestLiteralTraverser() (*BaseTraverser, error) {
	return CreateTraverserFromTestFile("declaration.go")
}

func TestBaseTraverser_Literal(t *testing.T) {
	traverser, err := setupTestLiteralTraverser()
	assert.NoError(t, err)

	node := traverser.ActiveFile().AST.Decls[1].(*ast.GenDecl).Specs[0].(*ast.ValueSpec).Values[0]
	literalTraverser, err := traverser.Literal(node, 0)
	assert.NoError(t, err)
	assert.NotNil(t, literalTraverser)
	assert.Equal(t, node, literalTraverser.Node)
}
