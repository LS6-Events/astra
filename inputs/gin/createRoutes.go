package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/gengo"
	"reflect"
	"runtime"
)

// createRoutes creates routes from a gin routes
// It will only create the routes and refer to the handler function by name, file and line number
// The routes will be populated later by parseRoutes
// It will individually call createRoute for each route
func createRoutes(router *gin.Engine) gengo.ServiceFunction {
	return func(s *gengo.Service) error {
		s.Log.Debug().Msg("Populating service with gin routes")
		for _, route := range router.Routes() {
			s.Log.Debug().Str("path", route.Path).Str("method", route.Method).Msg("Populating route")

			denied := false
			for _, denyFunc := range s.PathDenyList {
				if denyFunc(route.Path) {
					s.Log.Debug().Str("path", route.Path).Str("method", route.Method).Msg("Path is blacklisted")
					denied = true
					break
				}
			}
			if denied {
				continue
			}

			pc := reflect.ValueOf(route.HandlerFunc).Pointer()
			file, line := runtime.FuncForPC(pc).FileLine(pc)

			s.Log.Debug().Str("path", route.Path).Str("method", route.Method).Str("file", file).Int("line", line).Msg("Found route handler")

			s.Log.Debug().Str("path", route.Path).Str("method", route.Method).Str("file", file).Int("line", line).Msg("Parsing route")
			err := createRoute(s, file, line, route)
			if err != nil {
				s.Log.Error().Str("path", route.Path).Str("method", route.Method).Str("file", file).Int("line", line).Err(err).Msg("Failed to parse route")
				return err
			}
		}
		s.Log.Debug().Msg("Populated service with gin routes")

		return nil
	}
}
