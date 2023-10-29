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

func CreatePostJSON(c *gin.Context) {
	var postDTO types.PostDTO
	err := c.ShouldBindJSON(&postDTO)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid post")
		return
	}

	c.JSON(http.StatusOK, createOperation())
}

func CreatePostYAML(c *gin.Context) {
	var postDTO types.PostDTO
	err := c.ShouldBindYAML(&postDTO)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid post")
		return
	}

	c.YAML(http.StatusOK, createOperation())
}
