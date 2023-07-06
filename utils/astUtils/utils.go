package astUtils

import "go/ast"

// SplitIdentSelectorExpr splits an expression into a package name and a variable name
// The package name is the default package name if the expression is an identifier
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

// FuncDeclToFuncLit converts a function declaration to a function literal
// Useful to convert named functions to anonymous functions, as we don't need the function names for now
func FuncDeclToFuncLit(expr *ast.FuncDecl) *ast.FuncLit {
	return &ast.FuncLit{
		Type: expr.Type,
		Body: expr.Body,
	}
}
