package gin

import (
	"errors"
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/astTraversal"
	"github.com/ls6-events/astra/utils"
	"go/ast"
	"go/types"
	"strings"
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
func parseFunction(s *astra.Service, funcTraverser *astTraversal.FunctionTraverser, currRoute *astra.Route, activeFile *astTraversal.FileNode, level int) error {
	traverser := funcTraverser.Traverser

	traverser.SetActiveFile(activeFile)
	traverser.SetAddComponentFunction(addComponent(s))

	if level == 0 {
		funcDoc, err := funcTraverser.GoDoc()
		if err != nil {
			return err
		}
		if funcDoc != nil {
			currRoute.Doc = strings.TrimSpace(funcDoc.Doc)
		}
	}

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

					returnType := astra.ReturnType{
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

					returnType := astra.ReturnType{
						StatusCode: statusCode,
						Field: astra.Field{
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

					returnType := astra.ReturnType{
						StatusCode: statusCode,
						Field: astra.Field{
							Type: "nil",
						},
					}
					currRoute.ReturnTypes = utils.AddReturnType(currRoute.ReturnTypes, returnType)

				// Query Param methods
				case "GetQuery":
					fallthrough
				case "Query":
					var queryParam astra.Param
					queryParam, err = extractSingleRequestParam(traverser, callExpr.Args()[0], astra.Param{})
					if err != nil {
						return false
					}

					currRoute.QueryParams = append(currRoute.QueryParams, queryParam)

				case "GetQueryArray":
					fallthrough
				case "QueryArray":
					var queryParam astra.Param
					queryParam, err = extractSingleRequestParam(traverser, callExpr.Args()[0], astra.Param{
						IsArray: true,
					})
					if err != nil {
						return false
					}

					currRoute.QueryParams = append(currRoute.QueryParams, queryParam)

				case "GetQueryMap":
					fallthrough
				case "QueryMap":
					var queryParam astra.Param
					queryParam, err = extractSingleRequestParam(traverser, callExpr.Args()[0], astra.Param{
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
					var queryParam astra.Param
					queryParam, err = extractBoundRequestParam(traverser, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.QueryParams = append(currRoute.QueryParams, queryParam)

				// Body Param methods
				case "ShouldBind":
					fallthrough
				case "Bind":
					var bodyParam astra.Param
					bodyParam, err = extractBoundRequestParam(traverser, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.BodyType = "form"

					currRoute.Body = append(currRoute.Body, bodyParam)

				case "ShouldBindJSON":
					fallthrough
				case "BindJSON":
					var bodyParam astra.Param
					bodyParam, err = extractBoundRequestParam(traverser, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/json"

					currRoute.Body = append(currRoute.Body, bodyParam)

				case "ShouldBindXML":
					fallthrough
				case "BindXML":
					var bodyParam astra.Param
					bodyParam, err = extractBoundRequestParam(traverser, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/xml"

					currRoute.Body = append(currRoute.Body, bodyParam)

				case "ShouldBindYAML":
					fallthrough
				case "BindYAML":
					var bodyParam astra.Param
					bodyParam, err = extractBoundRequestParam(traverser, callExpr.Args()[0])
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/yaml"

					currRoute.Body = append(currRoute.Body, bodyParam)

				case "GetPostForm":
					fallthrough
				case "PostForm":
					var bodyParam astra.Param
					bodyParam, err = extractSingleRequestParam(traverser, callExpr.Args()[0], astra.Param{})
					if err != nil {
						return false
					}

					currRoute.BodyType = "application/x-www-form-urlencoded"
					currRoute.Body = append(currRoute.Body, bodyParam)

				case "GetPostFormArray":
					fallthrough
				case "PostFormArray":
					var bodyParam astra.Param
					bodyParam, err = extractSingleRequestParam(traverser, callExpr.Args()[0], astra.Param{
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
					var bodyParam astra.Param
					bodyParam, err = extractSingleRequestParam(traverser, callExpr.Args()[0], astra.Param{
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

func extractSingleRequestParam(traverser *astTraversal.BaseTraverser, node ast.Node, baseParam astra.Param) (astra.Param, error) {
	expr := traverser.Expression(node)

	name, err := expr.Value()
	if err != nil {
		traverser.Log.Error().Err(err).Msg("failed to parse expression")
		return astra.Param{}, err
	}

	exprType, err := expr.Type()
	if err != nil {
		traverser.Log.Error().Err(err).Msg("failed to parse expression type")
		return astra.Param{}, err
	}

	return astra.Param{
		Name: name,
		Field: astra.Field{
			Type: exprType.String(),
		},
		IsArray:    baseParam.IsArray,
		IsMap:      baseParam.IsMap,
		IsRequired: baseParam.IsRequired,
	}, nil
}

func extractBoundRequestParam(traverser *astTraversal.BaseTraverser, node ast.Node) (astra.Param, error) {
	exprType, err := traverser.Expression(node).Type()
	if err != nil {
		return astra.Param{}, err
	}

	result, err := traverser.Type(exprType, traverser.ActiveFile().Package).Result()
	if err != nil {
		return astra.Param{}, err
	}

	bodyParam := astra.Param{
		IsBound: true,
		Field:   parseResultToField(result),
	}

	return bodyParam, nil
}

func parseResultToField(result astTraversal.Result) astra.Field {
	field := astra.Field{
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

	// If the godoc is populated, we need to parse the response
	if result.Doc != nil {
		field.Doc = strings.TrimSpace(result.Doc.Doc)
	}

	// If the type is not a primitive type, we need to get the package path
	// If the type is named, it is referring to a type
	// If the slice type is populated and not a primitive type, we need to get the package path for the slice
	// If the array type is populated and not a primitive type, we need to get the package path for the array
	// If the map value type is populated and not a primitive type, we need to get the package path for the map value
	if !astra.IsAcceptedType(result.Type) || result.Name != "" ||
		(result.SliceType != "" && !astra.IsAcceptedType(result.SliceType)) ||
		(result.ArrayType != "" && !astra.IsAcceptedType(result.ArrayType)) ||
		(result.MapValueType != "" && !astra.IsAcceptedType(result.MapValueType)) {
		field.Package = result.Package.Path()
	}

	// If the map key type is populated and not a primitive type, we need to get the package path for the map key
	if result.MapKeyType != "" && !astra.IsAcceptedType(result.MapKeyType) {
		field.MapKeyPackage = result.MapKeyPackage.Path()
	}

	// If the struct fields are populated, we need to parse them
	if result.StructFields != nil {
		field.StructFields = make(map[string]astra.Field)
		for name, value := range result.StructFields {

			// Check if the result's Doc and Decl are populated.
			// If they are, iterate over each spec in the Decl's specs.
			// If the spec is a TypeSpec and its Type is a StructType,
			// iterate over each field in the StructType's fields.
			// If the field's name matches the given name and its Doc is not empty,
			// store the trimmed Doc text in the 'doc' variable and break out of the loops.
			var doc string
			if result.Doc != nil && result.Doc.Decl != nil {
				for _, spec := range result.Doc.Decl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							for _, structField := range structType.Fields.List {
								var fieldName string
								var isShown bool
								if structField.Tag == nil && len(structField.Names) == 0 {
									isShown = false
								} else if structField.Tag == nil {
									fieldName = structField.Names[0].Name
									isShown = true
								} else {
									fieldName, _, isShown = astTraversal.ParseStructTag(strings.Trim(structField.Tag.Value, "`"))
								}

								if !isShown {
									continue
								}

								if fieldName == name {
									fieldDoc := strings.TrimSpace(structField.Doc.Text())
									if fieldDoc != "" {
										doc = fieldDoc
										break
									}
								}
							}

							if doc != "" {
								break
							}
						}
					}
				}
			}

			structField := parseResultToField(value)
			structField.Doc = doc
			field.StructFields[name] = structField
		}
	}

	return field
}

func addComponent(s *astra.Service) func(astTraversal.Result) error {
	return func(result astTraversal.Result) error {
		field := parseResultToField(result)

		if field.Package != "" {
			s.Components = utils.AddComponent(s.Components, field)
		}
		return nil
	}
}
