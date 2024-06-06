package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
	"withvalidation/types"
)

func GetPost(c *gin.Context) {
	postId := c.Param("id")

	if postId == "" {
		c.String(http.StatusBadRequest, "Missing post id")
	}

	c.JSON(http.StatusOK, types.Post{
		ID:          uuid.New().String(),
		Name:        "Second post",
		Body:        "This is the second post",
		PublishedAt: time.Now(),
		Author: types.Author{
			ID:        uuid.New().String(),
			FirstName: "Jane",
			LastName:  "Doe",
		},
	})
}
