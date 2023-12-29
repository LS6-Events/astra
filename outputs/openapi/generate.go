package openapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/utils"

	"gopkg.in/yaml.v3"
)

// Generate the OpenAPI output.
// It will marshal the OpenAPI struct and write it to a file.
// It will also generate the paths and their operations.
// It will also generate the components and their schemas.
func Generate(filePath string) astra.ServiceFunction {
	return func(s *astra.Service) error {
		s.Log.Debug().Msg("Generating OpenAPI output")
		if s.Config == nil {
			s.Log.Error().Msg("No config found")
			return astra.ErrConfigNotFound
		}

		s.Log.Debug().Msg("Making collision safe struct names")
		makeCollisionSafeNamesFromComponents(s.Components)

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
					In:       "path",
					Required: pathParam.IsRequired,
					Schema:   mapParamToSchema(pathParam),
				})
			}

			for _, requestHeader := range endpoint.RequestHeaders {
				s.Log.Debug().Str("endpointPath", endpoint.Path).Str("method", endpoint.Method).Str("param", requestHeader.Name).Msg("Adding request header")
				parameter := Parameter{
					Name:     requestHeader.Name,
					In:       "header",
					Required: requestHeader.IsRequired,
					Schema:   mapParamToSchema(requestHeader),
				}

				operation.Parameters = append(operation.Parameters, parameter)
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

			var responseHeaders map[string]Header
			if len(endpoint.ResponseHeaders) > 0 {
				responseHeaders = make(map[string]Header)
				for _, responseHeader := range endpoint.ResponseHeaders {
					s.Log.Debug().Str("endpointPath", endpoint.Path).Str("method", endpoint.Method).Str("param", responseHeader.Name).Msg("Adding response header")
					responseHeaders[responseHeader.Name] = Header{
						Schema:   mapParamToSchema(responseHeader),
						Required: responseHeader.IsRequired,
					}
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
					Headers:     responseHeaders,
					Content:     content,
					Links:       nil,
				}
			}

			if endpoint.Doc != "" {
				operation.Description = endpoint.Doc
			}

			if endpoint.OperationID != "" {
				operation.OperationID = endpoint.OperationID
			}

			var endpointPath Path
			if _, ok := paths[endpoint.Path]; !ok {
				endpointPath = Path{}
			} else {
				endpointPath = paths[endpoint.Path]
			}
			switch endpoint.Method {
			case http.MethodGet:
				endpointPath.Get = &operation
			case http.MethodPost:
				endpointPath.Post = &operation
			case http.MethodPut:
				endpointPath.Put = &operation
			case http.MethodPatch:
				endpointPath.Patch = &operation
			case http.MethodDelete:
				endpointPath.Delete = &operation
			case http.MethodHead:
				endpointPath.Head = &operation
			case http.MethodOptions:
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

			schema := componentToSchema(component)

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
