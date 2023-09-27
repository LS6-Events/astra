package astTraversal

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
)

type PackageNode struct {
	Parent  *PackageNode
	Name    string
	Package *packages.Package
	Edges   []*PackageNode
	Files   []*FileNode
	Doc     *doc.Package
}

func (p *PackageNode) Path() string {
	if p == nil {
		return ""
	}

	var path string
	current := p
	for current.Parent != nil {
		if path == "" {
			path = current.Name
		} else {
			path = fmt.Sprintf("%s/%s", current.Name, path)
		}

		current = current.Parent
	}

	return path
}

func (p *PackageNode) AddFile(file *FileNode) {
	for _, f := range p.Files {
		if f.FileName == file.FileName {
			f.AST = file.AST
			return
		}
	}

	p.Files = append(p.Files, file)
}

func (p *PackageNode) FindASTForType(typeName string) (ast.Node, *FileNode, error) {
	if p.Package == nil {
		return nil, nil, fmt.Errorf("package %s not populated", p.Path())
	}
	t := p.Package.Types.Scope().Lookup(typeName)
	if t == nil {
		return nil, nil, fmt.Errorf("type %s not found in package %s", typeName, p.Path())
	}

	pos := p.Package.Fset.Position(t.Pos())

	for _, f := range p.Package.Syntax {
		if pos.Filename == p.Package.Fset.Position(f.Pos()).Filename {
			nT := f.Scope.Lookup(typeName)
			if nT == nil {
				return nil, nil, fmt.Errorf("type %s not found in package %s", typeName, p.Path())
			}

			result, ok := nT.Decl.(ast.Node)
			if !ok {
				return nil, nil, fmt.Errorf("type %s not found in package %s", typeName, p.Path())
			}

			var fileNode *FileNode
			for _, fN := range p.Files {
				if fN.FileName == pos.Filename {
					fileNode = fN
				}
			}

			return result, fileNode, nil
		}
	}

	return nil, nil, fmt.Errorf("type %s not found in package %s", typeName, p.Path())
}

func (p *PackageNode) FindTypeForExpr(expr ast.Expr) (types.Type, error) {
	if p.Package == nil {
		return nil, fmt.Errorf("package %s not populated", p.Path())
	}

	t := p.Package.TypesInfo.TypeOf(expr)
	if t == nil {
		return nil, fmt.Errorf("type for %s not found in package %s", expr, p.Path())
	}

	return t, nil
}

func (p *PackageNode) FindObjectForName(name string) (types.Object, error) {
	if p.Package == nil {
		return nil, fmt.Errorf("package %s not populated", p.Path())
	}

	obj := p.Package.Types.Scope().Lookup(name)
	if obj == nil {
		return nil, fmt.Errorf("object %s not found in package %s", name, p.Path())
	}

	return obj, nil
}

func (p *PackageNode) FindObjectForIdentFuzzy(ident *ast.Ident) (types.Object, error) {
	if p.Package == nil {
		return nil, fmt.Errorf("package %s not populated", p.Path())
	}

	obj := p.Package.TypesInfo.ObjectOf(ident)
	if obj == nil {
		// Try to find the object in the package
		for k, v := range p.Package.TypesInfo.Defs {
			if k.Name == ident.Name {
				obj = v
				break
			}
		}

		if obj == nil {
			return nil, fmt.Errorf("object %s not found in package %s", ident.Name, p.Path())
		}
	}

	return obj, nil
}

func (p *PackageNode) FindObjectForIdent(ident *ast.Ident) (types.Object, error) {
	if p.Package == nil {
		return nil, fmt.Errorf("package %s not populated", p.Path())
	}

	obj := p.Package.TypesInfo.ObjectOf(ident)
	if obj == nil {
		return nil, fmt.Errorf("object %s not found in package %s", ident.Name, p.Path())
	}

	return obj, nil
}

func (p *PackageNode) FindUsesForIdent(ident *ast.Ident) (types.Object, error) {
	if p.Package == nil {
		return nil, fmt.Errorf("package %s not populated", p.Path())
	}

	uses := p.Package.TypesInfo.Uses[ident]
	if uses == nil {
		return nil, fmt.Errorf("uses for %s not found in package %s", ident.Name, p.Path())
	}

	return uses, nil
}

func (p *PackageNode) ASTAtPos(pos token.Pos) (ast.Node, error) {
	if p.Package == nil {
		return nil, fmt.Errorf("package %s not populated", p.Path())
	}

	node := p.Package.Fset.Position(pos)
	for _, f := range p.Package.Syntax {
		if node.Filename == p.Package.Fset.Position(f.Pos()).Filename {
			// Find the right node
			var result ast.Node
			ast.Inspect(f, func(n ast.Node) bool {
				if n == nil {
					return true
				}

				if n.Pos() == pos {
					result = n
					return false
				}

				return true
			})

			if result == nil {
				return nil, fmt.Errorf("node at %s not found in package %s", node, p.Path())
			}

			return result, nil
		}
	}

	return nil, fmt.Errorf("node at %s not found in package %s", node, p.Path())
}
