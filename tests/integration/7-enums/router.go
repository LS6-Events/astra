package integration

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/enums/string", getStringEnum)
	r.GET("/enums/string-struct", getStringStructWithEnum)

	r.GET("/enums/int", getIntEnum)
	r.GET("/enums/int-struct", getIntStructWithEnum)

	return r
}
