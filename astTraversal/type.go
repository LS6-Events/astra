package astTraversal

import (
	"go/ast"
	"go/types"
)

type TypeTraverser struct {
	Traverser *BaseTraverser
	Node      types.Type
	Package   *PackageNode
}

func (t *BaseTraverser) Type(node types.Type, packageNode *PackageNode) *TypeTraverser {
	return &TypeTraverser{
		Traverser: t,
		Node:      node,
		Package:   packageNode,
	}
}

func (t *TypeTraverser) Result() (Result, error) {
	var result Result
	switch n := t.Node.(type) {
	case *types.Basic:
		result = Result{
			Type:    n.Name(),
			Package: t.Package,
		}
	case *types.Named:
		pkg := t.Traverser.Packages.FindOrAdd(n.Obj().Pkg().Path())
		_, err := t.Traverser.Packages.Get(pkg)
		if err != nil {
			return Result{}, err
		}

		if t.Traverser.shouldAddComponent {
			namedUnderlyingResult, err := t.Traverser.Type(n.Underlying(), pkg).Result()
			if err != nil {
				return Result{}, err
			}

			namedUnderlyingResult.Name = n.Obj().Name()

			namedUnderlyingResult.Doc, err = t.Doc()
			if err != nil {
				return Result{}, err
			}

			err = t.Traverser.addComponent(namedUnderlyingResult)
			if err != nil {
				return Result{}, err
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
			isEmbedded := f.Embedded()

			// Get if "binding:required" tag is present and json/yaml/xml/form as well
			tag := n.Tag(i)
			name, isRequired, isShown := ParseStructTag(tag)

			if name == "" {
				name = f.Id()
			}

			if !f.Exported() {
				isShown = false
			}

			if !isShown {
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

			structFieldResult.IsRequired = isRequired
			structFieldResult.IsEmbedded = isEmbedded
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
