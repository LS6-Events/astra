package gin

import (
	"errors"
	"github.com/ls6-events/gengo"
	"github.com/ls6-events/gengo/astTraversal"
	"github.com/ls6-events/gengo/utils"
	"go/ast"
	"go/types"
)

const (
	// GinPackagePath is the import path of the gin package
	GinPackagePath = "github.com/gin-gonic/gin"
	// GinContextType is the type of the context variable
	GinContextType = "Context"
	// GinContextIsPointer is whether the context variable is a pointer for the handler functions
	GinContextIsPointer = true
)

// parseFunction parses a function and adds it to the service
// It is designed to be called recursively should it be required
// The level parameter is used to determine the depth of recursion
// And the package name and path are used to determine the package of the currently analysed function
// The currRoute reference is used to manipulate the current route being analysed
// The imports are used to determine the package of the context variable
func parseFunction(s *gengo.Service, funcTraverser *astTraversal.FunctionTraverser, currRoute *gengo.Route, activeFile *astTraversal.FileNode, level int) error {
	traverser := funcTraverser.Traverser

	traverser.SetActiveFile(activeFile)
	traverser.SetAddComponentFunction(addComponent(s))

	ctxName := funcTraverser.FindArgumentNameByType(GinContextType, GinPackagePath, GinContextIsPointer)
	if ctxName == "" {
		return errors.New("failed to find context variable name")
	}

	var err error
	// Loop over every statement in the function
	ast.Inspect(funcTraverser.Node.Body, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		// If a function is called
		var callExpr *astTraversal.CallExpressionTraverser
		callExpr, err = traverser.CallExpression(n)
		if errors.Is(err, astTraversal.ErrInvalidNodeType) {
			err = nil
			return true
		} else if err != nil {
			return true
		}

		// If the function takes the context as any argument, traverse it
		_, ok := callExpr.ArgIndex(ctxName)
		if ok {
			var function *astTraversal.FunctionTraverser
			function, err = callExpr.Function()
			if err != nil {
				traverser.Log.Error().Err(err).Msg("failed to get function")
				return false
			}

			err = parseFunction(s, function, currRoute, function.Traverser.ActiveFile(), level+1)
			if err != nil {
				traverser.Log.Error().Err(err).Msg("error parsing function")
				return false
			}

			traverser.SetActiveFile(activeFile)
		} else {
			var funcType *types.Func
			funcType, err = callExpr.Type()
			if err != nil {
				return false
			}

			signature := funcType.Type().(*types.Signature)

			signaturePath := GinPackagePath + "." + GinContextType
			if GinContextIsPointer {
				signaturePath = "*" + signaturePath
			}

			if signature.Recv() != nil && signature.Recv().Type().String() == signaturePath {
				switch funcType.Name() {
				case "JSON":
					fallthrough
				case "XML":
					fallthrough
				case "YAML":
					fallthrough
				case "ProtoBuf":
					fallthrough
				case "Data":
					switch funcType.Name() {
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
					statusCode, err = traverser.ExtractStatusCode(callExpr.Args()[0])
					if err != nil {
						return true
					}

					argNo := 1
					if funcType.Name() == "Data" {
						argNo = 2
					}

					var exprType types.Type
					exprType, err = traverser.Expression(callExpr.Args()[argNo]).Type()
					if err != nil {
						traverser.Log.Error().Err(err).Msg("failed to get type for expression")
						return false
					}

					var result astTraversal.Result
					result, err = traverser.Type(exprType, traverser.ActiveFile().Package).Result()

					returnType := gengo.ReturnType{
						StatusCode: statusCode,
						Field:      parseResultToField(result),
					}

					currRoute.ReturnTypes = utils.AddReturnType(currRoute.ReturnTypes, returnType)

				case "String": // c.String
					currRoute.ContentType = "text/plain"

					var statusCode int
					statusCode, err = traverser.ExtractStatusCode(callExpr.Args()[0])
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

				case "Status": // c.Status
					var statusCode int
					statusCode, err = traverser.ExtractStatusCode(callExpr.Args()[0])
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

				// Query Param methods
				case "GetQuery":
					fallthrough
				case "Query":
					var queryParam gengo.Param
					queryParam, err = extractSingleRequestParam(traverser, callExpr.Args()[0], gengo.Param{})
					if err != nil {
						return false
					}

					currRoute.QueryParams = append(currRoute.QueryParams, queryParam)

				case "GetQueryArray":
					fallthrough
				case "QueryArray":
					var queryParam gengo.Param
					queryParam, err = extractSingleRequestParam(traverser, callExpr.Args()[0], gengo.Param{
						IsArray: true,
					})
					if err != nil {
						return false
					}

					currRoute.QueryParams = append(currRoute.QueryParams, queryParam)

				case "GetQueryMap":
					fallthrough
				case "QueryMap":
					var queryParam gengo.Param
					queryParam, err = extractSingleRequestParam(traverser, callExpr.Args()[0], gengo.Param{
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
					queryParam, err = extractBoundRequestParam(traverser, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.QueryParams = append(currRoute.QueryParams, queryParam)

				// Body Param methods
				case "ShouldBind":
					fallthrough
				case "Bind":
					var bodyParam gengo.Param
					bodyParam, err = extractBoundRequestParam(traverser, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.BodyType = "form"

					currRoute.Body = append(currRoute.Body, bodyParam)

				case "ShouldBindJSON":
					fallthrough
				case "BindJSON":
					var bodyParam gengo.Param
					bodyParam, err = extractBoundRequestParam(traverser, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/json"

					currRoute.Body = append(currRoute.Body, bodyParam)

				case "ShouldBindXML":
					fallthrough
				case "BindXML":
					var bodyParam gengo.Param
					bodyParam, err = extractBoundRequestParam(traverser, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/xml"

					currRoute.Body = append(currRoute.Body, bodyParam)

				case "ShouldBindYAML":
					fallthrough
				case "BindYAML":
					var bodyParam gengo.Param
					bodyParam, err = extractBoundRequestParam(traverser, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/yaml"

					currRoute.Body = append(currRoute.Body, bodyParam)

				case "GetPostForm":
					fallthrough
				case "PostForm":
					var bodyParam gengo.Param
					bodyParam, err = extractSingleRequestParam(traverser, callExpr.Args()[0], gengo.Param{})
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/x-www-form-urlencoded"
					currRoute.Body = append(currRoute.Body, bodyParam)

				case "GetPostFormArray":
					fallthrough
				case "PostFormArray":
					var bodyParam gengo.Param
					bodyParam, err = extractSingleRequestParam(traverser, callExpr.Args()[0], gengo.Param{
						IsArray: true,
					})
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/x-www-form-urlencoded"
					currRoute.Body = append(currRoute.Body, bodyParam)

				case "GetPostFormMap":
					fallthrough
				case "PostFormMap":
					var bodyParam gengo.Param
					bodyParam, err = extractSingleRequestParam(traverser, callExpr.Args()[0], gengo.Param{
						IsMap: true,
					})
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/x-www-form-urlencoded"
					currRoute.Body = append(currRoute.Body, bodyParam)
				}
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

func extractSingleRequestParam(traverser *astTraversal.BaseTraverser, node ast.Node, baseParam gengo.Param) (gengo.Param, error) {
	expr := traverser.Expression(node)

	name, err := expr.Value()
	if err != nil {
		traverser.Log.Error().Err(err).Msg("failed to parse expression")
		return gengo.Param{}, err
	}

	exprType, err := expr.Type()
	if err != nil {
		traverser.Log.Error().Err(err).Msg("failed to parse expression type")
		return gengo.Param{}, err
	}

	return gengo.Param{
		Name: name,
		Field: gengo.Field{
			Type: exprType.String(),
		},
		IsArray:    baseParam.IsArray,
		IsMap:      baseParam.IsMap,
		IsRequired: baseParam.IsRequired,
	}, nil
}

func extractBoundRequestParam(traverser *astTraversal.BaseTraverser, node ast.Node) (gengo.Param, error) {
	exprType, err := traverser.Expression(node).Type()
	if err != nil {
		return gengo.Param{}, err
	}

	result, err := traverser.Type(exprType, traverser.ActiveFile().Package).Result()

	bodyParam := gengo.Param{
		IsBound: true,
		Field:   parseResultToField(result),
	}

	return bodyParam, nil
}

func parseResultToField(result astTraversal.Result) gengo.Field {
	field := gengo.Field{
		Type:         result.Type,
		Name:         result.Name,
		IsRequired:   result.IsRequired,
		IsEmbedded:   result.IsEmbedded,
		SliceType:    result.SliceType,
		ArrayType:    result.ArrayType,
		ArrayLength:  result.ArrayLength,
		MapKeyType:   result.MapKeyType,
		MapValueType: result.MapValueType,
	}

	// If the type is not a primitive type, we need to get the package path
	// If the type is named, it is referring to a type
	// If the slice type is populated and not a primitive type, we need to get the package path for the slice
	// If the array type is populated and not a primitive type, we need to get the package path for the array
	// If the map value type is populated and not a primitive type, we need to get the package path for the map value
	if !gengo.IsAcceptedType(result.Type) || result.Name != "" ||
		(result.SliceType != "" && !gengo.IsAcceptedType(result.SliceType)) ||
		(result.ArrayType != "" && !gengo.IsAcceptedType(result.ArrayType)) ||
		(result.MapValueType != "" && !gengo.IsAcceptedType(result.MapValueType)) {
		field.Package = result.Package.Path()
	}

	// If the map key type is populated and not a primitive type, we need to get the package path for the map key
	if result.MapKeyType != "" && !gengo.IsAcceptedType(result.MapKeyType) {
		field.MapKeyPackage = result.MapKeyPackage.Path()
	}

	// If the struct fields are populated, we need to parse them
	if result.StructFields != nil {
		field.StructFields = make(map[string]gengo.Field)
		for name, value := range result.StructFields {
			field.StructFields[name] = parseResultToField(value)
		}
	}

	return field
}

func addComponent(s *gengo.Service) func(astTraversal.Result) error {
	return func(result astTraversal.Result) error {
		field := parseResultToField(result)

		if field.Package != "" {
			s.Components = utils.AddComponent(s.Components, field)
		}
		return nil
	}
}
