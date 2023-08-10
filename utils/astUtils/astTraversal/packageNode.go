package astTraversal

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/packages"
)

type PackageNode struct {
	Parent  *PackageNode
	Name    string
	Package *packages.Package
	Edges   []*PackageNode
	Files   []*FileNode
}

func (p *PackageNode) Path() string {
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
