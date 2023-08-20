package astTraversal

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestExpressionTraverser_SetReturnNum(t *testing.T) {
	bt := &BaseTraverser{}
	et := bt.Expression(&ast.BasicLit{})
	assert.Equal(t, 0, et.ReturnNum)

	et.SetReturnNum(5)
	assert.Equal(t, 5, et.ReturnNum)
}
