package astUtils

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"
)

// ExtractContext extracts the name of the context parameter from a function literal
// It returns an error if the context parameter is not found
func ExtractContext(pkgPath, typeName string, node *ast.FuncLit, imports []*ast.ImportSpec) (string, error) {
	var pkgName string
	for _, im := range imports {
		if im.Path.Value == fmt.Sprintf("\"%s\"", pkgPath) {
			if im.Name == nil {
				pkgName = strings.Split(pkgPath, "/")[len(strings.Split(pkgPath, "/"))-1]
			} else {
				pkgName = im.Name.Name
			}
			break
		}
	}
	if pkgName == "" {
		return "", fmt.Errorf("package %s not imported", pkgPath)
	}

	var ctxName string
	for _, param := range node.Type.Params.List {
		if len(param.Names) == 0 {
			continue
		}

		if strings.HasPrefix(typeName, "*") {
			starExpr, ok := param.Type.(*ast.StarExpr)
			if !ok {
				continue
			}

			selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
			if !ok {
				continue
			}

			ident, ok := selectorExpr.X.(*ast.Ident)
			if !ok || ident.Name != pkgName {
				continue
			}

			if selectorExpr.Sel.Name == typeName[1:] {
				ctxName = param.Names[0].Name
				break
			}
		} else {
			ident, ok := param.Type.(*ast.Ident)
			if !ok {
				continue
			}

			if ident.Name == typeName {
				ctxName = param.Names[0].Name
				break
			}
		}

	}

	if ctxName == "" {
		return "", errors.New("context parameter not found")
	}

	return ctxName, nil
}
