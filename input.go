package gengo

type Input struct {
	Mode         InputMode       `json:"mode"`
	CreateRoutes ServiceFunction `json:"-"`
	ParseRoutes  ServiceFunction `json:"-"`
}

type InputMode string

const (
	InputModeGin InputMode = "gin"
)
