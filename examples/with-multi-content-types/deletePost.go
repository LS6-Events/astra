package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"withmulticontenttypes/types"
)

func deleteOperation(postID types.PostID) types.Operation {
	return types.Operation{
		Type:       types.OperationDelete,
		EntityType: "post",
		EntityID:   postID.ID,
	}
}

func DeletePost(c *gin.Context) {
	var postURI types.PostID
	err := c.ShouldBindUri(postURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Error{
			Error:   "cannot bind uri",
			Details: err.Error(),
		})
		return
	}

	switch c.ContentType() {
	case "application/json":
		{
			c.JSON(http.StatusOK, deleteOperation(postURI))
		}
	case "application/yaml":
		{
			c.YAML(http.StatusOK, deleteOperation(postURI))
		}
	default:
		{
			c.String(http.StatusUnsupportedMediaType, "unsupported media type")
		}
	}
}
