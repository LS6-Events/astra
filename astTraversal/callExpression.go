package astTraversal

import (
	"errors"
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
		return nil, errors.Join(ErrInvalidNodeType, errors.New("expected *ast.CallExpr"))
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
		return nil, errors.Join(ErrInvalidNodeType, errors.New("expected *ast.CallExpr.Fun to not be nil"))
	}

	var obj types.Object
	var err error
	switch nodeFun := c.Node.Fun.(type) {
	case *ast.Ident:
		obj, err = c.File.Package.FindObjectForIdent(nodeFun)
	case *ast.SelectorExpr:
		obj, err = c.File.Package.FindObjectForIdent(nodeFun.Sel)
	default:
		err = errors.Join(ErrInvalidNodeType, errors.New("expected *ast.CallExpr.Fun to be *ast.Ident or *ast.SelectorExpr"))
	}
	if err != nil {
		return nil, err
	}

	switch objType := obj.(type) {
	case *types.Func:
		return objType, nil
	case *types.Builtin:
		return nil, ErrBuiltInFunction
	}

	return nil, errors.Join(ErrInvalidNodeType, errors.New("expected *types.Func"))
}

func (c *CallExpressionTraverser) ReturnType(returnNum int) (types.Type, error) {
	funcType, err := c.Type()
	if err != nil {
		return nil, err
	}

	signature, ok := funcType.Type().(*types.Signature)
	if !ok {
		return nil, errors.Join(ErrInvalidNodeType, errors.New("expected *types.Signature"))
	}

	if signature.Results().Len() <= returnNum {
		return nil, errors.Join(ErrInvalidIndex, errors.New("expected returnNum to be less than signature.Results().Len()"))
	}

	return signature.Results().At(returnNum).Type(), nil
}

func (c *CallExpressionTraverser) ArgType(argNum int) (types.Object, error) {
	funcType, err := c.Type()
	if err != nil {
		return nil, err
	}

	signature, ok := funcType.Type().(*types.Signature)
	if !ok {
		return nil, errors.Join(ErrInvalidNodeType, errors.New("expected *types.Signature"))
	}

	if signature.Params().Len() <= argNum {
		return nil, errors.Join(ErrInvalidIndex, errors.New("expected argNum to be less than signature.Params().Len()"))
	}

	return signature.Params().At(argNum), nil
}
