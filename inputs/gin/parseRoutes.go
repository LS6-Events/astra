package gin

import (
	"github.com/ls6-events/gengo"
)

// ParseRoutes parses routes from a gin routes
// It will populate the routes with the handler function
// It will individually call parseRoute for each route
// createRoutes must be called before this
func ParseRoutes() gengo.ServiceFunction {
	return func(s *gengo.Service) error {
		s.Log.Debug().Msg("Populating routes from gin routes")
		for _, route := range s.Routes {
			s.Log.Debug().Str("path", route.Path).Str("method", route.Method).Msg("Populating route")

			s.Log.Debug().Str("path", route.Path).Str("method", route.Method).Str("file", route.File).Int("line", route.LineNo).Msg("Parsing route")
			err := parseRoute(s, &route)
			if err != nil {
				s.Log.Error().Str("path", route.Path).Str("method", route.Method).Str("file", route.File).Int("line", route.LineNo).Err(err).Msg("Failed to parse route")
				return err
			}

			s.ReplaceRoute(route)
		}
		s.Log.Debug().Msg("Populated service with gin routes")

		return nil
	}
}
