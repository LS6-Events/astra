package astra

import (
	"github.com/ls6-events/astra/astTraversal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseResultToField(t *testing.T) {
	t.Run("parses results to field", func(t *testing.T) {
		result := astTraversal.Result{
			Type: "string",
			Name: "test",
		}

		field := ParseResultToField(result)

		require.Equal(t, "string", field.Type)
		require.Equal(t, "test", field.Name)
	})

	t.Run("trims whitespace from godoc", func(t *testing.T) {
		result := astTraversal.Result{
			Type: "string",
			Name: "test",
			Doc:  "  This is a test  ",
		}

		field := ParseResultToField(result)

		require.Equal(t, "This is a test", field.Doc)
	})

	t.Run("gets the package path for non-primitive types", func(t *testing.T) {
		makePackage := func(bottomName string) *astTraversal.PackageNode {
			return &astTraversal.PackageNode{
				Parent: &astTraversal.PackageNode{
					Parent: &astTraversal.PackageNode{
						Name: "github.com",
						Parent: &astTraversal.PackageNode{
							Parent: nil,
						},
					},
					Name: "ls6-events",
				},
				Name: bottomName,
			}
		}

		t.Run("non-primitive type and name is set", func(t *testing.T) {
			result := astTraversal.Result{
				Type:    "string",
				Name:    "test",
				Package: makePackage("astra"),
			}

			field := ParseResultToField(result)

			require.Equal(t, "github.com/ls6-events/astra", field.Package)
		})

		t.Run("non-primitive slice type", func(t *testing.T) {
			result := astTraversal.Result{
				Type:      "string",
				SliceType: "test",
				Package:   makePackage("slice"),
			}

			field := ParseResultToField(result)

			require.Equal(t, "github.com/ls6-events/slice", field.Package)
		})

		t.Run("non-primitive array type", func(t *testing.T) {
			result := astTraversal.Result{
				Type:      "string",
				ArrayType: "test",
				Package:   makePackage("array"),
			}

			field := ParseResultToField(result)

			require.Equal(t, "github.com/ls6-events/array", field.Package)
		})

		t.Run("non-primitive map value type", func(t *testing.T) {
			result := astTraversal.Result{
				Type:         "string",
				MapValueType: "test",
				Package:      makePackage("map"),
			}

			field := ParseResultToField(result)

			require.Equal(t, "github.com/ls6-events/map", field.Package)
		})

		t.Run("non-primitive map key type", func(t *testing.T) {
			result := astTraversal.Result{
				Type:          "string",
				MapKeyType:    "test",
				MapKeyPackage: makePackage("mapKey"),
			}

			field := ParseResultToField(result)

			require.Equal(t, "github.com/ls6-events/mapKey", field.MapKeyPackage)
		})
	})

	t.Run("parses struct fields", func(t *testing.T) {
		result := astTraversal.Result{
			Type: "struct",
			StructFields: map[string]astTraversal.Result{
				"test": {
					Type: "string",
					Name: "test",
				},
			},
		}

		field := ParseResultToField(result)

		require.Equal(t, "string", field.StructFields["test"].Type)
		require.Equal(t, "test", field.StructFields["test"].Name)
	})
}
