package astTraversal

import "go/ast"

type FileNode struct {
	Package  *PackageNode
	FileName string
	Imports  []FileImport
	AST      *ast.File
}

type FileImport struct {
	Package *PackageNode
	Name    string
}

func (f *FileNode) IsImportedPackage(ident string) bool {
	_, ok := f.FindImport(ident)
	return ok
}

func (f *FileNode) FindImport(ident string) (FileImport, bool) {
	for _, im := range f.Imports {
		if im.Name == ident || im.Package.Name == ident {
			return im, true
		}
	}

	return FileImport{}, false
}

func (t *BaseTraverser) ASTFileToNode(file *ast.File, packageNode *PackageNode) (*FileNode, error) {
	pkg, err := t.Packages.Get(packageNode)
	if err != nil {
		return nil, err
	}
	pos := pkg.Fset.Position(file.Pos())

	imports := t.Packages.MapImportSpecs(file.Imports)

	node := &FileNode{
		Package:  packageNode,
		FileName: pos.Filename,
		Imports:  imports,
		AST:      file,
	}

	packageNode.AddFile(node)

	return node, nil
}
