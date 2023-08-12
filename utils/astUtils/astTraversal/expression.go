package astTraversal

import (
	"errors"
	"go/ast"
	"strings"
)

type ExpressionTraverser struct {
	Traverser *BaseTraverser
	Node      ast.Expr
	File      *FileNode
	ReturnNum int
}

func (t *BaseTraverser) Expression(node ast.Node) *ExpressionTraverser {
	return &ExpressionTraverser{
		Traverser: t,
		Node:      node.(ast.Expr),
		File:      t.ActiveFile(),
		ReturnNum: 0,
	}
}

func (e *ExpressionTraverser) SetReturnNum(returnNum int) *ExpressionTraverser {
	e.ReturnNum = returnNum
	return e
}

func (e *ExpressionTraverser) ReservedFunctions(callExpr *ast.CallExpr) (Result, error) {
	if ident, ok := callExpr.Fun.(*ast.Ident); ok {
		switch ident.Name {
		case "new":
			fallthrough
		case "make":
			arg := callExpr.Args[0]
			result, err := e.Traverser.Expression(arg).Result()
			if err != nil {
				return Result{}, err
			}

			return result, nil
		case "len":
			return Result{
				Type: "int",
			}, nil
		}
	}

	return Result{}, nil
}

func (e *ExpressionTraverser) DoesNeedTracing() bool {
	switch n := e.Node.(type) {
	case *ast.StarExpr:
		return e.Traverser.Expression(n.X).DoesNeedTracing()
	case *ast.UnaryExpr:
		return e.Traverser.Expression(n.X).DoesNeedTracing()
	case *ast.BasicLit:
		return false
	case *ast.CompositeLit:
		return false
	}

	return true
}

func (e *ExpressionTraverser) Result() (Result, error) {
	switch n := e.Node.(type) {
	case *ast.StarExpr:
		return e.Traverser.Expression(n.X).Result()
	case *ast.UnaryExpr:
		return e.Traverser.Expression(n.X).Result()
	case *ast.Ident:
		return e.Traverser.ExtractVarName(n), nil
	case *ast.SelectorExpr:
		return e.Traverser.ExtractVarName(n), nil
	case *ast.ArrayType:
		embeddedType := e.Traverser.ExtractVarName(n.Elt)
		return Result{
			Type:      "slice",
			Package:   embeddedType.Package,
			SliceType: embeddedType.Type,
		}, nil
	case *ast.MapType:
		keyType := e.Traverser.ExtractVarName(n.Key)
		valueType := e.Traverser.ExtractVarName(n.Value)
		return Result{
			Type:          "map",
			Package:       valueType.Package,
			MapValType:    valueType.Type,
			MapKeyType:    keyType.Type,
			MapKeyPackage: keyType.Package,
		}, nil
	case *ast.CompositeLit:
		return e.Traverser.Expression(n.Type).Result()
	case *ast.BasicLit:
		result := Result{
			Type:          strings.ToLower(n.Kind.String()),
			Package:       e.File.Package,
			ConstantValue: n.Value,
		}

		if result.Type == "string" {
			result.ConstantValue = strings.Trim(result.ConstantValue, "\"")
		}

		return result, nil
	case *ast.CallExpr:
		result, err := e.ReservedFunctions(n)
		if err != nil {
			return Result{}, err
		}
		if result.Type != "" {
			return result, nil
		}

		callExpr, err := e.Traverser.CallExpression(n)
		if err != nil {
			return Result{}, err
		}

		return callExpr.ReturnResult(e.ReturnNum)
	default:
		return Result{}, errors.New("unsupported expression type")
	}

}
