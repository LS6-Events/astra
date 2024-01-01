package astra

import (
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestService_ParseRoutes(t *testing.T) {
	t.Run("parse routes is called and create routes is not", func(t *testing.T) {
		var createRoutesCalled bool
		var parseRoutesCalled bool
		service := &Service{
			Inputs: []Input{
				{
					Mode: "",
					CreateRoutes: func(service *Service) error {
						createRoutesCalled = true
						return nil
					},
					ParseRoutes: func(service *Service) error {
						parseRoutesCalled = true
						return nil
					},
				},
			},
		}

		err := service.ParseRoutes()
		require.NoError(t, err)
		require.False(t, createRoutesCalled)
		require.True(t, parseRoutesCalled)
	})

	t.Run("caches service if enabled", func(t *testing.T) {
		t.Run("enabled", func(t *testing.T) {
			service := &Service{
				CacheEnabled: true,
				Inputs: []Input{
					{
						ParseRoutes: func(service *Service) error {
							return nil
						},
					},
				},
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
				Inputs: []Input{
					{
						ParseRoutes: func(service *Service) error {
							return nil
						},
					},
				},
			}

			err := service.Setup()
			require.NoError(t, err)

			err = service.ParseRoutes()
			require.NoError(t, err)

			stat, err := os.Stat(path.Join(service.getAstraDirPath(), cacheFileName))
			require.Error(t, err)
			require.Nil(t, stat)
		})
	})
}
