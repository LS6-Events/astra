package routes

import "github.com/gin-gonic/gin"

func GetRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/posts", GetPosts)
	r.GET("/posts/:id", GetPost)
	r.POST("/posts", CreatePost)
	r.PUT("/posts/:id", UpdatePost)
	r.DELETE("/posts/:id", DeletePost)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return r
}
