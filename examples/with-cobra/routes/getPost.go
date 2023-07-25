package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"withcobra/types"
)

func GetPost(c *gin.Context) {
	postId := c.Param("id")

	if postId == "" {
		c.String(http.StatusBadRequest, "Missing post id")
	}

	c.JSON(http.StatusOK, types.Post{
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
}
