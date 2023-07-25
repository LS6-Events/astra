package astUtils

import (
	"go/ast"
	"strings"
)

// ParseInputPath parses the input path from the imports, package name and package path
// It returns the package path if the package name is not found in the imports
// It returns the package name if the package name is not found in the imports and the package path is empty
// If the package isn't found in the imports, we assume that the package name is the same as the last part of the package path
func ParseInputPath(imports []*ast.ImportSpec, pkgName, pkgPath string) string {
	var p string
	if pkgName != "" {
		for _, imp := range imports {
			if imp.Name != nil && imp.Name.Name == pkgName {
				p = strings.ReplaceAll(imp.Path.Value, "\"", "")
				break
			} else if imp.Name == nil {
				split := strings.Split(strings.ReplaceAll(imp.Path.Value, "\"", ""), "/")
				if len(split) > 0 && split[len(split)-1] == pkgName {
					p = strings.ReplaceAll(imp.Path.Value, "\"", "")
					break
				}
			}
		}
	}

	if p == "" {
		if pkgPath == "" {
			p = pkgName
		} else {
			p = pkgPath
		}
	}

	return p
}
