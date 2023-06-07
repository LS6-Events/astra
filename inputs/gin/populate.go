package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/gengo"
	"reflect"
	"runtime"
)

func populate(router *gin.Engine) gengo.PopulateFunction {
	return func(s *gengo.Service) error {
		s.Log.Debug().Msg("Populating service with gin routes")
		for _, route := range router.Routes() {
			s.Log.Debug().Str("path", route.Path).Str("method", route.Method).Msg("Populating route")
			pc := reflect.ValueOf(route.HandlerFunc).Pointer()
			file, _ := runtime.FuncForPC(pc).FileLine(pc)
			s.Log.Debug().Str("path", route.Path).Str("method", route.Method).Str("file", file).Msg("Found route handler")

			s.Log.Debug().Str("path", route.Path).Str("method", route.Method).Str("file", file).Msg("Parsing route")
			err := parseRoute(s, file, route)
			if err != nil {
				s.Log.Error().Str("path", route.Path).Str("method", route.Method).Str("file", file).Err(err).Msg("Failed to parse route")
				return err
			}
		}
		s.Log.Debug().Msg("Populated service with gin routes")

		return nil
	}
}
