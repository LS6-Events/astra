package gin

import (
	"github.com/ls6-events/astra"
)

// ParseRoutes parses routes from a gin routes.
// It will populate the routes with the handler function.
// It will individually call parseHandler for each route.
// createRoutes must be called before this.
func ParseRoutes() astra.ServiceFunction {
	return func(s *astra.Service) error {
		s.Log.Debug().Msg("Populating routes from gin routes")
		for _, route := range s.Routes {
			s.Log.Debug().Str("path", route.Path).Str("method", route.Method).Msg("Parsing route")

			for i := range route.Handlers {
				s.Log.Debug().Str("path", route.Path).Str("method", route.Method).Str("file", route.Handlers[i].File).Int("line", route.Handlers[i].LineNo).Msg("Parsing handler")
				err := parseHandler(s, &route, i)
				if err != nil {
					s.Log.Error().Str("path", route.Path).Str("method", route.Method).Str("file", route.Handlers[i].File).Int("line", route.Handlers[i].LineNo).Err(err).Msg("Failed to parse handler")
					return err
				}
			}

			s.ReplaceRoute(route)
		}
		s.Log.Debug().Msg("Populated service with gin routes")

		return nil
	}
}
