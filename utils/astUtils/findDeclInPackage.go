package astUtils

import "go/ast"

// FindDeclInPackage finds the declaration for a type in a package
// It will return an error if the type is not found
// It will return the first type found if there are multiple types with the same name
func FindDeclInPackage(expr ast.Expr, imports []*ast.ImportSpec, pkgName, pkgPath, workDir string, onPkgName func(ParseResult) (ParseResult, error)) (ParseResult, ast.Node, []*ast.ImportSpec, error) {
	res := SplitIdentSelectorExpr(expr, pkgName)

	res, err := onPkgName(res)
	if err != nil {
		return ParseResult{}, nil, nil, err
	}

	nPkgPath := ParseInputPath(imports, res.PkgName, pkgPath)
	pkg, err := LoadPackage(nPkgPath, workDir)
	if err != nil {
		return ParseResult{}, nil, nil, err
	}

	node, nImports, err := FindASTForType(pkg, res.VarName)
	if err != nil {
		return ParseResult{}, nil, nil, err
	}

	return res, node, nImports, nil
}
