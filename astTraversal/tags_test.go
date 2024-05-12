package astTraversal

import (
	"go/types"
	"reflect"
	"testing"
)

func TestParseStructTag(t *testing.T) {
	testCases := []struct {
		field               string
		tag                 string
		node                types.Type
		expectedBindingTags BindingTagMap
	}{
		{
			field: "Field1",
			tag:   `json:"field1"`,
			node:  &types.Basic{},
			expectedBindingTags: BindingTagMap{
				JSONBindingTag: {
					Name:           "field1",
					NotShown:       false,
					ReturnOptional: false,
				},
			},
		},
		{
			field: "Field2",
			tag:   `json:"field2,omitempty"`,
			node:  &types.Basic{},
			expectedBindingTags: BindingTagMap{
				JSONBindingTag: {
					Name:           "field2",
					NotShown:       false,
					ReturnOptional: true,
				},
			},
		},
		{
			field: "Field3",
			tag:   `json:""`,
			node:  &types.Basic{},
			expectedBindingTags: BindingTagMap{
				JSONBindingTag: {
					Name:           "Field3",
					NotShown:       false,
					ReturnOptional: false,
				},
			},
		},
		{
			field: "Field4",
			tag:   `json:",omitempty"`,
			node:  &types.Basic{},
			expectedBindingTags: BindingTagMap{
				JSONBindingTag: {
					Name:           "Field4",
					NotShown:       false,
					ReturnOptional: true,
				},
			},
		},
		{
			field: "Field5",
			tag:   `json:"-"`,
			node:  &types.Basic{},
			expectedBindingTags: BindingTagMap{
				JSONBindingTag: {
					Name:           "",
					NotShown:       true,
					ReturnOptional: false,
				},
			},
		},
		{
			field: "Field7",
			tag:   ``,
			node:  &types.Basic{},
			expectedBindingTags: BindingTagMap{
				NoBindingTag: {
					Name:           "Field7",
					NotShown:       false,
					ReturnOptional: false,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run("field="+testCase.field, func(t *testing.T) {
			bindingTags, _, _ := ParseStructTag(testCase.field, testCase.node, testCase.tag)

			if !reflect.DeepEqual(bindingTags, testCase.expectedBindingTags) {
				t.Errorf("Expected BindingTags: %v, but got: %v", testCase.expectedBindingTags, bindingTags)
			}
		})
	}
}
