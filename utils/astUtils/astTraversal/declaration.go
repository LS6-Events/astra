package astTraversal

import (
	"errors"
	"go/ast"
)

type DeclarationTraverser struct {
	Traverser *BaseTraverser
	Decl      ast.Node
	File      *FileNode
	VarName   string
}

func (t *BaseTraverser) Declaration(node ast.Node, varName string) (*DeclarationTraverser, error) {
	return &DeclarationTraverser{
		Traverser: t,
		Decl:      node,
		File:      t.ActiveFile(),
		VarName:   varName, // The name of the variable on the LHS of the arrangement
	}, nil
}

func (dt *DeclarationTraverser) AssignStmt(assignStmt *ast.AssignStmt, varName string) (Result, error) {
	var varNum int
	for i, expr := range assignStmt.Lhs {
		if expr.(*ast.Ident).Name == varName {
			varNum = i
			break
		}
	}

	var node ast.Node

	if len(assignStmt.Lhs) == len(assignStmt.Rhs) {
		node = assignStmt.Rhs[varNum]
	} else {
		node = assignStmt.Rhs[0]
	}

	literal, err := dt.Traverser.Literal(node, varNum)
	if err != nil {
		return Result{}, err
	}

	return literal.Result()
}

func (dt *DeclarationTraverser) Result(varName string) (Result, error) {
	switch decl := dt.Decl.(type) {
	case *ast.AssignStmt:
		return dt.AssignStmt(decl, varName)
	case *ast.GenDecl:
		specIndex := -1
		nameIndex := -1
		for i, spec := range decl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				return Result{}, errors.New("not parsed as value type for gendecl")
			}

			for j, name := range valueSpec.Names {
				if name.Name == dt.VarName {
					specIndex = i
					nameIndex = j
					break
				}
			}
		}

		if specIndex != -1 && nameIndex != -1 {
			valueSpec := decl.Specs[specIndex].(*ast.ValueSpec)

			result, err := dt.ValueSpecResult(valueSpec, nameIndex)
			if err != nil {
				return Result{}, err
			}

			if decl.Tok.String() == "const" {
				literal, err := dt.Traverser.Literal(valueSpec.Values[nameIndex], -1)
				if err != nil {
					return Result{}, err
				}

				literalResult, err := literal.Result()
				if err != nil {
					return Result{}, err
				}

				result.ConstantValue = literalResult.ConstantValue
			}

			return result, nil
		} else {
			return Result{}, errors.New("cannot find declaration line for gendecl")
		}
	case *ast.ValueSpec:
		nameIndex := 0 // Default value of 0 for single value spec
		for i, name := range decl.Names {
			if name.Name == dt.VarName {
				nameIndex = i
				break
			}
		}
		// TODO Figure out way to trace value spec back to Gen Decl AST - ident.Obj.Decl gives you a value spec
		return dt.ValueSpecResult(decl, nameIndex)
	default:
		return Result{}, errors.New("unsupported declaration type")
	}
}

func (dt *DeclarationTraverser) ValueSpecResult(valueSpec *ast.ValueSpec, nameIndex int) (Result, error) {
	if valueSpec.Type == nil { // Type inference
		lt, err := dt.Traverser.Literal(valueSpec.Values[nameIndex], -1)
		if err != nil {
			return Result{}, err
		}

		return lt.Result()
	}
	return dt.Traverser.Expression(valueSpec.Type).Result()
}
