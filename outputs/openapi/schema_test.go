package openapi

import (
	"github.com/ls6-events/validjsonator"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetQueryParamStyle(t *testing.T) {
	t.Run("it returns the correct style for object types", func(t *testing.T) {
		style, explode := getQueryParamStyle(validjsonator.Schema{
			Type: "object",
		})

		require.Equal(t, "deepObject", style)
		require.True(t, explode)
	})

	t.Run("it returns the correct style for a primitives or array types", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			style, explode := getQueryParamStyle(validjsonator.Schema{
				Type: "string",
			})

			require.Equal(t, "form", style)
			require.False(t, explode)

		})

		t.Run("int", func(t *testing.T) {
			style, explode := getQueryParamStyle(validjsonator.Schema{
				Type: "int",
			})

			require.Equal(t, "form", style)
			require.False(t, explode)
		})

		t.Run("array", func(t *testing.T) {
			style, explode := getQueryParamStyle(validjsonator.Schema{
				Type: "array",
			})

			require.Equal(t, "form", style)
			require.False(t, explode)
		})
	})
}
