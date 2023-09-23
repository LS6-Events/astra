package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/utils"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strconv"
	"strings"
)

// Generate the OpenAPI output
// It will marshal the OpenAPI struct and write it to a file
// It will also generate the paths and their operations
// It will also generate the components and their schemas
func Generate(filePath string) astra.ServiceFunction {
	return func(s *astra.Service) error {
		s.Log.Debug().Msg("Generating OpenAPI output")
		if s.Config == nil {
			s.Log.Error().Msg("No config found")
			return astra.ErrConfigNotFound
		}

		protocol := "http"
		if s.Config.Secure {
			protocol += "s"
		}

		paths := make(Paths)
		s.Log.Debug().Msg("Adding paths")
		for _, endpoint := range s.Routes {
			s.Log.Debug().Str("endpointPath", endpoint.Path).Str("method", endpoint.Method).Msg("Generating endpoint")

			endpoint.Path = utils.MapPathParams(endpoint.Path, func(param string) string {
				if param[0] == ':' {
					return fmt.Sprintf("{%s}", param[1:])
				} else {
					return fmt.Sprintf("{%s*}", param[1:])
				}
			})

			operation := Operation{
				Responses: make(map[string]Response),
			}

			for _, pathParam := range endpoint.PathParams {
				s.Log.Debug().Str("endpointPath", endpoint.Path).Str("method", endpoint.Method).Str("param", pathParam.Name).Msg("Adding endpointPath parameter")
				operation.Parameters = append(operation.Parameters, Parameter{
					Name:     pathParam.Name,
					In:       "endpointPath",
					Required: pathParam.IsRequired,
				})
			}

			for _, queryParam := range endpoint.QueryParams {
				s.Log.Debug().Str("endpointPath", endpoint.Path).Str("method", endpoint.Method).Str("param", queryParam.Name).Msg("Adding query parameter")
				parameter := Parameter{
					Name:     queryParam.Name,
					In:       "query",
					Required: queryParam.IsRequired,
					Explode:  true,
					Style:    "form",
					Schema:   mapParamToSchema(queryParam),
				}

				operation.Parameters = append(operation.Parameters, parameter)
			}

			for _, bodyParam := range endpoint.Body {
				s.Log.Debug().Str("endpointPath", endpoint.Path).Str("method", endpoint.Method).Str("param", bodyParam.Name).Msg("Adding body parameter")
				var mediaType MediaType
				schema := mapFieldToSchema(bodyParam.Field)
				if bodyParam.Name != "" {
					mediaType.Schema = Schema{
						Type: "object",
						Properties: map[string]Schema{
							bodyParam.Name: schema,
						},
					}
				} else {
					mediaType.Schema = schema
				}

				operation.RequestBody = &RequestBody{
					Content: map[string]MediaType{
						endpoint.BodyType: mediaType,
					},
				}
			}

			for _, returnType := range endpoint.ReturnTypes {
				s.Log.Debug().Str("endpointPath", endpoint.Path).Str("method", endpoint.Method).Str("return", returnType.Field.Name).Msg("Adding return type")
				var mediaType MediaType
				mediaType.Schema = mapFieldToSchema(returnType.Field)

				var content map[string]MediaType
				if mediaType.Schema.Type != "" || mediaType.Schema.Ref != "" {
					content = map[string]MediaType{
						endpoint.ContentType: mediaType,
					}
				} else {
					content = nil
				}

				operation.Responses[strconv.Itoa(returnType.StatusCode)] = Response{
					Description: "",
					Headers:     nil,
					Content:     content,
					Links:       nil,
				}
			}

			if endpoint.Doc != "" {
				operation.Description = endpoint.Doc
			}

			var endpointPath Path
			if _, ok := paths[endpoint.Path]; !ok {
				endpointPath = Path{}
			} else {
				endpointPath = paths[endpoint.Path]
			}
			switch endpoint.Method {
			case "GET":
				endpointPath.Get = &operation
			case "POST":
				endpointPath.Post = &operation
			case "PUT":
				endpointPath.Put = &operation
			case "PATCH":
				endpointPath.Patch = &operation
			case "DELETE":
				endpointPath.Delete = &operation
			case "HEAD":
				endpointPath.Head = &operation
			case "OPTIONS":
				endpointPath.Options = &operation
			}

			paths[endpoint.Path] = endpointPath
			s.Log.Debug().Str("path", endpoint.Path).Str("method", endpoint.Method).Msg("Added path")
		}
		s.Log.Debug().Msg("Added paths")

		components := Components{
			Schemas: make(map[string]Schema),
		}

		s.Log.Debug().Msg("Adding components")
		for _, component := range s.Components {
			s.Log.Debug().Str("name", component.Name).Msg("Adding component")
			var schema Schema

			if component.Type == "struct" {
				embeddedProperties := make([]Schema, 0)
				schema = Schema{
					Type:       "object",
					Properties: make(map[string]Schema),
				}
				for key, field := range component.StructFields {
					// We should aim to use doc comments in the future
					// However https://github.com/OAI/OpenAPI-Specification/issues/1514
					if field.IsEmbedded {
						embeddedProperties = append(embeddedProperties, Schema{
							Ref: makeComponentRef(field.Type, field.Package),
						})
						continue
					}
					if !astra.IsAcceptedType(field.Type) {
						schema.Properties[key] = Schema{
							Ref: makeComponentRef(field.Type, field.Package),
						}
					} else {
						schema.Properties[key] = mapAcceptedType(field.Type)
					}
				}

				if len(embeddedProperties) > 0 {
					if len(schema.Properties) == 0 {
						schema.AllOf = embeddedProperties
					} else {
						schema.AllOf = append(embeddedProperties, Schema{
							Properties: schema.Properties,
						})

						schema.Properties = nil
					}
				}
			} else if component.Type == "slice" {
				itemSchema := mapAcceptedType(component.SliceType)
				schema = Schema{
					Type:  "array",
					Items: &itemSchema,
				}
				if !astra.IsAcceptedType(component.SliceType) {
					schema.Items = &Schema{
						Ref: makeComponentRef(component.SliceType, component.Package),
					}
				}
			} else if component.Type == "map" {
				var additionalProperties Schema

				if !astra.IsAcceptedType(component.MapValueType) {
					additionalProperties.Ref = makeComponentRef(component.MapValueType, component.Package)
				} else {
					additionalProperties = mapAcceptedType(component.MapValueType)
				}

				schema = Schema{
					Type:                 "object",
					AdditionalProperties: &additionalProperties,
				}
			} else {
				schema = mapAcceptedType(component.Type)
			}

			if component.Doc != "" {
				schema.Description = component.Doc
			}

			components.Schemas[makeComponentRefName(component.Name, component.Package)] = schema
		}
		s.Log.Debug().Msg("Added components")

		if s.Config.Description == "" {
			s.Config.Description = "Generated by astra"
		}

		s.Log.Debug().Msg("Generating OpenAPI schema file")
		output := OpenAPISchema{
			OpenAPI: "3.0.0",
			Info: Info{
				Title:       s.Config.Title,
				Description: s.Config.Description,
				Contact:     Contact(s.Config.Contact),
				License:     License(s.Config.License),
				Version:     s.Config.Version,
			},
			Servers: []Server{
				{
					URL: fmt.Sprintf("%s://%s:%d%s", protocol, s.Config.Host, s.Config.Port, s.Config.BasePath),
				},
			},
			Paths:      paths,
			Components: components,
		}

		if !strings.HasSuffix(filePath, ".json") && !strings.HasSuffix(filePath, ".yaml") && !strings.HasSuffix(filePath, ".yml") {
			s.Log.Debug().Msg("No file extension provided, defaulting to .json")
			filePath += ".json"
		}

		var file []byte
		var err error
		if strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml") {
			s.Log.Debug().Msg("Writing YAML file")
			file, err = yaml.Marshal(output)
		} else {
			s.Log.Debug().Msg("Writing JSON file")
			file, err = json.MarshalIndent(output, "", "  ")
		}
		if err != nil {
			s.Log.Error().Err(err).Msg("Failed to marshal OpenAPI schema")
			return err
		}

		filePath = path.Join(s.WorkDir, filePath)
		err = os.WriteFile(filePath, file, 0644)
		if err != nil {
			s.Log.Error().Err(err).Msg("Failed to write OpenAPI schema file")
			return err
		}

		s.Log.Debug().Str("filePath", filePath).Msg("Successfully generated OpenAPI schema file")

		return nil
	}
}
