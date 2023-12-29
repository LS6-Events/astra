package gin

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/astTraversal"
	"github.com/ls6-events/astra/utils"
	"go/ast"
	"path"
)

// parseRoute parses a route from a gin routes
// It will populate the route with the handler function
// createRoute must be called before this
// It will open the file as an AST and find the handler function using the line number and function name
// It can also find the path parameters from the handler function
// It calls the parseFunction function to parse the handler function
func parseRoute(s *astra.Service, baseRoute *astra.Route) error {
	log := s.Log.With().Str("path", baseRoute.Path).Str("method", baseRoute.Method).Str("file", baseRoute.File).Logger()

	traverser := astTraversal.New(s.WorkDir).SetLog(&log)

	traverser.Packages.AddPathLoader(func(path string) (string, error) {
		if path == "main" {
			return s.GetMainPackageName()
		}
		return path, nil
	})

	handler := utils.SplitHandlerPath(baseRoute.Handler)

	pkgPath := handler.PackagePath()
	pkgName := handler.PackageName()

	if len(handler.HandlerParts) < 1 {
		err := fmt.Errorf("invalid handler name for file: %s", baseRoute.Handler)
		log.Error().Err(err).Msg("Failed to parse handler name")
		return err
	}

	funcName := handler.FuncName()

	pkgNode := traverser.Packages.AddPackage(pkgPath)

	log.Debug().Str("pkgName", pkgName).Str("funcName", funcName).Msg("Found handler name")

	log.Debug().Msg("Parsing file")

	_, err := traverser.Packages.Get(pkgNode)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get package")
		return err
	}

	for _, file := range pkgNode.Files {
		if path.Base(file.FileName) == path.Base(baseRoute.File) {
			log.Debug().Str("fileName", file.FileName).Msg("Found file")
			traverser.SetActiveFile(file)
			break
		}
	}

	if traverser.ActiveFile() == nil {
		err := fmt.Errorf("could not find file: %s", baseRoute.File)
		log.Error().Err(err).Msg("Failed to find file")
		return err
	}

	baseRoute.PathParams = utils.ExtractParamsFromPath(baseRoute.Path)
	if len(baseRoute.PathParams) > 0 {
		log.Debug().Interface("pathParams", baseRoute.PathParams).Msg("Found path params")
	} else {
		log.Debug().Msg("No path params found")
	}

	ast.Inspect(traverser.ActiveFile().AST, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		funcDecl, ok := n.(*ast.FuncDecl)

		if ok && funcDecl.Name.Name == funcName {
			log.Debug().Str("funcName", funcName).Msg("Found handler function")

			startPos := traverser.ActiveFile().Package.Package.Fset.Position(funcDecl.Pos())

			if baseRoute.LineNo != startPos.Line {
				// This means that the function is set inline in the route definition
				log.Debug().Str("funcName", funcName).Msg("Function is inline")

				ast.Inspect(funcDecl, func(n ast.Node) bool {
					if n == nil {
						return true
					}

					funcLit, ok := n.(*ast.FuncLit)

					if ok {
						inlineStartPos := traverser.ActiveFile().Package.Package.Fset.Position(funcLit.Pos())

						if baseRoute.LineNo == inlineStartPos.Line {
							log.Debug().Str("funcName", funcName).Msg("Found inline handler function")

							function, err := traverser.Function(funcLit)
							if err != nil {
								log.Error().Err(err).Msg("Failed to get function")
								return false
							}

							err = parseFunction(s, function, baseRoute, traverser.ActiveFile(), 0)
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

			// If the function is not inline, we can just parse it normally
			function, err := traverser.Function(funcDecl)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get function")
				return false
			}

			// And define the function name as the operation ID
			baseRoute.OperationID = strcase.ToLowerCamel(funcName)

			err = parseFunction(s, function, baseRoute, traverser.ActiveFile(), 0)
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
