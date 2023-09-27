package astTraversal

import (
	"go/ast"
	"go/doc"
	"go/types"
)

type FunctionTraverser struct {
	Traverser *BaseTraverser
	Node      *ast.FuncLit
	File      *FileNode
	DeclNode  *ast.FuncDecl
}

func (t *BaseTraverser) Function(node ast.Node) (*FunctionTraverser, error) {
	var funcLit *ast.FuncLit
	var funcDecl *ast.FuncDecl
	switch n := node.(type) {
	case *ast.FuncLit:
		funcLit = n
	case *ast.FuncDecl:
		funcLit = &ast.FuncLit{
			Type: n.Type,
			Body: n.Body,
		}
		funcDecl = n
	default:
		return nil, ErrInvalidNodeType
	}

	return &FunctionTraverser{
		Traverser: t,
		Node:      funcLit,
		File:      t.ActiveFile(),
		DeclNode:  funcDecl,
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
	for _, arg := range f.Arguments() {
		argType, err := f.File.Package.FindTypeForExpr(arg.Type)
		if err != nil {
			continue
		}

		if isPointer {
			if named, ok := argType.(*types.Pointer); ok {
				argType = named.Elem()
			} else {
				continue
			}
		}

		if named, ok := argType.(*types.Named); ok && named.Obj().Pkg().Path() == packagePath && named.Obj().Name() == typeName {
			return arg.Names[0].Name
		} else if basic, ok := argType.(*types.Basic); ok && basic.Name() == typeName {
			return arg.Names[0].Name
		}
	}

	return ""
}

func (f *FunctionTraverser) GoDoc() (*doc.Func, error) {
	pkgDoc, err := f.Traverser.Packages.GoDoc(f.File.Package)
	if err != nil {
		return nil, err
	}

	for _, fun := range pkgDoc.Funcs {
		if f.DeclNode != nil && fun.Name == f.DeclNode.Name.Name {
			return fun, nil
		}
	}

	return nil, nil
}
