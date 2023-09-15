package utils

import (
	"github.com/ls6-events/astra"
	"regexp"
)

// getPathParamRegex returns a regex to match path parameters
// This regex matches both :param and *param (gin style)
// It will need to be updated if other frameworks are supported that contain different syntax
func getPathParamRegex() *regexp.Regexp {
	return regexp.MustCompile(`:[^\/]+|\*[^\/]+`)
}

// ExtractParamsFromPath extracts the parameters from a path
func ExtractParamsFromPath(path string) []astra.Param {
	resultParams := make([]astra.Param, 0)

	paramRegex := getPathParamRegex()
	if paramRegex.MatchString(path) {
		params := paramRegex.FindAllString(path, -1)
		for _, param := range params {
			resultParams = append(resultParams, astra.Param{
				Name: param[1:],
				Field: astra.Field{
					Type: "string",
				},
				IsRequired: param[0] == ':',
			})
		}
	}

	return resultParams
}

// MapPathParams maps the path parameters to a new path
// Useful when converting a Gin path to an OpenAPI path
func MapPathParams(path string, repl func(string) string) string {
	return getPathParamRegex().ReplaceAllStringFunc(path, repl)
}
