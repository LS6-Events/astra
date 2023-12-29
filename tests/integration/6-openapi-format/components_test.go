package petstore

import (
	"github.com/ls6-events/astra/tests/integration/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFormatSchemas(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := setupRouter()

	testAstra, err := helpers.SetupTestAstraWithDefaultConfig(t, r)
	require.NoError(t, err)

	schemas := testAstra.Path("components.schemas")

	properties := schemas.Search("6-openapi-format.TestStructFormatter", "properties")

	t.Run("String", func(t *testing.T) {
		require.Equal(t, "string", properties.Path("string.type").Data().(string))
	})

	t.Run("Int", func(t *testing.T) {
		require.Equal(t, "integer", properties.Path("int.type").Data().(string))
		require.Equal(t, "int32", properties.Path("int.format").Data().(string))
	})

	t.Run("Int8", func(t *testing.T) {
		require.Equal(t, "integer", properties.Path("int8.type").Data().(string))
		require.Equal(t, "int8", properties.Path("int8.format").Data().(string))
	})

	t.Run("Int16", func(t *testing.T) {
		require.Equal(t, "integer", properties.Path("int16.type").Data().(string))
		require.Equal(t, "int16", properties.Path("int16.format").Data().(string))
	})

	t.Run("Int32", func(t *testing.T) {
		require.Equal(t, "integer", properties.Path("int32.type").Data().(string))
		require.Equal(t, "int32", properties.Path("int32.format").Data().(string))
	})

	t.Run("Int64", func(t *testing.T) {
		require.Equal(t, "integer", properties.Path("int64.type").Data().(string))
		require.Equal(t, "int64", properties.Path("int64.format").Data().(string))
	})

	t.Run("Uint", func(t *testing.T) {
		require.Equal(t, "integer", properties.Path("uint.type").Data().(string))
		require.Equal(t, "uint", properties.Path("uint.format").Data().(string))
	})

	t.Run("Uint8", func(t *testing.T) {
		require.Equal(t, "integer", properties.Path("uint8.type").Data().(string))
		require.Equal(t, "uint8", properties.Path("uint8.format").Data().(string))
	})

	t.Run("Uint16", func(t *testing.T) {
		require.Equal(t, "integer", properties.Path("uint16.type").Data().(string))
		require.Equal(t, "uint16", properties.Path("uint16.format").Data().(string))
	})

	t.Run("Uint32", func(t *testing.T) {
		require.Equal(t, "integer", properties.Path("uint32.type").Data().(string))
		require.Equal(t, "uint32", properties.Path("uint32.format").Data().(string))
	})

	t.Run("Uint64", func(t *testing.T) {
		require.Equal(t, "integer", properties.Path("uint64.type").Data().(string))
		require.Equal(t, "uint64", properties.Path("uint64.format").Data().(string))
	})

	t.Run("Float32", func(t *testing.T) {
		require.Equal(t, "number", properties.Path("float32.type").Data().(string))
		require.Equal(t, "float32", properties.Path("float32.format").Data().(string))
	})

	t.Run("Float64", func(t *testing.T) {
		require.Equal(t, "number", properties.Path("float64.type").Data().(string))
		require.Equal(t, "float64", properties.Path("float64.format").Data().(string))
	})

	t.Run("Bool", func(t *testing.T) {
		require.Equal(t, "boolean", properties.Path("bool.type").Data().(string))
	})

	t.Run("Byte", func(t *testing.T) {
		require.Equal(t, "string", properties.Path("byte.type").Data().(string))
		require.Equal(t, "byte", properties.Path("byte.format").Data().(string))
	})

	t.Run("Rune", func(t *testing.T) {
		require.Equal(t, "string", properties.Path("rune.type").Data().(string))
		require.Equal(t, "rune", properties.Path("rune.format").Data().(string))
	})

	t.Run("Struct", func(t *testing.T) {
		require.Equal(t, "object", properties.Path("struct.type").Data().(string))
	})

	t.Run("Map", func(t *testing.T) {
		require.Equal(t, "object", properties.Path("map.type").Data().(string))
	})

	t.Run("Slice", func(t *testing.T) {
		require.Equal(t, "array", properties.Path("slice.type").Data().(string))
	})

	t.Run("Any", func(t *testing.T) {
		require.False(t, properties.Exists("any", "type"))
	})

	t.Run("time.Time", func(t *testing.T) {
		require.Equal(t, "#/components/schemas/time.Time", properties.Path("time.$ref").Data().(string))
		require.Equal(t, "string", schemas.Search("time.Time", "type").Data().(string))
		require.Equal(t, "date-time", schemas.Search("time.Time", "format").Data().(string))
	})

	t.Run("github.com/google/uuid.UUID", func(t *testing.T) {
		require.Equal(t, "#/components/schemas/uuid.UUID", properties.Path("uuid.$ref").Data().(string))
		require.Equal(t, "string", schemas.Search("uuid.UUID", "type").Data().(string))
		require.Equal(t, "uuid", schemas.Search("uuid.UUID", "format").Data().(string))
	})
}
