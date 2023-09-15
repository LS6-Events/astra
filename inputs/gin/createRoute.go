package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/astra"
	"os"
	"path/filepath"
)

// createRoute creates a route from a gin RouteInfo
// It will only create the route and refer to the handler function by name, file and line number
// The route will be populated later by parseRoute
func createRoute(s *astra.Service, file string, line int, info gin.RouteInfo) error {
	log := s.Log.With().Str("path", info.Path).Str("method", info.Method).Str("handler", info.Handler).Logger()

	cwd, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get working directory")
		return err
	}

	relativePath, err := filepath.Rel(cwd, file)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get relative path")
		return err
	}

	baseRoute := astra.Route{
		Handler:     info.Handler,
		File:        relativePath,
		LineNo:      line,
		Path:        info.Path,
		Method:      info.Method,
		PathParams:  make([]astra.Param, 0),
		Body:        make([]astra.Param, 0),
		QueryParams: make([]astra.Param, 0),
		ReturnTypes: make([]astra.ReturnType, 0),
	}

	s.AddRoute(baseRoute)

	log.Debug().Msg("Populated route")

	return nil
}
