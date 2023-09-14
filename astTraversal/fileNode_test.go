package astTraversal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileNode_IsImportedPackage(t *testing.T) {
	file := &FileNode{
		Imports: []FileImport{
			{Package: &PackageNode{Name: "fmt"}},
			{Name: "mytime", Package: &PackageNode{Name: "time"}},
		},
	}

	assert.True(t, file.IsImportedPackage("fmt"))
	assert.True(t, file.IsImportedPackage("mytime"))
	assert.False(t, file.IsImportedPackage("math"))
}

func TestFileNode_FindImport(t *testing.T) {
	pkg := &PackageNode{Name: "fmt"}
	file := &FileNode{
		Imports: []FileImport{
			{Package: pkg, Name: "fmtAlias"},
		},
	}

	im, found := file.FindImport("fmtAlias")
	assert.True(t, found)
	assert.Equal(t, "fmtAlias", im.Name)

	_, found = file.FindImport("nonexistent")
	assert.False(t, found)
}
