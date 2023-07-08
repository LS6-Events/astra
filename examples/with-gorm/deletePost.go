package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"withgorm/types"
)

func DeletePost(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	postId := c.Param("id")

	if postId == "" {
		c.String(http.StatusBadRequest, "Missing post id")
	}

	err := db.Delete(&types.Post{}, "id = ?", postId).Error
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not delete post")
		return
	}

	c.Status(http.StatusOK)
}
