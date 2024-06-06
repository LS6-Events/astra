package petstore

import (
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMiddleware(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r, astra.UnstableWithMiddleware())
	require.NoError(t, err)

	paths := testAstra.Path("paths")

	// /pets middleware
	// GET /pets
	require.Equal(t, "string", paths.Search("/pets", "get", "parameters", "0", "schema", "type").Data().(string))
	require.Equal(t, "api_key", paths.Search("/pets", "get", "parameters", "0", "name").Data().(string))
	require.Equal(t, "query", paths.Search("/pets", "get", "parameters", "0", "in").Data().(string))

	// GET /pets/{id}
	require.Equal(t, "string", paths.Search("/pets/{id}", "get", "parameters", "0", "schema", "type").Data().(string))
	require.Equal(t, "api_key", paths.Search("/pets/{id}", "get", "parameters", "0", "name").Data().(string))
	require.Equal(t, "query", paths.Search("/pets/{id}", "get", "parameters", "0", "in").Data().(string))

	// POST /pets
	require.Equal(t, "string", paths.Search("/pets", "post", "parameters", "0", "schema", "type").Data().(string))
	require.Equal(t, "api_key", paths.Search("/pets", "post", "parameters", "0", "name").Data().(string))
	require.Equal(t, "query", paths.Search("/pets", "post", "parameters", "0", "in").Data().(string))

	// DELETE /pets/{id}
	require.Equal(t, "string", paths.Search("/pets/{id}", "delete", "parameters", "0", "schema", "type").Data().(string))
	require.Equal(t, "api_key", paths.Search("/pets/{id}", "delete", "parameters", "0", "name").Data().(string))
	require.Equal(t, "query", paths.Search("/pets/{id}", "delete", "parameters", "0", "in").Data().(string))

	// /no-middleware
	require.True(t, paths.Exists("/no-middleware", "get"))
	require.False(t, paths.Exists("/no-middleware", "parameters", "0"))

	// /middleware
	require.Equal(t, "string", paths.Search("/middleware", "get", "parameters", "0", "schema", "type").Data().(string))
	require.Equal(t, "Authorization", paths.Search("/middleware", "get", "parameters", "0", "name").Data().(string))
	require.Equal(t, "header", paths.Search("/middleware", "get", "parameters", "0", "in").Data().(string))

	// /wrapper-func-middleware
	// Higher order functions don't seem to play nicely with runtime.FuncForPC when not in JetBrains debugger (unknown reason, will require further investigation)
	//require.Equal(t, "string", paths.Search("/wrapper-func-middleware", "get", "parameters", "0", "schema", "type").Data().(string))
	//require.Equal(t, "inline_api_key", paths.Search("/wrapper-func-middleware", "get", "parameters", "0", "name").Data().(string))
	//require.Equal(t, "query", paths.Search("/wrapper-func-middleware", "get", "parameters", "0", "in").Data().(string))
}
