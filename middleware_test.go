package astra

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_UnstableWithMiddleware(t *testing.T) {
	service := &Service{
		UnstableEnableMiddleware: false,
	}

	UnstableWithMiddleware()(service)

	require.True(t, service.UnstableEnableMiddleware)
}
