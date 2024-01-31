package petstore

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func uploadFile(c *gin.Context) {
	_, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusOK)
}
