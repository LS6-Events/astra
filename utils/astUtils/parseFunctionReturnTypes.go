package astUtils

import (
	"github.com/rs/zerolog"
	"go/ast"
)

// ParseFunctionReturnTypes parses the return types of a function
// It will extract the package name and type of the element on the return types
// And map them appropriately for the different types
func ParseFunctionReturnTypes(log zerolog.Logger, node ast.Node, argType *ast.Ident) (ParseResult, bool) {
	switch fieldType := node.(type) {
	case *ast.StarExpr:
		return ParseFunctionReturnTypes(log, fieldType.X, argType)
	case *ast.ArrayType:
		arrayResult, ok := ParseFunctionReturnTypes(log, fieldType.Elt, argType)
		if ok {
			return ParseResult{
				VarName:   "slice",
				PkgName:   arrayResult.PkgName,
				SliceType: arrayResult.VarName,
			}, true
		} else {
			return ParseResult{}, false
		}
	case *ast.MapType:
		keyResult, keyOk := ParseFunctionReturnTypes(log, fieldType.Key, argType)
		valResult, valOk := ParseFunctionReturnTypes(log, fieldType.Value, argType)
		if keyOk && valOk {
			return ParseResult{
				VarName:   "map",
				MapKey:    keyResult.VarName,
				MapKeyPkg: keyResult.PkgName,
				MapVal:    valResult.VarName,
				PkgName:   valResult.PkgName,
			}, true
		} else {
			return ParseResult{}, false
		}
	case *ast.SelectorExpr:
		return SplitIdentSelectorExpr(fieldType, argType.Name), true
	case *ast.Ident:
		return SplitIdentSelectorExpr(fieldType, argType.Name), true
	default:
		return ParseResult{}, false
	}
}
