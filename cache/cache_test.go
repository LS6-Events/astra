package cache

import (
	"github.com/ls6-events/astra"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithCache(t *testing.T) {
	service := &astra.Service{}

	require.False(t, service.CacheEnabled)

	WithCache()(service)

	require.True(t, service.CacheEnabled)
}

func TestWithCustomCachePath(t *testing.T) {
	service := &astra.Service{}

	require.False(t, service.CacheEnabled)
	require.Empty(t, service.CachePath)

	WithCustomCachePath("test")(service)

	require.True(t, service.CacheEnabled)
	require.Equal(t, "test", service.CachePath)
}
