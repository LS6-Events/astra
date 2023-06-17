package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/ls6-events/gengo"
	"github.com/ls6-events/gengo/utils"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

func generate(filePath string) gengo.GenerateFunction {
	return func(s *gengo.Service) error {
		s.Log.Debug().Msg("Generating OpenAPI output")
		if s.Config == nil {
			s.Log.Error().Msg("No config found")
			return gengo.ErrConfigNotFound
		}

		protocol := "http"
		if s.Config.Secure {
			protocol += "s"
		}

		paths := make(Paths)
		s.Log.Debug().Msg("Adding paths")
		for _, endpoint := range s.Routes {
			s.Log.Debug().Str("path", endpoint.Path).Str("method", endpoint.Method).Msg("Generating endpoint")

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
				s.Log.Debug().Str("path", endpoint.Path).Str("method", endpoint.Method).Str("param", pathParam.Name).Msg("Adding path parameter")
				operation.Parameters = append(operation.Parameters, Parameter{
					Name:     pathParam.Name,
					In:       "path",
					Required: pathParam.IsRequired,
				})
			}

			for _, queryParam := range endpoint.QueryParams {
				s.Log.Debug().Str("path", endpoint.Path).Str("method", endpoint.Method).Str("param", queryParam.Name).Msg("Adding query parameter")
				parameter := Parameter{
					Name:     queryParam.Name,
					In:       "query",
					Required: queryParam.IsRequired,
					Explode:  true,
					Style:    "form",
				}

				if queryParam.IsBound {
					parameter.Schema = Schema{
						Ref: makeComponentRef(queryParam.Type, queryParam.Package),
					}
				} else if queryParam.IsArray {
					itemSchema := Schema{
						Type: mapAcceptedType(queryParam.Type).Type,
					}
					if !gengo.IsAcceptedType(queryParam.Type) {
						itemSchema = Schema{
							Ref: makeComponentRef(queryParam.Type, queryParam.Package),
						}
					}
					parameter.Schema = Schema{
						Type:  "array",
						Items: &itemSchema,
					}
				} else if queryParam.IsMap {
					var additionalProperties Schema
					if !gengo.IsAcceptedType(queryParam.Type) {
						additionalProperties.Ref = makeComponentRef(queryParam.Type, queryParam.Package)
					} else {
						additionalProperties.Type = mapAcceptedType(queryParam.Type).Type
					}
					parameter.Schema = Schema{
						Type:                 "object",
						AdditionalProperties: &additionalProperties,
					}
				} else {
					parameter.Schema = Schema{
						Type: queryParam.Type,
					}
				}

				operation.Parameters = append(operation.Parameters, parameter)
			}

			for _, bodyParam := range endpoint.Body {
				s.Log.Debug().Str("path", endpoint.Path).Str("method", endpoint.Method).Str("param", bodyParam.Name).Msg("Adding body parameter")
				var mediaType MediaType
				if bodyParam.IsBound {
					mediaType.Schema = Schema{
						Ref: makeComponentRef(bodyParam.Type, bodyParam.Package),
					}
				} else if bodyParam.IsArray {
					itemSchema := Schema{
						Type: mapAcceptedType(bodyParam.Type).Type,
					}
					if !gengo.IsAcceptedType(bodyParam.Type) {
						itemSchema = Schema{
							Ref: makeComponentRef(bodyParam.Type, bodyParam.Package),
						}
					}
					mediaType.Schema = Schema{
						Type:  "array",
						Items: &itemSchema,
					}
				} else if bodyParam.IsMap {
					var additionalProperties Schema
					if !gengo.IsAcceptedType(bodyParam.Type) {
						additionalProperties.Ref = makeComponentRef(bodyParam.Type, bodyParam.Package)
					} else {
						additionalProperties.Type = mapAcceptedType(bodyParam.Type).Type
					}
					mediaType.Schema = Schema{
						Type:                 "object",
						AdditionalProperties: &additionalProperties,
					}
				} else {
					mediaType.Schema = Schema{
						Type: mapAcceptedType(bodyParam.Type).Type,
					}
				}

				operation.RequestBody = &RequestBody{
					Content: map[string]MediaType{
						endpoint.BodyType: mediaType,
					},
				}
			}

			for _, returnType := range endpoint.ReturnTypes {
				s.Log.Debug().Str("path", endpoint.Path).Str("method", endpoint.Method).Str("return", returnType.Field.Name).Msg("Adding return type")
				var mediaType MediaType
				if !gengo.IsAcceptedType(returnType.Field.Type) {
					mediaType.Schema = Schema{
						Ref: makeComponentRef(returnType.Field.Type, returnType.Field.Package),
					}
				} else {
					mediaType.Schema = Schema{
						Type: mapAcceptedType(returnType.Field.Type).Type,
					}
					if returnType.Field.Type == "slice" {
						itemSchema := Schema{
							Type: mapAcceptedType(returnType.Field.SliceType).Type,
						}
						if !gengo.IsAcceptedType(returnType.Field.SliceType) {
							itemSchema = Schema{
								Ref: makeComponentRef(returnType.Field.SliceType, returnType.Field.Package),
							}
						}
						mediaType.Schema.Items = &itemSchema
					} else if returnType.Field.Type == "map" {
						var additionalProperties Schema
						if !gengo.IsAcceptedType(returnType.Field.MapValue) {
							additionalProperties.Ref = makeComponentRef(returnType.Field.MapValue, returnType.Field.Package)
						} else {
							additionalProperties.Type = mapAcceptedType(returnType.Field.MapValue).Type
						}
						mediaType.Schema = Schema{
							Type:                 "object",
							AdditionalProperties: &additionalProperties,
						}
					}
				}

				operation.Responses[strconv.Itoa(returnType.StatusCode)] = Response{
					Description: "",
					Headers:     nil,
					Content: map[string]MediaType{
						endpoint.ContentType: mediaType,
					},
					Links: nil,
				}
			}

			var path Path
			if _, ok := paths[endpoint.Path]; !ok {
				path = Path{}
			} else {
				path = paths[endpoint.Path]
			}
			switch endpoint.Method {
			case "GET":
				path.Get = &operation
			case "POST":
				path.Post = &operation
			case "PUT":
				path.Put = &operation
			case "PATCH":
				path.Patch = &operation
			case "DELETE":
				path.Delete = &operation
			case "HEAD":
				path.Head = &operation
			case "OPTIONS":
				path.Options = &operation
			}

			paths[endpoint.Path] = path
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
					if field.IsEmbedded {
						embeddedProperties = append(embeddedProperties, Schema{
							Ref: makeComponentRef(field.Type, field.Package),
						})
						continue
					}
					if !gengo.IsAcceptedType(field.Type) {
						schema.Properties[key] = Schema{
							Ref: makeComponentRef(field.Type, field.Package),
						}
					} else {
						schema.Properties[key] = Schema{
							Type: mapAcceptedType(field.Type).Type,
						}
					}
				}

				if len(embeddedProperties) > 0 {
					schema.AllOf = append(embeddedProperties, Schema{
						Properties: schema.Properties,
					})

					schema.Properties = nil
				}
			} else if component.Type == "slice" {
				schema = Schema{
					Type: "array",
					Items: &Schema{
						Type: mapAcceptedType(component.SliceType).Type,
					},
				}
				if !gengo.IsAcceptedType(component.SliceType) {
					schema.Items = &Schema{
						Ref: makeComponentRef(component.SliceType, component.Package),
					}
				}
			} else if component.Type == "map" {
				var additionalProperties Schema

				if !gengo.IsAcceptedType(component.MapValue) {
					additionalProperties.Ref = makeComponentRef(component.MapValue, component.Package)
				} else {
					additionalProperties.Type = mapAcceptedType(component.MapValue).Type
				}

				schema = Schema{
					Type:                 "object",
					AdditionalProperties: &additionalProperties,
				}
			} else {
				schema = Schema{
					Type: mapAcceptedType(component.Type).Type,
				}
			}
			components.Schemas[makeComponentRefName(component.Name, component.Package)] = schema
		}
		s.Log.Debug().Msg("Added components")

		if s.Config.Description == "" {
			s.Config.Description = "Generated by gengo"
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

		err = os.WriteFile(filePath, file, 0644)
		if err != nil {
			s.Log.Error().Err(err).Msg("Failed to write OpenAPI schema file")
			return err
		}

		s.Log.Debug().Str("filePath", filePath).Msg("Successfully generated OpenAPI schema file")

		return nil
	}
}
