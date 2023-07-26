package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/gengo"
	"github.com/ls6-events/gengo/cli"
	"github.com/ls6-events/gengo/inputs"
	"github.com/ls6-events/gengo/outputs"
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

	gen := gengo.New(inputs.WithGinInput(r), outputs.WithOpenAPIOutput("openapi.generated.yaml"), cli.WithCLI())

	config := gengo.Config{
		Title:   "Example API with Cache",
		Version: "1.0.0",
		Host:    "localhost",
		Port:    8000,
	}

	gen.SetConfig(&config)

	err := gen.SetupParse() // NOTE: This is SetupParse instead of Parse
	if err != nil {
		panic(err)
	}

	err = r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
