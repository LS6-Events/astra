package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/astTraversal"
	"github.com/ls6-events/astra/inputs"
	"github.com/ls6-events/astra/outputs"
)

func main() {
	r := gin.Default()

	r.GET("/posts", GetPosts)
	r.GET("/posts/:id", GetPost)
	r.POST("/posts", CreatePost)
	r.PUT("/posts/:id", UpdatePost)
	r.DELETE("/posts/:id", DeletePost)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	gen := astra.New(
		inputs.WithGinInput(r),
		outputs.WithOpenAPIOutput("openapi.generated.yaml"),
		astra.WithCustomFunc(func(contextVarName string, contextFuncBuilder *astra.ContextFuncBuilder) (*astra.Route, error) {
			funcType, err := contextFuncBuilder.Traverser.Type()
			if err != nil {
				return nil, err
			}

			// handleError(c, statusCode, err)
			if funcType.Name() == "handleError" {
				// Because we know explicitly that the first argument is the context, we can ignore it
				// We only need to concern ourselves with the status code, which is a unique case
				// We can also ignore the error, as it will be parsed as a string
				return contextFuncBuilder.Ignored().StatusCode().ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
					// Params is a list of the arguments returned from the function
					// [0] is ignored
					// [1] is the status code (int)
					// [2] is the error (type)
					statusCode := params[1].(int)

					// Create the return type for this explicit error code
					returnType := astra.ReturnType{
						ContentType: "text/plain",
						StatusCode:  statusCode,
						Field: astra.Field{
							Type: "string",
						},
					}

					// A custom utility function that prevents duplicate return types
					route.ReturnTypes = astra.AddReturnType(route.ReturnTypes, returnType)
					return route, nil
				})
			}

			// handleSuccess(c, statusCode, data)
			if funcType.Name() == "handleSuccess" {
				// Because we know explicitly that the first argument is the context, we can ignore it
				// We only need to concern ourselves with the status code, which is a unique case
				// We can also ignore the error, as it will be parsed as a string
				return contextFuncBuilder.Ignored().StatusCode().ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
					// Params is a list of the arguments returned from the function
					// [0] is ignored
					// [1] is the status code (int)
					// [2] is the data (type), this comes as a result type from the astTraversal package
					statusCode := params[1].(int)
					dataType := params[2].(astTraversal.Result)

					returnType := astra.ReturnType{
						ContentType: "application/json",
						StatusCode:  statusCode,
						Field:       astra.ParseResultToField(dataType),
					}

					route.ReturnTypes = astra.AddReturnType(route.ReturnTypes, returnType)

					return route, nil
				})
			}

			// When the function is not one of our custom functions, we return nil, nil
			// This will tell Astra to continue with the default behaviour for this function
			return nil, nil
		}),
	)

	config := astra.Config{
		Title:   "Example API with Custom Functions",
		Version: "1.0.0",
		Host:    "localhost",
		Port:    8000,
	}

	gen.SetConfig(&config)

	err := gen.Parse()
	if err != nil {
		panic(err)
	}

	err = r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
