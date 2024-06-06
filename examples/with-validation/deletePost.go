package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeletePost(c *gin.Context) {
	postId := c.Param("id")

	if postId == "" {
		c.String(http.StatusBadRequest, "Missing post id")
	}

	c.Status(http.StatusOK)
}
