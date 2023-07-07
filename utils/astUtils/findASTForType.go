package astUtils

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/packages"
)

// FindASTForType finds the AST node for a type in a package
// It will return an error if the type is not found
// It will return the first type found if there are multiple types with the same name
func FindASTForType(pkg *packages.Package, typeName string) (ast.Node, []*ast.ImportSpec, error) {
	t := pkg.Types.Scope().Lookup(typeName)
	if t == nil {
		return nil, nil, fmt.Errorf("type %s not found in package %s", typeName, pkg.PkgPath)
	}

	pos := pkg.Fset.Position(t.Pos())

	for _, f := range pkg.Syntax {
		if pos.Filename == pkg.Fset.Position(f.Pos()).Filename {
			nT := f.Scope.Lookup(typeName)
			if nT == nil {
				return nil, nil, fmt.Errorf("type %s not found in package %s", typeName, pkg.PkgPath)
			}

			result, ok := nT.Decl.(ast.Node)
			if !ok {
				return nil, nil, fmt.Errorf("type %s not found in package %s", typeName, pkg.PkgPath)
			}

			return result, f.Imports, nil
		}
	}

	return nil, nil, fmt.Errorf("type %s not found in package %s", typeName, pkg.PkgPath)
}
