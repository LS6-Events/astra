package petstore

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/top-level", topLevelHandler)
	r.GET("/nested", nestedHandler)

	return r
}
