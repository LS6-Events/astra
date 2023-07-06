package gengo

import (
	"fmt"
)

// specialTypes is a map of types that need to be handled specially
var specialTypes = map[string]string{
	"time.Time": "string", // We can assume that all time.Time definitions will be serialized as strings
}

// handleSpecialType handles special types
func (s *Service) handleSpecialType(component *Field) {
	if _, ok := specialTypes[fmt.Sprintf("%s.%s", component.Package, component.Name)]; ok {
		s.Log.Debug().Type("pkg", component.Package).Str("type", component.Name).Msg("Mapping special case")

		component.Type = specialTypes[fmt.Sprintf("%s.%s", component.Package, component.Name)]

		component.StructFields = nil
		component.SliceType = ""
		component.MapValue = ""
	}
}
