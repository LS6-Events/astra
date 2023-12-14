package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"withmulticontenttypes/types"
)

var posts = []types.Post{
	{
		PostID:      types.PostID{ID: 1},
		Name:        "First post",
		Body:        "This is the first post",
		PublishedAt: time.Now(),
		Author: types.Author{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
		},
	},
	{
		PostID:      types.PostID{ID: 2},
		Name:        "Second post",
		Body:        "This is the second post",
		PublishedAt: time.Now(),
		Author: types.Author{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Doe",
		},
	},
}

func GetPosts(c *gin.Context) {
	switch c.ContentType() {
	case "application/json":
		c.JSON(http.StatusOK, posts)
	case "application/yaml":
		c.YAML(http.StatusOK, posts)
	default:
		c.String(http.StatusUnsupportedMediaType, "unsupported media type")
	}
}
