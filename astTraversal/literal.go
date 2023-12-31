package astTraversal

import (
	"go/ast"
	"go/types"
)

type LiteralTraverser struct {
	Traverser *BaseTraverser
	Node      ast.Node
	ReturnNum int
}

func (t *BaseTraverser) Literal(node ast.Node, returnNum int) (*LiteralTraverser, error) {
	return &LiteralTraverser{
		Traverser: t,
		Node:      node,
		ReturnNum: returnNum,
	}, nil
}

func (lt *LiteralTraverser) Type() (types.Type, error) {
	exprNode, ok := lt.Node.(ast.Expr)
	if !ok {
		return nil, ErrInvalidNodeType
	}

	return lt.Traverser.Expression(exprNode).SetReturnNum(lt.ReturnNum).Type()
}
