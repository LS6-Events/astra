package astUtils

import (
	"github.com/rs/zerolog"
	"go/ast"
	"strings"
)

// ParseAssignStatement parses an assignment statement
// It will extract the package name and type of the element on the right-hand side of the assignment
func ParseAssignStatement(log zerolog.Logger, expr ast.Expr, assignStmt *ast.AssignStmt, pkgPath string, pkgName string, imports []*ast.ImportSpec, argType *ast.Ident, onExternalPkg func(funcName, pkgName, pkgPath string) error) (ParseResult, error, bool) {
	var err error
	var res ParseResult
	switch rhs := expr.(type) {
	case *ast.UnaryExpr:
		return ParseAssignStatement(log, rhs.X, assignStmt, pkgPath, pkgName, imports, argType, onExternalPkg)
	case *ast.CompositeLit:
		res = HandleExpr(rhs.Type, pkgName)
	case *ast.BasicLit:
		res = ParseResult{
			VarName: strings.ToLower(rhs.Kind.String()),
			PkgName: pkgName,
			Value:   rhs.Value,
		}
	case *ast.Ident:
		assignStmt, ok := rhs.Obj.Decl.(*ast.AssignStmt)
		if !ok {
			return ParseResult{}, nil, false
		}

		var assignedIndex int
		for i, expr := range assignStmt.Lhs {
			if expr.(*ast.Ident).Name == rhs.Name {
				assignedIndex = i
				break
			}
		}

		var assignedExpr ast.Expr
		if len(assignStmt.Lhs) == len(assignStmt.Rhs) { // If the number of variables and values are the same
			assignedExpr = assignStmt.Rhs[assignedIndex]
		} else { // If the number of variables and values are different (i.e. a function call)
			assignedExpr = assignStmt.Rhs[0]
		}

		return ParseAssignStatement(log, assignedExpr, assignStmt, pkgPath, pkgName, imports, argType, onExternalPkg)
	case *ast.CallExpr:
		res, err = HandleReservedFunctions(rhs, pkgName)
		if err != nil {
			return ParseResult{}, err, false
		}
		if res.VarName != "" {
			return res, nil, true
		}

		switch fun := rhs.Fun.(type) {
		case *ast.SelectorExpr: // foo.Bar()
			ident, ok := fun.X.(*ast.Ident)
			if !ok {
				return ParseResult{}, nil, false
			}

			err = onExternalPkg(fun.Sel.Name, ident.Name, pkgPath)
			if err != nil {
				return ParseResult{}, err, false
			} else {
				return ParseResult{}, nil, false
			}
		case *ast.Ident: // Bar()
			if fun.Obj == nil { // If the function is in the current package but not declared in the same file
				err = onExternalPkg(fun.Name, pkgName, pkgPath)
				if err != nil {
					return ParseResult{}, err, false
				} else {
					return ParseResult{}, nil, false
				}
			} else {
				funcDecl, ok := fun.Obj.Decl.(*ast.FuncDecl)
				if !ok {
					return ParseResult{}, nil, false
				}

				var funcReturnIndex int
				for i, field := range assignStmt.Lhs {
					if f, ok := field.(*ast.Ident); ok {
						if f.Name == argType.Name {
							funcReturnIndex = i
						}
					}
				}

				field := funcDecl.Type.Results.List[funcReturnIndex]

				res, ok = ParseFunctionReturnTypes(log, field.Type, argType.Name)
				if !ok {
					return ParseResult{}, nil, false
				}
			}
		default:
			return ParseResult{}, nil, false
		}
	default:
		return ParseResult{}, nil, false
	}

	return res, nil, true
}
