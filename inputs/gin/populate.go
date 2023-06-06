package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/gengo"
	"reflect"
	"runtime"
)

func populate(router *gin.Engine) gengo.PopulateFunction {
	return func(s *gengo.Service) error {
		for _, route := range router.Routes() {

			pc := reflect.ValueOf(route.HandlerFunc).Pointer()
			file, _ := runtime.FuncForPC(pc).FileLine(pc)

			err := parseRoute(s, file, route)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
