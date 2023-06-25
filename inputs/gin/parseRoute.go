package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/gengo"
	"github.com/ls6-events/gengo/utils/astUtils"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

func parseRoute(s *gengo.Service, file string, line int, info gin.RouteInfo) error {
	fset := token.NewFileSet()

	log := s.Log.With().Str("path", info.Path).Str("method", info.Method).Str("file", file).Logger()

	log.Debug().Msg("Parsing file")
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse file")
		return err
	}

	pkgPath := strings.Split(info.Handler, "/")
	names := strings.Split(pkgPath[len(pkgPath)-1], ".")

	if len(names) < 2 {
		err := fmt.Errorf("invalid handler name for file: %s", info.Handler)
		log.Error().Err(err).Msg("Failed to parse handler name")
		return err
	}

	pkgName := names[0]
	funcName := names[1]
	log.Debug().Str("pkgName", pkgName).Str("funcName", funcName).Msg("Found handler name")

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
		log.Debug().Str("path", baseRoute.Path).Msg("Found path params")
		params := paramRegex.FindAllString(baseRoute.Path, -1)
		for _, param := range params {
			baseRoute.PathParams = append(baseRoute.PathParams, gengo.Param{
				Name:       param[1:],
				Type:       "string",
				IsRequired: param[0] == ':',
			})
		}
	} else {
		log.Debug().Str("path", baseRoute.Path).Msg("No path params found")
	}

	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)

		if ok && funcDecl.Name.Name == funcName {
			log.Debug().Str("funcName", funcName).Msg("Found handler function")

			startPos := fset.Position(funcDecl.Pos())

			if line != startPos.Line {
				// This means that the function is set inline in the route definition
				log.Debug().Str("funcName", funcName).Msg("Function is inline")

				ast.Inspect(funcDecl, func(n ast.Node) bool {
					funcLit, ok := n.(*ast.FuncLit)

					if ok {
						inlineStartPos := fset.Position(funcLit.Pos())

						if line == inlineStartPos.Line {
							log.Debug().Str("funcName", funcName).Msg("Found inline handler function")

							err = parseFunction(s, log, &baseRoute, funcLit, node.Imports, pkgName, strings.Join(pkgPath[:len(pkgPath)-1], "/"), 0)
							if err != nil {
								log.Error().Err(err).Msg("Failed to parse inline function")
								return false
							}

							log.Debug().Str("funcName", funcName).Interface("route", baseRoute).Msg("Adding route")
							s.AddRoute(baseRoute)

							return false
						}
					}

					return true
				})

				return false
			}

			err = parseFunction(s, log, &baseRoute, astUtils.FuncDeclToFuncLit(funcDecl), node.Imports, pkgName, strings.Join(pkgPath[:len(pkgPath)-1], "/"), 0)
			if err != nil {
				log.Error().Err(err).Msg("Failed to parse function")
				return false
			}

			log.Debug().Str("funcName", funcName).Interface("route", baseRoute).Msg("Adding route")
			s.AddRoute(baseRoute)

			return false
		}

		return true
	})

	if err != nil {
		return err
	}

	return nil
}
