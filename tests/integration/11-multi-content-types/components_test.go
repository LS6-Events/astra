package petstore

import (
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMultiContentTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r)
	require.NoError(t, err)

	schemas := testAstra.Path("components.schemas")

	require.False(t, schemas.Exists("11-multi-content-types.MultiContentType"))

	// JSON
	require.True(t, schemas.Exists("11-multi-content-types.json.MultiContentType"))
	require.True(t, schemas.Exists("11-multi-content-types.json.MultiContentType", "properties", "json-test"))

	// XML
	require.True(t, schemas.Exists("11-multi-content-types.xml.MultiContentType"))
	require.True(t, schemas.Exists("11-multi-content-types.xml.MultiContentType", "properties", "xml-test"))

	// YAML
	require.True(t, schemas.Exists("11-multi-content-types.yaml.MultiContentType"))
	require.True(t, schemas.Exists("11-multi-content-types.yaml.MultiContentType", "properties", "yaml-test"))
}
