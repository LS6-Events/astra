package astra

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestWithCustomTypeMapping(t *testing.T) {
	t.Run("adds custom type mapping", func(t *testing.T) {
		service := &Service{
			CustomTypeMapping: make(map[string]TypeFormat),
		}

		WithCustomTypeMapping(map[string]TypeFormat{
			"old": {
				Type:   "string",
				Format: "new",
			},
		})(service)

		require.Equal(t, TypeFormat{
			Type:   "string",
			Format: "new",
		}, service.CustomTypeMapping["old"])
	})

	t.Run("overwrites existing type mapping", func(t *testing.T) {
		service := &Service{
			CustomTypeMapping: map[string]TypeFormat{
				"old": {
					Type:   "string",
					Format: "old",
				},
			},
		}

		WithCustomTypeMapping(map[string]TypeFormat{
			"old": {
				Type:   "string",
				Format: "new",
			},
		})(service)

		require.Equal(t, TypeFormat{
			Type:   "string",
			Format: "new",
		}, service.CustomTypeMapping["old"])
	})

	t.Run("overwrites existing type mapping set from function", func(t *testing.T) {
		service := &Service{
			CustomTypeMapping: make(map[string]TypeFormat),
		}

		WithCustomTypeMapping(map[string]TypeFormat{
			"old": {
				Type:   "string",
				Format: "old",
			},
		})(service)

		require.Equal(t, TypeFormat{
			Type:   "string",
			Format: "old",
		}, service.CustomTypeMapping["old"])

		WithCustomTypeMapping(map[string]TypeFormat{
			"old": {
				Type:   "string",
				Format: "new",
			},
		})(service)

		require.Equal(t, TypeFormat{
			Type:   "string",
			Format: "new",
		}, service.CustomTypeMapping["old"])
	})
}

func TestWithCustomTypeMappingSingle(t *testing.T) {
	t.Run("adds custom type mapping", func(t *testing.T) {
		service := &Service{
			CustomTypeMapping: make(map[string]TypeFormat),
		}

		WithCustomTypeMappingSingle("old", "string", "new")(service)

		require.Equal(t, TypeFormat{
			Type:   "string",
			Format: "new",
		}, service.CustomTypeMapping["old"])
	})

	t.Run("overwrites existing type mapping", func(t *testing.T) {
		service := &Service{
			CustomTypeMapping: map[string]TypeFormat{
				"old": {
					Type:   "string",
					Format: "old",
				},
			},
		}

		WithCustomTypeMappingSingle("old", "string", "new")(service)

		require.Equal(t, TypeFormat{
			Type:   "string",
			Format: "new",
		}, service.CustomTypeMapping["old"])
	})
}

func TestService_GetTypeMapping(t *testing.T) {
	t.Run("sets full type mapping if it doesn't exist initially", func(t *testing.T) {
		service := &Service{
			CustomTypeMapping: make(map[string]TypeFormat),
		}

		require.Empty(t, service.fullTypeMapping)

		service.GetTypeMapping("test", "test")

		require.NotEmpty(t, service.fullTypeMapping)
	})

	t.Run("returns the type mapping for the given custom key", func(t *testing.T) {
		service := &Service{
			fullTypeMapping: map[string]TypeFormat{
				"test.test": {
					Type:   "string",
					Format: "test",
				},
			},
		}

		result, ok := service.GetTypeMapping("test", "test")

		require.True(t, ok)
		require.Equal(t, TypeFormat{
			Type:   "string",
			Format: "test",
		}, result)
	})

	t.Run("returns predefined types for type maps", func(t *testing.T) {
		service := &Service{}

		for k, v := range PredefinedTypeMap {
			t.Run(k, func(t *testing.T) {
				splitKey := strings.Split(k, ".")
				var result TypeFormat
				var ok bool
				if len(splitKey) > 1 {
					result, ok = service.GetTypeMapping(splitKey[len(splitKey)-1], strings.Join(splitKey[0:len(splitKey)-1], "."))
				} else {
					result, ok = service.GetTypeMapping(splitKey[0], "")
				}
				require.True(t, ok)
				require.Equal(t, v, result)
			})
		}
	})
}
