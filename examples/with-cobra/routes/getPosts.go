package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"withcobra/types"
)

func GetPosts(c *gin.Context) {
	posts := make([]types.Post, 0)

	posts = append(posts, types.Post{
		ID:          1,
		Name:        "First post",
		Body:        "This is the first post",
		PublishedAt: time.Now(),
		Author: types.Author{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
		},
	})

	posts = append(posts, types.Post{
		ID:          2,
		Name:        "Second post",
		Body:        "This is the second post",
		PublishedAt: time.Now(),
		Author: types.Author{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Doe",
		},
	})

	c.JSON(http.StatusOK, posts)
}
