package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/gengo"
)

func WithGinInput(router *gin.Engine) gengo.Option {
	return func(s *gengo.Service) {
		s.Inputs = append(s.Inputs, gengo.Input{
			Mode:         gengo.InputModeGin,
			CreateRoutes: createRoutes(router),
			ParseRoutes:  parseRoutes(),
		})
	}
}
