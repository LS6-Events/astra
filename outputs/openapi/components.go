package openapi

import (
	"github.com/ls6-events/astra"
	"strconv"
	"strings"
)

// makeComponentRef creates a reference to the component in the OpenAPI specification
func makeComponentRef(name, pkg string) string {
	return "#/components/schemas/" + makeComponentRefName(name, pkg)
}

// makeComponentRefName converts the component and package name to a valid OpenAPI reference name (to avoid collisions)
func makeComponentRefName(name, pkg string) string {
	return strings.ReplaceAll(pkg, "/", "_") + "." + name
}

func mapComponentEnums(component astra.Field) []interface{} {
	var enums []interface{}
	for _, enum := range component.EnumValues {
		schema := mapAcceptedType(component.Type)
		switch schema.Type {
		case "string":
			enums = append(enums, strings.Trim(enum, "\""))
		case "integer":
			i, err := strconv.Atoi(enum)
			if err != nil {
				continue
			}

			enums = append(enums, i)
		case "number":
			f, err := strconv.ParseFloat(enum, 64)
			if err != nil {
				continue
			}

			enums = append(enums, f)
		case "boolean":
			b, err := strconv.ParseBool(enum)
			if err != nil {
				continue
			}

			enums = append(enums, b)
		}
	}

	return enums
}
