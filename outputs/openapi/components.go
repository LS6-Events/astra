package openapi

import (
	"github.com/ls6-events/astra"
	"slices"
	"strings"
)

// collisionSafeNames is a map of a full name package path to a collision safe name
var collisionSafeNames = make(map[string]string)

// collisionSafeKey creates a key for the collisionSafeNames map
func collisionSafeKey(name, pkg string) string {
	return pkg + "." + name
}

// getPackageName gets the package name from the package path (i.e. github.com/ls6-events/astra -> astra)
func getPackageName(pkg string) string {
	return pkg[strings.LastIndex(pkg, "/")+1:]
}

// makeCollisionSafeNamesFromComponents creates collision safe names for the components
// This needs to be run before any routes or components are generated
// As the makeComponentRefName function relies on the collisionSafeNames map
func makeCollisionSafeNamesFromComponents(components []astra.Field) {
	// Group the components by package name
	packageNames := make(map[string][]astra.Field)
	for _, component := range components {
		packageName := getPackageName(component.Package)

		// If the package name doesn't exist in the map, create it
		if _, ok := packageNames[packageName]; !ok {
			packageNames[packageName] = make([]astra.Field, 0)
		}

		// Append the component to the package name
		packageNames[packageName] = append(packageNames[packageName], component)
	}

	for _, components := range packageNames {
		// Iterate over every component and see if there is ever a case where the full package path doesn't match up
		sameUntil := 0
		for i := 0; i < len(components)-1; i++ {
			for j := i + 1; j < len(components)-i; j++ {
				// If the packages don't match, we need to find the first point where they don't match
				if components[i].Package != components[j].Package {
					// Create a slice of the package path split by "/"
					iComponentPackageSplit := strings.Split(components[i].Package, "/")
					jComponentPackageSplit := strings.Split(components[j].Package, "/")

					// Reverse the package path so we can iterate from the end first
					slices.Reverse(iComponentPackageSplit)
					slices.Reverse(jComponentPackageSplit)

					// Iterate over the package path slices and find the first point where they don't match
					for k := 0; k < len(iComponentPackageSplit) && k < len(jComponentPackageSplit); k++ {
						if iComponentPackageSplit[k] != jComponentPackageSplit[k] {
							// We've found the first point where they don't match, set sameUntil to k and break out of the loop
							sameUntil = k
							break
						}
					}
					break
				}
			}
		}

		// Iterate over every component and create a collision safe name
		for _, component := range components {
			// If sameUntil is greater than 0, we need to remove the package path up to the point where they first don't match
			if sameUntil > 0 {
				// Split the package path by "/"
				splitPackage := strings.Split(component.Package, "/")

				// Pick the final part of the package path, guided by sameUntil
				// We add 1 because we want to access the first different part of the package path
				// e.g. github.com/ls6-events/astra and github.com/different/astra would give us sameUntil = 1
				// and split into "ls6-events" and "different"

				splitPackage = splitPackage[len(splitPackage)-(sameUntil+1):]

				// Join the path back together
				collisionSafeNames[collisionSafeKey(component.Name, component.Package)] = strings.Join(splitPackage, ".") + "." + component.Name
			} else {
				// If sameUntil is 0, we can just use the package name
				collisionSafeNames[collisionSafeKey(component.Name, component.Package)] = getPackageName(component.Package) + "." + component.Name
			}
		}
	}
}

// makeComponentRef creates a reference to the component in the OpenAPI specification
func makeComponentRef(name, pkg string) string {
	return "#/components/schemas/" + makeComponentRefName(name, pkg)
}

// makeComponentRefName converts the component and package name to a valid OpenAPI reference name (to avoid collisions)
func makeComponentRefName(name, pkg string) string {
	return collisionSafeNames[collisionSafeKey(name, pkg)]
}

// componentToSchema converts a component to a schema
func componentToSchema(component astra.Field) Schema {
	var schema Schema

	if component.Type == "struct" {
		embeddedProperties := make([]Schema, 0)
		schema = Schema{
			Type:       "object",
			Properties: make(map[string]Schema),
		}
		for key, field := range component.StructFields {
			// We should aim to use doc comments in the future
			// However https://github.com/OAI/OpenAPI-Specification/issues/1514
			if field.IsEmbedded {
				embeddedProperties = append(embeddedProperties, Schema{
					Ref: makeComponentRef(field.Type, field.Package),
				})
				continue
			}

			schema.Properties[key] = componentToSchema(field)
		}

		if len(embeddedProperties) > 0 {
			if len(schema.Properties) == 0 {
				schema.AllOf = embeddedProperties
			} else {
				schema.AllOf = append(embeddedProperties, Schema{
					Properties: schema.Properties,
				})

				schema.Properties = nil
			}
		}
	} else if component.Type == "slice" {
		itemSchema := mapAcceptedType(component.SliceType)

		if itemSchema.Type == "" && !astra.IsAcceptedType(component.SliceType) {
			itemSchema = Schema{
				Ref: makeComponentRef(component.SliceType, component.Package),
			}
		}

		schema = Schema{
			Type:  "array",
			Items: &itemSchema,
		}
	} else if component.Type == "map" {
		additionalProperties := mapAcceptedType(component.MapValueType)

		if additionalProperties.Type == "" && !astra.IsAcceptedType(component.MapValueType) {
			additionalProperties.Ref = makeComponentRef(component.MapValueType, component.Package)
		}

		schema = Schema{
			Type:                 "object",
			AdditionalProperties: &additionalProperties,
		}
	} else {
		schema = mapAcceptedType(component.Type)
		if schema.Type == "" && !astra.IsAcceptedType(component.Type) {
			schema = Schema{
				Ref: makeComponentRef(component.Type, component.Package),
			}
		} else {
			schema.Enum = component.EnumValues
		}
	}

	return schema
}
