package astUtils

import (
	"go/ast"
)

// HandleExpr handles an expression
// It can be an identifier, a selector expression, an array type or a map type
// It will return the package name and the variable name as part of the ParseResult
func HandleExpr(expr ast.Expr, pkgName string) ParseResult {
	switch n := expr.(type) {
	case *ast.StarExpr:
		return HandleExpr(n.X, pkgName)
	case *ast.UnaryExpr:
		return HandleExpr(n.X, pkgName)
	case *ast.Ident:
		return SplitIdentSelectorExpr(n, pkgName)
	case *ast.SelectorExpr:
		return SplitIdentSelectorExpr(n, pkgName)
	case *ast.ArrayType:
		embeddedType := SplitIdentSelectorExpr(n.Elt, pkgName)
		return ParseResult{
			VarName:   "slice",
			PkgName:   embeddedType.PkgName,
			SliceType: embeddedType.VarName,
		}
	case *ast.MapType:
		keyType := SplitIdentSelectorExpr(n.Key, pkgName)
		valueType := SplitIdentSelectorExpr(n.Value, pkgName)
		return ParseResult{
			VarName:   "map",
			MapKey:    keyType.VarName,
			MapKeyPkg: keyType.PkgName,
			MapVal:    valueType.VarName,
			PkgName:   valueType.PkgName,
		}
	}

	return ParseResult{}
}
