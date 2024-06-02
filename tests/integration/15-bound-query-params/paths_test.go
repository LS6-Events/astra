package petstore

import (
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPaths(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r)
	require.NoError(t, err)

	require.NotNil(t, testAstra)

	paths := testAstra.Search("paths")

	// Inline query params
	require.Equal(t, "decimal", paths.Search("/inlineQueryParams", "get", "parameters", "0", "name").Data().(string))
	require.Equal(t, "form", paths.Search("/inlineQueryParams", "get", "parameters", "0", "style").Data().(string))
	require.Nil(t, paths.Search("/inlineQueryParams", "get", "parameters", "0", "explode").Data())

	require.Equal(t, "limit", paths.Search("/inlineQueryParams", "get", "parameters", "1", "name").Data().(string))
	require.Equal(t, "form", paths.Search("/inlineQueryParams", "get", "parameters", "1", "style").Data().(string))
	require.Nil(t, paths.Search("/inlineQueryParams", "get", "parameters", "1", "explode").Data())

	require.Equal(t, "name", paths.Search("/inlineQueryParams", "get", "parameters", "2", "name").Data().(string))
	require.Equal(t, "form", paths.Search("/inlineQueryParams", "get", "parameters", "2", "style").Data().(string))
	require.Nil(t, paths.Search("/inlineQueryParams", "get", "parameters", "2", "explode").Data())

	require.Equal(t, "tag", paths.Search("/inlineQueryParams", "get", "parameters", "3", "name").Data().(string))
	require.Equal(t, "form", paths.Search("/inlineQueryParams", "get", "parameters", "3", "style").Data().(string))
	require.Nil(t, paths.Search("/inlineQueryParams", "get", "parameters", "3", "explode").Data())

	// Bound query params
	require.Equal(t, "decimal", paths.Search("/boundQueryParams", "get", "parameters", "0", "name").Data().(string))
	require.Equal(t, "form", paths.Search("/boundQueryParams", "get", "parameters", "0", "style").Data().(string))
	require.Nil(t, paths.Search("/boundQueryParams", "get", "parameters", "0", "explode").Data())
	require.Equal(t, "number", paths.Search("/boundQueryParams", "get", "parameters", "0", "schema", "type").Data().(string))
	require.Equal(t, "float64", paths.Search("/boundQueryParams", "get", "parameters", "0", "schema", "format").Data().(string))

	require.Equal(t, "limit", paths.Search("/boundQueryParams", "get", "parameters", "1", "name").Data().(string))
	require.Equal(t, "form", paths.Search("/boundQueryParams", "get", "parameters", "1", "style").Data().(string))
	require.Nil(t, paths.Search("/boundQueryParams", "get", "parameters", "1", "explode").Data())
	require.Equal(t, "integer", paths.Search("/boundQueryParams", "get", "parameters", "1", "schema", "type").Data().(string))
	require.Equal(t, "int32", paths.Search("/boundQueryParams", "get", "parameters", "1", "schema", "format").Data().(string))

	require.Equal(t, "map", paths.Search("/boundQueryParams", "get", "parameters", "2", "name").Data().(string))
	require.Equal(t, "deepObject", paths.Search("/boundQueryParams", "get", "parameters", "2", "style").Data().(string))
	require.Equal(t, true, paths.Search("/boundQueryParams", "get", "parameters", "2", "explode").Data().(bool))
	require.Equal(t, "object", paths.Search("/boundQueryParams", "get", "parameters", "2", "schema", "type").Data().(string))
	require.Equal(t, "string", paths.Search("/boundQueryParams", "get", "parameters", "2", "schema", "additionalProperties", "type").Data().(string))

	require.Equal(t, "name", paths.Search("/boundQueryParams", "get", "parameters", "3", "name").Data().(string))
	require.Equal(t, "form", paths.Search("/boundQueryParams", "get", "parameters", "3", "style").Data().(string))
	require.Nil(t, paths.Search("/boundQueryParams", "get", "parameters", "3", "explode").Data())
	require.Equal(t, "string", paths.Search("/boundQueryParams", "get", "parameters", "3", "schema", "type").Data().(string))

	require.Equal(t, "slice", paths.Search("/boundQueryParams", "get", "parameters", "4", "name").Data().(string))
	require.Equal(t, "form", paths.Search("/boundQueryParams", "get", "parameters", "4", "style").Data().(string))
	require.Nil(t, paths.Search("/boundQueryParams", "get", "parameters", "4", "explode").Data())
	require.Equal(t, "array", paths.Search("/boundQueryParams", "get", "parameters", "4", "schema", "type").Data().(string))
	require.Equal(t, "string", paths.Search("/boundQueryParams", "get", "parameters", "4", "schema", "items", "type").Data().(string))

	require.Equal(t, "tag", paths.Search("/boundQueryParams", "get", "parameters", "5", "name").Data().(string))
	require.Equal(t, "form", paths.Search("/boundQueryParams", "get", "parameters", "5", "style").Data().(string))
	require.Nil(t, paths.Search("/boundQueryParams", "get", "parameters", "5", "explode").Data())
	require.Equal(t, "string", paths.Search("/boundQueryParams", "get", "parameters", "5", "schema", "type").Data().(string))
}
