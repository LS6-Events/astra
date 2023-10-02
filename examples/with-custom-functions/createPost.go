package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"withcustomfunctions/types"
)

func CreatePost(c *gin.Context) {
	var postDTO types.PostDTO
	err := c.ShouldBindJSON(&postDTO)
	if err != nil {
		handleError(c, http.StatusBadRequest, err)
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

	handleSuccess(c, http.StatusCreated, post)
}
