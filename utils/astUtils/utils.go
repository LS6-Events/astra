package astUtils

import (
	"go/ast"
)

// FuncDeclToFuncLit converts a function declaration to a function literal
// Useful to convert named functions to anonymous functions, as we don't need the function names for now
func FuncDeclToFuncLit(expr *ast.FuncDecl) *ast.FuncLit {
	return &ast.FuncLit{
		Type: expr.Type,
		Body: expr.Body,
	}
}
