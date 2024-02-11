package integration

import (
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEnums(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r)
	require.NoError(t, err)

	schemas := testAstra.Path("components.schemas")

	t.Run("String Enum", func(t *testing.T) {
		stringEnum := schemas.Search("7-enums.TestStringEnum")

		require.Equal(t, "string", stringEnum.Path("type").Data().(string))
		require.Equal(t, []any{"available", "sold"}, stringEnum.Path("enum").Data().([]any))

		stringStruct := schemas.Search("7-enums.TestStructWithStringEnum")

		require.Equal(t, "#/components/schemas/7-enums.TestStringEnum", stringStruct.Path("properties.enum.$ref").Data().(string))
	})

	t.Run("Int Enum", func(t *testing.T) {
		intEnum := schemas.Search("7-enums.TestIntEnum")

		require.Equal(t, "integer", intEnum.Path("type").Data().(string))
		require.Equal(t, []any{1.0, 2.0}, intEnum.Path("enum").Data().([]any))

		intStruct := schemas.Search("7-enums.TestStructWithIntEnum")

		require.Equal(t, "#/components/schemas/7-enums.TestIntEnum", intStruct.Path("properties.enum.$ref").Data().(string))
	})
}
