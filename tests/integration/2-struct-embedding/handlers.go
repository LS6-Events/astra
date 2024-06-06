package integration

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ls6-events/astra/tests/petstore"
)

func getCatByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	pet, err := petstore.PetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Cat{
		Pet:           *pet,
		Breed:         "Persian",
		IsIndependent: false,
	})
}

func getDogByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	pet, err := petstore.PetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Dog{
		Pet:       *pet,
		Breed:     "Labrador",
		IsTrained: true,
	})
}
