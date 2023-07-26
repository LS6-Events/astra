package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/gengo"
	"github.com/ls6-events/gengo/inputs"
	"github.com/ls6-events/gengo/outputs"
	"regexp"
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

	r.Static("/swaggerui", "./swaggerui")

	gen := gengo.New(inputs.WithGinInput(r), outputs.WithOpenAPIOutput("./swaggerui/swagger.json"), gengo.WithPathDenyListRegex(regexp.MustCompile("^/swaggerui.*")))

	config := gengo.Config{
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
