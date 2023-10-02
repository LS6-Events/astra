package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"withcustomfunctions/types"
)

func UpdatePost(c *gin.Context) {
	postId := c.Param("id")

	if postId == "" {
		handleError(c, http.StatusBadRequest, errors.New("post id is required"))
		return
	}

	var postDTO types.PostDTO
	err := c.ShouldBindJSON(&postDTO)
	if err != nil {
		handleError(c, http.StatusBadRequest, err)
		return
	}

	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		handleError(c, http.StatusBadRequest, err)
		return
	}

	post := types.Post{
		ID:   postIdInt,
		Name: postDTO.Name,
		Body: postDTO.Body,
	}

	handleSuccess(c, http.StatusOK, post)
}
