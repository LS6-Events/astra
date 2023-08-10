package astTraversal

import (
	"go/ast"
)

type LiteralTraverser struct {
	Traverser *Traverser
	Node      ast.Node
	ReturnNum int
}

func (t *Traverser) Literal(node ast.Node, returnNum int) (*LiteralTraverser, error) {
	return &LiteralTraverser{
		Traverser: t,
		Node:      node,
		ReturnNum: returnNum,
	}, nil
}

func (lt *LiteralTraverser) Result() (Result, error) {
	return lt.Traverser.Expression(lt.Node).SetReturnNum(lt.ReturnNum).Result()
}
