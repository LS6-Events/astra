package astTraversal

import (
	"fmt"
	"golang.org/x/tools/go/packages"
)

var cachedPackages = make(map[string]*packages.Package)

// LoadPackage loads a package from a path
// Because of the way the packages.Load function works, we cache the packages to avoid loading the same package multiple times
// As we load these packages one at a time
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

	for _, pkg := range pkgs {
		for _, pkgErr := range pkg.Errors {
			switch pkgErr.Kind {
			case packages.ListError:
				return nil, fmt.Errorf("package %s has list errors", pkgPath)
			case packages.TypeError:
				return nil, fmt.Errorf("package %s has type errors", pkgPath)
			case packages.ParseError:
				return nil, fmt.Errorf("package %s has parse errors", pkgPath)
			case packages.UnknownError:
				return nil, fmt.Errorf("package %s has unknown errors", pkgPath)
			}
		}
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("package %s not found", pkgPath)
	}

	cachedPackages[pkgPath] = pkgs[0]

	return pkgs[0], nil
}
