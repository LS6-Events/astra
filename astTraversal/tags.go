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
	Name           string `json:"name,omitempty" yaml:"name,omitempty"`
	NotShown       bool   `json:"not_shown,omitempty" yaml:"not_shown,omitempty"`
	ReturnOptional bool   `json:"return_optional,omitempty" yaml:"return_optional,omitempty"`
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
	IsRequired bool `json:"is_required,omitempty" yaml:"is_required,omitempty"`
}

type ValidationTagMap map[ValidationTagType]ValidationTag

func ParseStructTag(field string, tag string) (BindingTagMap, ValidationTagMap) {
	bindingTags := make(BindingTagMap)
	for _, bindingTag := range BindingTags {
		tagValue, tagOk := reflect.StructTag(tag).Lookup(string(bindingTag))
		if !tagOk {
			continue
		}

		newBindingTag := BindingTag{}

		tagItems := strings.Split(tagValue, ",")
		if tagItems[0] == "" {
			newBindingTag.Name = field
		} else if tagItems[0] != "-" {
			newBindingTag.Name = tagItems[0]
		} else {
			newBindingTag.NotShown = true
		}

		if len(tagItems) > 1 && tagItems[1] == "omitempty" {
			newBindingTag.ReturnOptional = true
		} else {
			newBindingTag.ReturnOptional = false
		}

		bindingTags[bindingTag] = newBindingTag
	}
	if len(bindingTags) == 0 {
		bindingTags[NoBindingTag] = BindingTag{Name: field, NotShown: false, ReturnOptional: false}
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
