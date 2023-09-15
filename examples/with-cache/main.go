package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/cache"
	"github.com/ls6-events/astra/inputs"
	"github.com/ls6-events/astra/outputs"
)

func main() {
	r := gin.Default()

	r.GET("/posts", GetPosts)
	r.GET("/posts/:id", GetPost)
	r.POST("/posts", CreatePost)
	r.PUT("/posts/:id", UpdatePost)
	r.DELETE("/posts/:id", DeletePost)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Normally, you would use the following option to add the cache to the generated code:
	// gen := astra.New(astraGin.WithGinInput(r), openapi.WithOpenAPIOutput("openapi.generated.yaml"), cache.WithCache())
	// However, the folder by default is git ignored, so we will use the following option instead:
	gen := astra.New(inputs.WithGinInput(r), outputs.WithOpenAPIOutput("openapi.generated.yaml"), cache.WithCustomCachePath("cache.json"))

	config := astra.Config{
		Title:   "Example API with Cache",
		Version: "1.0.0",
		Host:    "localhost",
		Port:    8000,
	}

	gen.SetConfig(&config)

	err := gen.Parse()
	if err != nil {
		panic(err)
	}

	err = r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
