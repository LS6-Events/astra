package outputs

import (
	"github.com/ls6-events/astra"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithAzureFunctionsOutput(t *testing.T) {
	service := &astra.Service{}

	require.Len(t, service.Outputs, 0)

	WithAzureFunctionsOutput("./")(service)

	require.Len(t, service.Outputs, 1)
}

func TestWithJSONOutput(t *testing.T) {
	service := &astra.Service{}

	require.Len(t, service.Outputs, 0)

	WithJSONOutput("./")(service)

	require.Len(t, service.Outputs, 1)
}

func TestWithOpenAPIOutput(t *testing.T) {
	service := &astra.Service{}

	require.Len(t, service.Outputs, 0)

	WithOpenAPIOutput("./")(service)

	require.Len(t, service.Outputs, 1)
}
