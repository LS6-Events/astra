package integration

import (
	"github.com/gin-gonic/gin"
)

type MultiContentType struct {
	Test string `json:"json-test" xml:"xml-test" yaml:"yaml-test"`
}

func getJSON(c *gin.Context) {
	c.JSON(200, MultiContentType{
		Test: "json",
	})
}

func getXML(c *gin.Context) {
	c.XML(200, MultiContentType{
		Test: "xml",
	})
}

func getYAML(c *gin.Context) {
	c.YAML(200, MultiContentType{
		Test: "yaml",
	})
}
