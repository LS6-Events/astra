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
