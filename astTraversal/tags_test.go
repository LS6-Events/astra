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
			tag:   `json:"field1,omitempty" form:"field1" binding:"required"`,
			expectedBindingTags: BindingTagMap{
				JSONBindingTag: {Name: "field1", IsShown: true, IsOptional: true},
				FormBindingTag: {Name: "field1", IsShown: true, IsOptional: false},
			},
			expectedValidationTags: ValidationTagMap{
				GinValidationTag: {IsRequired: true},
			},
		},
		{
			field: "Field2",
			tag:   `xml:"field2" form:",omitempty"`,
			expectedBindingTags: BindingTagMap{
				XMLBindingTag:  {Name: "field2", IsShown: true, IsOptional: false},
				FormBindingTag: {Name: "Field2", IsShown: true, IsOptional: true},
			},
			expectedValidationTags: ValidationTagMap{},
		},
		{
			field: "Field3",
			tag:   `xml:"-"`,
			expectedBindingTags: BindingTagMap{
				XMLBindingTag: {IsShown: false, IsOptional: false},
			},
			expectedValidationTags: ValidationTagMap{},
		},
		{
			field: "Field4",
			tag:   `json:"field4,omitempty" form:"-"`,
			expectedBindingTags: BindingTagMap{
				JSONBindingTag: {Name: "field4", IsShown: true, IsOptional: true},
				FormBindingTag: {IsShown: false, IsOptional: false},
			},
			expectedValidationTags: ValidationTagMap{},
		},
		{
			field: "Field5",
			tag:   `binding:"required"`,
			expectedBindingTags: BindingTagMap{
				NoBindingTag: {Name: "Field5", IsShown: true, IsOptional: false},
			},
			expectedValidationTags: ValidationTagMap{
				GinValidationTag: {IsRequired: true},
			},
		},
		{
			field: "Field6",
			tag:   `json:""`,
			expectedBindingTags: BindingTagMap{
				JSONBindingTag: {Name: "Field6", IsShown: true, IsOptional: false},
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
