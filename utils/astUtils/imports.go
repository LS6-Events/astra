package astUtils

import (
	"go/ast"
	"strings"
)

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
