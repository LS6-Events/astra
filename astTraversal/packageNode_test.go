package astTraversal

import (
	"go/ast"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestPackageNode_Path(t *testing.T) {
	root := &PackageNode{}

	pkg1 := &PackageNode{
		Name: "pkg1",
	}

	pkg2 := &PackageNode{
		Name: "pkg2",
	}

	pkg3 := &PackageNode{
		Name: "pkg3",
	}

	pkg4 := &PackageNode{
		Name: "pkg4",
	}

	pkg5 := &PackageNode{
		Name: "pkg5",
	}

	pkg6 := &PackageNode{
		Name: "pkg6",
	}

	pkg1.Parent = root
	pkg2.Parent = pkg1
	pkg1.Edges = append(pkg1.Edges, pkg2)
	pkg3.Parent = pkg2
	pkg2.Edges = append(pkg2.Edges, pkg3)

	pkg4.Parent = root
	pkg4.Edges = append(pkg4.Edges, pkg5)
	pkg5.Parent = pkg4

	pkg6.Parent = root

	root.Edges = append(root.Edges, pkg1, pkg4, pkg6)

	t.Run("Test a single root package with children", func(t *testing.T) {
		assert.Equal(t, "pkg1", pkg1.Path())
	})

	t.Run("Test a child of a node with a child", func(t *testing.T) {
		assert.Equal(t, "pkg1/pkg2", pkg2.Path())
	})

	t.Run("Test a child of a 2-layered node with no children", func(t *testing.T) {
		assert.Equal(t, "pkg1/pkg2/pkg3", pkg3.Path())
	})

	t.Run("Test a child of root with a single edge", func(t *testing.T) {
		assert.Equal(t, "pkg4", pkg4.Path())
	})

	t.Run("Test a child of a single-layered node with no children", func(t *testing.T) {
		assert.Equal(t, "pkg4/pkg5", pkg5.Path())
	})

	t.Run("Test a child of a root with no children", func(t *testing.T) {
		assert.Equal(t, "pkg6", pkg6.Path())
	})
}

func TestPackageNode_AddFile(t *testing.T) {
	p := &PackageNode{Name: "root"}

	file1 := &FileNode{FileName: "main.go"}
	file2 := &FileNode{FileName: "util.go"}

	// Adding a file to an empty PackageNode
	p.AddFile(file1)
	assert.Len(t, p.Files, 1)
	assert.Equal(t, "main.go", p.Files[0].FileName)

	// Adding a different file to the PackageNode
	p.AddFile(file2)
	assert.Len(t, p.Files, 2)

	// Adding the same file should not increase the length but should update the AST
	newAST := &ast.File{Name: &ast.Ident{Name: "NewAST"}}
	file1.AST = newAST
	p.AddFile(file1)
	assert.Len(t, p.Files, 2)
	assert.Equal(t, newAST, p.Files[0].AST)
}

func TestPackageNode_FindASTForType(t *testing.T) {
	traverser, err := createTraverserFromTestFile("usefulTypes.go")
	assert.NoError(t, err)

	assert.NotNil(t, traverser.ActiveFile().Package)

	_, err = traverser.Packages.Get(traverser.ActiveFile().Package)
	assert.NoError(t, err)

	t.Run("Test finding a type in the same file", func(t *testing.T) {
		// Find the AST for the type "MyStruct"
		node, file, err := traverser.ActiveFile().Package.FindASTForType("MyStruct")
		assert.NoError(t, err)
		assert.NotNil(t, node)
		assert.NotNil(t, file)
		assert.Contains(t, file.FileName, "usefulTypes.go")
	})

	t.Run("Test finding a type in another file in the same package", func(t *testing.T) {
		// Find the AST for the type "YourStruct"
		node, file, err := traverser.ActiveFile().Package.FindASTForType("YourStruct")
		assert.NoError(t, err)
		assert.NotNil(t, node)
		assert.NotNil(t, file)
		assert.Contains(t, file.FileName, "usefulTypes2.go")
	})

	t.Run("Test finding a type in another package", func(t *testing.T) {
		// Find the AST for the type "otherpkg1.Foo"
		_, _, err := traverser.ActiveFile().Package.FindASTForType("otherpkg1.Foo")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "type otherpkg1.Foo not found in package github.com/ls6-events/astra/astTraversal/testfiles")
	})
}
