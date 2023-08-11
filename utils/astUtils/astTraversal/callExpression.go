package astTraversal

import (
	"go/ast"
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

func (c *CallExpressionTraverser) FuncName() string {
	result := c.Traverser.ExtractVarName(c.Node.Fun)

	return result.Type
}

func (c *CallExpressionTraverser) FuncResult() Result {
	return c.Traverser.ExtractVarName(c.Node.Fun)
}

func (c *CallExpressionTraverser) IsExternal() bool {
	node, ok := c.Node.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := node.X.(*ast.Ident)
	if !ok {
		return false
	}

	return c.File.IsImportedPackage(ident.Name)
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

func (c *CallExpressionTraverser) ReturnResult(returnNum int) (Result, error) {
	function, err := c.Function()
	if err != nil {
		return Result{}, err
	}

	if len(function.Results()) <= returnNum {
		return Result{}, ErrInvalidIndex
	}

	resultType := function.Results()[returnNum]

	result, err := c.Traverser.Expression(resultType.Type).Result()

	return result, nil
}

func (c *CallExpressionTraverser) ArgResult(argNum int) (Result, error) {
	if len(c.Node.Args) <= argNum {
		return Result{}, ErrInvalidIndex
	}
	arg := c.Node.Args[argNum]

	result, err := c.Traverser.Expression(arg).Result()
	if err != nil {
		return Result{}, err
	}

	return result, nil
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
