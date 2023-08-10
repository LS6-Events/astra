package gin

import (
	"errors"
	"github.com/ls6-events/gengo"
	"github.com/ls6-events/gengo/utils"
	"github.com/ls6-events/gengo/utils/astUtils"
	"github.com/ls6-events/gengo/utils/astUtils/astTraversal"
	"go/ast"
	"strings"
)

// parseFunction parses a function and adds it to the service
// It is designed to be called recursively should it be required
// The level parameter is used to determine the depth of recursion
// And the package name and path are used to determine the package of the currently analysed function
// The currRoute reference is used to manipulate the current route being analysed
// The imports are used to determine the package of the context variable
func parseFunction(traverser *astTraversal.Traverser, s *gengo.Service, currRoute *gengo.Route, node ast.Node, level int) error {
	// Get the variable name of the context parameter
	funcExpr, err := traverser.Function(node)
	if err != nil {
		return err
	}

	ctxName := funcExpr.FindArgumentNameByType("Context", "github.com/gin-gonic/gin", true)
	if ctxName == "" {
		return errors.New("failed to find context variable name")
	}

	// Loop over every statement in the function
	ast.Inspect(funcExpr.Node.Body, func(n ast.Node) bool {
		// If a function is called
		var callExpr *astTraversal.CallExpressionTraverser
		callExpr, err = traverser.CallExpression(n)
		if errors.Is(err, astTraversal.ErrInvalidNodeType) {
			err = nil
			return true
		} else if err != nil {
			return true
		}

		_, ok := callExpr.ArgIndex(ctxName)
		if !ok {
			result := callExpr.FuncResult()
			if len(result.StructNames) > 0 && result.StructNames[0] == ctxName {
				switch callExpr.FuncName() {
				case "JSON":
					fallthrough
				case "XML":
					fallthrough
				case "YAML":
					fallthrough
				case "ProtoBuf":
					fallthrough
				case "Data":
					switch callExpr.FuncName() {
					case "JSON":
						currRoute.ContentType = "application/json"
					case "XML":
						currRoute.ContentType = "application/xml"
					case "YAML":
						currRoute.ContentType = "application/yaml"
					case "ProtoBuf":
						currRoute.ContentType = "application/protobuf"
					}

					// Get the status code
					var statusCode int
					statusCode, err = astUtils.ExtractStatusCode(callExpr.Args()[0])
					if err != nil {
						return true
					}

					argNo := 1
					if callExpr.FuncName() == "Data" {
						argNo = 2
					}

					expr := traverser.Expression(callExpr.Args()[argNo])

					result, err = expr.Result()
					if err != nil {
						traverser.Log.Error().Err(err).Msg("failed to get result for expression")
						return false
					}
					if expr.DoesNeedTracing() {
						var decl *astTraversal.DeclarationTraverser
						decl, err = traverser.FindDeclarationForNode(expr.Node)
						if err != nil {
							return false
						}

						result, err = decl.Result(result.Type)
						if err != nil {
							traverser.Log.Error().Err(err).Msg("failed to get result for declaration")
							return false
						}
					}

					returnType := gengo.ReturnType{
						StatusCode: statusCode,
						Field:      parseResultToField(s, result),
					}

					currRoute.ReturnTypes = utils.AddReturnType(currRoute.ReturnTypes, returnType)

					return true
				case "String": // c.String
					currRoute.ContentType = "text/plain"

					var statusCode int
					statusCode, err = astUtils.ExtractStatusCode(callExpr.Args()[0])
					if err != nil {
						return false
					}

					returnType := gengo.ReturnType{
						StatusCode: statusCode,
						Field: gengo.Field{
							Type: "string",
						},
					}
					currRoute.ReturnTypes = utils.AddReturnType(currRoute.ReturnTypes, returnType)

					return true
				case "Status": // c.Status
					var statusCode int
					statusCode, err = astUtils.ExtractStatusCode(callExpr.Args()[0])
					if err != nil {
						return false
					}

					returnType := gengo.ReturnType{
						StatusCode: statusCode,
						Field: gengo.Field{
							Type: "nil",
						},
					}
					currRoute.ReturnTypes = utils.AddReturnType(currRoute.ReturnTypes, returnType)

					return true
				// Query Param methods
				case "GetQuery":
					fallthrough
				case "Query":
					var queryParam gengo.Param
					queryParam, err = extractSingleRequestParam(traverser, s, callExpr.Args()[0], gengo.Param{})
					if err != nil {
						return false
					}

					currRoute.QueryParams = append(currRoute.QueryParams, queryParam)

					return true
				case "GetQueryArray":
					fallthrough
				case "QueryArray":
					var queryParam gengo.Param
					queryParam, err = extractSingleRequestParam(traverser, s, callExpr.Args()[0], gengo.Param{
						IsArray: true,
					})
					if err != nil {
						return false
					}

					currRoute.QueryParams = append(currRoute.QueryParams, queryParam)

					return true
				case "GetQueryMap":
					fallthrough
				case "QueryMap":
					var queryParam gengo.Param
					queryParam, err = extractSingleRequestParam(traverser, s, callExpr.Args()[0], gengo.Param{
						IsMap: true,
					})
					if err != nil {
						return false
					}

					currRoute.QueryParams = append(currRoute.QueryParams, queryParam)

					return true
				case "ShouldBindQuery":
					fallthrough
				case "BindQuery":
					var queryParam gengo.Param
					queryParam, err = extractBoundRequestParam(traverser, s, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.QueryParams = append(currRoute.QueryParams, queryParam)

					return true
				// Body Param methods
				case "ShouldBind":
					fallthrough
				case "Bind":
					var bodyParam gengo.Param
					bodyParam, err = extractBoundRequestParam(traverser, s, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.BodyType = "form"

					currRoute.Body = append(currRoute.Body, bodyParam)

					return true
				case "ShouldBindJSON":
					fallthrough
				case "BindJSON":
					var bodyParam gengo.Param
					bodyParam, err = extractBoundRequestParam(traverser, s, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/json"

					currRoute.Body = append(currRoute.Body, bodyParam)

					return true
				case "ShouldBindXML":
					fallthrough
				case "BindXML":
					var bodyParam gengo.Param
					bodyParam, err = extractBoundRequestParam(traverser, s, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/xml"

					currRoute.Body = append(currRoute.Body, bodyParam)

					return true
				case "ShouldBindYAML":
					fallthrough
				case "BindYAML":
					var bodyParam gengo.Param
					bodyParam, err = extractBoundRequestParam(traverser, s, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/yaml"

					currRoute.Body = append(currRoute.Body, bodyParam)

					return true
				case "GetPostForm":
					fallthrough
				case "PostForm":
					var bodyParam gengo.Param
					bodyParam, err = extractSingleRequestParam(traverser, s, callExpr.Args()[0], gengo.Param{})
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/x-www-form-urlencoded"
					currRoute.Body = append(currRoute.Body, bodyParam)

					return true
				case "GetPostFormArray":
					fallthrough
				case "PostFormArray":
					var bodyParam gengo.Param
					bodyParam, err = extractSingleRequestParam(traverser, s, callExpr.Args()[0], gengo.Param{
						IsArray: true,
					})
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/x-www-form-urlencoded"
					currRoute.Body = append(currRoute.Body, bodyParam)

					return true
				case "GetPostFormMap":
					fallthrough
				case "PostFormMap":
					var bodyParam gengo.Param
					bodyParam, err = extractSingleRequestParam(traverser, s, callExpr.Args()[0], gengo.Param{
						IsMap: true,
					})
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/x-www-form-urlencoded"
					currRoute.Body = append(currRoute.Body, bodyParam)

					return true
				default:
					return true
				}
			}
		} else {
			var function *astTraversal.FunctionTraverser
			function, err = callExpr.Function()
			if err != nil {
				traverser.Log.Error().Err(err).Msg("failed to get function")
				return false
			}

			err = parseFunction(traverser, s, currRoute, function.Node, level+1)
			if err != nil {
				traverser.Log.Error().Err(err).Msg("error parsing function")
				return false
			}
		}

		return true
	})

	if err != nil {
		return err
	}

	if len(currRoute.ReturnTypes) == 0 && level == 0 {
		return errors.New("return type not found")
	}

	return nil
}

func extractSingleRequestParam(traverser *astTraversal.Traverser, s *gengo.Service, node ast.Node, baseParam gengo.Param) (gengo.Param, error) {
	expr := traverser.Expression(node)
	result, err := expr.Result()
	if err != nil {
		traverser.Log.Error().Err(err).Msg("failed to parse expression")
		return gengo.Param{}, err
	}

	if expr.DoesNeedTracing() {
		var decl *astTraversal.DeclarationTraverser
		decl, err = traverser.FindDeclarationForNode(expr.Node)
		if err != nil {
			traverser.Log.Error().Err(err).Msg("failed to find declaration")
			return gengo.Param{}, err
		}

		result, err = decl.Result(result.Type)
		if err != nil {
			traverser.Log.Error().Err(err).Msg("failed to get result for declaration")
			return gengo.Param{}, err
		}
	}

	return gengo.Param{
		Name:       strings.ReplaceAll(result.ConstantValue, "\"", ""),
		Field:      parseResultToField(s, result),
		IsArray:    baseParam.IsArray,
		IsMap:      baseParam.IsMap,
		IsRequired: baseParam.IsRequired,
	}, nil
}

func extractBoundRequestParam(traverser *astTraversal.Traverser, s *gengo.Service, node ast.Node) (gengo.Param, error) {
	expr := traverser.Expression(node)
	result, err := expr.Result()
	if err != nil {
		traverser.Log.Error().Err(err).Msg("failed to parse expression")
		return gengo.Param{}, err
	}

	if expr.DoesNeedTracing() {
		var decl *astTraversal.DeclarationTraverser
		decl, err = traverser.FindDeclarationForNode(expr.Node)
		if err != nil {
			traverser.Log.Error().Err(err).Msg("failed to find declaration")
			return gengo.Param{}, err
		}

		result, err = decl.Result(result.Type)
		if err != nil {
			traverser.Log.Error().Err(err).Msg("failed to get result for declaration")
			return gengo.Param{}, err
		}
	}

	bodyParam := gengo.Param{
		IsBound: true,
		Field:   parseResultToField(s, result),
	}

	return bodyParam, nil
}

func parseResultToField(s *gengo.Service, result astTraversal.Result) gengo.Field {
	field := gengo.Field{
		Type:      result.Type,
		SliceType: result.SliceType,
		MapKey:    result.MapKeyType,
		MapValue:  result.MapValType,
	}

	if !gengo.IsAcceptedType(result.Type) {
		s.AddToBeProcessed(result.Package.Path(), result.Type)
		field.Package = result.Package.Path()
	}
	if result.SliceType != "" && !gengo.IsAcceptedType(result.SliceType) {
		s.AddToBeProcessed(result.Package.Path(), result.SliceType)
		field.Package = result.Package.Path()
	}
	if result.MapKeyType != "" && !gengo.IsAcceptedType(result.MapKeyType) {
		s.AddToBeProcessed(result.MapKeyPackage.Path(), result.MapKeyType)
		field.MapKeyPkg = result.MapKeyPackage.Path()
	}
	if result.MapValType != "" && !gengo.IsAcceptedType(result.MapValType) {
		s.AddToBeProcessed(result.Package.Path(), result.MapValType)
		field.Package = result.Package.Path()
	}

	return field
}
