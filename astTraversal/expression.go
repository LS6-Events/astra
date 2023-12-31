package astTraversal

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

type ExpressionTraverser struct {
	Traverser *BaseTraverser
	Node      ast.Expr
	File      *FileNode
	ReturnNum int
}

func (t *BaseTraverser) Expression(node ast.Expr) *ExpressionTraverser {
	return &ExpressionTraverser{
		Traverser: t,
		Node:      node,
		File:      t.ActiveFile(),
		ReturnNum: 0,
	}
}

func (e *ExpressionTraverser) SetReturnNum(returnNum int) *ExpressionTraverser {
	e.ReturnNum = returnNum
	return e
}

func (e *ExpressionTraverser) BuiltInFunctions(callExpr *ast.CallExpr) (types.Type, error) {
	if ident, ok := callExpr.Fun.(*ast.Ident); ok {
		switch ident.Name {
		case "new":
			fallthrough
		case "make":
			arg := callExpr.Args[0]
			return e.Traverser.Expression(arg).Type()
		case "len":
			return types.Typ[types.Int], nil
		}
	}

	return nil, nil
}

func (e *ExpressionTraverser) Value() (string, error) {
	switch n := e.Node.(type) {
	case *ast.BasicLit:
		value := n.Value
		if n.Kind == token.STRING {
			value = strings.Trim(value, "\"")
		}

		return value, nil
	case *ast.Ident:
		obj, err := e.File.Package.FindObjectForIdent(n)
		if err != nil {
			return "", err
		}

		if constant, ok := obj.(*types.Const); ok {
			return constant.Val().String(), nil
		}

		node, err := e.File.Package.ASTAtPos(obj.Pos())
		if err != nil {
			return "", err
		}

		declaration, err := e.Traverser.Declaration(node, n.Name)
		if err != nil {
			return "", err
		}

		astNode, err := declaration.Value()
		if err != nil {
			return "", err
		}

		exprNode, ok := astNode.(ast.Expr)
		if !ok {
			return "", errors.New("astNode is not of type ast.Expr")
		}

		return e.Traverser.Expression(exprNode).Value()
	case *ast.SelectorExpr:
		obj, err := e.File.Package.FindObjectForIdent(n.Sel)
		if err != nil {
			return "", err
		}

		pkg := e.Traverser.Packages.FindOrAdd(obj.Pkg().Path())
		_, err = e.Traverser.Packages.Get(pkg)
		if err != nil {
			return "", err
		}

		obj, err = pkg.FindObjectForIdentFuzzy(n.Sel)
		if err != nil {
			return "", err
		}

		if constant, ok := obj.(*types.Const); ok {
			return constant.Val().String(), nil
		}

		node, err := pkg.ASTAtPos(obj.Pos())
		if err != nil {
			return "", err
		}

		declaration, err := e.Traverser.Declaration(node, n.Sel.Name)
		if err != nil {
			return "", err
		}

		astNode, err := declaration.Value()
		if err != nil {
			return "", err
		}

		exprNode, ok := astNode.(ast.Expr)
		if !ok {
			return "", errors.New("astNode is not of type ast.Expr")
		}

		return e.Traverser.Expression(exprNode).Value()
	}

	return "", errors.New("value not retrievable")
}

func (e *ExpressionTraverser) Type() (types.Type, error) {
	switch n := e.Node.(type) {
	case *ast.StarExpr:
		return e.Traverser.Expression(n.X).Type()
	case *ast.UnaryExpr:
		return e.Traverser.Expression(n.X).Type()
	case *ast.Ident:
		obj, err := e.File.Package.FindObjectForIdent(n)
		if err != nil {
			return nil, err
		}
		return obj.Type(), nil
	case *ast.SelectorExpr:
		obj, err := e.File.Package.FindObjectForIdent(n.Sel)
		if err != nil {
			return nil, err
		}
		return obj.Type(), nil
	case *ast.ArrayType:
		return e.File.Package.FindTypeForExpr(n)
	case *ast.MapType:
		return e.File.Package.FindTypeForExpr(n)
	case *ast.CompositeLit:
		return e.Traverser.Expression(n.Type).Type()
	case *ast.BasicLit:
		return e.File.Package.FindTypeForExpr(n)
	case *ast.CallExpr:
		result, err := e.BuiltInFunctions(n)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, err
		}

		callExpr, err := e.Traverser.CallExpression(n)
		if err != nil {
			return nil, err
		}

		return callExpr.ReturnType(e.ReturnNum)
	}

	return nil, errors.New("unsupported expression type")
}
