package astTraversal

import (
	"reflect"
	"strings"
)

type BindingTagType string

const (
	HeaderBindingTag BindingTagType = "header"
	FormBindingTag   BindingTagType = "form"
	URIBindingTag    BindingTagType = "uri"
	JSONBindingTag   BindingTagType = "json"
	XMLBindingTag    BindingTagType = "xml"
	YAMLBindingTag   BindingTagType = "yaml"
	NoBindingTag     BindingTagType = ""
)

var BindingTags = []BindingTagType{HeaderBindingTag, FormBindingTag, URIBindingTag, JSONBindingTag, XMLBindingTag, YAMLBindingTag}

type BindingTag struct {
	Name       string
	IsShown    bool
	IsOptional bool
}

type BindingTagMap map[BindingTagType]BindingTag

type ValidationTagType string

const (
	GinValidationTag       ValidationTagType = "binding"
	ValidatorValidationTag ValidationTagType = "validate"
	NoValidationTag        ValidationTagType = ""
)

var ValidationTags = []ValidationTagType{GinValidationTag, ValidatorValidationTag}

type ValidationTag struct {
	IsRequired bool
}

type ValidationTagMap map[ValidationTagType]ValidationTag

func ParseStructTag(field string, tag string) (BindingTagMap, ValidationTagMap) {
	bindingTags := make(BindingTagMap)
	for _, bindingTag := range BindingTags {
		newBindingTag := BindingTag{}

		tagValue := reflect.StructTag(tag).Get(string(bindingTag))
		if tagValue == "" {
			continue
		}

		tagItems := strings.Split(tagValue, ",")
		if tagItems[0] == "" {
			newBindingTag.Name = field
			newBindingTag.IsShown = true
		} else if tagItems[0] != "-" {
			newBindingTag.Name = tagItems[0]
			newBindingTag.IsShown = true
		} else {
			newBindingTag.IsShown = false
		}

		if len(tagItems) > 1 && tagItems[1] == "omitempty" {
			newBindingTag.IsOptional = true
		} else {
			newBindingTag.IsOptional = false
		}

		bindingTags[bindingTag] = newBindingTag
	}
	if len(bindingTags) == 0 {
		bindingTags[NoBindingTag] = BindingTag{Name: field, IsShown: true, IsOptional: false}
	}

	validationTags := make(ValidationTagMap)
	for _, validationTag := range ValidationTags {
		tagValue := reflect.StructTag(tag).Get(string(validationTag))
		if tagValue == "" {
			continue
		}

		validationTags[validationTag] = ValidationTag{
			IsRequired: strings.Contains(tagValue, "required"),
		}
	}

	return bindingTags, validationTags
}
