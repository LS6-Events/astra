package petstore

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/inlineQueryParams", inlineQueryParams)
	r.GET("/boundQueryParams", boundQueryParams)

	return r
}
