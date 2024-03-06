package integration

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, TestStructFormatter{})
	})

	return r
}
