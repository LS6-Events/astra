package gin

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/ls6-events/astra"

	"github.com/gin-gonic/gin"
)

var deniedHandlers = []string{
	"github.com/gin-gonic/gin.LoggerWithConfig.func1",
	"github.com/gin-gonic/gin.CustomRecoveryWithWriter.func1",
}

// createRoute creates a route from a gin RouteInfo.
// It will only create the route and refer to the handler function by name, file and line number.
// The route will be populated later by parseHandler.
func createRoute(s *astra.Service, handlers []uintptr, info gin.RouteInfo) error {
	log := s.Log.With().Str("path", info.Path).Str("method", info.Method).Str("handler", info.Handler).Logger()

	cwd, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get working directory")
		return err
	}

	log.Debug().Msg("Found route handlers")
	resultHandlers := make([]astra.Handler, 0)
	for _, handler := range handlers {
		funcForPC := runtime.FuncForPC(handler)
		name := funcForPC.Name()
		file, line := funcForPC.FileLine(handler)

		var found bool
		for _, deniedHandler := range deniedHandlers {
			if deniedHandler == name {
				found = true
				break
			}
		}

		if found {
			log.Debug().Str("name", name).Str("file", file).Int("line", line).Msg("Denied handler")
			continue
		}

		log.Debug().Str("name", name).Str("file", file).Int("line", line).Msg("Found handler")

		relativePath, err := filepath.Rel(cwd, file)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get relative path")
			return err
		}

		resultHandlers = append(resultHandlers, astra.Handler{
			Name:   name,
			File:   relativePath,
			LineNo: line,
		})
	}

	baseRoute := astra.Route{
		Handlers:    resultHandlers,
		Path:        info.Path,
		Method:      info.Method,
		PathParams:  make([]astra.Param, 0),
		Body:        make([]astra.BodyParam, 0),
		QueryParams: make([]astra.Param, 0),
		ReturnTypes: make([]astra.ReturnType, 0),
	}

	s.AddRoute(baseRoute)

	log.Debug().Msg("Populated route")

	return nil
}
