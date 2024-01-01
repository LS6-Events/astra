package astTraversal

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestNew(t *testing.T) {
	workDir := "/path/to/dir"
	traverser := New(workDir)
	assert.NotNil(t, traverser.Packages)
}

func TestTraverser_SetLog(t *testing.T) {
	traverser := &BaseTraverser{}
	logger := zerolog.Nop()
	traverser.SetLog(&logger)

	assert.NotNil(t, traverser.Log)
}

func TestTraverser_SetActiveFile(t *testing.T) {
	traverser := &BaseTraverser{}
	file := &FileNode{}
	traverser.SetActiveFile(file)

	assert.Equal(t, file, traverser.activeFile)
	assert.Equal(t, file, traverser.baseActiveFile)
}

func TestTraverser_Reset(t *testing.T) {
	traverser := &BaseTraverser{}
	file1 := &FileNode{}
	file2 := &FileNode{}

	traverser.SetActiveFile(file1).activeFile = file2
	traverser.Reset()

	assert.Equal(t, file1, traverser.activeFile)
}

func TestTraverser_ExtractVarName(t *testing.T) {
	traverser, err := CreateTraverserFromTestFile("usefulTypes.go")
	assert.NoError(t, err)

	assert.NotNil(t, traverser.ActiveFile().Package)

	t.Run("should return the name of an Ident", func(t *testing.T) {
		ident := &ast.Ident{
			Name: "MyStruct",
		}
		result := traverser.ExtractVarName(ident)
		assert.Equal(t, "MyStruct", result.Type)
		assert.Contains(t, result.Package.Path(), "testfiles")
	})

	t.Run("should return the name of a SelectorExpr", func(t *testing.T) {
		selectorExpr := &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "otherpkg1",
			},
			Sel: &ast.Ident{
				Name: "MyField",
			},
		}
		result := traverser.ExtractVarName(selectorExpr)
		assert.Equal(t, "MyField", result.Type)
		assert.Contains(t, result.Package.Path(), "testfiles")
	})

	t.Run("should return the name of a SelectorExpr with multiple X", func(t *testing.T) {
		// MyStruct.MyField.MyField2
		selectorExpr := &ast.SelectorExpr{
			X: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "MyStruct",
				},
				Sel: &ast.Ident{
					Name: "MyField",
				},
			},
			Sel: &ast.Ident{
				Name: "MyField2",
			},
		}
		result := traverser.ExtractVarName(selectorExpr)
		assert.Equal(t, "MyField2", result.Type)
		assert.Equal(t, []string{"MyStruct", "MyField"}, result.Names)
	})

	t.Run("should return the name of a SelectorExpr with multiplx X and a package", func(t *testing.T) {
		// otherpkg1.MyStruct.MyField.MyField2
		selectorExpr := &ast.SelectorExpr{
			X: &ast.SelectorExpr{
				X: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "otherpkg1",
					},
					Sel: &ast.Ident{
						Name: "MyStruct",
					},
				},
				Sel: &ast.Ident{
					Name: "MyField",
				},
			},
			Sel: &ast.Ident{
				Name: "MyField2",
			},
		}
		result := traverser.ExtractVarName(selectorExpr)
		assert.Equal(t, "MyField2", result.Type)
		assert.Equal(t, []string{"MyStruct", "MyField"}, result.Names)
		assert.Contains(t, result.Package.Path(), "testfiles") // We don't change the active file here
	})
}
