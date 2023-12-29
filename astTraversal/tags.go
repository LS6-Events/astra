package astTraversal

import (
	"reflect"
	"strings"
)

func ParseStructTag(tag string) (name string, isRequired bool, isShown bool) {
	isShown = true
	if tag != "" {
		binding := reflect.StructTag(tag).Get("binding")
		if binding != "" {
			isRequired = strings.Contains(binding, "required")
		}

		yaml := reflect.StructTag(tag).Get("yaml")
		if yaml != "" && yaml != "-" {
			isShown = true
			name = strings.Split(yaml, ",")[0]
		} else if yaml == "-" && isShown {
			isShown = false
		}

		xml := reflect.StructTag(tag).Get("xml")
		if xml != "" && xml != "-" {
			isShown = true
			name = strings.Split(xml, ",")[0]
		} else if xml == "-" && isShown {
			isShown = false
		}

		form := reflect.StructTag(tag).Get("form")
		if form != "" && form != "-" {
			isShown = true
			name = strings.Split(form, ",")[0]
		} else if form == "-" && isShown {
			isShown = false
		}

		json := reflect.StructTag(tag).Get("json")
		if json != "" && json != "-" {
			isShown = true
			name = strings.Split(json, ",")[0]
		} else if json == "-" && isShown {
			isShown = false
		}
	}

	return name, isRequired, isShown
}
