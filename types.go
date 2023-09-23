package astra

// These are types that are used throughout the astra package

// Route is a route in the service and all of its potential options
type Route struct {
	Handler     string       `json:"handler" yaml:"handler"`
	File        string       `json:"file" yaml:"file"`
	LineNo      int          `json:"lineNo" yaml:"lineNo"`
	Method      string       `json:"method" yaml:"method"`
	Path        string       `json:"path" yaml:"path"`
	ContentType string       `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	BodyType    string       `json:"bodyType,omitempty" yaml:"bodyType,omitempty"`
	PathParams  []Param      `json:"params,omitempty" yaml:"params,omitempty"` // for now, we use :param in the path to denote a required path param, and *param to denote an optional path param
	QueryParams []Param      `json:"queryParams,omitempty" yaml:"queryParams,omitempty"`
	Body        []Param      `json:"body,omitempty" yaml:"body,omitempty"`
	ReturnTypes []ReturnType `json:"returnTypes,omitempty" yaml:"returnTypes,omitempty"`
	Doc         string       `json:"doc,omitempty" yaml:"doc,omitempty"`
}

// ReturnType is a return type for a route
// It contains the status code and the field that is returned
type ReturnType struct {
	StatusCode int   `json:"statusCode,omitempty" yaml:"statusCode,omitempty"`
	Field      Field `json:"field,omitempty" yaml:"field,omitempty"`
}

// Param is a parameter for a route
// It contains the name, type, and whether it is required
// It also contains an IsBound field, which is used to denote whether the param is a struct reference
type Param struct {
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Field      Field  `json:"type,omitempty" yaml:"type,omitempty"`
	IsRequired bool   `json:"isRequired,omitempty" yaml:"isRequired,omitempty"`
	IsArray    bool   `json:"isArray,omitempty" yaml:"isArray,omitempty"`
	IsMap      bool   `json:"isMap,omitempty" yaml:"isMap,omitempty"`

	IsBound bool `json:"isBound,omitempty" yaml:"isBound,omitempty"` // I.e. is a struct reference
}

// Processable is a struct that is processable by the astra package
// It just contains the name of the type and the package it came from
type Processable struct {
	Name string
	Pkg  string
}

// Field is a field in a struct
// It contains the package, type, and name of the field (the type is slice, map or struct in the case of a slice, map or struct)
// It also contains whether the field is required and whether it is embedded
// If the field is a slice, it contains the type of the slice (package is the package of the type)
// If the field is a map, it contains the key and value types of the map (and the key package, we treat the package as the value package)
// If the field is a struct, it contains the fields of the struct
type Field struct {
	Package string `json:"package,omitempty" yaml:"package,omitempty"`
	Type    string `json:"type,omitempty" yaml:"type,omitempty"`
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`

	IsRequired bool `json:"isRequired,omitempty" yaml:"isRequired,omitempty"`
	IsEmbedded bool `json:"isEmbedded,omitempty" yaml:"isEmbedded,omitempty"`

	SliceType string `json:"sliceType,omitempty" yaml:"sliceType,omitempty"`

	ArrayType   string `json:"arrayType,omitempty" yaml:"arrayType,omitempty"`
	ArrayLength int64  `json:"arrayLength,omitempty" yaml:"arrayLength,omitempty"`

	MapKeyPackage string `json:"mapKeyPackage,omitempty" yaml:"mapKeyPackage,omitempty"`
	MapKeyType    string `json:"mapKeyType,omitempty" yaml:"mapKeyType,omitempty"`
	MapValueType  string `json:"mapValueType,omitempty" yaml:"mapValueType,omitempty"`

	StructFields map[string]Field `json:"structFields,omitempty" yaml:"structFields,omitempty"`

	Doc string `json:"doc,omitempty" yaml:"doc,omitempty"`
}
