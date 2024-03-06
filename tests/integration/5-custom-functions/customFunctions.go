package integration

import "github.com/gin-gonic/gin"

func handleError(c *gin.Context, statusCode int, err error) {
	c.String(statusCode, err.Error())
}
