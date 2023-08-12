package gengo

// Input is the input for the generator
type Input struct {
	Mode         InputMode       `json:"mode"`
	CreateRoutes ServiceFunction `json:"-"`
	ParseRoutes  ServiceFunction `json:"-"`
}

type InputMode string
