package astTraversal

import (
	"errors"
	"github.com/rs/zerolog"
	"go/ast"
)

type BaseTraverser struct {
	baseActiveFile *FileNode
	activeFile     *FileNode
	Log            *zerolog.Logger
	Packages       *PackageManager
}

func New(workDir string) *BaseTraverser {
	packages := NewPackageManager(workDir)
	return &BaseTraverser{
		Packages: packages,
	}
}

func (t *BaseTraverser) SetLog(log *zerolog.Logger) *BaseTraverser {
	t.Log = log
	return t
}

func (t *BaseTraverser) ActiveFile() *FileNode {
	return t.activeFile
}

func (t *BaseTraverser) SetActiveFile(file *FileNode) *BaseTraverser {
	t.baseActiveFile = file
	t.activeFile = file
	return t
}

func (t *BaseTraverser) Reset() {
	t.activeFile = t.baseActiveFile
}

func (t *BaseTraverser) ExtractVarName(node ast.Node) Result {
	packageNode := t.ActiveFile().Package

	switch nodeType := node.(type) {
	case *ast.Ident:
		return Result{
			Type:    nodeType.Name,
			Package: packageNode,
		}
	case *ast.SelectorExpr:
		isSelectorExpr := true
		names := []string{nodeType.Sel.Name}
		for isSelectorExpr {
			isSelectorExpr = false
			switch n := nodeType.X.(type) {
			case *ast.Ident:
				name := n.Name
				fileImport, ok := t.ActiveFile().FindImport(name) // We don't change the active file here
				if ok {
					packageNode = fileImport.Package
				} else {
					names = append([]string{name}, names...)
				}

				return Result{
					Package:     packageNode,
					Type:        names[len(names)-1],
					StructNames: names[:len(names)-1],
				}
			case *ast.SelectorExpr:
				isSelectorExpr = true
				names = append([]string{n.Sel.Name}, names...)
				nodeType = n
			default:
				return Result{}
			}
		}
	}

	return Result{}
}

func (t *BaseTraverser) FindDeclarationForNode(node ast.Node) (*DeclarationTraverser, error) {
	switch nodeType := node.(type) {
	case *ast.UnaryExpr:
		return t.FindDeclarationForNode(nodeType.X)
	case *ast.Ident:
		if nodeType.Obj != nil { // Defined in file
			return t.Declaration(nodeType.Obj.Decl.(ast.Node), nodeType.Name)
		} else { // Defined in package
			_, err := t.Packages.Get(t.ActiveFile().Package)
			if err != nil {
				return nil, err
			}

			newNode, file, err := t.ActiveFile().Package.FindASTForType(nodeType.Name)
			t.activeFile = file

			return t.Declaration(newNode, nodeType.Name)
		}
	case *ast.SelectorExpr:
		packageName := nodeType.X.(*ast.Ident).Name
		if fileImport, ok := t.ActiveFile().FindImport(packageName); ok {
			_, err := t.Packages.Get(fileImport.Package)
			if err != nil {
				return nil, err
			}

			newNode, file, err := fileImport.Package.FindASTForType(nodeType.Sel.Name)
			t.activeFile = file

			return t.Declaration(newNode, nodeType.Sel.Name)
		} else { // Property of a struct - need to recursively
			// TODO: Implement struct property lookup and recursive traversal
			return nil, errors.New("not implemented yet")
		}
	}

	return nil, errors.New("unsupported type")
}
