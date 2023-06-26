package gengo

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
}

type ReturnType struct {
	StatusCode int   `json:"statusCode,omitempty" yaml:"statusCode,omitempty"`
	Field      Field `json:"field,omitempty" yaml:"field,omitempty"`
}

type Param struct {
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Type       string `json:"type,omitempty" yaml:"type,omitempty"`
	Package    string `json:"package,omitempty" yaml:"package,omitempty"`
	IsRequired bool   `json:"isRequired,omitempty" yaml:"isRequired,omitempty"`
	IsArray    bool   `json:"isArray,omitempty" yaml:"isArray,omitempty"`
	IsMap      bool   `json:"isMap,omitempty" yaml:"isMap,omitempty"`

	IsBound bool `json:"isBound,omitempty" yaml:"isBound,omitempty"` // I.e. is a struct reference
}

type Processable struct {
	Name string
	Pkg  string
}

type Field struct {
	Package string `json:"package,omitempty" yaml:"package,omitempty"`
	Type    string `json:"type,omitempty" yaml:"type,omitempty"`
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`

	IsRequired bool `json:"isRequired,omitempty" yaml:"isRequired,omitempty"`
	IsEmbedded bool `json:"isEmbedded,omitempty" yaml:"isEmbedded,omitempty"`

	SliceType string `json:"sliceType,omitempty" yaml:"sliceType,omitempty"`

	MapKeyPkg string `json:"mapKeyPkg,omitempty" yaml:"mapKeyPkg,omitempty"`
	MapKey    string `json:"mapKey,omitempty" yaml:"mapKey,omitempty"`
	MapValue  string `json:"mapValue,omitempty" yaml:"mapValue,omitempty"`

	StructFields map[string]Field `json:"structFields,omitempty" yaml:"structFields,omitempty"`
}
