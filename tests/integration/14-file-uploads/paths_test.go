package integration

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

	require.Equal(t, "string", testAstra.Search("paths", "/upload", "post", "requestBody", "content", "multipart/form-data", "schema", "properties", "file", "type").Data().(string))
	require.Equal(t, "binary", testAstra.Search("paths", "/upload", "post", "requestBody", "content", "multipart/form-data", "schema", "properties", "file", "format").Data().(string))
}
