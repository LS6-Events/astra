package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/gengo"
)

// WithGinInput adds gin as an input to the service
// CreateRoutes is called before ParseRoutes
// CreateRoutes is the only function that will have access to the router - it will create the routes and refer to the handler function by name, file and line number
// ParseRoutes will populate the routes with the handler function, should not need access to the router because there will be cases where the router is nil (CLI)
func WithGinInput(router *gin.Engine) gengo.Option {
	return func(s *gengo.Service) {
		s.Inputs = append(s.Inputs, gengo.Input{
			Mode:         gengo.InputModeGin,
			CreateRoutes: createRoutes(router),
			ParseRoutes:  parseRoutes(),
		})
	}
}
