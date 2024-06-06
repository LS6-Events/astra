package integration

import (
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCustomFunction(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r, astra.WithCustomFunc(func(contextVarName string, contextFuncBuilder *astra.ContextFuncBuilder) (*astra.Route, error) {
		funcType, err := contextFuncBuilder.Traverser.Type()
		if err != nil {
			return nil, err
		}

		if funcType.Name() == "handleError" {
			return contextFuncBuilder.Ignored().StatusCode().ExpressionResult().Build(func(route *astra.Route, params []any) (*astra.Route, error) {
				statusCode := params[1].(int)

				returnType := astra.ReturnType{
					ContentType: "text/plain",
					StatusCode:  statusCode,
					Field: astra.Field{
						Type: "string",
					},
				}

				route.ReturnTypes = astra.AddReturnType(route.ReturnTypes, returnType)
				return route, nil
			})
		}
		return nil, nil
	}))
	require.NoError(t, err)

	paths := testAstra.Path("paths")

	// POST /pets
	require.Equal(t, "string", paths.Path("/pets.post.responses.400.content.text/plain.schema.type").Data().(string))

	// GET /pets/{id}
	require.Equal(t, "string", paths.Path("/pets/{id}.get.responses.400.content.text/plain.schema.type").Data().(string))

	// DELETE /pets/{id}
	require.Equal(t, "string", paths.Path("/pets/{id}.delete.responses.400.content.text/plain.schema.type").Data().(string))
}
