package gengo

type Route struct {
	Method      string       `json:"method"`
	Path        string       `json:"path"`
	ContentType string       `json:"contentType,omitempty"`
	BodyType    string       `json:"bodyType,omitempty"`
	PathParams  []Param      `json:"params,omitempty"` // for now, we use :param in the path to denote a required path param, and *param to denote an optional path param
	QueryParams []Param      `json:"queryParams,omitempty"`
	Body        []Param      `json:"body,omitempty"`
	ReturnTypes []ReturnType `json:"returnTypes,omitempty"`
}

type ReturnType struct {
	StatusCode int   `json:"statusCode,omitempty"`
	Field      Field `json:"field,omitempty"`
}

type Param struct {
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	Package    string `json:"package,omitempty"`
	IsRequired bool   `json:"isRequired,omitempty"`
	IsArray    bool   `json:"isArray,omitempty"`
	IsMap      bool   `json:"isMap,omitempty"`

	IsBound bool `json:"isBound,omitempty"` // I.e. is a struct reference
}

type Processable struct {
	Name string
	Pkg  string
}

type Field struct {
	Package string `json:"package,omitempty"`
	Type    string `json:"type,omitempty"`
	Name    string `json:"name,omitempty"`

	IsRequired bool `json:"isRequired,omitempty"`
	IsEmbedded bool `json:"isEmbedded,omitempty"`

	SliceType string `json:"sliceType,omitempty"`

	MapKeyPkg string `json:"mapKeyPkg,omitempty"`
	MapKey    string `json:"mapKey,omitempty"`
	MapValue  string `json:"mapValue,omitempty"`

	StructFields map[string]Field `json:"structFields,omitempty"`
}
