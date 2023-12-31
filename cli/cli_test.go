package cli

import (
	"github.com/ls6-events/astra"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithCLI(t *testing.T) {
	service := &astra.Service{}

	require.False(t, service.CacheEnabled)
	require.Empty(t, service.CLIMode)

	WithCLI()(service)

	require.True(t, service.CacheEnabled)
	require.Equal(t, astra.CLIModeSetup, service.CLIMode)
}

func TestWithCLIBuilder(t *testing.T) {
	service := &astra.Service{}

	require.Empty(t, service.CLIMode)

	WithCLIBuilder()(service)

	require.Equal(t, astra.CLIModeBuilder, service.CLIMode)
}
