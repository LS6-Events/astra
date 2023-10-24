package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"withmulticontenttypes/types"
)

func updateOperation(postID types.PostID) types.Operation {
	return types.Operation{
		Type:       types.OperationUpdate,
		EntityType: "post",
		EntityID:   postID.ID,
	}
}

func UpdatePostJSON(c *gin.Context) {
	var postURI types.PostID
	err := c.ShouldBindUri(postURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Error{
			Error:   "cannot bind uri",
			Details: err.Error(),
		})
		return
	}

	var postDTO types.PostDTO
	err = c.ShouldBindJSON(&postDTO)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid post")
		return
	}

	c.JSON(http.StatusOK, updateOperation(postURI))
}

func UpdatePostYAML(c *gin.Context) {
	var postURI types.PostID
	err := c.ShouldBindUri(postURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Error{
			Error:   "cannot bind uri",
			Details: err.Error(),
		})
		return
	}

	var postDTO types.PostDTO
	err = c.ShouldBindYAML(&postDTO)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid post")
		return
	}

	c.JSON(http.StatusOK, updateOperation(postURI))
}
