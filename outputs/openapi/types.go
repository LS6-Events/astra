package openapi

// OpenAPISchema is the OpenAPI schema
type OpenAPISchema struct {
	OpenAPI           string        `json:"openapi" yaml:"openapi"`
	Info              Info          `json:"info" yaml:"info"`
	JSONSchemaDialect string        `json:"jsonSchemaDialect,omitempty" yaml:"jsonSchemaDialect,omitempty"`
	Servers           []Server      `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths             Paths         `json:"paths,omitempty" yaml:"paths,omitempty"`
	Webhooks          Paths         `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
	Components        Components    `json:"components,omitempty" yaml:"components,omitempty"`
	Security          []Security    `json:"security,omitempty" yaml:"security,omitempty"`
	Tags              []Tag         `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs      *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// Info is the OpenAPI info
type Info struct {
	Title          string  `json:"title" yaml:"title"`
	Description    string  `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string  `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        Contact `json:"contact" yaml:"contact"`
	License        License `json:"license" yaml:"license"`
	Version        string  `json:"version" yaml:"version"`
}

// Contact is the OpenAPI contact
type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// License is the OpenAPI license
type License struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

type Server struct {
	URL         string `json:"url,omitempty" yaml:"url,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// Paths is the OpenAPI paths
type Paths map[string]Path

// Path is the OpenAPI path
type Path struct {
	Ref         string      `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Summary     string      `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *Operation  `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *Operation  `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *Operation  `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *Operation  `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options     *Operation  `json:"options,omitempty" yaml:"options,omitempty"`
	Head        *Operation  `json:"head,omitempty" yaml:"head,omitempty"`
	Patch       *Operation  `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace       *Operation  `json:"trace,omitempty" yaml:"trace,omitempty"`
	Servers     []Server    `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters  []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// Operation is the OpenAPI operation
type Operation struct {
	Ref          string        `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Tags         []string      `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary      string        `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  string        `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	OperationID  string        `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   []Parameter   `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  *RequestBody  `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses    Responses     `json:"responses,omitempty" yaml:"responses,omitempty"`
	Callbacks    Callbacks     `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	Deprecated   bool          `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security     []Security    `json:"security,omitempty" yaml:"security,omitempty"`
	Servers      []Server      `json:"servers,omitempty" yaml:"servers,omitempty"`
}

// Parameter is the OpenAPI parameter
type Parameter struct {
	Ref         string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	In          string `json:"in,omitempty" yaml:"in,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool   `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated  bool   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmpty  bool   `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Style       string `json:"style,omitempty" yaml:"style,omitempty"`
	Explode     bool   `json:"explode,omitempty" yaml:"explode,omitempty"`
	Schema      Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
}

// RequestBody is the OpenAPI request body
type RequestBody struct {
	Ref         string               `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Description string               `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty" yaml:"content,omitempty"`
	Required    bool                 `json:"required,omitempty" yaml:"required,omitempty"`
}

// MediaType is the OpenAPI media type
type MediaType struct {
	Schema   Schema              `json:"schema,omitempty" yaml:"schema,omitempty"`
	Encoding map[string]Encoding `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

// Encoding is the OpenAPI encoding
type Encoding struct {
	ContentType string            `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Headers     map[string]Header `json:"headers,omitempty" yaml:"headers,omitempty"`
}

// Header is the OpenAPI header
type Header struct {
	Ref         string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool   `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated  bool   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmpty  bool   `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
}

// Responses is the OpenAPI responses
type Responses map[string]Response

// Response is the OpenAPI response
type Response struct {
	Description string               `json:"description" yaml:"description"`
	Headers     map[string]Header    `json:"headers,omitempty" yaml:"headers,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty" yaml:"content,omitempty"`
	Links       map[string]Link      `json:"links,omitempty" yaml:"links,omitempty"`
}

// Link is the OpenAPI link
type Link struct {
	OperationRef string               `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	OperationID  string               `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   map[string]Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  interface{}          `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Description  string               `json:"description,omitempty" yaml:"description,omitempty"`
	Server       Server               `json:"server,omitempty" yaml:"server,omitempty"`
}

// Callbacks is the OpenAPI callbacks
type Callbacks map[string]Callback

// Callback is the OpenAPI callback
type Callback map[string]Path

// Components is the OpenAPI components
type Components struct {
	Schemas         map[string]Schema         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses       map[string]Response       `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameters      map[string]Parameter      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBodies   map[string]RequestBody    `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	Headers         map[string]Header         `json:"headers,omitempty" yaml:"headers,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Links           map[string]Link           `json:"links,omitempty" yaml:"links,omitempty"`
	Callbacks       map[string]Callback       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	PathItems       map[string]Path           `json:"pathItems,omitempty" yaml:"pathItems,omitempty"`
}

// Schema is JSON Schema utilised by OpenAPI
type Schema struct {
	Ref                  string            `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Title                string            `json:"title,omitempty" yaml:"title,omitempty"`
	MultipleOf           float64           `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	Maximum              float64           `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	ExclusiveMaximum     bool              `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	Minimum              float64           `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	ExclusiveMinimum     bool              `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`
	MaxLength            int               `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength            int               `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	Pattern              string            `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MaxItems             int               `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems             int               `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems          bool              `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	MaxProperties        int               `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties        int               `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	Required             []string          `json:"required,omitempty" yaml:"required,omitempty"`
	Enum                 []interface{}     `json:"enum,omitempty" yaml:"enum,omitempty"`
	Type                 string            `json:"type,omitempty" yaml:"type,omitempty"`
	Format               string            `json:"format,omitempty" yaml:"format,omitempty"`
	AllOf                []Schema          `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	OneOf                []Schema          `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOf                []Schema          `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	Not                  *Schema           `json:"not,omitempty" yaml:"not,omitempty"`
	Items                *Schema           `json:"items,omitempty" yaml:"items,omitempty"`
	Properties           map[string]Schema `json:"properties,omitempty" yaml:"properties,omitempty"`
	PatternProperties    map[string]Schema `json:"patternProperties,omitempty" yaml:"patternProperties,omitempty"`
	AdditionalProperties *Schema           `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Description          string            `json:"description,omitempty" yaml:"description,omitempty"`
}

// SecurityScheme is the OpenAPI security scheme
type SecurityScheme struct {
	Ref         string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	In          string `json:"in,omitempty" yaml:"in,omitempty"`
	Scheme      string `json:"scheme,omitempty" yaml:"scheme,omitempty"`
}

// Security is the OpenAPI security
type Security map[string][]string

// Tag is the OpenAPI tag
type Tag struct {
	Name         string        `json:"name,omitempty" yaml:"name,omitempty"`
	Description  string        `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// ExternalDocs is the OpenAPI external documentation
type ExternalDocs struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string `json:"url,omitempty" yaml:"url,omitempty"`
}
