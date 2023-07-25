package gin

import (
	"fmt"
	"github.com/ls6-events/gengo"
	"github.com/ls6-events/gengo/utils"
	"github.com/ls6-events/gengo/utils/astUtils"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"
)

// parseRoute parses a route from a gin routes
// It will populate the route with the handler function
// createRoute must be called before this
// It will open the file as an AST and find the handler function using the line number and function name
// It can also find the path parameters from the handler function
// It calls the parseFunction function to parse the handler function
func parseRoute(s *gengo.Service, baseRoute *gengo.Route) error {
	fset := token.NewFileSet()

	log := s.Log.With().Str("path", baseRoute.Path).Str("method", baseRoute.Method).Str("file", baseRoute.File).Logger()

	pkgPath := strings.Split(baseRoute.Handler, "/")
	names := strings.Split(pkgPath[len(pkgPath)-1], ".")

	if len(names) < 2 {
		err := fmt.Errorf("invalid handler name for file: %s", baseRoute.Handler)
		log.Error().Err(err).Msg("Failed to parse handler name")
		return err
	}

	pkgName := names[0]
	funcName := names[1]
	log.Debug().Str("pkgName", pkgName).Str("funcName", funcName).Msg("Found handler name")

	log.Debug().Msg("Parsing file")
	filePath := path.Join(s.WorkDir, baseRoute.File)
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse file")
		return err
	}

	baseRoute.PathParams = utils.ExtractParamsFromPath(baseRoute.Path)
	if len(baseRoute.PathParams) > 0 {
		log.Debug().Interface("pathParams", baseRoute.PathParams).Msg("Found path params")
	} else {
		log.Debug().Msg("No path params found")
	}

	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)

		if ok && funcDecl.Name.Name == funcName {
			log.Debug().Str("funcName", funcName).Msg("Found handler function")

			startPos := fset.Position(funcDecl.Pos())

			if baseRoute.LineNo != startPos.Line {
				// This means that the function is set inline in the route definition
				log.Debug().Str("funcName", funcName).Msg("Function is inline")

				ast.Inspect(funcDecl, func(n ast.Node) bool {
					funcLit, ok := n.(*ast.FuncLit)

					if ok {
						inlineStartPos := fset.Position(funcLit.Pos())

						if baseRoute.LineNo == inlineStartPos.Line {
							log.Debug().Str("funcName", funcName).Msg("Found inline handler function")

							err = parseFunction(s, log, baseRoute, funcLit, node.Imports, pkgName, strings.Join(pkgPath[:len(pkgPath)-1], "/"), 0)
							if err != nil {
								log.Error().Err(err).Msg("Failed to parse inline function")
								return false
							}

							log.Debug().Str("funcName", funcName).Interface("route", *baseRoute).Msg("Adding route")

							return false
						}
					}

					return true
				})

				return false
			}

			err = parseFunction(s, log, baseRoute, astUtils.FuncDeclToFuncLit(funcDecl), node.Imports, pkgName, strings.Join(pkgPath[:len(pkgPath)-1], "/"), 0)
			if err != nil {
				log.Error().Err(err).Msg("Failed to parse function")
				return false
			}

			log.Debug().Str("funcName", funcName).Interface("route", *baseRoute).Msg("Adding route")

			return false
		}

		return true
	})

	if err != nil {
		return err
	}

	return nil
}
