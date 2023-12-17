package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"withcustomtypemapping/types"
)

func GetPost(c *gin.Context) {
	postId := c.Param("id")

	if postId == "" {
		c.String(http.StatusBadRequest, "Missing post id")
	}

	c.JSON(http.StatusOK, types.Post{
		ID:   2,
		Name: "Second post",
		Body: sql.NullString{
			String: "This is the second post",
			Valid:  true,
		},
		PublishedAt: time.Now(),
		Author: types.Author{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Doe",
		},
	})
}
