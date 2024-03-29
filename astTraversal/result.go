package astTraversal

type Result struct {
	// Type is the type of the result
	Type string

	// Package is the package of the result
	// It's used to reference the external package of a types.Named type
	// It's used as the package of a Slice, Array or Map Value
	Package *PackageNode

	// Name is the name of the result if it's a types.Named type
	Name string

	// Names is a list of names that are associated with the result (e.g. for a struct field)
	Names []string

	// IsEmbedded is true if the result is embedded in a struct
	IsEmbedded bool

	// ConstantValue is the constant value of the result (e.g. for a string)
	ConstantValue string

	// EnumValues is a list of enum values (e.g. for an enum { Foo, Bar })
	// This is used for when the result of a types.Named becomes a types.Basic and has constant values defined in the package
	EnumValues []any

	// MapKeyPackage is the package of the map key (e.g. for a map[string]string)
	MapKeyPackage *PackageNode

	// MapKeyType is the type of the map key (e.g. for a map[string]string)
	MapKeyType string

	// MapValueType is the type of the map value (e.g. for a map[string]string)
	MapValueType string

	// SliceType is the type of the slice (e.g. for a []string)
	SliceType string

	// ArrayType is the type of the array (e.g. for a [5]string)
	ArrayType string

	// ArrayLength is the length of the array (e.g. for a [5]string)
	ArrayLength int64

	// StructFields is a map of struct fields (e.g. for a struct { Foo string })
	StructFields map[string]Result

	StructFieldBindingTags BindingTagMap

	StructFieldValidationTags ValidationTagMap

	// Doc is the documentation of the result
	Doc string
}
