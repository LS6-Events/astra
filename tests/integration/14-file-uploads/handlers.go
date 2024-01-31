package petstore

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func uploadFile(c *gin.Context) {
	_, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusOK)
}
