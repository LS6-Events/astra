package utils

import (
	"errors"
	"go/types"
	"golang.org/x/tools/go/packages"
	"os"
	"strconv"
	"strings"
)

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
				// Return the integer value of the status code constant
				return strconv.Atoi(pkg.Types.Scope().Lookup(statusCode).(*types.Const).Val().ExactString())
			}
		}
	}

	return 0, errors.New("status code not found")
}
