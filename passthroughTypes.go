package astra

import "fmt"

// passthroughTypes is a list of types that need to be passed through.
var passthroughTypes = []string{
	"time.Time",
	"github.com/google/uuid.UUID",
}

// HandlePassthroughTypes handles passthrough types.
func (s *Service) HandlePassthroughTypes(component *Field) {
	qualifiedName := fmt.Sprintf("%s.%s", component.Package, component.Name)

	if strSliceContains(passthroughTypes, qualifiedName) {
		s.Log.Debug().Type("pkg", component.Package).Str("type", component.Name).Msg("Handling passthrough types")

		component.Type = qualifiedName

		component.StructFields = nil
		component.SliceType = ""
		component.ArrayType = ""
		component.MapValueType = ""
	}
}
