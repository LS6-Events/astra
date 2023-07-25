package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"withgorm/types"
)

func GetPost(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	postId := c.Param("id")

	if postId == "" {
		c.String(http.StatusBadRequest, "Missing post id")
	}

	var post types.Post
	err := db.First(&post, "id = ?", postId).Preload("Author").Preload("Comments").Error
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not get post")
		return
	}

	c.JSON(http.StatusOK, post)
}
