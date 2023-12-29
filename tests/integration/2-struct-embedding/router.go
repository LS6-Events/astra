package petstore

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/cats/:id", getCatByID)
	r.GET("/dogs/:id", getDogByID)

	return r
}
