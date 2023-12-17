package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"withcustomtypemapping/types"
)

func GetPosts(c *gin.Context) {
	posts := make([]types.Post, 0)

	posts = append(posts, types.Post{
		ID:   1,
		Name: "First post",
		Body: sql.NullString{
			String: "This is the first post",
			Valid:  true,
		},
		PublishedAt: time.Now(),
		Author: types.Author{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
		},
	})

	posts = append(posts, types.Post{
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

	c.JSON(http.StatusOK, posts)
}
