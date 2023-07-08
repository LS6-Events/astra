package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"withcli/types"
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

	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid post id")
		return
	}

	post := types.Post{
		ID:   postIdInt,
		Name: postDTO.Name,
		Body: postDTO.Body,
	}

	c.JSON(http.StatusOK, post)
}
