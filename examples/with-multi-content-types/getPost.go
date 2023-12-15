package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"withmulticontenttypes/types"
)

func post(postID types.PostID) types.Post {
	return types.Post{
		PostID:      postID,
		Name:        "Selected post",
		Body:        "This is the post you selected",
		PublishedAt: time.Now(),
		Author: types.Author{
			ID:        1,
			FirstName: "Jane",
			LastName:  "Doe",
		},
	}
}

func GetPost(c *gin.Context) {
	var postURI types.PostID
	err := c.ShouldBindUri(postURI)
	if err != nil {
		c.String(http.StatusBadRequest, "cannot bind uri")
		return
	}

	switch c.ContentType() {
	case "application/json":
		c.JSON(http.StatusOK, post(postURI))
	case "application/yaml":
		c.YAML(http.StatusOK, post(postURI))
	default:
		c.String(http.StatusUnsupportedMediaType, "unsupported media type")
	}
}
