package astra

import "maps"

// TypeFormat is the types of the standard go types that are accepted by OpenAPI.
type TypeFormat struct {
	Type   string
	Format string
}

// PredefinedTypeMap is the map of the standard go types that are accepted by OpenAPI.
// It contains the go type as a string and the corresponding OpenAPI type as the value - also including the format.
var PredefinedTypeMap = map[string]TypeFormat{
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

	// Custom handlers
	"file": {
		Type:   "string",
		Format: "binary",
	},
}

// WithCustomTypeMapping adds a custom type mapping to the predefined type map.
func WithCustomTypeMapping(customTypeMap map[string]TypeFormat) Option {
	return func(service *Service) {
		for k, v := range customTypeMap {
			service.CustomTypeMapping[k] = v
		}
	}
}

// WithCustomTypeMappingSingle adds a custom type mapping to the predefined type map.
func WithCustomTypeMappingSingle(key string, valueType string, valueFormat string) Option {
	return func(service *Service) {
		service.CustomTypeMapping[key] = TypeFormat{
			Type:   valueType,
			Format: valueFormat,
		}
	}
}

// GetTypeMapping returns the type mapping for the given key.
func (s *Service) GetTypeMapping(key string, pkg string) (TypeFormat, bool) {
	if s.fullTypeMapping == nil {
		s.fullTypeMapping = make(map[string]TypeFormat)
		maps.Copy(s.fullTypeMapping, PredefinedTypeMap)
		maps.Copy(s.fullTypeMapping, s.CustomTypeMapping)
	}

	if !IsAcceptedType(key) {
		key = pkg + "." + key
	}

	value, exists := s.fullTypeMapping[key]

	return value, exists
}
