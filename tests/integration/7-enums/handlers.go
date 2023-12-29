package petstore

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getStringEnum(c *gin.Context) {
	c.JSON(http.StatusOK, TestStringEnumAvailable)
}

func getStringStructWithEnum(c *gin.Context) {
	c.JSON(http.StatusOK, TestStructWithStringEnum{
		Enum: TestStringEnumSold,
	})
}

func getIntEnum(c *gin.Context) {
	c.JSON(http.StatusOK, TestIntEnumAvailable)
}

func getIntStructWithEnum(c *gin.Context) {
	c.JSON(http.StatusOK, TestStructWithIntEnum{
		Enum: TestIntEnumSold,
	})
}
