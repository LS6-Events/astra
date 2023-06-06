package openapi

import "strings"

func makeComponentRef(name, pkg string) string {
	return "#/components/schemas/" + makeComponentRefName(name, pkg)
}

func makeComponentRefName(name, pkg string) string {
	return strings.ReplaceAll(pkg, "/", "_") + "." + name
}
