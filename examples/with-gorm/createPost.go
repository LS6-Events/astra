package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
	"withgorm/types"
)

func CreatePost(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var postDTO types.PostDTO
	err := c.ShouldBindJSON(&postDTO)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid post")
		return
	}

	post := types.Post{
		Name:        postDTO.Name,
		Body:        postDTO.Body,
		PublishedAt: time.Now(),
		AuthorID:    postDTO.AuthorID,
	}

	err = db.Create(&post).Error
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not create post")
		return
	}

	c.JSON(http.StatusOK, post)
}
