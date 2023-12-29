package utils

import (
	"errors"
	"go/types"
	"os"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

// ConvertStatusCodeTypeToInt converts a status code type to an integer.
// If it uses the net/http package to get the status code, it will return the integer value of the status code constant.
func ConvertStatusCodeTypeToInt(statusCode string) (int, error) {
	if strings.HasPrefix(statusCode, "http.") {
		statusCode = strings.Replace(statusCode, "http.", "", 1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return 0, err
	}

	pkgs, err := packages.Load(&packages.Config{
		Dir:  cwd,
		Mode: packages.NeedTypes,
	}, "net/http")
	if err != nil {
		return 0, err
	}

	for _, pkg := range pkgs {
		for _, typ := range pkg.Types.Scope().Names() {
			if typ == statusCode {
				// Return the integer value of the status code constant.

				constNode, ok := pkg.Types.Scope().Lookup(statusCode).(*types.Const)
				if !ok {
					return 0, errors.New("node is not a constant")
				}

				return strconv.Atoi(constNode.Val().ExactString())
			}
		}
	}

	return 0, errors.New("status code not found")
}
