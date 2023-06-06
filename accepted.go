package gengo

var AcceptedTypes = []string{
	"nil",
	"string",
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"float",
	"float32",
	"float64",
	"bool",
	"byte",
	"rune",
	"struct",
	"map",
	"slice",
	"any",
}

func IsAcceptedType(t string) bool {
	return contains(AcceptedTypes, t)
}
