package astTraversal

import (
	"reflect"
	"testing"
)

func TestParseStructTag(t *testing.T) {
	testCases := []struct {
		field                  string
		tag                    string
		expectedBindingTags    BindingTagMap
		expectedValidationTags ValidationTagMap
	}{
		{
			field: "Field1",
			tag:   `json:"field1"`,
			expectedBindingTags: BindingTagMap{
				JSONBindingTag: {
					Name:           "field1",
					NotShown:       false,
					ReturnOptional: false,
				},
			},
			expectedValidationTags: ValidationTagMap{},
		},
		{
			field: "Field2",
			tag:   `json:"field2,omitempty"`,
			expectedBindingTags: BindingTagMap{
				JSONBindingTag: {
					Name:           "field2",
					NotShown:       false,
					ReturnOptional: true,
				},
			},
			expectedValidationTags: ValidationTagMap{},
		},
		{
			field: "Field3",
			tag:   `json:""`,
			expectedBindingTags: BindingTagMap{
				JSONBindingTag: {
					Name:           "Field3",
					NotShown:       false,
					ReturnOptional: false,
				},
			},
			expectedValidationTags: ValidationTagMap{},
		},
		{
			field: "Field4",
			tag:   `json:",omitempty"`,
			expectedBindingTags: BindingTagMap{
				JSONBindingTag: {
					Name:           "Field4",
					NotShown:       false,
					ReturnOptional: true,
				},
			},
			expectedValidationTags: ValidationTagMap{},
		},
		{
			field: "Field5",
			tag:   `json:"-"`,
			expectedBindingTags: BindingTagMap{
				JSONBindingTag: {
					Name:           "",
					NotShown:       true,
					ReturnOptional: false,
				},
			},
			expectedValidationTags: ValidationTagMap{},
		},
		{
			field: "Field6",
			tag:   `validate:"required"`,
			expectedBindingTags: BindingTagMap{
				NoBindingTag: {
					Name:           "Field6",
					NotShown:       false,
					ReturnOptional: false,
				},
			},
			expectedValidationTags: ValidationTagMap{
				ValidatorValidationTag: {
					IsRequired: true,
				},
			},
		},
		{
			field: "Field7",
			tag:   ``,
			expectedBindingTags: BindingTagMap{
				NoBindingTag: {
					Name:           "Field7",
					NotShown:       false,
					ReturnOptional: false,
				},
			},
			expectedValidationTags: ValidationTagMap{},
		},
	}

	for _, testCase := range testCases {
		t.Run("field="+testCase.field, func(t *testing.T) {
			bindingTags, validationTags := ParseStructTag(testCase.field, testCase.tag)

			if !reflect.DeepEqual(bindingTags, testCase.expectedBindingTags) {
				t.Errorf("Expected BindingTags: %v, but got: %v", testCase.expectedBindingTags, bindingTags)
			}

			if !reflect.DeepEqual(validationTags, testCase.expectedValidationTags) {
				t.Errorf("Expected ValidationTags: %v, but got: %v", testCase.expectedValidationTags, validationTags)
			}
		})
	}
}
