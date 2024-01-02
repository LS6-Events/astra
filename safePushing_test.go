package astra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddReturnType(t *testing.T) {
	t.Run("AddingToEmptySlice", func(t *testing.T) {
		var emptySlice []ReturnType
		newReturn := ReturnType{
			StatusCode: 200,
			Field: Field{
				Package: "test",
				Type:    "TestType",
			},
		}
		result := AddReturnType(emptySlice, newReturn)
		expected := []ReturnType{newReturn}
		assert.Equal(t, expected, result)
	})

	t.Run("AddingToExistingSliceWithSameField", func(t *testing.T) {
		existingReturn := ReturnType{
			StatusCode: 200,
			Field: Field{
				Package: "test",
				Type:    "TestType",
			},
		}
		existingSlice := []ReturnType{existingReturn}
		newReturn := ReturnType{
			StatusCode: 200,
			Field: Field{
				Package: "test",
				Type:    "TestType",
			},
		}
		result := AddReturnType(existingSlice, newReturn)
		assert.Equal(t, existingSlice, result)
	})

	t.Run("AddingToExistingSliceWithDifferentField", func(t *testing.T) {
		existingReturn := ReturnType{
			StatusCode: 200,
			Field: Field{
				Package: "test",
				Type:    "TestType",
			},
		}
		existingSlice := []ReturnType{existingReturn}
		differentFieldReturn := ReturnType{
			StatusCode: 400,
			Field: Field{
				Package: "other",
				Type:    "OtherType",
			},
		}
		result := AddReturnType(existingSlice, differentFieldReturn)
		expected := append(existingSlice, differentFieldReturn)
		assert.Equal(t, expected, result)
	})
}

func TestAddComponent(t *testing.T) {
	t.Run("AddingToEmptySlice", func(t *testing.T) {
		var emptySlice []Field
		newComponent := Field{
			Name:    "ComponentA",
			Package: "test",
		}
		result := AddComponent(emptySlice, newComponent)
		expected := []Field{newComponent}
		assert.Equal(t, expected, result)
	})

	t.Run("AddingToExistingSliceWithSameComponent", func(t *testing.T) {
		existingComponent := Field{
			Name:    "ComponentA",
			Package: "test",
		}
		existingSlice := []Field{existingComponent}
		newComponent := Field{
			Name:    "ComponentA",
			Package: "test",
		}
		result := AddComponent(existingSlice, newComponent)
		assert.Equal(t, existingSlice, result)
	})

	t.Run("AddingToExistingSliceWithDifferentComponent", func(t *testing.T) {
		existingComponent := Field{
			Name:    "ComponentA",
			Package: "test",
		}
		existingSlice := []Field{existingComponent}
		differentComponent := Field{
			Name:    "ComponentB",
			Package: "other",
		}
		result := AddComponent(existingSlice, differentComponent)
		expected := append(existingSlice, differentComponent)
		assert.Equal(t, expected, result)
	})
}
