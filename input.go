package gengo

type PopulateFunction func(service *Service) error

type Input struct {
	Mode     InputMode
	Populate PopulateFunction
}

type InputMode string

const (
	InputModeGin InputMode = "gin"
)
