package openapi

import (
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/astTraversal"
	"github.com/ls6-events/validjsonator"
)

func mapParamToSchema(bindingType astTraversal.BindingTagType, param astra.Param) (validjsonator.Schema, bool) {
	if param.IsBound {
		return mapFieldToSchema(bindingType, param.Field)
	} else if param.IsArray {
		itemSchema := mapPredefinedTypeFormat(param.Field.Type)
		if !astra.IsAcceptedType(param.Field.Type) {
			componentRef, bound := makeComponentRef(bindingType, param.Field.Type, param.Field.Package)
			if bound {
				itemSchema = validjsonator.Schema{
					Ref: componentRef,
				}
			}
		}
		return validjsonator.Schema{
			Type:  "array",
			Items: &itemSchema,
		}, true
	} else if param.IsMap {
		var additionalProperties validjsonator.Schema
		if !astra.IsAcceptedType(param.Field.Type) {
			componentRef, bound := makeComponentRef(bindingType, param.Field.Type, param.Field.Package)
			if bound {
				additionalProperties.Ref = componentRef
			}
		} else {
			additionalProperties = mapPredefinedTypeFormat(param.Field.Type)
		}
		return validjsonator.Schema{
			Type:                 "object",
			AdditionalProperties: &additionalProperties,
		}, true
	} else {
		return mapPredefinedTypeFormat(param.Field.Type), true
	}
}

func mapFieldToSchema(bindingType astTraversal.BindingTagType, field astra.Field) (validjsonator.Schema, bool) {
	if !astra.IsAcceptedType(field.Type) {
		componentRef, bound := makeComponentRef(bindingType, field.Type, field.Package)
		if bound {
			return validjsonator.Schema{
				Ref: componentRef,
			}, true
		}

		return validjsonator.Schema{}, false
	} else {
		schema := mapPredefinedTypeFormat(field.Type)
		if field.Type == "slice" {
			itemSchema := validjsonator.Schema{
				Type: mapPredefinedTypeFormat(field.SliceType).Type,
			}
			if !astra.IsAcceptedType(field.SliceType) {
				componentRef, bound := makeComponentRef(bindingType, field.SliceType, field.Package)
				if bound {
					itemSchema = validjsonator.Schema{
						Ref: componentRef,
					}
				}
			}
			schema.Items = &itemSchema
		} else if field.Type == "map" {
			var additionalProperties validjsonator.Schema
			if !astra.IsAcceptedType(field.MapValueType) {
				componentRef, bound := makeComponentRef(bindingType, field.MapValueType, field.Package)
				if bound {
					additionalProperties.Ref = componentRef
				}
			} else {
				additionalProperties = mapPredefinedTypeFormat(field.MapValueType)
			}
			schema.AdditionalProperties = &additionalProperties
		}

		return schema, true
	}
}

// mapTypeFormat maps the type with the list of types from the service.
// This should be primarily used for custom types in components that need to be mapped.
func mapTypeFormat(service *astra.Service, acceptedType string, pkg string) validjsonator.Schema {
	if acceptedType, ok := service.GetTypeMapping(acceptedType, pkg); ok {
		if acceptedType.Type == "" {
			return validjsonator.Schema{}
		}
		return validjsonator.Schema{
			Type:   acceptedType.Type,
			Format: acceptedType.Format,
		}
	}

	return validjsonator.Schema{}
}

// mapPredefinedTypeFormat maps the type with the list of types that are predefined.
// This should be primarily used for types that are not custom types, i.e. everywhere except top level components.
func mapPredefinedTypeFormat(acceptedType string) validjsonator.Schema {
	if acceptedType, ok := astra.PredefinedTypeMap[acceptedType]; ok {
		if acceptedType.Type == "" {
			return validjsonator.Schema{}
		}
		return validjsonator.Schema{
			Type:   acceptedType.Type,
			Format: acceptedType.Format,
		}
	}

	return validjsonator.Schema{}
}

// getQueryParamStyle returns the style of the query parameter, based on the schema.
func getQueryParamStyle(schema validjsonator.Schema) (style string, explode bool) {
	if schema.Type == "object" {
		return "deepObject", true
	}

	// The default behavior is to use the form style.
	// For arrays, we want comma separated values.
	return "form", false
}

// findComponentByPackageAndType finds the schema by the package and type.
func findComponentByPackageAndType(fields []astra.Field, pkg string, typeName string) (astra.Field, bool) {
	for _, field := range fields {
		if field.Package == pkg && field.Name == typeName {
			return field, true
		}
	}

	return astra.Field{}, false
}
