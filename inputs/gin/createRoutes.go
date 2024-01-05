package gin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/astra"
	"reflect"
)

var (
	ErrRouteNotFound = errors.New("route not found")
)

// CreateRoutes creates routes from a gin routes.
// It will only create the routes and refer to the handler function by name, file and line number.
// The routes will be populated later by parseRoutes.
// It will individually call createRoute for each route.
func CreateRoutes(router *gin.Engine) astra.ServiceFunction {
	return func(s *astra.Service) error {
		s.Log.Debug().Msg("Populating service with gin routes")

		// To prevent performance issues, only find the gin.Engine tree if middleware is enabled
		var trees reflect.Value
		if s.UnstableEnableMiddleware {
			trees = reflect.ValueOf(router).Elem().FieldByName("trees")
		}

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

			var handlers []uintptr

			// If middleware is enabled, find the handlers for the route using the gin.Engine tree
			if s.UnstableEnableMiddleware {
				var found bool
				handlers, found = findHandlersForRoute(trees, route)
				if !found {
					s.Log.Error().Str("path", route.Path).Str("method", route.Method).Msg("Route not found")
					return ErrRouteNotFound
				}
			} else {
				// If middleware is disabled, use the gin.Engine handlers
				handlers = []uintptr{reflect.ValueOf(route.HandlerFunc).Pointer()}
			}

			err := createRoute(s, s.Context, handlers, route)
			if err != nil {
				s.Log.Error().Str("path", route.Path).Err(err).Msg("Failed to parse route")
				return err
			}
		}
		s.Log.Debug().Msg("Populated service with gin routes")

		return nil
	}
}
