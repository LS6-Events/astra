package astTraversal

import (
	"errors"
	"github.com/ls6-events/gengo/utils/astUtils"
	"go/ast"
)

type FunctionTraverser struct {
	Traverser *Traverser
	Node      *ast.FuncLit
}

func (t *Traverser) Function(node ast.Node) (*FunctionTraverser, error) {
	var funcLit *ast.FuncLit
	switch n := node.(type) {
	case *ast.FuncLit:
		funcLit = n
	case *ast.FuncDecl:
		funcLit = astUtils.FuncDeclToFuncLit(n)
	default:
		return nil, ErrInvalidNodeType
	}

	return &FunctionTraverser{
		Traverser: t,
		Node:      funcLit,
	}, nil
}

func (f *FunctionTraverser) Arguments() []*ast.Field {
	if f.Node.Type.Params == nil {
		return []*ast.Field{}
	}
	return f.Node.Type.Params.List
}

func (f *FunctionTraverser) Results() []*ast.Field {
	if f.Node.Type.Results == nil {
		return []*ast.Field{}
	}
	return f.Node.Type.Results.List
}

func (f *FunctionTraverser) FindArgumentNameByType(typeName string, packagePath string, isPointer bool) string {
	var packageIdentifier string
	if packagePath != "" {
		for _, im := range f.Traverser.ActiveFile().Imports {
			if im.Package.Path() == packagePath {
				if im.Name == "" {
					packageIdentifier = im.Package.Name
				} else {
					packageIdentifier = im.Name
				}
				break
			}
		}
	}

	var varName string
	for _, param := range f.Arguments() {
		if len(param.Names) == 0 {
			continue
		}

		if isPointer {
			starExpr, ok := param.Type.(*ast.StarExpr)
			if !ok {
				continue
			}

			if packageIdentifier != "" {
				selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
				if !ok {
					continue
				}

				ident, ok := selectorExpr.X.(*ast.Ident)
				if !ok || ident.Name != packageIdentifier {
					continue
				}

				if selectorExpr.Sel.Name == typeName {
					return param.Names[0].Name
				}
			} else {
				ident, ok := starExpr.X.(*ast.Ident)
				if !ok {
					continue
				}

				if ident.Name == typeName {
					return param.Names[0].Name
				}
			}

		} else {
			if packageIdentifier != "" {
				selectorExpr, ok := param.Type.(*ast.SelectorExpr)
				if !ok {
					continue
				}

				ident, ok := selectorExpr.X.(*ast.Ident)
				if !ok || ident.Name != packageIdentifier {
					continue
				}

				if selectorExpr.Sel.Name == typeName {
					return param.Names[0].Name
				}
			} else {
				ident, ok := param.Type.(*ast.Ident)
				if !ok {
					continue
				}

				if ident.Name == typeName {
					return param.Names[0].Name
				}
			}
		}
	}

	return varName
}

func (f *FunctionTraverser) ReturnTypeResult(varNum int) (Result, error) {
	results := f.Results()

	if len(results) < varNum+1 {
		return Result{}, errors.New("too few return variables")
	}

	result, err := f.Traverser.Expression(results[varNum].Type).Result()
	if err != nil {
		return Result{}, err
	}

	return result, nil
}
