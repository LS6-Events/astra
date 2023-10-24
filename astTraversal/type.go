package astTraversal

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"
)

type TypeTraverser struct {
	Traverser *BaseTraverser
	Node      types.Type
	Package   *PackageNode
	// name is the name of the type, if it comes from a types.Named
	// If it's not a types.Named, it's empty
	name string
}

func (t *BaseTraverser) Type(node types.Type, packageNode *PackageNode) *TypeTraverser {
	return &TypeTraverser{
		Traverser: t,
		Node:      node,
		Package:   packageNode,
	}
}

func (t *TypeTraverser) SetName(name string) *TypeTraverser {
	t.name = name
	return t
}

func (t *TypeTraverser) Result() (Result, error) {
	var result Result
	switch n := t.Node.(type) {
	case *types.Basic:
		result = Result{
			Type:    n.Name(),
			Package: t.Package,
		}

		// If the name isn't empty, it's a named type
		// Therefore it has the potential to be an enum
		if t.name != "" {
			// Get the package
			_, err := t.Traverser.Packages.Get(t.Package)
			if err != nil {
				return Result{}, err
			}

			// Iterate through the package's AST to find the enum values
			// We start by iterating over every file in the package
			for _, file := range t.Package.Package.Syntax {
				// Then we iterate over every declaration in the file
				for _, decl := range file.Decls {
					// If the declaration is a GenDecl, it's a const/var declaration
					if genDecl, ok := decl.(*ast.GenDecl); ok {
						// If the declaration isn't a const, we skip it (we're only looking for constants)
						if genDecl.Tok != token.CONST {
							continue
						}

						// If the declaration is a const, we iterate over every spec
						for _, spec := range genDecl.Specs {
							// If the spec is a ValueSpec, we check if the type is the same as the named type
							if valueSpec, ok := spec.(*ast.ValueSpec); ok {
								// If the type is the same as the named type, we iterate over every value
								if valueSpec.Type != nil {
									// We check this by comparing the name of the type to the name of the named type
									// It must be an Ident, otherwise it's not a named type, or it's from another package, not the one we're looking for
									if ident, ok := valueSpec.Type.(*ast.Ident); ok {
										if ident.Name == t.name {
											// We iterate over every value in the value spec
											for _, value := range valueSpec.Values {
												// If the value is a basic literal, we add it to the enum values
												if basicLit, ok := value.(*ast.BasicLit); ok {
													// Switch over the basic literal's kind to determine the type of the value
													// And format it accordingly
													switch n.Kind() {
													case types.String:
														result.EnumValues = append(result.EnumValues, strings.Trim(basicLit.Value, "\""))
													case types.Int:
														i, err := strconv.Atoi(basicLit.Value)
														if err != nil {
															continue
														}

														result.EnumValues = append(result.EnumValues, i)
													case types.Float32, types.Float64:
														f, err := strconv.ParseFloat(basicLit.Value, 64)
														if err != nil {
															continue
														}

														result.EnumValues = append(result.EnumValues, f)
													case types.Bool:
														b, err := strconv.ParseBool(basicLit.Value)
														if err != nil {
															continue
														}

														result.EnumValues = append(result.EnumValues, b)
													case types.Int8, types.Int16, types.Int32, types.Int64:
														i, err := strconv.ParseInt(basicLit.Value, 10, 64)
														if err != nil {
															continue
														}

														result.EnumValues = append(result.EnumValues, i)
													case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
														i, err := strconv.ParseUint(basicLit.Value, 10, 64)
														if err != nil {
															continue
														}

														result.EnumValues = append(result.EnumValues, i)
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

	case *types.Named:
		var pkg *PackageNode
		if n.Obj().Pkg() != nil {
			pkg = t.Traverser.Packages.FindOrAdd(n.Obj().Pkg().Path())
			_, err := t.Traverser.Packages.Get(pkg)
			if err != nil {
				return Result{}, err
			}

			if t.Traverser.shouldAddComponent {
				namedUnderlyingResult, err := t.Traverser.Type(n.Underlying(), pkg).SetName(n.Obj().Name()).Result()
				if err != nil {
					return Result{}, err
				}

				namedUnderlyingResult.Doc, err = t.Doc()
				if err != nil {
					return Result{}, err
				}

				err = t.Traverser.addComponent(namedUnderlyingResult)
				if err != nil {
					return Result{}, err
				}
			}
		}

		result = Result{
			Type:    n.Obj().Name(),
			Package: pkg,
		}
	case *types.Pointer:
		return t.Traverser.Type(n.Elem(), t.Package).Result()
	case *types.Slice:
		sliceElemResult, err := t.Traverser.Type(n.Elem(), t.Package).Result()
		if err != nil {
			return Result{}, err
		}

		result = Result{
			Type:      "slice",
			SliceType: sliceElemResult.Type,
			Package:   sliceElemResult.Package,
		}
	case *types.Array:
		arrayElemResult, err := t.Traverser.Type(n.Elem(), t.Package).Result()
		if err != nil {
			return Result{}, err
		}

		result = Result{
			Type:        "array",
			ArrayType:   arrayElemResult.Type,
			ArrayLength: n.Len(),
			Package:     arrayElemResult.Package,
		}
	case *types.Map:
		keyResult, err := t.Traverser.Type(n.Key(), t.Package).Result()
		if err != nil {
			return Result{}, err
		}

		valueResult, err := t.Traverser.Type(n.Elem(), t.Package).Result()
		if err != nil {
			return Result{}, err
		}

		result = Result{
			Type:          "map",
			MapKeyType:    keyResult.Type,
			MapKeyPackage: keyResult.Package,
			MapValueType:  valueResult.Type,
			Package:       valueResult.Package,
		}
	case *types.Struct:
		fields := make(map[string]Result)
		for i := 0; i < n.NumFields(); i++ {
			f := n.Field(i)
			name := f.Id()
			isExported := f.Exported()
			isEmbedded := f.Embedded()

			var bindingTag BindingTagMap
			var validationTags ValidationTagMap
			if isExported {
				bindingTag, validationTags = ParseStructTag(name, n.Tag(i))
			} else {
				continue
			}

			structFieldResult, err := t.Traverser.Type(f.Type(), t.Package).Result()
			if err != nil {
				return Result{}, err
			}

			if structFieldResult.Package != nil {
				_, err = t.Traverser.Packages.Get(structFieldResult.Package)
				if err != nil {
					return Result{}, err
				}

				pos := f.Pos()

				// TODO - this isn't working entirely well for external packages
				// Needs investigation
				node, err := structFieldResult.Package.ASTAtPos(pos)
				if err != nil || node == nil {
					t.Traverser.Log.Warn().Err(err).Msgf("failed to get AST at position %d for field %s", pos, f.Id())
				} else if field, ok := node.(*ast.Field); ok {
					structFieldResult.Doc = FormatDoc(field.Doc.Text())
				}
			}

			structFieldResult.IsEmbedded = isEmbedded
			structFieldResult.StructFieldBindingTags = bindingTag
			structFieldResult.StructFieldValidationTags = validationTags

			fields[name] = structFieldResult
		}

		result = Result{
			Type:         "struct",
			StructFields: fields,
			Package:      t.Package,
		}
	case *types.Interface:
		result = Result{
			Type:    "any",
			Package: t.Package,
		}
	}

	if t.name != "" {
		result.Name = t.name
	}

	if result.Type != "" {
		return result, nil
	} else {
		return Result{}, ErrInvalidNodeType
	}
}

func (t *TypeTraverser) Doc() (string, error) {
	if named, ok := t.Node.(*types.Named); ok {

		pkg := t.Traverser.Packages.AddPackage(named.Obj().Pkg().Path())

		_, err := t.Traverser.Packages.Get(pkg)
		if err != nil {
			return "", err
		}

		node, err := pkg.ASTAtPos(named.Obj().Pos())
		if err != nil || node == nil {
			t.Traverser.Log.Warn().Err(err).Msgf("failed to get AST for type %s", t.Node.String())
			for expr, info := range pkg.Package.TypesInfo.Types {
				if info.Type.String() == named.String() {
					node = expr
					break
				}
			}
		}

		for node != nil {
			switch n := node.(type) {
			case *ast.TypeSpec:
				doc := n.Doc.Text()
				if doc != "" {
					return FormatDoc(doc), nil
				}

				// If the doc doesn't exist, we need to find the declaration of the type
				// and get the doc from there
				// This is because the doc is attached to the declaration, not the AST type
				for _, file := range pkg.Package.Syntax {
					if pkg.Package.Fset.Position(file.Pos()).Filename == pkg.Package.Fset.Position(n.Pos()).Filename {
						for _, decl := range file.Decls {
							if genDecl, ok := decl.(*ast.GenDecl); ok {
								for _, spec := range genDecl.Specs {
									if typeSpec, ok := spec.(*ast.TypeSpec); ok {
										if typeSpec.Name.Name == n.Name.Name {
											return FormatDoc(genDecl.Doc.Text()), nil
										}
									}
								}
							}
						}
					}
				}

				node = nil
			case *ast.CompositeLit:
				node = n.Type
			case *ast.Ident:
				if n.Obj != nil {
					node = n.Obj.Decl.(ast.Node)
				} else {
					node = nil
				}
			default:
				node = nil
			}
		}

	}

	return "", nil
}
