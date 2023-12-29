package petstore

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getHeader(c *gin.Context) {
	header := c.GetHeader("X-Test-Header")
	c.String(http.StatusOK, header)
}

func setHeader(c *gin.Context) {
	c.Header("X-Test-Header", "test")
	c.Status(http.StatusOK)
}

func abortWithError(c *gin.Context) {
	c.AbortWithError(http.StatusBadRequest, nil)
}

func abortWithStatus(c *gin.Context) {
	c.AbortWithStatus(http.StatusUnauthorized)
}

func abortWithStatusJSON(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{
		"message": "unauthorized",
	})
}
