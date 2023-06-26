package gengo

type Output struct {
	Mode          OutputMode      `json:"mode"`
	Generate      ServiceFunction `json:"-"`
	Configuration IOConfiguration `json:"configuration"`
}

type OutputMode string

const (
	OutputModeJSON    OutputMode = "json"
	OutputModeOpenAPI OutputMode = "openapi"
)
