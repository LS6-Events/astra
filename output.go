package gengo

// Output is the output for the generator
// It takes in a configuration object to allow the caching mechanism to store the file paths for the outputs
type Output struct {
	Mode          OutputMode      `json:"mode"`
	Generate      ServiceFunction `json:"-"`
	Configuration IOConfiguration `json:"configuration"`
}

type OutputMode string

const (
	OutputModeJSON    OutputMode = "json"    // JSON file - primarily used for debugging
	OutputModeOpenAPI OutputMode = "openapi" // OpenAPI 3.0 file
)
