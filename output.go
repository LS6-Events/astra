package gengo

type GenerateFunction func(service *Service) error

type Output struct {
	Mode     OutputMode       `json:"mode"`
	Generate GenerateFunction `json:"-"`
}

type OutputMode string

const (
	OutputModeJSON    OutputMode = "json"
	OutputModeOpenAPI OutputMode = "openapi"
)
