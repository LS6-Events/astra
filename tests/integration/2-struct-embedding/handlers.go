package petstore

import (
	"github.com/gin-gonic/gin"
	petstore2 "github.com/ls6-events/astra/tests/petstore"
	"net/http"
	"strconv"
)

func getCatByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pet, err := petstore2.PetByID(int64(id))
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pet, err := petstore2.PetByID(int64(id))
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
