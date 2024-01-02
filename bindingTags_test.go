package astra

import (
	"github.com/ls6-events/astra/astTraversal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExtractBindingTags(t *testing.T) {
	t.Run("Single Binding Tags", func(t *testing.T) {
		fields := map[string]Field{
			"Name": {
				StructFieldBindingTags: map[astTraversal.BindingTagType]astTraversal.BindingTag{
					astTraversal.JSONBindingTag: {
						Name: "name",
					},
				},
			},
			"Age": {
				StructFieldBindingTags: map[astTraversal.BindingTagType]astTraversal.BindingTag{
					astTraversal.JSONBindingTag: {
						Name: "age",
					},
				},
			},
		}

		bindingTags, uniqueBindings := ExtractBindingTags(fields)
		require.ElementsMatch(t, []astTraversal.BindingTagType{astTraversal.JSONBindingTag}, bindingTags)
		require.False(t, uniqueBindings)
	})

	t.Run("Multiple Binding Tags", func(t *testing.T) {
		fields := map[string]Field{
			"Name": {
				StructFieldBindingTags: map[astTraversal.BindingTagType]astTraversal.BindingTag{
					astTraversal.JSONBindingTag: {
						Name: "name",
					},
					astTraversal.XMLBindingTag: {
						Name: "name",
					},
				},
			},
			"Age": {
				StructFieldBindingTags: map[astTraversal.BindingTagType]astTraversal.BindingTag{
					astTraversal.JSONBindingTag: {
						Name: "age",
					},
					astTraversal.XMLBindingTag: {
						Name: "age",
					},
				},
			},
		}

		bindingTags, uniqueBindings := ExtractBindingTags(fields)
		require.ElementsMatch(t, []astTraversal.BindingTagType{astTraversal.JSONBindingTag, astTraversal.XMLBindingTag}, bindingTags)
		require.False(t, uniqueBindings)
	})

	t.Run("Different Binding Tags", func(t *testing.T) {
		fields := map[string]Field{
			"Name": {
				StructFieldBindingTags: map[astTraversal.BindingTagType]astTraversal.BindingTag{
					astTraversal.JSONBindingTag: {
						Name: "name",
					},
				},
			},
			"Age": {
				StructFieldBindingTags: map[astTraversal.BindingTagType]astTraversal.BindingTag{
					astTraversal.XMLBindingTag: {
						Name: "age",
					},
				},
			},
		}

		bindingTags, uniqueBindings := ExtractBindingTags(fields)
		require.ElementsMatch(t, []astTraversal.BindingTagType{astTraversal.JSONBindingTag, astTraversal.XMLBindingTag}, bindingTags)
		require.False(t, uniqueBindings)
	})

	t.Run("No Binding Tags", func(t *testing.T) {
		fields := map[string]Field{
			"Name": {
				StructFieldBindingTags: map[astTraversal.BindingTagType]astTraversal.BindingTag{},
			},
			"Age": {
				StructFieldBindingTags: map[astTraversal.BindingTagType]astTraversal.BindingTag{},
			},
		}

		bindingTags, uniqueBindings := ExtractBindingTags(fields)
		require.ElementsMatch(t, astTraversal.BindingTags, bindingTags)
		require.False(t, uniqueBindings)
	})

	t.Run("Unique Binding Tag Names", func(t *testing.T) {
		fields := map[string]Field{
			"Name": {
				StructFieldBindingTags: map[astTraversal.BindingTagType]astTraversal.BindingTag{
					astTraversal.JSONBindingTag: {
						Name: "json-name",
					},
					astTraversal.XMLBindingTag: {
						Name: "xml-name",
					},
				},
			},
			"Age": {
				StructFieldBindingTags: map[astTraversal.BindingTagType]astTraversal.BindingTag{
					astTraversal.JSONBindingTag: {
						Name: "json-age",
					},
					astTraversal.XMLBindingTag: {
						Name: "xml-age",
					},
				},
			},
		}

		bindingTags, uniqueBindings := ExtractBindingTags(fields)
		require.ElementsMatch(t, []astTraversal.BindingTagType{astTraversal.JSONBindingTag, astTraversal.XMLBindingTag}, bindingTags)
		require.True(t, uniqueBindings)
	})
}

func TestContentTypeToBindingTag(t *testing.T) {
	t.Run("application/json", func(t *testing.T) {
		require.Equal(t, astTraversal.JSONBindingTag, ContentTypeToBindingTag("application/json"))
	})

	t.Run("application/xml", func(t *testing.T) {
		require.Equal(t, astTraversal.XMLBindingTag, ContentTypeToBindingTag("application/xml"))
	})

	t.Run("application/x-www-form-urlencoded", func(t *testing.T) {
		require.Equal(t, astTraversal.FormBindingTag, ContentTypeToBindingTag("application/x-www-form-urlencoded"))
	})

	t.Run("multipart/form-data", func(t *testing.T) {
		require.Equal(t, astTraversal.FormBindingTag, ContentTypeToBindingTag("multipart/form-data"))
	})

	t.Run("application/yaml", func(t *testing.T) {
		require.Equal(t, astTraversal.YAMLBindingTag, ContentTypeToBindingTag("application/yaml"))
	})
}

func TestBindingTagToContentTypes(t *testing.T) {
	t.Run("JSONBindingTag", func(t *testing.T) {
		require.Equal(t, []string{"application/json"}, BindingTagToContentTypes(astTraversal.JSONBindingTag))
	})

	t.Run("XMLBindingTag", func(t *testing.T) {
		require.Equal(t, []string{"application/xml"}, BindingTagToContentTypes(astTraversal.XMLBindingTag))
	})

	t.Run("FormBindingTag", func(t *testing.T) {
		require.Equal(t, []string{"application/x-www-form-urlencoded", "multipart/form-data"}, BindingTagToContentTypes(astTraversal.FormBindingTag))
	})

	t.Run("YAMLBindingTag", func(t *testing.T) {
		require.Equal(t, []string{"application/yaml"}, BindingTagToContentTypes(astTraversal.YAMLBindingTag))
	})
}
