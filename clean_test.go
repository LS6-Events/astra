package astra

import (
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestService_Clean(t *testing.T) {
	t.Run("cleans up components", func(t *testing.T) {
		service := &Service{
			Components: []Field{
				{
					Name:    "Duration",
					Package: "time",
				},
			},
		}

		err := service.Setup()
		require.NoError(t, err)

		require.Empty(t, service.Components[0].Type)

		err = service.Clean()
		require.NoError(t, err)

		require.Equal(t, "int", service.Components[0].Type)
	})

	t.Run("cleans up routes", func(t *testing.T) {
		t.Run("cleans up return types", func(t *testing.T) {
			service := &Service{
				Routes: []Route{
					{
						ReturnTypes: []ReturnType{
							{
								StatusCode:  200,
								ContentType: "application/json",
								Field: Field{
									Name:    "Duration",
									Package: "time",
								},
							},
						},
					},
				},
			}

			err := service.Setup()
			require.NoError(t, err)

			require.Empty(t, service.Routes[0].ReturnTypes[0].Field.Type)

			err = service.Clean()
			require.NoError(t, err)

			require.Equal(t, "int", service.Routes[0].ReturnTypes[0].Field.Type)
		})

		t.Run("cleans up path parameters", func(t *testing.T) {
			service := &Service{
				Routes: []Route{
					{
						PathParams: []Param{
							{
								Name: "duration",
								Field: Field{
									Name:    "Duration",
									Package: "time",
								},
							},
						},
					},
				},
			}

			err := service.Setup()
			require.NoError(t, err)

			require.Empty(t, service.Routes[0].PathParams[0].Field.Type)

			err = service.Clean()
			require.NoError(t, err)

			require.Equal(t, "int", service.Routes[0].PathParams[0].Field.Type)
		})

		t.Run("cleans up query parameters", func(t *testing.T) {
			service := &Service{
				Routes: []Route{
					{
						QueryParams: []Param{
							{
								Name: "duration",
								Field: Field{
									Name:    "Duration",
									Package: "time",
								},
							},
						},
					},
				},
			}

			err := service.Setup()
			require.NoError(t, err)

			require.Empty(t, service.Routes[0].QueryParams[0].Field.Type)

			err = service.Clean()
			require.NoError(t, err)

			require.Equal(t, "int", service.Routes[0].QueryParams[0].Field.Type)
		})

		t.Run("cleans up body parameters", func(t *testing.T) {
			service := &Service{
				Routes: []Route{
					{
						Body: []BodyParam{
							{
								Name: "duration",
								Field: Field{
									Name:    "Duration",
									Package: "time",
								},
							},
						},
					},
				},
			}

			err := service.Setup()
			require.NoError(t, err)

			require.Empty(t, service.Routes[0].Body[0].Field.Type)

			err = service.Clean()
			require.NoError(t, err)

			require.Equal(t, "int", service.Routes[0].Body[0].Field.Type)
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

			err = service.Clean()
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

			err = service.Clean()
			require.NoError(t, err)

			_, err = os.Stat(path.Join(service.getAstraDirPath(), cacheFileName))
			require.Error(t, err)
		})
	})
}

func TestService_cleanField(t *testing.T) {
	t.Run("does not do anything if the field is not a substitute type nor the main field doesn't need changing", func(t *testing.T) {
		service := &Service{}

		field := Field{
			Name:    "TestType",
			Package: "test",
			Type:    "string",
		}

		newField := service.cleanField(field, "main")

		require.Equal(t, "test", newField.Package)
		require.Equal(t, "TestType", newField.Name)
		require.Equal(t, "string", newField.Type)
	})

	t.Run("clean substitute types", func(t *testing.T) {
		service := &Service{}

		field := Field{
			Name:    "Duration",
			Package: "time",
		}

		newField := service.cleanField(field, "main")

		require.Equal(t, "int", newField.Type)
		require.Equal(t, "time", newField.Package)
		require.Equal(t, "Duration", newField.Name)
	})

	t.Run("handles main package name for package", func(t *testing.T) {
		service := &Service{}

		field := Field{
			Name:    "TestType",
			Package: "not-main",
			Type:    "string",
		}

		newField := service.cleanField(field, "not-main")

		require.Equal(t, "main", newField.Package)
		require.Equal(t, "TestType", newField.Name)
		require.Equal(t, "string", newField.Type)
	})

	t.Run("handles main package name for map key package", func(t *testing.T) {
		service := &Service{}

		field := Field{
			Name:          "TestType",
			MapKeyPackage: "not-main",
			Type:          "map",
			MapKeyType:    "TestKey",
			MapValueType:  "string",
		}

		newField := service.cleanField(field, "not-main")

		require.Equal(t, "main", newField.MapKeyPackage)
		require.Equal(t, "TestType", newField.Name)
		require.Equal(t, "map", newField.Type)
		require.Equal(t, "TestKey", newField.MapKeyType)
		require.Equal(t, "string", newField.MapValueType)
	})

	t.Run("runs recursively for struct fields", func(t *testing.T) {
		t.Run("does not do anything if the field is not a substitute type nor the main field doesn't need changing", func(t *testing.T) {
			service := &Service{}

			field := Field{
				Name:    "TestStruct",
				Package: "example",
				Type:    "struct",
				StructFields: map[string]Field{
					"Test": {
						Name:    "TestType",
						Package: "test",
						Type:    "string",
					},
				},
			}

			newField := service.cleanField(field, "main")

			require.Equal(t, "example", newField.Package)
			require.Equal(t, "TestStruct", newField.Name)
			require.Equal(t, "struct", newField.Type)
			require.Equal(t, "test", newField.StructFields["Test"].Package)
			require.Equal(t, "TestType", newField.StructFields["Test"].Name)
			require.Equal(t, "string", newField.StructFields["Test"].Type)
		})

		t.Run("cleans substitute types", func(t *testing.T) {
			service := &Service{}

			field := Field{
				Name:    "TestStruct",
				Package: "example",
				Type:    "struct",
				StructFields: map[string]Field{
					"Test": {
						Name:    "Duration",
						Package: "time",
					},
				},
			}

			newField := service.cleanField(field, "main")

			require.Equal(t, "struct", newField.Type)
			require.Equal(t, "example", newField.Package)
			require.Equal(t, "TestStruct", newField.Name)
			require.Equal(t, "int", newField.StructFields["Test"].Type)
			require.Equal(t, "time", newField.StructFields["Test"].Package)
			require.Equal(t, "Duration", newField.StructFields["Test"].Name)
		})

		t.Run("handles main package name for package", func(t *testing.T) {
			service := &Service{}

			field := Field{
				Name:    "TestStruct",
				Package: "example",
				Type:    "struct",
				StructFields: map[string]Field{
					"Test": {
						Name:    "TestType",
						Package: "not-main",
						Type:    "string",
					},
				},
			}

			newField := service.cleanField(field, "not-main")

			require.Equal(t, "main", newField.StructFields["Test"].Package)
			require.Equal(t, "TestType", newField.StructFields["Test"].Name)
			require.Equal(t, "string", newField.StructFields["Test"].Type)
		})

		t.Run("handles main package name for map key package", func(t *testing.T) {
			service := &Service{}

			field := Field{
				Name:    "TestStruct",
				Package: "example",
				Type:    "struct",
				StructFields: map[string]Field{
					"Test": {
						Name:          "TestType",
						MapKeyPackage: "not-main",
						Type:          "map",
						MapKeyType:    "TestKey",
						MapValueType:  "string",
					},
				},
			}

			newField := service.cleanField(field, "not-main")

			require.Equal(t, "main", newField.StructFields["Test"].MapKeyPackage)
			require.Equal(t, "TestType", newField.StructFields["Test"].Name)
			require.Equal(t, "map", newField.StructFields["Test"].Type)
			require.Equal(t, "TestKey", newField.StructFields["Test"].MapKeyType)
			require.Equal(t, "string", newField.StructFields["Test"].MapValueType)
		})
	})
}
