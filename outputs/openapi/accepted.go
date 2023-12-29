package openapi

// openAPIJSONType is the types of the standard go types that are accepted by OpenAPI.
type openAPIJSONType struct {
	Type   string
	Format string
}

// acceptedTypeMap is the map of the standard go types that are accepted by OpenAPI.
// It contains the go type as a string and the corresponding OpenAPI type as the value - also including the format.
var acceptedTypeMap = map[string]openAPIJSONType{
	"string": {
		Type: "string",
	},
	"int": {
		Type:   "integer",
		Format: "int32",
	},
	"int8": {
		Type:   "integer",
		Format: "int8",
	},
	"int16": {
		Type:   "integer",
		Format: "int16",
	},
	"int32": {
		Type:   "integer",
		Format: "int32",
	},
	"int64": {
		Type:   "integer",
		Format: "int64",
	},
	"uint": {
		Type:   "integer",
		Format: "uint",
	},
	"uint8": {
		Type:   "integer",
		Format: "uint8",
	},
	"uint16": {
		Type:   "integer",
		Format: "uint16",
	},
	"uint32": {
		Type:   "integer",
		Format: "uint32",
	},
	"uint64": {
		Type:   "integer",
		Format: "uint64",
	},
	"float": {
		Type:   "number",
		Format: "float",
	},
	"float32": {
		Type:   "number",
		Format: "float32",
	},
	"float64": {
		Type:   "number",
		Format: "float64",
	},
	"bool": {
		Type: "boolean",
	},
	"byte": {
		Type:   "string",
		Format: "byte",
	},
	"rune": {
		Type:   "string",
		Format: "rune",
	},
	"struct": {
		Type: "object",
	},
	"map": {
		Type: "object",
	},
	"slice": {
		Type: "array",
	},
	"any": {
		Type: "",
	},
	"nil": {
		Type: "",
	},

	// Passthrough types
	"time.Time": {
		Type:   "string",
		Format: "date-time",
	},
	"github.com/google/uuid.UUID": {
		Type:   "string",
		Format: "uuid",
	},
}

// mapAcceptedType maps the accepted type to the OpenAPI type.
// It returns a Schema with the type set to the OpenAPI type.
// If not found, it returns an empty Schema.
func mapAcceptedType(acceptedType string) Schema {
	if acceptedType, ok := acceptedTypeMap[acceptedType]; ok {
		if acceptedType.Type == "" {
			return Schema{}
		}
		return Schema{
			Type:   acceptedType.Type,
			Format: acceptedType.Format,
		}
	}

	return Schema{}
}
