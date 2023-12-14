package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"withmulticontenttypes/types"
)

func createOperation() types.Operation {
	return types.Operation{
		Type:       types.OperationCreate,
		EntityType: "post",
		EntityID:   rand.Int(),
	}
}

func CreatePost(c *gin.Context) {
	switch c.ContentType() {
	case "application/json":
		var postDTO types.PostDTO
		err := c.ShouldBindJSON(&postDTO)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid post")
			return
		}

		c.JSON(http.StatusOK, createOperation())
	case "application/yaml":
		var postDTO types.PostDTO
		err := c.ShouldBindYAML(&postDTO)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid post")
			return
		}

		c.YAML(http.StatusOK, createOperation())
	default:
		c.String(http.StatusUnsupportedMediaType, "unsupported media type")
	}
}
