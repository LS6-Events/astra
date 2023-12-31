package astra

import (
	"encoding/json"
	"github.com/Jeffail/gabs/v2"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"testing"
)

func TestService_Cache(t *testing.T) {
	t.Run("creates a cache in json format", func(t *testing.T) {
		s := Service{
			CacheEnabled: true,
			CachePath:    "./test-cache.json",
			Components: []Field{
				{
					Name:    "TestComponent",
					Package: "test",
				},
			},
		}

		err := s.Cache()
		require.NoError(t, err)
		defer func() {
			err := os.Remove("./test-cache.json")
			require.NoError(t, err)
		}()

		jsonBlob, err := gabs.ParseJSONFile("./test-cache.json")
		require.NoError(t, err)

		require.Equal(t, "TestComponent", jsonBlob.Path("components.0.name").Data().(string))
		require.Equal(t, "test", jsonBlob.Path("components.0.package").Data().(string))
	})

	t.Run("creates a cache file in yaml format", func(t *testing.T) {
		s := Service{
			CacheEnabled: true,
			CachePath:    "./test-cache.yaml",
			Components: []Field{
				{
					Name:    "TestComponent",
					Package: "test",
				},
			},
		}

		err := s.Cache()
		require.NoError(t, err)
		defer func() {
			err := os.Remove("./test-cache.yaml")
			require.NoError(t, err)
		}()

		yamlContent, err := os.ReadFile("./test-cache.yaml")
		require.NoError(t, err)

		var decodedYaml *Service
		err = yaml.Unmarshal(yamlContent, &decodedYaml)
		require.NoError(t, err)

		require.Equal(t, "TestComponent", decodedYaml.Components[0].Name)
		require.Equal(t, "test", decodedYaml.Components[0].Package)
	})

	t.Run("creates cache in default path (.astra/cache.json)", func(t *testing.T) {
		s := Service{
			CacheEnabled: true,
			Components: []Field{
				{
					Name:    "TestComponent",
					Package: "test",
				},
			},
		}

		err := s.setupAstraDir()
		require.NoError(t, err)

		err = s.Cache()
		require.NoError(t, err)

		jsonBlob, err := gabs.ParseJSONFile("./.astra/cache.json")
		require.NoError(t, err)

		require.Equal(t, "TestComponent", jsonBlob.Path("components.0.name").Data().(string))
		require.Equal(t, "test", jsonBlob.Path("components.0.package").Data().(string))
	})
}

func TestService_LoadCache(t *testing.T) {
	setupCache := func(path string, t *testing.T) {
		t.Helper()

		s := Service{
			Components: []Field{
				{
					Name:    "TestComponent",
					Package: "test",
				},
			},
		}

		var buff []byte
		var err error
		if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
			buff, err = yaml.Marshal(s)
		} else {
			buff, err = json.Marshal(s)
		}

		require.NoError(t, err)

		err = os.WriteFile(path, buff, 0644)
		require.NoError(t, err)
	}

	cleanupCache := func(path string, t *testing.T) {
		t.Helper()

		err := os.Remove(path)
		require.NoError(t, err)
	}

	t.Run("loads cache from specified json path", func(t *testing.T) {
		service := &Service{
			CacheEnabled: true,
			CachePath:    "./test-cache.json",
		}

		setupCache("./test-cache.json", t)
		defer cleanupCache("./test-cache.json", t)

		require.Empty(t, service.Components)

		err := service.LoadCache()
		require.NoError(t, err)

		require.Equal(t, "TestComponent", service.Components[0].Name)
		require.Equal(t, "test", service.Components[0].Package)
	})

	t.Run("loads cache from specified yaml path", func(t *testing.T) {
		service := &Service{
			CacheEnabled: true,
			CachePath:    "./test-cache.yaml",
		}

		setupCache("./test-cache.yaml", t)
		defer cleanupCache("./test-cache.yaml", t)

		require.Empty(t, service.Components)

		err := service.LoadCache()
		require.NoError(t, err)

		require.Equal(t, "TestComponent", service.Components[0].Name)
		require.Equal(t, "test", service.Components[0].Package)
	})

	t.Run("loads cache from default path (.astra/cache.json)", func(t *testing.T) {
		service := &Service{
			CacheEnabled: true,
		}

		err := service.setupAstraDir()
		require.NoError(t, err)

		setupCache("./.astra/cache.json", t)
		defer cleanupCache("./.astra/cache.json", t)

		require.Empty(t, service.Components)

		err = service.LoadCache()
		require.NoError(t, err)

		require.Equal(t, "TestComponent", service.Components[0].Name)
		require.Equal(t, "test", service.Components[0].Package)
	})
}

func TestService_LoadCacheFromCustomPath(t *testing.T) {
	setupCache := func(path string, toCache Service, t *testing.T) {
		t.Helper()

		var buff []byte
		var err error
		if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
			buff, err = yaml.Marshal(toCache)
		} else {
			buff, err = json.Marshal(toCache)
		}

		require.NoError(t, err)

		err = os.WriteFile(path, buff, 0644)
		require.NoError(t, err)
	}

	cleanupCache := func(path string, t *testing.T) {
		t.Helper()

		err := os.Remove(path)
		require.NoError(t, err)
	}

	t.Run("loads cache from specified json path", func(t *testing.T) {
		topLevelService := &Service{}

		cachedService := Service{
			Components: []Field{
				{
					Name:    "TestComponent",
					Package: "test",
				},
			},
		}

		setupCache("./test-cache.json", cachedService, t)
		defer cleanupCache("./test-cache.json", t)

		require.Empty(t, topLevelService.Components)

		err := topLevelService.LoadCacheFromCustomPath("./test-cache.json")
		require.NoError(t, err)

		require.Equal(t, "TestComponent", topLevelService.Components[0].Name)
		require.Equal(t, "test", topLevelService.Components[0].Package)
	})

	t.Run("loads cache from specified yaml path", func(t *testing.T) {
		topLevelService := &Service{}

		cachedService := Service{
			Components: []Field{
				{
					Name:    "TestComponent",
					Package: "test",
				},
			},
		}

		setupCache("./test-cache.yaml", cachedService, t)
		defer cleanupCache("./test-cache.yaml", t)

		require.Empty(t, topLevelService.Components)

		err := topLevelService.LoadCacheFromCustomPath("./test-cache.yaml")
		require.NoError(t, err)

		require.Equal(t, "TestComponent", topLevelService.Components[0].Name)
		require.Equal(t, "test", topLevelService.Components[0].Package)
	})

	t.Run("returns an error if the file format is not yaml or json", func(t *testing.T) {
		topLevelService := &Service{}

		cachedService := Service{
			Components: []Field{
				{
					Name:    "TestComponent",
					Package: "test",
				},
			},
		}

		setupCache("./test-cache.txt", cachedService, t)
		defer cleanupCache("./test-cache.txt", t)

		err := topLevelService.LoadCacheFromCustomPath("./test-cache.txt")
		require.Error(t, err)
	})

	t.Run("returns an error if the file does not exist", func(t *testing.T) {
		topLevelService := &Service{}

		err := topLevelService.LoadCacheFromCustomPath("./test-cache-should-not-exist.txt")
		require.Error(t, err)
	})

	t.Run("sets the correct services fields", func(t *testing.T) {
		t.Run("Inputs", func(t *testing.T) {
			topLevelService := &Service{}

			cachedService := Service{
				Inputs: []Input{
					{
						Mode: InputMode("gin"),
					},
				},
			}

			setupCache("./test-cache.json", cachedService, t)
			defer cleanupCache("./test-cache.json", t)

			require.Empty(t, topLevelService.Inputs)

			err := topLevelService.LoadCacheFromCustomPath("./test-cache.json")
			require.NoError(t, err)

			require.Equal(t, InputMode("gin"), topLevelService.Inputs[0].Mode)
		})

		t.Run("Outputs", func(t *testing.T) {
			topLevelService := &Service{}

			cachedService := Service{
				Outputs: []Output{
					{
						Mode: OutputMode("openapi"),
						Configuration: IOConfiguration{
							IOConfigurationKeyFilePath: "test-openapi.json",
						},
					},
				},
			}

			setupCache("./test-cache.json", cachedService, t)
			defer cleanupCache("./test-cache.json", t)

			require.Empty(t, topLevelService.Outputs)

			err := topLevelService.LoadCacheFromCustomPath("./test-cache.json")
			require.NoError(t, err)

			require.Equal(t, OutputMode("openapi"), topLevelService.Outputs[0].Mode)
			require.Equal(t, "test-openapi.json", topLevelService.Outputs[0].Configuration[IOConfigurationKeyFilePath])
		})

		t.Run("Config", func(t *testing.T) {
			topLevelService := &Service{}

			cachedService := Service{
				Config: &Config{
					BasePath: "/test-base-path",
					Port:     4000,
					Host:     "google.com",
				},
			}

			setupCache("./test-cache.json", cachedService, t)
			defer cleanupCache("./test-cache.json", t)

			require.Nil(t, topLevelService.Config)

			err := topLevelService.LoadCacheFromCustomPath("./test-cache.json")
			require.NoError(t, err)

			require.Equal(t, "/test-base-path", topLevelService.Config.BasePath)
			require.Equal(t, 4000, topLevelService.Config.Port)
			require.Equal(t, "google.com", topLevelService.Config.Host)
		})

		t.Run("Routes", func(t *testing.T) {
			topLevelService := &Service{}

			cachedService := Service{
				Routes: []Route{
					{
						Method: "GET",
						Path:   "/test-path",
					},
				},
			}

			setupCache("./test-cache.json", cachedService, t)
			defer cleanupCache("./test-cache.json", t)

			require.Empty(t, topLevelService.Routes)

			err := topLevelService.LoadCacheFromCustomPath("./test-cache.json")
			require.NoError(t, err)

			require.Equal(t, "GET", topLevelService.Routes[0].Method)
			require.Equal(t, "/test-path", topLevelService.Routes[0].Path)
		})

		t.Run("Components", func(t *testing.T) {
			topLevelService := &Service{}

			cachedService := Service{
				Components: []Field{
					{
						Name:    "TestComponent",
						Package: "test",
					},
				},
			}

			setupCache("./test-cache.json", cachedService, t)
			defer cleanupCache("./test-cache.json", t)

			require.Empty(t, topLevelService.Components)

			err := topLevelService.LoadCacheFromCustomPath("./test-cache.json")
			require.NoError(t, err)

			require.Equal(t, "TestComponent", topLevelService.Components[0].Name)
			require.Equal(t, "test", topLevelService.Components[0].Package)
		})
	})
}

func TestService_ClearCache(t *testing.T) {
	setupCache := func(t *testing.T, path string) {
		t.Helper()

		service := Service{}

		buff, err := json.Marshal(service)
		require.NoError(t, err)

		err = os.WriteFile(path, buff, 0644)
	}

	t.Run("clears the specified cache file", func(t *testing.T) {
		service := &Service{
			CacheEnabled: true,
			CachePath:    "./test-cache.json",
		}

		setupCache(t, "./test-cache.json")

		err := service.ClearCache()
		require.NoError(t, err)

		_, err = os.Stat("./test-cache.json")
		require.Error(t, err)
	})

	t.Run("clears the default cache file", func(t *testing.T) {
		service := &Service{
			CacheEnabled: true,
		}

		err := service.setupAstraDir()
		require.NoError(t, err)

		setupCache(t, "./.astra/cache.json")

		err = service.ClearCache()
		require.NoError(t, err)

		_, err = os.Stat("./.astra/cache.json")
		require.Error(t, err)
	})
}
