package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"withvalidation/types"
)

func UpdatePost(c *gin.Context) {
	postId := c.Param("id")

	if postId == "" {
		c.String(http.StatusBadRequest, "Missing post id")
	}

	var postDTO types.PostDTO
	err := c.ShouldBindJSON(&postDTO)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid post")
		return
	}

	post := types.Post{
		ID:   postId,
		Name: postDTO.Name,
		Body: postDTO.Body,
	}

	c.JSON(http.StatusOK, post)
}
