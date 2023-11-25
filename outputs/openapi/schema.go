package openapi

import (
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/astTraversal"
)

func mapParamToSchema(bindingType astTraversal.BindingTagType, param astra.Param) (Schema, bool) {
	if param.IsBound {
		return mapFieldToSchema(bindingType, param.Field)
	} else if param.IsArray {
		itemSchema := mapAcceptedType(param.Field.Type)
		if !astra.IsAcceptedType(param.Field.Type) {
			componentRef, ok := makeComponentRef(bindingType, param.Field.Type, param.Field.Package)
			if ok {
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
			componentRef, ok := makeComponentRef(bindingType, param.Field.Type, param.Field.Package)
			if ok {
				additionalProperties.Ref = componentRef
			}
		} else {
			additionalProperties = mapAcceptedType(param.Field.Type)
		}
		return Schema{
			Type:                 "object",
			AdditionalProperties: &additionalProperties,
		}, true
	} else {
		return mapAcceptedType(param.Field.Type), true
	}
}

func mapFieldToSchema(bindingType astTraversal.BindingTagType, field astra.Field) (Schema, bool) {
	if !astra.IsAcceptedType(field.Type) {
		componentRef, ok := makeComponentRef(bindingType, field.Type, field.Package)
		if ok {
			return Schema{
				Ref: componentRef,
			}, true
		}

		return Schema{}, false
	} else {
		schema := mapAcceptedType(field.Type)
		if field.Type == "slice" {
			itemSchema := Schema{
				Type: mapAcceptedType(field.Type).Type,
			}
			if !astra.IsAcceptedType(field.SliceType) {
				componentRef, ok := makeComponentRef(bindingType, field.SliceType, field.Package)
				if ok {
					itemSchema = Schema{
						Ref: componentRef,
					}
				}
			}
			schema.Items = &itemSchema
		} else if field.Type == "map" {
			var additionalProperties Schema
			if !astra.IsAcceptedType(field.MapValueType) {
				componentRef, ok := makeComponentRef(bindingType, field.MapValueType, field.Package)
				if ok {
					additionalProperties.Ref = componentRef
				}
			} else {
				additionalProperties = mapAcceptedType(field.MapValueType)
			}
			schema.AdditionalProperties = &additionalProperties
		}

		return schema, true
	}
}
