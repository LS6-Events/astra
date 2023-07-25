package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"withgorm/types"
)

func GetPosts(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var posts []types.Post
	err := db.Preload("Author").Find(&posts).Error
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not get posts")
		return
	}

	c.JSON(http.StatusOK, posts)
}
