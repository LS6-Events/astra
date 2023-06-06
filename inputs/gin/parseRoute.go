package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/gengo"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

func parseRoute(s *gengo.Service, file string, info gin.RouteInfo) error {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	pkgPath := strings.Split(info.Handler, "/")
	names := strings.Split(pkgPath[len(pkgPath)-1], ".")

	if len(names) < 2 {
		return fmt.Errorf("invalid handler name: %s", info.Handler)
	}

	pkgName := names[0]
	funcName := names[1]

	baseRoute := gengo.Route{
		Path:        info.Path,
		Method:      info.Method,
		PathParams:  make([]gengo.Param, 0),
		Body:        make([]gengo.Param, 0),
		QueryParams: make([]gengo.Param, 0),
		ReturnTypes: make([]gengo.ReturnType, 0),
	}

	paramRegex := regexp.MustCompile(`:[^\/]+|\*[^\/]+`)
	if paramRegex.MatchString(baseRoute.Path) {
		params := paramRegex.FindAllString(baseRoute.Path, -1)
		for _, param := range params {
			baseRoute.PathParams = append(baseRoute.PathParams, gengo.Param{
				Name:       param[1:],
				Type:       "string",
				IsRequired: param[0] == ':',
			})
		}
	}

	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if ok && funcDecl.Name.Name == funcName {
			err = parseFunction(s, &baseRoute, funcDecl, node.Imports, pkgName, strings.Join(pkgPath[:len(pkgPath)-1], "/"), 0)
			if err != nil {
				return false
			}

			s.AddRoute(baseRoute)
		}

		return true
	})

	if err != nil {
		return err
	}

	return nil
}
