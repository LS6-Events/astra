package astra

import (
	"fmt"
)

// substituteTypes is a map of types that need to be replaced with other types.
var substituteTypes = map[string]string{
	"time.Duration": "int",    // We can assume that all time.Duration definitions will be serialized as int.
	"error":         "string", // We can assume that all errors will be serialized as strings.
}

// HandleSubstituteTypes handles substitute types.
func (s *Service) HandleSubstituteTypes(component *Field) {
	qualifiedName := fmt.Sprintf("%s.%s", component.Package, component.Name)

	if _, ok := substituteTypes[qualifiedName]; ok {
		s.Log.Debug().Type("pkg", component.Package).Str("type", component.Name).Msg("Handling substitute types")

		component.Type = substituteTypes[qualifiedName]

		component.StructFields = nil
		component.SliceType = ""
		component.MapValueType = ""
	}
}
