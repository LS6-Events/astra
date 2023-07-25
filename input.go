package gengo

// Input is the input for the generator
type Input struct {
	Mode         InputMode       `json:"mode"`
	CreateRoutes ServiceFunction `json:"-"`
	ParseRoutes  ServiceFunction `json:"-"`
}

type InputMode string

const (
	InputModeGin InputMode = "gin" // github.com/gin-gonic/gin web framework
)
