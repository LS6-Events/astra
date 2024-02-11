package petstore

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func inlineQueryParams(c *gin.Context) {
	_ = c.Query("name")
	_ = c.Query("tag")
	_ = c.Query("limit")
	_ = c.Query("decimal")

	c.Status(http.StatusOK)
}

type BoundQueryParams struct {
	Name    string  `form:"name"`
	Tag     string  `form:"tag"`
	Limit   int     `form:"limit"`
	Decimal float64 `form:"decimal"`

	Slice []string `form:"slice"`

	Map map[string]string `form:"map"`
}

func boundQueryParams(c *gin.Context) {
	var params BoundQueryParams
	_ = c.BindQuery(&params)

	c.Status(http.StatusOK)
}
