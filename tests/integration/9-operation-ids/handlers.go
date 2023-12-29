package petstore

import (
	"github.com/gin-gonic/gin"
	petstore2 "github.com/ls6-events/astra/tests/petstore"
	"net/http"
	"strconv"
)

func getAllPets(c *gin.Context) {
	allPets := petstore2.Pets

	c.JSON(http.StatusOK, allPets)
}

func getPetByID(c *gin.Context) {
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

	c.JSON(http.StatusOK, pet)
}

func createPet(c *gin.Context) {
	var pet petstore2.PetDTO
	err := c.BindJSON(&pet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	petstore2.AddPet(petstore2.Pet{
		Name:      pet.Name,
		PhotoURLs: pet.PhotoURLs,
		Status:    pet.Status,
		Tags:      pet.Tags,
	})

	c.JSON(http.StatusOK, pet)
}

func deletePet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	petstore2.RemovePet(int64(id))

	c.Status(http.StatusOK)
}
