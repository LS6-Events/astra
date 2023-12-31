package astra

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsAcceptedType(t *testing.T) {
	for _, acceptedType := range AcceptedTypes {
		require.True(t, IsAcceptedType(acceptedType))
	}
}
