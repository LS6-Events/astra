package astTraversal

import (
	"github.com/stretchr/testify/assert"
	"go/token"
	"go/types"
	"strings"
	"testing"
)

func TestTypeTraverser_Result(t *testing.T) {
	baseTraverser, err := createTraverserFromTestFile("usefulTypes.go")
	assert.NoError(t, err)

	_, err = baseTraverser.Packages.Get(baseTraverser.ActiveFile().Package)
	assert.NoError(t, err)

	t.Run("Basic", func(t *testing.T) {
		tt := baseTraverser.Type(types.Typ[types.Int], nil)
		res, err := tt.Result()
		assert.Nil(t, err)
		assert.Equal(t, "int", res.Type)
	})

	t.Run("Pointer", func(t *testing.T) {
		tt := baseTraverser.Type(types.NewPointer(types.Typ[types.Int]), nil)
		res, err := tt.Result()
		assert.Nil(t, err)
		assert.Equal(t, "int", res.Type)
	})

	t.Run("Slice", func(t *testing.T) {
		tt := baseTraverser.Type(types.NewSlice(types.Typ[types.Int]), nil)
		res, err := tt.Result()
		assert.Nil(t, err)
		assert.Equal(t, "slice", res.Type)
		assert.Equal(t, "int", res.SliceType)
	})

	t.Run("Array", func(t *testing.T) {
		tt := baseTraverser.Type(types.NewArray(types.Typ[types.Int], 5), nil)
		res, err := tt.Result()
		assert.Nil(t, err)
		assert.Equal(t, "array", res.Type)
		assert.Equal(t, "int", res.ArrayType)
		assert.Equal(t, int64(5), res.ArrayLength)
	})

	t.Run("Map", func(t *testing.T) {
		tt := baseTraverser.Type(types.NewMap(types.Typ[types.String], types.Typ[types.Int]), nil)
		res, err := tt.Result()
		assert.Nil(t, err)
		assert.Equal(t, "map", res.Type)
		assert.Equal(t, "string", res.MapKeyType)
		assert.Equal(t, "int", res.MapValueType)
	})

	t.Run("Named", func(t *testing.T) {
		// The only test that requires us having to use the testfiles package - looking at MyStruct in usefulTypes.go
		namedType, err := baseTraverser.ActiveFile().Package.FindObjectForName("MyStruct")
		assert.NoError(t, err)

		tt := baseTraverser.Type(namedType.Type(), baseTraverser.ActiveFile().Package)
		res, err := tt.Result()
		assert.Nil(t, err)
		assert.Equal(t, "MyStruct", res.Type)
	})

	t.Run("Struct", func(t *testing.T) {
		// Creating a simple struct with a field "Age" of type int
		fields := []*types.Var{
			types.NewField(token.NoPos, nil, "Age", types.Typ[types.Int], false),
		}
		structType := types.NewStruct(fields, []string{"json:\"age\" binding:\"required\""})

		tt := baseTraverser.Type(structType, nil)
		res, err := tt.Result()
		assert.Nil(t, err)
		assert.Equal(t, "struct", res.Type)

		// Check if the struct fields are parsed correctly
		ageField, ok := res.StructFields["Age"]
		assert.True(t, ok)
		assert.Equal(t, "int", ageField.Type)
		assert.True(t, ageField.StructFieldValidationTags[GinValidationTag].IsRequired)
	})
}

func TestTypeTraverser_Doc(t *testing.T) {
	baseTraverser, err := createTraverserFromTestFile("usefulTypes.go")
	assert.NoError(t, err)

	_, err = baseTraverser.Packages.Get(baseTraverser.ActiveFile().Package)
	assert.NoError(t, err)

	t.Run("Basic", func(t *testing.T) {
		tt := baseTraverser.Type(types.Typ[types.Int], baseTraverser.ActiveFile().Package)
		doc, err := tt.Doc()
		assert.Nil(t, err)
		assert.Empty(t, doc)
	})

	t.Run("Named", func(t *testing.T) {
		// The only test that requires us having to use the testfiles package - looking at MyStruct in usefulTypes.go
		namedType, err := baseTraverser.ActiveFile().Package.FindObjectForName("MyStruct")
		assert.NoError(t, err)

		tt := baseTraverser.Type(namedType.Type(), baseTraverser.ActiveFile().Package)
		doc, err := tt.Doc()
		assert.Nil(t, err)
		assert.Equal(t, "MyStruct is a struct", strings.TrimSpace(doc))
	})
}
