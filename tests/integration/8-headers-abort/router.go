package integration

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/headers", getHeader)
	r.POST("/headers", setHeader)
	r.GET("/abort-with-error", abortWithError)
	r.GET("/abort-with-status", abortWithStatus)
	r.GET("/abort-with-status-json", abortWithStatusJSON)

	return r
}
