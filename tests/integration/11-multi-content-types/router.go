package petstore

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/json", getJSON)
	r.GET("/xml", getXML)
	r.GET("/yaml", getYAML)

	return r
}
