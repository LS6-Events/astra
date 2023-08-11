package astTraversal

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestNewPackageManager(t *testing.T) {
	workDir := "/path/to/dir"
	pm := NewPackageManager(workDir)

	assert.NotNil(t, pm.tree)
	assert.Equal(t, workDir, pm.workDir)
}

func TestPackageManager_AddPathLoader(t *testing.T) {
	pm := NewPackageManager("/path/to/dir")
	loader := func(path string) (string, error) {
		return "/new/path", nil
	}

	pm.AddPathLoader(loader)

	for _, pathLoader := range pm.pathLoaders {
		startPath := "/path/to/dir"
		path, err := pathLoader(startPath)
		assert.NoError(t, err)

		assert.Equal(t, "/new/path", path)
	}
}

func TestPackageManager_AddPackage(t *testing.T) {
	pm := NewPackageManager("/path/to/dir")
	pkgPath := "github.com/user/project"
	node := pm.AddPackage(pkgPath)

	assert.NotNil(t, node)
	assert.Equal(t, "project", node.Name)
	assert.Equal(t, "user", node.Parent.Name)
	assert.Equal(t, "github.com", node.Parent.Parent.Name)
	assert.Equal(t, pkgPath, node.Path())
}

func TestPackageManager_Get(t *testing.T) {
	pm := NewPackageManager(".")
	pkgPath := "time"
	node := pm.AddPackage(pkgPath)

	assert.Nil(t, node.Package)

	pkg, err := pm.Get(node)
	assert.Nil(t, err)

	assert.NotNil(t, pkg)
	assert.Equal(t, pkgPath, pkg.PkgPath)
}

func TestPackageManager_Find(t *testing.T) {
	pm := NewPackageManager(".")
	pkgPath := "github.com/user/project"
	node := pm.AddPackage(pkgPath)

	assert.Nil(t, node.Package)

	pkg := pm.Find(pkgPath)
	assert.NotNil(t, pkg)
	assert.Equal(t, pkgPath, pkg.Path())

	pkg = pm.Find("github.com/user/other")
	assert.Nil(t, pkg)
}

func TestPackageManager_FindOrAdd(t *testing.T) {
	pm := NewPackageManager(".")
	pkgPath := "github.com/project"
	node := pm.AddPackage(pkgPath)

	assert.Nil(t, node.Package)

	pkg := pm.FindOrAdd(pkgPath)
	assert.NotNil(t, pkg)
	assert.Equal(t, pkgPath, pkg.Path())
	assert.Len(t, pm.tree.Edges, 1)

	pkg = pm.FindOrAdd(pkgPath)
	assert.NotNil(t, pkg)
	assert.Equal(t, pkgPath, pkg.Path())
	assert.Len(t, pm.tree.Edges, 1)
}

func TestPackageManager_MapImportSpecs(t *testing.T) {
	pm := NewPackageManager(".")
	imports := []*ast.ImportSpec{
		{
			Path: &ast.BasicLit{Value: "\"github.com/user/project\""},
		},
	}
	fileImports := pm.MapImportSpecs(imports)

	assert.Equal(t, 1, len(fileImports))
	assert.Equal(t, "github.com/user/project", fileImports[0].Package.Path())
}
