package astUtils

// ParseResult is the result of parsing a value
// It contains the variable name, package name and value
// If the value is a map, it contains the key and value types
// If the value is a slice, it contains the slice type
// If it is a constant value, it contains the value
type ParseResult struct {
	VarName string
	PkgName string

	Value string

	MapKeyPkg string
	MapKey    string
	MapVal    string

	SliceType string
}
