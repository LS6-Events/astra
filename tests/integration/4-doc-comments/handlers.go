package integration

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ls6-events/astra/tests/petstore"
)

// getAllPets returns all pets.
func getAllPets(c *gin.Context) {
	allPets := petstore.Pets

	c.JSON(http.StatusOK, allPets)
}

// getPetByID returns a pet by its ID.
// It takes in the ID as a path parameter.
func getPetByID(c *gin.Context) {
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

	c.JSON(http.StatusOK, pet)
}

// createPet creates a pet.
// It takes in a Pet without an ID in the request body.
func createPet(c *gin.Context) {
	var pet petstore.PetDTO
	err := c.BindJSON(&pet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	petstore.AddPet(petstore.Pet{
		Name:      pet.Name,
		PhotoURLs: pet.PhotoURLs,
		Status:    pet.Status,
		Tags:      pet.Tags,
	})

	c.JSON(http.StatusOK, pet)
}

// deletePet deletes a pet by its ID.
// It takes in the ID as a path parameter.
func deletePet(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	petstore.RemovePet(id)

	c.Status(http.StatusOK)
}

type NoDescription struct {
	Foo string `json:"foo"`
}

func noDescriptionRoute(c *gin.Context) {
	c.JSON(http.StatusOK, NoDescription{Foo: "bar"})
}
