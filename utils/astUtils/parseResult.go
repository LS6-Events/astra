package astUtils

type ParseResult struct {
	VarName string
	PkgName string

	Value string

	MapKeyPkg string
	MapKey    string
	MapVal    string

	SliceType string
}
