package astra

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithCustomWorkDir(t *testing.T) {
	service := &Service{}

	require.NotEqual(t, "/test-work-dir", service.WorkDir)

	WithCustomWorkDir("/test-work-dir")(service)

	require.Equal(t, "/test-work-dir", service.WorkDir)
}
