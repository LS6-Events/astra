package gin

import (
	"errors"
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/astTraversal"
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
		funcDoc, err := funcTraverser.Doc()
		if err != nil {
			return err
		}
		if funcDoc != "" {
			currRoute.Doc = strings.TrimSpace(funcDoc)
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

		funcBuilder := astra.NewContextFuncBuilder(currRoute, callExpr)

		// Loop over every custom function
		// If the custom function returns a route, use that route instead of the current route
		// And break out of this AST traversal for this call expression
		// Otherwise, continue on
		var shouldBreak bool
		for _, customFunc := range s.CustomFuncs {
			var newRoute *astra.Route
			newRoute, err = customFunc(ctxName, funcBuilder)
			if err != nil {
				return false
			}
			if newRoute != nil {
				currRoute = newRoute
				shouldBreak = true
				break
			}
		}
		if shouldBreak {
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
					currRoute, err = funcBuilder.StatusCode().ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						statusCode := params[0].(int)
						result := params[1].(astTraversal.Result)

						returnType := astra.ReturnType{
							StatusCode:  statusCode,
							ContentType: "application/json",
							Field:       astra.ParseResultToField(result),
						}

						route.ReturnTypes = astra.AddReturnType(route.ReturnTypes, returnType)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "XML":
					currRoute, err = funcBuilder.StatusCode().ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						statusCode := params[0].(int)
						result := params[1].(astTraversal.Result)

						returnType := astra.ReturnType{
							StatusCode:  statusCode,
							ContentType: "application/xml",
							Field:       astra.ParseResultToField(result),
						}

						route.ReturnTypes = astra.AddReturnType(route.ReturnTypes, returnType)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "YAML":
					currRoute, err = funcBuilder.StatusCode().ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						statusCode := params[0].(int)
						result := params[1].(astTraversal.Result)

						returnType := astra.ReturnType{
							StatusCode:  statusCode,
							ContentType: "application/yaml",
							Field:       astra.ParseResultToField(result),
						}

						route.ReturnTypes = astra.AddReturnType(route.ReturnTypes, returnType)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "ProtoBuf":
					currRoute, err = funcBuilder.StatusCode().ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						statusCode := params[0].(int)
						result := params[1].(astTraversal.Result)

						returnType := astra.ReturnType{
							StatusCode:  statusCode,
							ContentType: "application/protobuf",
							Field:       astra.ParseResultToField(result),
						}

						route.ReturnTypes = astra.AddReturnType(route.ReturnTypes, returnType)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "Data":
					currRoute, err = funcBuilder.StatusCode().Ignored().ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						statusCode := params[0].(int)
						result := params[2].(astTraversal.Result)

						returnType := astra.ReturnType{
							StatusCode: statusCode,
							Field:      astra.ParseResultToField(result),
						}

						route.ReturnTypes = astra.AddReturnType(route.ReturnTypes, returnType)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "String": // c.String
					currRoute, err = funcBuilder.StatusCode().Ignored().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						statusCode := params[0].(int)

						returnType := astra.ReturnType{
							StatusCode:  statusCode,
							ContentType: "text/plain",
							Field: astra.Field{
								Type: "string",
							},
						}

						route.ReturnTypes = astra.AddReturnType(route.ReturnTypes, returnType)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "Status": // c.Status
					currRoute, err = funcBuilder.StatusCode().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						statusCode := params[0].(int)

						returnType := astra.ReturnType{
							StatusCode: statusCode,
							Field: astra.Field{
								Type: "nil",
							},
						}

						route.ReturnTypes = astra.AddReturnType(route.ReturnTypes, returnType)

						return route, nil
					})
				// Query Param methods
				case "GetQuery", "Query":
					currRoute, err = funcBuilder.Value().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						name := params[0].(string)

						param := astra.Param{
							Field: astra.Field{
								Type: "string",
							},
							Name: name,
						}

						route.QueryParams = append(route.QueryParams, param)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "GetQueryArray", "QueryArray":
					currRoute, err = funcBuilder.Value().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						name := params[0].(string)

						param := astra.Param{
							Field: astra.Field{
								Type: "string",
							},
							Name:    name,
							IsArray: true,
						}

						route.QueryParams = append(route.QueryParams, param)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "GetQueryMap", "QueryMap":
					currRoute, err = funcBuilder.Value().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						name := params[0].(string)

						param := astra.Param{
							Field: astra.Field{
								Type: "string",
							},
							Name:  name,
							IsMap: true,
						}

						route.QueryParams = append(route.QueryParams, param)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "ShouldBindQuery", "BindQuery":
					currRoute, err = funcBuilder.ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						field := astra.ParseResultToField(params[0].(astTraversal.Result))

						route.QueryParams = append(route.QueryParams, astra.Param{
							IsBound: true,
							Field:   field,
						})

						return route, nil
					})
					if err != nil {
						return false
					}

				// Body Param methods
				case "ShouldBind", "Bind":
					currRoute, err = funcBuilder.ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						field := astra.ParseResultToField(params[0].(astTraversal.Result))

						route.PathParams = append(route.PathParams, astra.Param{
							IsBound: true,
							Field:   field,
						})

						route.QueryParams = append(route.QueryParams, astra.Param{
							IsBound: true,
							Field:   field,
						})

						route.RequestHeaders = append(route.RequestHeaders, astra.Param{
							IsBound: true,
							Field:   field,
						})

						for _, bodyBindingTag := range []astTraversal.BindingTagType{astTraversal.FormBindingTag, astTraversal.JSONBindingTag, astTraversal.XMLBindingTag, astTraversal.YAMLBindingTag} {
							contentTypes := astra.BindingTagToContentTypes(bodyBindingTag)

							for _, contentType := range contentTypes {
								route.Body = append(route.Body, astra.BodyParam{
									ContentType: contentType,
									IsBound:     true,
									Field:       field,
								})
							}
						}

						return route, nil
					})
					if err != nil {
						return false
					}
				case "ShouldBindJSON", "BindJSON":
					currRoute, err = funcBuilder.ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						field := astra.ParseResultToField(params[0].(astTraversal.Result))

						route.Body = append(route.Body, astra.BodyParam{
							ContentType: "application/json",
							IsBound:     true,
							Field:       field,
						})

						return route, nil
					})
					if err != nil {
						return false
					}
				case "ShouldBindXML", "BindXML":
					currRoute, err = funcBuilder.ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						field := astra.ParseResultToField(params[0].(astTraversal.Result))

						route.Body = append(route.Body, astra.BodyParam{
							ContentType: "application/xml",
							IsBound:     true,
							Field:       field,
						})

						return route, nil
					})
					if err != nil {
						return false
					}
				case "ShouldBindYAML", "BindYAML":
					currRoute, err = funcBuilder.ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						field := astra.ParseResultToField(params[0].(astTraversal.Result))

						route.Body = append(route.Body, astra.BodyParam{
							ContentType: "application/yaml",
							IsBound:     true,
							Field:       field,
						})

						return route, nil
					})
					if err != nil {
						return false
					}
				case "GetPostForm", "PostForm":
					currRoute, err = funcBuilder.Value().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						name := params[0].(string)

						param := astra.BodyParam{
							ContentType: "application/x-www-form-urlencoded",
							Field: astra.Field{
								Type: "string",
							},
							Name: name,
						}

						route.Body = append(route.Body, param)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "GetPostFormArray", "PostFormArray":
					currRoute, err = funcBuilder.Value().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						name := params[0].(string)

						param := astra.BodyParam{
							ContentType: "application/x-www-form-urlencoded",
							Field: astra.Field{
								Type: "string",
							},
							Name:    name,
							IsArray: true,
						}

						route.Body = append(route.Body, param)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "GetPostFormMap", "PostFormMap":
					currRoute, err = funcBuilder.Value().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						name := params[0].(string)

						param := astra.BodyParam{
							ContentType: "application/x-www-form-urlencoded",
							Field: astra.Field{
								Type: "string",
							},
							Name:  name,
							IsMap: true,
						}

						route.Body = append(route.Body, param)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "GetHeader":
					currRoute, err = funcBuilder.Value().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						name := params[0].(string)

						param := astra.Param{
							Field: astra.Field{
								Type: "string",
							},
							Name: name,
						}

						route.RequestHeaders = append(route.RequestHeaders, param)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "ShouldBindHeader", "BindHeader":
					currRoute, err = funcBuilder.ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						field := astra.ParseResultToField(params[0].(astTraversal.Result))

						route.RequestHeaders = append(route.RequestHeaders, astra.Param{
							IsBound: true,
							Field:   field,
						})

						return route, nil
					})
					if err != nil {
						return false
					}
				case "Header":
					currRoute, err = funcBuilder.Value().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						name := params[0].(string)

						param := astra.Param{
							Field: astra.Field{
								Type: "string",
							},
							Name: name,
						}

						route.ResponseHeaders = append(route.ResponseHeaders, param)

						return route, nil
					})
				case "AbortWithError":
					currRoute, err = funcBuilder.StatusCode().Ignored().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						statusCode := params[0].(int)

						returnType := astra.ReturnType{
							StatusCode: statusCode,
							Field: astra.Field{
								Type: "nil",
							},
						}

						route.ReturnTypes = astra.AddReturnType(route.ReturnTypes, returnType)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "AbortWithStatus":
					currRoute, err = funcBuilder.StatusCode().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						statusCode := params[0].(int)

						returnType := astra.ReturnType{
							StatusCode: statusCode,
							Field: astra.Field{
								Type: "nil",
							},
						}

						route.ReturnTypes = astra.AddReturnType(route.ReturnTypes, returnType)

						return route, nil
					})
					if err != nil {
						return false
					}
				case "AbortWithStatusJSON":
					currRoute, err = funcBuilder.StatusCode().ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
						statusCode := params[0].(int)
						result := params[1].(astTraversal.Result)

						returnType := astra.ReturnType{
							ContentType: "application/json",
							StatusCode:  statusCode,
							Field:       astra.ParseResultToField(result),
						}

						route.ReturnTypes = astra.AddReturnType(route.ReturnTypes, returnType)

						return route, nil
					})
					if err != nil {
						return false
					}
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

func addComponent(s *astra.Service) func(astTraversal.Result) error {
	return func(result astTraversal.Result) error {
		field := astra.ParseResultToField(result)

		if field.Package != "" {
			s.Components = astra.AddComponent(s.Components, field)
		}
		return nil
	}
}
