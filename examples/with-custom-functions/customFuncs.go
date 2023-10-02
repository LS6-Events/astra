package main

import "github.com/gin-gonic/gin"

func handleError(c *gin.Context, statusCode int, err error) {
	c.Error(err)
	c.String(statusCode, err.Error())
}

func handleSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}
