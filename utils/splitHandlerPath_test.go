package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSplitHandlerPath(t *testing.T) {
	testCases := []struct {
		name   string
		path   string
		result HandlerPath
	}{
		{
			name: "simple",
			path: "hello.world",
			result: HandlerPath{
				PathParts:    []string{"hello"},
				HandlerParts: []string{"world"},
			},
		},
		{
			name: "nested path",
			path: "foo/bar.hello",
			result: HandlerPath{
				PathParts:    []string{"foo", "bar"},
				HandlerParts: []string{"hello"},
			},
		},
		{
			name: "nested handler",
			path: "foo/bar/hello.world.func1",
			result: HandlerPath{
				PathParts:    []string{"foo", "bar", "hello"},
				HandlerParts: []string{"world", "func1"},
			},
		},
		{
			name: "nested path and handler",
			path: "foo/bar/hello.world.func1.func2",
			result: HandlerPath{
				PathParts:    []string{"foo", "bar", "hello"},
				HandlerParts: []string{"world", "func1", "func2"},
			},
		},
		{
			name: "dot in path",
			path: "foo.bar/hello.world",
			result: HandlerPath{
				PathParts:    []string{"foo.bar", "hello"},
				HandlerParts: []string{"world"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SplitHandlerPath(tc.path)
			require.ElementsMatch(t, tc.result.PathParts, result.PathParts)
			require.ElementsMatch(t, tc.result.HandlerParts, result.HandlerParts)
		})
	}
}
