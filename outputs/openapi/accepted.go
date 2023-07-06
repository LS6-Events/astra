package openapi

// openAPIJSONType is the types of the standard go types that are accepted by OpenAPI
type openAPIJSONType struct {
	Type   string
	Format string
}

// acceptedTypeMap is the map of the standard go types that are accepted by OpenAPI
// It contains the go type as a string and the corresponding OpenAPI type as the value - also including the format
var acceptedTypeMap = map[string]openAPIJSONType{
	"string": openAPIJSONType{
		Type: "string",
	},
	"int": openAPIJSONType{
		Type:   "integer",
		Format: "int32",
	},
	"int8": openAPIJSONType{
		Type:   "integer",
		Format: "int8",
	},
	"int16": openAPIJSONType{
		Type:   "integer",
		Format: "int16",
	},
	"int32": openAPIJSONType{
		Type:   "integer",
		Format: "int32",
	},
	"int64": openAPIJSONType{
		Type:   "integer",
		Format: "int64",
	},
	"uint": openAPIJSONType{
		Type:   "integer",
		Format: "uint",
	},
	"uint8": openAPIJSONType{
		Type:   "integer",
		Format: "uint8",
	},
	"uint16": openAPIJSONType{
		Type:   "integer",
		Format: "uint16",
	},
	"uint32": openAPIJSONType{
		Type:   "integer",
		Format: "uint32",
	},
	"uint64": openAPIJSONType{
		Type:   "integer",
		Format: "uint64",
	},
	"float": openAPIJSONType{
		Type:   "number",
		Format: "float",
	},
	"float32": openAPIJSONType{
		Type:   "number",
		Format: "float32",
	},
	"float64": openAPIJSONType{
		Type:   "number",
		Format: "float64",
	},
	"bool": openAPIJSONType{
		Type: "boolean",
	},
	"byte": openAPIJSONType{
		Type:   "string",
		Format: "byte",
	},
	"rune": openAPIJSONType{
		Type:   "string",
		Format: "rune",
	},
	"struct": openAPIJSONType{
		Type: "object",
	},
	"map": openAPIJSONType{
		Type: "object",
	},
	"slice": openAPIJSONType{
		Type: "array",
	},
	"any": openAPIJSONType{
		Type: "",
	},
	"nil": openAPIJSONType{
		Type: "",
	},
}

// mapAcceptedType maps the accepted type to the OpenAPI type
// It returns a Schema with the type set to the OpenAPI type
// If not found, it returns an empty Schema
// TODO: Add support for format field in Schema
func mapAcceptedType(acceptedType string) Schema {
	if acceptedType, ok := acceptedTypeMap[acceptedType]; ok {
		if acceptedType.Type == "" {
			return Schema{}
		}
		return Schema{
			Type: acceptedType.Type,
		}
	}

	return Schema{}
}
