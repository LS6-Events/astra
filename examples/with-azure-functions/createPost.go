package main

import (
	"azurefunctions/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CreatePost(c *gin.Context) {
	var postDTO types.PostDTO
	err := c.ShouldBindJSON(&postDTO)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid post")
		return
	}

	post := types.Post{
		ID:          1,
		Name:        postDTO.Name,
		Body:        postDTO.Body,
		PublishedAt: time.Now(),
		Author: types.Author{
			ID:        postDTO.AuthorID,
			FirstName: "John",
			LastName:  "Doe",
		},
	}

	c.JSON(http.StatusOK, post)
}
