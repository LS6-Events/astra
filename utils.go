package astra

import (
	"strings"

	"github.com/ls6-events/astra/astTraversal"
)

// ParseResultToField changes a result from the AST traversal to a local field.
func ParseResultToField(result astTraversal.Result) Field {
	field := Field{
		Type:                          result.Type,
		Name:                          result.Name,
		EnumValues:                    result.EnumValues,
		IsEmbedded:                    result.IsEmbedded,
		SliceType:                     result.SliceType,
		ArrayType:                     result.ArrayType,
		ArrayLength:                   result.ArrayLength,
		MapKeyType:                    result.MapKeyType,
		MapValueType:                  result.MapValueType,
		StructFieldBindingTags:        result.StructFieldBindingTags,
		StructFieldValidationTags:     result.StructFieldValidationTags,
		StructFieldValidationRequired: result.StructFieldValidationRequired,
	}

	// If the godoc is populated, we need to parse the response.s
	if result.Doc != "" {
		field.Doc = strings.TrimSpace(result.Doc)
	}

	// If the type is not a primitive type, we need to get the package path.
	// If the type is named, it is referring to a type.
	// If the slice type is populated and not a primitive type, we need to get the package path for the slice.
	// If the array type is populated and not a primitive type, we need to get the package path for the array.
	// If the map value type is populated and not a primitive type, we need to get the package path for the map value.
	if !IsAcceptedType(result.Type) || result.Name != "" ||
		(result.SliceType != "" && !IsAcceptedType(result.SliceType)) ||
		(result.ArrayType != "" && !IsAcceptedType(result.ArrayType)) ||
		(result.MapValueType != "" && !IsAcceptedType(result.MapValueType)) {
		field.Package = result.Package.Path()
	}

	// If the map key type is populated and not a primitive type, we need to get the package path for the map key.
	if result.MapKeyType != "" && !IsAcceptedType(result.MapKeyType) {
		field.MapKeyPackage = result.MapKeyPackage.Path()
	}

	// If the struct fields are populated, we need to parse them.
	if result.StructFields != nil {
		field.StructFields = make(map[string]Field)
		for name, value := range result.StructFields {
			structField := ParseResultToField(value)
			field.StructFields[name] = structField
		}
	}

	return field
}
