package utils

import (
	"github.com/ls6-events/astra"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExtractParamsFromPath(t *testing.T) {
	testCases := []struct {
		name   string
		path   string
		result []astra.Param
	}{
		{
			name:   "simple",
			path:   "/pets",
			result: []astra.Param{},
		},
		{
			name: "with param",
			path: "/pets/:id",
			result: []astra.Param{
				{
					Name: "id",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
			},
		},
		{
			name: "with multiple params",
			path: "/pets/:id/:name",
			result: []astra.Param{
				{
					Name: "id",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
				{
					Name: "name",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
			},
		},
		{
			name: "with multiple params and static paths",
			path: "/pets/:id/:name/owner",
			result: []astra.Param{
				{
					Name: "id",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
				{
					Name: "name",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
			},
		},
		{
			name: "with multiple params and static paths",
			path: "/pets/:id/:name/owner/:ownerID",
			result: []astra.Param{
				{
					Name: "id",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
				{
					Name: "name",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
				{
					Name: "ownerID",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
			},
		},
		{
			name: "with multiple params and static paths",
			path: "/pets/:id/:name/owner/:ownerID/:ownerName",
			result: []astra.Param{
				{
					Name: "id",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
				{
					Name: "name",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
				{
					Name: "ownerID",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
				{
					Name: "ownerName",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
			},
		},
		{
			name: "with multiple params and optional static paths",
			path: "/pets/:id/:name/owner/:ownerID/*ownerName",
			result: []astra.Param{
				{
					Name: "id",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
				{
					Name: "name",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
				{
					Name: "ownerID",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: true,
				},
				{
					Name: "ownerName",
					Field: astra.Field{
						Type: "string",
					},
					IsRequired: false,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ExtractParamsFromPath(tc.path)
			require.ElementsMatch(t, tc.result, result)
		})
	}
}
