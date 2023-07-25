package astUtils

import (
	"go/ast"
)

// HandleReservedFunctions handles reserved functions
// It will return no error if the function is not reserved
// These functions are:
// - new
// - make
// - len
func HandleReservedFunctions(callExpr *ast.CallExpr, pkgName string) (ParseResult, error) {
	if ident, ok := callExpr.Fun.(*ast.Ident); ok {
		switch ident.Name {
		case "new":
			fallthrough
		case "make":
			arg := callExpr.Args[0]
			return HandleExpr(arg, pkgName), nil
		case "len":
			return ParseResult{
				VarName: "int",
			}, nil
		}
	}

	return ParseResult{}, nil
}
