package petstore

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/astra/tests/integration/10-struct-name-collision-avoidance/nested/types"
	topLevelTypes "github.com/ls6-events/astra/tests/integration/10-struct-name-collision-avoidance/types"
	"net/http"
)

func topLevelHandler(c *gin.Context) {
	c.JSON(http.StatusOK, topLevelTypes.TestType{
		TopLevelField: "topLevel",
	})
}

func nestedHandler(c *gin.Context) {
	c.JSON(http.StatusOK, types.TestType{
		NestedField: "nested",
	})
}
