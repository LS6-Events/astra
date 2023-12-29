package petstore

import (
	"github.com/gin-gonic/gin"
	petstore2 "github.com/ls6-events/astra/tests/petstore"
	"net/http"
	"strconv"
)

// getAllPets returns all pets
func getAllPets(c *gin.Context) {
	allPets := petstore2.Pets

	c.JSON(http.StatusOK, allPets)
}

// getPetByID returns a pet by its ID
// It takes in the ID as a path parameter
func getPetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, err)
		return
	}

	pet, err := petstore2.PetByID(int64(id))
	if err != nil {
		handleError(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, pet)
}

// createPet creates a pet
// It takes in a Pet without an ID in the request body
func createPet(c *gin.Context) {
	var pet petstore2.PetDTO
	err := c.BindJSON(&pet)
	if err != nil {
		handleError(c, http.StatusBadRequest, err)
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

// deletePet deletes a pet by its ID
// It takes in the ID as a path parameter
func deletePet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, err)
		return
	}

	petstore2.RemovePet(int64(id))

	c.Status(http.StatusOK)
}
