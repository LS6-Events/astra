package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
	"withvalidation/types"
)

func GetPosts(c *gin.Context) {
	posts := make([]types.Post, 0)

	posts = append(posts, types.Post{
		ID:          uuid.New().String(),
		Name:        "First post",
		Body:        "This is the first post",
		PublishedAt: time.Now(),
		Author: types.Author{
			ID:        uuid.New().String(),
			FirstName: "John",
			LastName:  "Doe",
		},
	})

	posts = append(posts, types.Post{
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

	c.JSON(http.StatusOK, posts)
}
