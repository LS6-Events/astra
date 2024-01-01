package astra

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_AddRoute(t *testing.T) {
	service := &Service{}
	route := Route{}

	require.Len(t, service.Routes, 0)

	service.AddRoute(route)

	require.Len(t, service.Routes, 1)
	require.Equal(t, route, service.Routes[0])
}

func TestService_ReplaceRoute(t *testing.T) {
	service := &Service{}
	route := Route{
		Path:   "/test",
		Method: "GET",
		Doc:    "Route 1",
	}

	service.AddRoute(route)

	require.Len(t, service.Routes, 1)
	require.Equal(t, route, service.Routes[0])

	route = Route{
		Path:   "/test",
		Method: "POST",
		Doc:    "Route 2",
	}

	service.ReplaceRoute(route)

	require.Len(t, service.Routes, 1)
	require.Equal(t, "Route 1", service.Routes[0].Doc)

	route = Route{
		Path:   "/test",
		Method: "GET",
		Doc:    "Route 2",
	}

	service.ReplaceRoute(route)

	require.Len(t, service.Routes, 1)
	require.Equal(t, "Route 2", service.Routes[0].Doc)
}
