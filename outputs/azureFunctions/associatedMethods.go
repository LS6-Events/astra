package azureFunctions

import "strings"

func associatedMethods(method string) []string {
	return []string{strings.ToLower(method), "options"}
}
