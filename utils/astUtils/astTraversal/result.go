package astTraversal

type Result struct {
	Type string

	// TODO handle structs better
	StructNames []string

	Package *PackageNode

	ConstantValue string

	MapKeyPackage *PackageNode
	MapKeyType    string
	MapValType    string

	SliceType string
}
