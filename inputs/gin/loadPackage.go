package gin

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"os"
)

var cachedPackages = make(map[string]*packages.Package)

func loadPackage(pkgPath string) (*packages.Package, error) {
	if pkg, ok := cachedPackages[pkgPath]; ok {
		return pkg, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedImports | packages.NeedDeps | packages.NeedName,
		Dir:  cwd,
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
