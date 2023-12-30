package astTraversal

import (
	"go/parser"
	"go/token"
	"os"
	"path"
)

const relativeTestFilePath = "astTraversal/testfiles"

func createTraverserFromTestFile(testFilePath string) (*BaseTraverser, error) {
	wd, err := os.Getwd() // Will work as tests are only run from the root of the project
	if err != nil {
		return nil, err
	}

	for path.Base(wd) != "astra" {
		wd = path.Dir(wd)
	}

	fullPath := path.Join(wd, relativeTestFilePath, testFilePath)

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, fullPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	traverser := New(wd)

	packagePath := "github.com/ls6-events/astra/" + relativeTestFilePath

	fileNode := &FileNode{
		AST:      file,
		FileName: path.Base(fullPath),
		Package:  traverser.Packages.AddPackage(packagePath),
		Imports:  traverser.Packages.MapImportSpecs(file.Imports),
	}

	traverser.SetActiveFile(fileNode)

	return traverser, nil
}
