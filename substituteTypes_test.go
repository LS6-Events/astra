package astra

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestService_HandleSubstituteTypes(t *testing.T) {
	service := &Service{}

	for typeName, typeReplacement := range substituteTypes {
		t.Run(typeName, func(t *testing.T) {
			typeValue := &Field{
				Name: typeName,
			}

			if strings.Contains(typeName, ".") {
				splitTypeName := strings.Split(typeName, ".")

				typeValue = &Field{
					Package: splitTypeName[0],
					Name:    splitTypeName[1],
				}
			}

			require.NotEqual(t, typeReplacement, typeValue.Type)

			service.HandleSubstituteTypes(typeValue)

			require.Equal(t, typeReplacement, typeValue.Type)
		})
	}

	t.Run("ignores types that are not in the substituteTypes map", func(t *testing.T) {
		typeValue := &Field{
			Package: "test",
			Name:    "test",
		}

		service.HandleSubstituteTypes(typeValue)

		require.Equal(t, "test", typeValue.Package)
		require.Equal(t, "test", typeValue.Name)
	})

	t.Run("sets the struct fields to nil", func(t *testing.T) {
		typeValue := &Field{
			Package: "time",
			Name:    "Duration",
			StructFields: map[string]Field{
				"test": {
					Name: "test",
				},
			},
		}

		service.HandleSubstituteTypes(typeValue)

		require.Nil(t, typeValue.StructFields)
	})

	t.Run("sets the slice type to an empty string", func(t *testing.T) {
		typeValue := &Field{
			Package:   "time",
			Name:      "Duration",
			SliceType: "test",
		}

		service.HandleSubstituteTypes(typeValue)

		require.Equal(t, "", typeValue.SliceType)
	})

	t.Run("sets the map value type to an empty string", func(t *testing.T) {
		typeValue := &Field{
			Package:      "time",
			Name:         "Duration",
			MapValueType: "test",
		}

		service.HandleSubstituteTypes(typeValue)

		require.Equal(t, "", typeValue.MapValueType)
	})
}
