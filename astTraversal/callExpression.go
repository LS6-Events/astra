package astTraversal

import (
	"go/ast"
	"go/types"
)

type CallExpressionTraverser struct {
	Traverser *BaseTraverser
	Node      *ast.CallExpr
	File      *FileNode
}

func (t *BaseTraverser) CallExpression(node ast.Node) (*CallExpressionTraverser, error) {
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return nil, ErrInvalidNodeType
	}

	return &CallExpressionTraverser{
		Traverser: t,
		Node:      callExpr,
		File:      t.ActiveFile(),
	}, nil
}

func (c *CallExpressionTraverser) Function() (*FunctionTraverser, error) {
	decl, err := c.Traverser.FindDeclarationForNode(c.Node.Fun)
	if err != nil {
		return nil, err
	}

	function, err := c.Traverser.Function(decl.Decl)
	if err != nil {
		return nil, err
	}

	return function, nil
}

func (c *CallExpressionTraverser) ArgIndex(argName string) (int, bool) {
	for i, arg := range c.Node.Args {
		ident, ok := arg.(*ast.Ident)
		if !ok {
			continue
		}

		if ident.Name == argName {
			return i, true
		}
	}

	return 0, false
}

func (c *CallExpressionTraverser) Args() []ast.Expr {
	return c.Node.Args
}

func (c *CallExpressionTraverser) Type() (*types.Func, error) {
	if c.Node.Fun == nil {
		return nil, ErrInvalidNodeType
	}

	var obj types.Object
	var err error
	switch c.Node.Fun.(type) {
	case *ast.Ident:
		obj, err = c.File.Package.FindObjectForIdent(c.Node.Fun.(*ast.Ident))
	case *ast.SelectorExpr:
		obj, err = c.File.Package.FindObjectForIdent(c.Node.Fun.(*ast.SelectorExpr).Sel)
	default:
		err = ErrInvalidNodeType
	}
	if err != nil {
		return nil, err
	}

	return obj.(*types.Func), nil
}

func (c *CallExpressionTraverser) ReturnType(returnNum int) (types.Type, error) {
	funcType, err := c.Type()
	if err != nil {
		return nil, err
	}

	signature := funcType.Type().(*types.Signature)

	if signature.Results().Len() <= returnNum {
		return nil, ErrInvalidIndex
	}

	return signature.Results().At(returnNum).Type(), nil
}

func (c *CallExpressionTraverser) ArgType(argNum int) (types.Object, error) {
	funcType, err := c.Type()
	if err != nil {
		return nil, err
	}

	signature := funcType.Type().(*types.Signature)

	if signature.Params().Len() <= argNum {
		return nil, ErrInvalidIndex
	}

	return signature.Params().At(argNum), nil
}
