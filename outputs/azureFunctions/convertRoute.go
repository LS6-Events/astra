package azureFunctions

import (
	"github.com/ls6-events/gengo"
	"strings"
)

// Adapted from https://learn.microsoft.com/en-us/aspnet/web-api/overview/web-api-routing-and-actions/attribute-routing-in-web-api-2#route-constraints
var acceptedTypeMap = map[string]string{
	"string":    "alpha",
	"int":       "long",
	"int32":     "int",
	"int64":     "long",
	"uint":      "long",
	"uint32":    "int",
	"uint64":    "long",
	"float":     "double",
	"float32":   "float",
	"float64":   "double",
	"bool":      "bool",
	"time.Time": "datetime",
	"uuid.UUID": "guid",
}

func convertRoute(route gengo.Route) string {
	routeString := route.Path

	for _, pathParams := range route.PathParams {
		if azureType, ok := acceptedTypeMap[pathParams.Type]; ok {
			if azureType == "" {
				return ""
			}

			if pathParams.IsRequired {
				routeString = strings.Replace(routeString, ":"+pathParams.Name, "{"+pathParams.Name+":"+azureType+"}", 1)
			} else {
				routeString = strings.Replace(routeString, "*"+pathParams.Name, "{"+pathParams.Name+":"+azureType+"?}", 1)
			}
		}
	}

	return routeString
}
