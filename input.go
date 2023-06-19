package gengo

type PopulateFunction func(service *Service) error

type Input struct {
	Mode     InputMode        `json:"mode"`
	Populate PopulateFunction `json:"-"`
}

type InputMode string

const (
	InputModeGin InputMode = "gin"
)
