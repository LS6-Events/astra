# [WIP] GenGo
## A Go type extractor from web services and outputter to various formats

GenGo is a tool that extracts types from web services and outputs them in various formats. It is intended to be used as a library with your code. The idea stems from not having to define any additional comments or code duplication to export your types to other formats, not violating the DRY (Don't Repeat Yourself) principle.

*The key features of Gengo are:*
* Extract types from web services
* Read types from files
* Follow through different functions to extract types
* Extract types from any packages from the dependency tree
* Adapt to different coding styles
* Read the entirety of the `main` package (by extracting to a temporary directory `$GOPATH/.gengo/gengomain`) to extract types
* Support for different input and output formats

## Supported Formats

### Currently supported input formats
* [Gin](https://www.github.com/gin-gonic/gin)
### Currently supported output formats
* JSON
* [OpenAPI](https://www.openapis.org/)

## Usage
If you have [Go module](https://github.com/golang/go/wiki/Modules) support, then simply add this import to your configuration:
```go
import "github.com/ls6-events/gengo"
```
and run the appropriate `go` commands to pull in the dependency.

Otherwise, install the package using the following command:
```bash
$ go get -u github.com/ls6-events/gengo
```

Now you can set it up with your existing `gin.Engine` instance:
```go
package main

import (
	"github.com/ls6-events/gengo"
	gengoGin "github.com/ls6-events/gengo/inputs/gin"
	"github.com/ls6-events/gengo/outputs/json"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// The handlers can be defined in any package including main
	r.GET("/posts", GetPosts)
	r.GET("/posts/:id", GetPosts)
	r.POST("/posts", CreatePost)
	r.PUT("/posts/:id", UpdatePost)
	r.DELETE("/posts/:id", DeletePost)
	
	r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
        })
    })

	// We need to pass in the gin.Engine instance to the input, this needs to be done before Parse()
	// We also need an output instance, in this case we are using JSON
	// You can have multiple inputs and outputs
	gen := gengo.New(gengoGin.WithGinInput(r), json.WithJSONOutput("example.json"))
	
	// This step isn't necessary to run in production, recommend to only run in development
	err := gen.Parse()
	if err != nil {
		panic(err)
	}
	
	err = r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
```

We also support OpenAPI output, which can be used as follows:
```go
package main

import (
	"github.com/ls6-events/gengo"
	gengoGin "github.com/ls6-events/gengo/inputs/gin"
	"github.com/ls6-events/gengo/outputs/openapi"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/posts", GetPosts)
	r.GET("/posts/:id", GetPosts)
	r.POST("/posts", CreatePost)
	r.PUT("/posts/:id", UpdatePost)
	r.DELETE("/posts/:id", DeletePost)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	
	gen := gengo.New(gengoGin.WithGinInput(r), openapi.WithOpenAPIOutput("openapi.yaml"))

	// For OpenAPI to work, we need to define a configuration, which contains the title, version and description amongst other important information
	config := gengo.Config{
		Title:   "Example API",
		Version: "1.0.0",
		Host:    "localhost",
		Port:    8000,
	}
	
	// Or you can use our config builder, which will set the host to "localhost" by default, and will validate the configuration to test if it is valid.
	config, err := gengo.NewConfigBuilder().SetTitle("Example API").SetVersion("1.0.0").SetPort(8000).SetSecure(false).Build()
	if err != nil {
		panic(err)
	}

	// Either way, you need to set the config for the OpenAPI generator to use
	gen.SetConfig(config)
	
	
	err = gen.Parse()
	if err != nil {
		panic(err)
	}
	
	err = r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
```
### Logging
We use [ZeroLog](https://www.github.com/rs/zerolog) for logging, which is a fast and lightweight logging library. By default we have `info` level logging configured, but to specify `debug`, you can add a configuration option to the `New` function
```go
gengo.New(...(previous configuration)..., gengo.WithCustomLogLevel("debug"))
```

Or if you wanted to set a new logger function:
```go
import "github.com/rs/zerolog/log"

logger := log.With().Caller().Logger()
gengo.New(...(previous configuration)..., gengo.WithCustomLogger(logger))
```

## How it works

GenGo utilises the inbuilt Go AST (Abstract Syntax Tree) to parse the source code of your project. It then extracts the types from the source code and outputs them in the desired format. The documentation for the AST can be found [here](https://golang.org/pkg/go/ast/). The AST is a tree representation of the source code, which can be traversed to extract the types. We can then utilise the [Go tools packages library](https://pkg.go.dev/golang.org/x/tools/go/packages) to extract the types from the source code of the packages wherever they are located on the dependency tree.

We have methods to extract types from the following:
* Functions in your code
* Functions in the dependency tree
* Functions in the `main` package
* Anywhere that utilises the `gin.Context` type
* Functions from inline functions (_but it has to be set inside the function where you specify your routes_)

### Upcoming features
* Allow custom status codes nested inside packages (here we only allow for preset constants (i.e. 200), or the `http.Status*` constants)
* Extract types from other web frameworks as inputs (e.g. [Echo](https://github.com/labstack/echo), [Fiber](https://github.com/gofiber/fiber), etc.)
* Add support for more output formats (e.g. TypeScript interfaces/classes etc.)
* Add unit tests and more documentation
* Test more edge cases (please report any issues you find!)

## Contributing
We aim to welcome contributions to this project, will provide further information in due course.
This project uses [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) for commit messages.

Developed and maintained by the [LS6 Events Team](https://www.ls6.events)