package openapi

import "strings"

// makeComponentRef creates a reference to the component in the OpenAPI specification
func makeComponentRef(name, pkg string) string {
	return "#/components/schemas/" + makeComponentRefName(name, pkg)
}

// makeComponentRefName converts the component and package name to a valid OpenAPI reference name (to avoid collisions)
func makeComponentRefName(name, pkg string) string {
	return strings.ReplaceAll(pkg, "/", "_") + "." + name
}
