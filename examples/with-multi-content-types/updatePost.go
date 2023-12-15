package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"withmulticontenttypes/types"
)

func updateOperation(postID types.PostID) types.Operation {
	return types.Operation{
		Type:       types.OperationUpdate,
		EntityType: "post",
		EntityID:   postID.ID,
	}
}

func UpdatePost(c *gin.Context) {
	var postURI types.PostID
	err := c.ShouldBindUri(postURI)
	if err != nil {
		c.String(http.StatusBadRequest, "cannot bind uri")
		return
	}

	switch c.ContentType() {
	case "application/json":
		var postDTO types.PostDTO
		err = c.ShouldBindJSON(&postDTO)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid post")
			return
		}

		c.JSON(http.StatusOK, updateOperation(postURI))
	case "application/yaml":
		var postDTO types.PostDTO
		err = c.ShouldBindYAML(&postDTO)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid post")
			return
		}

		c.YAML(http.StatusOK, updateOperation(postURI))
	default:
		c.String(http.StatusUnsupportedMediaType, "unsupported media type")
	}
}
