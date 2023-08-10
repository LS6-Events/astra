package openapi

import "github.com/ls6-events/gengo"

func mapParamToSchema(param gengo.Param) Schema {
	if param.IsBound {
		return mapFieldToSchema(param.Field)
	} else if param.IsArray {
		itemSchema := mapAcceptedType(param.Field.Type)
		if !gengo.IsAcceptedType(param.Field.Type) {
			itemSchema = Schema{
				Ref: makeComponentRef(param.Field.Type, param.Field.Package),
			}
		}
		return Schema{
			Type:  "array",
			Items: &itemSchema,
		}
	} else if param.IsMap {
		var additionalProperties Schema
		if !gengo.IsAcceptedType(param.Field.Type) {
			additionalProperties.Ref = makeComponentRef(param.Field.Type, param.Field.Package)
		} else {
			additionalProperties = mapAcceptedType(param.Field.Type)
		}
		return Schema{
			Type:                 "object",
			AdditionalProperties: &additionalProperties,
		}
	} else {
		return mapAcceptedType(param.Field.Type)
	}
}

func mapFieldToSchema(field gengo.Field) Schema {
	if !gengo.IsAcceptedType(field.Type) {
		return Schema{
			Ref: makeComponentRef(field.Type, field.Package),
		}
	} else {
		schema := mapAcceptedType(field.Type)
		if field.Type == "slice" {
			itemSchema := Schema{
				Type: mapAcceptedType(field.SliceType).Type,
			}
			if !gengo.IsAcceptedType(field.SliceType) {
				itemSchema = Schema{
					Ref: makeComponentRef(field.SliceType, field.Package),
				}
			}
			schema.Items = &itemSchema
		} else if field.Type == "map" {
			var additionalProperties Schema
			if !gengo.IsAcceptedType(field.MapValue) {
				additionalProperties.Ref = makeComponentRef(field.MapValue, field.Package)
			} else {
				additionalProperties = mapAcceptedType(field.MapValue)
			}
			schema.AdditionalProperties = &additionalProperties
		}

		return schema
	}
}
