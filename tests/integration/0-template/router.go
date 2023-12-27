package petstore

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/pets", getAllPets)
	r.GET("/pets/:id", getPetByID)
	r.POST("/pets", createPet)
	r.DELETE("/pets/:id", deletePet)

	return r
}
