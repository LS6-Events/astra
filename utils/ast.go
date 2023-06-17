package utils

import (
	"go/ast"
)

type ParseResult struct {
	VarName string
	PkgName string

	Value string

	MapKeyPkg string
	MapKey    string
	MapVal    string

	SliceType string
}

func SplitIdentSelectorExpr(expr ast.Expr, defaultPkgName string) ParseResult {
	switch e := expr.(type) {
	case *ast.Ident:
		return ParseResult{
			VarName: e.Name,
			PkgName: defaultPkgName,
		}
	case *ast.SelectorExpr:
		return ParseResult{
			PkgName: e.X.(*ast.Ident).Name,
			VarName: e.Sel.Name,
		}
	}

	return ParseResult{}
}

func FuncDeclToFuncLit(expr *ast.FuncDecl) *ast.FuncLit {
	return &ast.FuncLit{
		Type: expr.Type,
		Body: expr.Body,
	}
}
