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

func GetPostJSON(c *gin.Context) {
	var postURI types.PostID
	err := c.ShouldBindUri(postURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Error{
			Error:   "cannot bind uri",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, post(postURI))
}

func GetPostYAML(c *gin.Context) {
	var postURI types.PostID
	err := c.ShouldBindUri(postURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Error{
			Error:   "cannot bind uri",
			Details: err.Error(),
		})
		return
	}

	c.YAML(http.StatusOK, post(postURI))
}
