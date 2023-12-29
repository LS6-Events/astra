package petstore

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/astra/tests/petstore"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/inline", func(c *gin.Context) {
		c.JSON(200, petstore.Pet{
			Name:      "inline",
			PhotoURLs: []string{"inline"},
			Status:    "inline",
			Tags:      []petstore.Tag{{Name: "inline"}},
		})
	})

	r.GET("/inline/:param", func(c *gin.Context) {
		param := c.Param("param")

		name, ok := c.GetQuery("name")
		if !ok {
			c.JSON(400, gin.H{"error": "name is required"})
			return
		}

		c.JSON(200, petstore.Pet{
			Name:      name,
			PhotoURLs: []string{"inline"},
			Status:    param,
			Tags:      []petstore.Tag{{Name: "inline"}},
		})
	})

	r.POST("/inline", func(c *gin.Context) {
		var pet petstore.PetDTO
		err := c.BindJSON(&pet)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, petstore.Pet{
			Name:      pet.Name,
			PhotoURLs: pet.PhotoURLs,
			Status:    pet.Status,
			Tags:      pet.Tags,
		})
	})

	return r
}
