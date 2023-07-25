package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"withgorm/types"
)

func UpdatePost(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	postId := c.Param("id")

	if postId == "" {
		c.String(http.StatusBadRequest, "Missing post id")
	}

	var postDTO types.PostDTO
	err := c.ShouldBindJSON(&postDTO)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid post")
		return
	}

	var post types.Post
	err = db.First(&post, "id = ?", postId).Error
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not get post")
		return
	}

	post.Name = postDTO.Name
	post.Body = postDTO.Body
	post.AuthorID = postDTO.AuthorID

	err = db.Save(&post).Error
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not update post")
		return
	}

	c.JSON(http.StatusOK, post)
}
