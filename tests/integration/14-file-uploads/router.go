package petstore

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/upload", uploadFile)

	return r
}
