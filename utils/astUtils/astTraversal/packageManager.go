package astTraversal

import (
	"go/ast"
	"golang.org/x/tools/go/packages"
	"strings"
	"sync"
)

type PackagePathLoader func(path string) (string, error)

type PackageManager struct {
	mu          sync.Mutex
	tree        *PackageNode
	workDir     string
	pathLoaders []PackagePathLoader
}

func NewPackageManager(workDir string) *PackageManager {
	return &PackageManager{
		workDir: workDir,
		tree: &PackageNode{
			Edges: make([]*PackageNode, 0),
		},
	}
}

func (pm *PackageManager) AddPathLoader(loader PackagePathLoader) {
	pm.pathLoaders = append(pm.pathLoaders, loader)
}

func (pm *PackageManager) AddPackage(pkgPath string) *PackageNode {
	pathItems := strings.Split(pkgPath, "/")

	pm.mu.Lock()
	defer pm.mu.Unlock()

	currNode := pm.tree
	for _, step := range pathItems {
		found := false
		for _, node := range currNode.Edges {
			if node.Name == step {
				found = true
				currNode = node
				break
			}
		}

		if !found {
			newNode := &PackageNode{
				Name:   step,
				Parent: currNode,
				Edges:  make([]*PackageNode, 0),
			}
			currNode.Edges = append(currNode.Edges, newNode)

			currNode = newNode
		}
	}

	return currNode
}

func (pm *PackageManager) Get(n *PackageNode) (*packages.Package, error) {
	if n.Package == nil {
		path := n.Path()
		for _, loader := range pm.pathLoaders {
			newPath, err := loader(path)
			if err == nil {
				path = newPath
				break
			}
		}
		pkg, err := LoadPackage(path, pm.workDir)
		if err != nil {
			return nil, err
		}

		n.Package = pkg

		for _, file := range n.Package.Syntax {
			n.Files = append(n.Files, &FileNode{
				Package:  n,
				FileName: n.Package.Fset.Position(file.Pos()).Filename,
				AST:      file,
				Imports:  pm.MapImportSpecs(file.Imports),
			})
		}
	}

	return n.Package, nil
}

func (pm *PackageManager) Find(path string) *PackageNode {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	currNode := pm.tree
	pathItems := strings.Split(path, "/")
	for _, step := range pathItems {
		found := false
		for _, node := range currNode.Edges {
			if node.Name == step {
				found = true
				currNode = node
				break
			}
		}

		if !found {
			return nil
		}
	}

	if currNode == pm.tree {
		return nil
	}

	return currNode
}

func (pm *PackageManager) FindOrAdd(path string) *PackageNode {
	pkg := pm.Find(path)
	if pkg == nil {
		pkg = pm.AddPackage(path)
	}

	return pkg
}

func (pm *PackageManager) MapImportSpecs(imports []*ast.ImportSpec) []FileImport {
	fileImports := make([]FileImport, 0)
	for _, imp := range imports {
		fileImport := FileImport{
			Package: pm.FindOrAdd(strings.Trim(imp.Path.Value, "\"")),
		}

		if imp.Name != nil {
			fileImport.Name = imp.Name.Name
		}

		fileImports = append(fileImports, fileImport)
	}

	return fileImports
}
