package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeletePost(c *gin.Context) {
	postId := c.Param("id")

	if postId == "" {
		handleError(c, http.StatusBadRequest, errors.New("post id is required"))
		return
	}

	c.Status(http.StatusOK)
}
