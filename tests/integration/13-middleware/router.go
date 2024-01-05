package petstore

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()

	pets := r.Group("/pets")

	pets.Use(petsMiddleware)

	pets.GET("", getAllPets)
	pets.GET("/:id", getPetByID)
	pets.POST("", createPet)
	pets.DELETE("/:id", deletePet)

	r.GET("/no-middleware", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "no middleware",
		})
	})

	r.GET("/middleware", headerMiddleware, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "middleware",
		})
	})

	r.GET("/wrapper-func-middleware", handlerFunc(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "inline middleware",
		})
	})

	return r
}
