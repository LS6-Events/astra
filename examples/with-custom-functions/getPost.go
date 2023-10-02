package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"withcustomfunctions/types"
)

func GetPost(c *gin.Context) {
	postId := c.Param("id")

	if postId == "" {
		handleError(c, http.StatusBadRequest, errors.New("post id is required"))
		return
	}

	handleSuccess(c, http.StatusOK, types.Post{
		ID:          2,
		Name:        "Second post",
		Body:        "This is the second post",
		PublishedAt: time.Now(),
		Author: types.Author{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Doe",
		},
	})
}
