package astTraversal

import "testing"
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
