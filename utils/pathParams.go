package utils

import (
	"github.com/ls6-events/gengo"
	"regexp"
)

func getPathParamRegex() *regexp.Regexp {
	return regexp.MustCompile(`:[^\/]+|\*[^\/]+`)
}

func ExtractParamsFromPath(path string) []gengo.Param {
	resultParams := make([]gengo.Param, 0)

	paramRegex := getPathParamRegex()
	if paramRegex.MatchString(path) {
		params := paramRegex.FindAllString(path, -1)
		for _, param := range params {
			resultParams = append(resultParams, gengo.Param{
				Name:       param[1:],
				Type:       "string",
				IsRequired: param[0] == ':',
			})
		}
	}

	return resultParams
}

func MapPathParams(path string, repl func(string) string) string {
	return getPathParamRegex().ReplaceAllStringFunc(path, repl)
}
