package gin

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

// findHandlersForRoute finds all handlers for a given route.
// It uses reflection to access the private properties of the gin.Engine,specifically it's method tree.
// It returns the handler pointers and a boolean indicating if the route was found.
func findHandlersForRoute(tree reflect.Value, route gin.RouteInfo) ([]uintptr, bool) {
	for i := 0; i < tree.Len(); i++ {
		method := tree.Index(i)
		methodString := method.FieldByName("method").String()
		if methodString != route.Method {
			continue
		}

		node := method.FieldByName("root")

		foundRoute, found := searchNode(node.Elem(), route.Path)
		if found {
			handlersField := foundRoute.FieldByName("handlers")
			handlers := make([]uintptr, handlersField.Len())

			for j := 0; j < handlersField.Len(); j++ {
				handlers[j] = handlersField.Index(j).Pointer()
			}

			return handlers, true
		}
	}

	return nil, false
}

// searchNode searches a gin.Engine node for a given path.
// It uses recursion to search all children nodes.
// It returns the node and a boolean indicating if the path was found.
func searchNode(node reflect.Value, path string) (reflect.Value, bool) {
	if node.FieldByName("fullPath").String() == path && node.FieldByName("handlers").Len() > 0 {
		return node, true
	}

	for i := 0; i < node.FieldByName("children").Len(); i++ {
		child, found := searchNode(node.FieldByName("children").Index(i).Elem(), path)
		if found {
			return child, true
		}
	}

	return reflect.Value{}, false
}
