package astra

import (
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestService_Setup(t *testing.T) {
	t.Run("sets up astra directory", func(t *testing.T) {
		service := &Service{}

		stat, err := os.Stat("./.astra")
		require.Error(t, err)
		require.Nil(t, stat)

		err = service.Setup()
		require.NoError(t, err)
		defer func() {
			err = service.ClearCache()
			require.NoError(t, err)
		}()

		stat, err = os.Stat(service.getAstraDirPath())
		require.NoError(t, err)
		require.NotNil(t, stat)
	})

	t.Run("sets up working directory if empty", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			service := &Service{}

			require.Empty(t, service.WorkDir)

			err := service.Setup()
			require.NoError(t, err)
			defer func() {
				err = service.ClearCache()
				require.NoError(t, err)
			}()

			require.NotEmpty(t, service.WorkDir)

			cwd, err := os.Getwd()
			require.NoError(t, err)

			require.Equal(t, cwd, service.WorkDir)
		})

		t.Run("not empty", func(t *testing.T) {
			service := &Service{
				WorkDir: "test",
			}

			require.NotEmpty(t, service.WorkDir)

			err := service.Setup()
			require.NoError(t, err)
			defer func() {
				err = service.ClearCache()
				require.NoError(t, err)
			}()

			require.NotEmpty(t, service.WorkDir)
			require.Equal(t, "test", service.WorkDir)
		})
	})

	t.Run("caches service if enabled", func(t *testing.T) {
		t.Run("enabled", func(t *testing.T) {
			service := &Service{
				CacheEnabled: true,
			}

			err := service.Setup()
			require.NoError(t, err)
			defer func() {
				err = service.ClearCache()
				require.NoError(t, err)
			}()

			err = service.Generate()
			require.NoError(t, err)

			stat, err := os.Stat(path.Join(service.getAstraDirPath(), cacheFileName))
			require.NoError(t, err)
			require.NotNil(t, stat)
		})

		t.Run("disabled", func(t *testing.T) {
			service := &Service{
				CacheEnabled: false,
			}

			err := service.Setup()
			require.NoError(t, err)
			defer func() {
				err = service.ClearCache()
				require.NoError(t, err)
			}()

			stat, err := os.Stat(path.Join(service.getAstraDirPath(), cacheFileName))
			require.Error(t, err)
			require.Nil(t, stat)
		})
	})
}
