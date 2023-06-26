package astUtils

import (
	"fmt"
	"golang.org/x/tools/go/packages"
)

var cachedPackages = make(map[string]*packages.Package)

func LoadPackage(pkgPath string, workDir string) (*packages.Package, error) {
	if pkg, ok := cachedPackages[pkgPath]; ok {
		return pkg, nil
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedImports | packages.NeedDeps | packages.NeedName,
		Dir:  workDir,
	}, pkgPath)
	if err != nil {
		return nil, err
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("package %s not found", pkgPath)
	}

	cachedPackages[pkgPath] = pkgs[0]

	return pkgs[0], nil
}
