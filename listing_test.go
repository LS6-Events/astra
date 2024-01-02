package astra

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"regexp"
	"strings"
	"testing"
)

func TestWithPathDenyList(t *testing.T) {
	service := &Service{
		PathDenyList: []func(string) bool{},
	}

	WithPathDenyList("test")(service)

	require.Equal(t, 1, len(service.PathDenyList))
	require.True(t, service.PathDenyList[0]("test"))
	require.False(t, service.PathDenyList[0]("test2"))
}

func TestWithPathDenyListRegex(t *testing.T) {
	service := &Service{
		PathDenyList: []func(string) bool{},
	}

	reg := regexp.MustCompile("test-[0-9]")

	WithPathDenyListRegex(reg)(service)

	require.Equal(t, 1, len(service.PathDenyList))
	for i := 0; i < 10; i++ {
		require.True(t, service.PathDenyList[0](fmt.Sprintf("test-%d", i)))
	}
	require.False(t, service.PathDenyList[0]("test"))
	require.False(t, service.PathDenyList[0]("test-a"))
}

func TestWithPathDenyListFunc(t *testing.T) {
	service := &Service{
		PathDenyList: []func(string) bool{},
	}

	WithPathDenyListFunc(func(path string) bool {
		return strings.HasPrefix(path, "test")
	})(service)

	require.Equal(t, 1, len(service.PathDenyList))
	require.True(t, service.PathDenyList[0]("test"))
	require.True(t, service.PathDenyList[0]("test2"))
	require.False(t, service.PathDenyList[0]("2test"))
}
