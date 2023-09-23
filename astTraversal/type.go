package astTraversal

import (
	"go/doc"
	"go/types"
	"strings"
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

func (t *TypeTraverser) Doc() (*doc.Type, error) {
	pkgDoc, err := t.Package.GoDoc()
	if err != nil {
		return nil, err
	}

	splitNode := strings.Split(t.Node.String(), ".")
	nodeName := splitNode[len(splitNode)-1]

	for _, typ := range pkgDoc.Types {
		if typ.Name == nodeName {
			return typ, nil
		}
	}

	return nil, nil
}
