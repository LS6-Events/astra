package openapi

import (
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/astTraversal"
)

func mapParamToSchema(bindingType astTraversal.BindingTagType, param astra.Param) (Schema, bool) {
	if param.IsBound {
		return mapFieldToSchema(bindingType, param.Field)
	} else if param.IsArray {
		itemSchema := mapPredefinedTypeFormat(param.Field.Type)
		if !astra.IsAcceptedType(param.Field.Type) {
			componentRef, bound := makeComponentRef(bindingType, param.Field.Type, param.Field.Package)
			if bound {
				itemSchema = Schema{
					Ref: componentRef,
				}
			}
		}
		return Schema{
			Type:  "array",
			Items: &itemSchema,
		}, true
	} else if param.IsMap {
		var additionalProperties Schema
		if !astra.IsAcceptedType(param.Field.Type) {
			componentRef, bound := makeComponentRef(bindingType, param.Field.Type, param.Field.Package)
			if bound {
				additionalProperties.Ref = componentRef
			}
		} else {
			additionalProperties = mapPredefinedTypeFormat(param.Field.Type)
		}
		return Schema{
			Type:                 "object",
			AdditionalProperties: &additionalProperties,
		}, true
	} else {
		return mapPredefinedTypeFormat(param.Field.Type), true
	}
}

func mapFieldToSchema(bindingType astTraversal.BindingTagType, field astra.Field) (Schema, bool) {
	if !astra.IsAcceptedType(field.Type) {
		componentRef, bound := makeComponentRef(bindingType, field.Type, field.Package)
		if bound {
			return Schema{
				Ref: componentRef,
			}, true
		}

		return Schema{}, false
	} else {
		schema := mapPredefinedTypeFormat(field.Type)
		if field.Type == "slice" {
			itemSchema := Schema{
				Type: mapPredefinedTypeFormat(field.SliceType).Type,
			}
			if !astra.IsAcceptedType(field.SliceType) {
				componentRef, bound := makeComponentRef(bindingType, field.SliceType, field.Package)
				if bound {
					itemSchema = Schema{
						Ref: componentRef,
					}
				}
			}
			schema.Items = &itemSchema
		} else if field.Type == "map" {
			var additionalProperties Schema
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

// mapTypeFormat maps the type with the list of types from the service
// This should be primarily used for custom types in components that need to be mapped
func mapTypeFormat(service *astra.Service, acceptedType string, pkg string) Schema {
	if acceptedType, ok := service.GetTypeMapping(acceptedType, pkg); ok {
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

// mapPredefinedTypeFormat maps the type with the list of types that are predefined
// This should be primarily used for types that are not custom types, i.e. everywhere except top level components
func mapPredefinedTypeFormat(acceptedType string) Schema {
	if acceptedType, ok := astra.PredefinedTypeMap[acceptedType]; ok {
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
