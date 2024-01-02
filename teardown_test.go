package astra

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestService_Teardown(t *testing.T) {
	t.Run("tears down astra directory when cache is disabled and cli mode is none", func(t *testing.T) {
		service := &Service{
			CacheEnabled: false,
			CLIMode:      CLIModeNone,
		}

		err := service.Setup()
		require.NoError(t, err)

		stat, err := os.Stat(service.getAstraDirPath())
		require.NoError(t, err)
		require.NotNil(t, stat)

		err = service.Teardown()
		require.NoError(t, err)

		stat, err = os.Stat(service.getAstraDirPath())
		require.Error(t, err)
		require.Nil(t, stat)
	})

	t.Run("leaves astra directory when cache is enabled", func(t *testing.T) {
		service := &Service{
			CacheEnabled: true,
			CLIMode:      CLIModeNone,
		}

		err := service.Setup()
		require.NoError(t, err)

		stat, err := os.Stat(service.getAstraDirPath())
		require.NoError(t, err)
		require.NotNil(t, stat)

		err = service.Teardown()
		require.NoError(t, err)

		stat, err = os.Stat(service.getAstraDirPath())
		require.NoError(t, err)
		require.NotNil(t, stat)
	})
}
