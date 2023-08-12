package inputs

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/gengo"
	gengoGin "github.com/ls6-events/gengo/inputs/gin"
)

const (
	InputModeGin gengo.InputMode = "gin" // github.com/gin-gonic/gin web framework
)

func addInput(mode gengo.InputMode, createRoutes gengo.ServiceFunction, parseRoutes gengo.ServiceFunction) gengo.Option {
	return func(s *gengo.Service) {
		s.Inputs = append(s.Inputs, gengo.Input{
			Mode:         mode,
			CreateRoutes: createRoutes,
			ParseRoutes:  parseRoutes,
		})
	}
}

// WithGinInput adds gin as an input to the service
// CreateRoutes is called before ParseRoutes
// CreateRoutes is the only function that will have access to the routes - it will create the routes and refer to the handler function by name, file and line number
// ParseRoutes will populate the routes with the handler function, should not need access to the routes because there will be cases where the routes is nil (CLI)
func WithGinInput(router *gin.Engine) gengo.Option {
	return addInput(
		InputModeGin,
		gengoGin.CreateRoutes(router),
		gengoGin.ParseRoutes(),
	)
}
