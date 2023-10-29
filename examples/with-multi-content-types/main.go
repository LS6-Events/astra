package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/astra"
	"github.com/ls6-events/astra/inputs"
	"github.com/ls6-events/astra/outputs"
)

func main() {
	r := gin.Default()

	r.GET("/posts/json", GetPostsJSON)
	r.GET("/posts/yaml", GetPostsYAML)
	r.GET("/posts/:id/json", GetPostJSON)
	r.GET("/posts/:id/yaml", GetPostYAML)
	r.POST("/posts/json", CreatePostJSON)
	r.POST("/posts/yaml", CreatePostYAML)
	r.PUT("/posts/:id/json", UpdatePostJSON)
	r.PUT("/posts/:id/yaml", UpdatePostYAML)
	r.DELETE("/posts/:id/json", DeletePostJSON)
	r.DELETE("/posts/:id/yaml", DeletePostYAML)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	gen := astra.New(inputs.WithGinInput(r), outputs.WithOpenAPIOutput("openapi.generated.yaml"))

	config := astra.Config{
		Title:   "Example API",
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
