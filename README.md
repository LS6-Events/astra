# GenGo
## An automatic Go type extractor for web services with very little additional configuration

GenGo is a tool that extracts types from web services _without any additional configuration for the handlers_ and outputs them in various formats. It is intended to be used as a library with your code. 

The idea stems from not having to define any additional comments or code duplication to export your types to be used in other systems, not violating the DRY (Don't Repeat Yourself) principle. 

For example, if you had a Go backend that used the [Gin](https://www.github.com/gin-gonic/gin) framework, you would have to define the types in the handler function, and then define them again in the OpenAPI specification. This is not only time consuming, but also violates the DRY principle. Whereas our solution, you'd only have to define the types once, and then use GenGo to extract them from the handlers and output them to the OpenAPI specification, which you could host using [Swagger UI](https://swagger.io/tools/swagger-ui/) or use them to generate client code.

## Features

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
* [OpenAPI](https://www.openapis.org/)
* [JSON](https://www.json.org/json-en.html) (for debugging purposes)

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
We recommend using the OpenAPI specification as the output format, as it is the most widely used format for describing RESTful APIs. To use it, you need to import the following package:
```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/gengo"
	gengoGin "github.com/ls6-events/gengo/inputs/gin"
	"github.com/ls6-events/gengo/outputs/openapi"
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
### CLI (CI/CD)
GenGo also has a method for running the program from the command line, which is useful for CI/CD pipelines. To use it, follow the instructions in the [CLI documenation](./docs/cli.md).

### Blacklist

GenGo can also blacklist certain functions from being parsed. This is useful if you have a function that you don't want to be parsed, such as a function that is used for testing. To use this, follow the instructions in the [blacklist documentation](./docs/blacklisting.md).


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

More information can be found in the [logging documentation](./docs/logging.md).

## How it works

GenGo utilises the inbuilt Go AST (Abstract Syntax Tree) to parse the source code of your project. It then extracts the types from the source code and outputs them in the desired format. The documentation for the AST can be found [here](https://golang.org/pkg/go/ast/). The AST is a tree representation of the source code, which can be traversed to extract the types. We can then utilise the [Go tools packages library](https://pkg.go.dev/golang.org/x/tools/go/packages) to extract the types from the source code of the packages wherever they are located on the dependency tree.

We have methods to extract types from the following:
* Functions in your code
* Functions in the dependency tree
* Functions in the `main` package
* Anywhere that utilises the `gin.Context` type
* Functions from inline functions (_but it has to be set inside the function where you specify your routes_)

There is more information in the [how it works documentation](./docs/how-it-works.md)

### Upcoming features
* Allow custom status codes nested inside packages (here we only allow for preset constants (i.e. 200), or the `http.Status*` constants)
* Extract types from other web frameworks as inputs (e.g. [Echo](https://github.com/labstack/echo), [Fiber](https://github.com/gofiber/fiber), etc.)
* Add support for more output formats (e.g. TypeScript interfaces/classes etc.)
* Add unit tests and more documentation
* Test more edge cases (please report any issues you find!)

## Issue Reporting

Naturally a project that has to adapt to so many people's coding styles will have some issues. If you find any, please report them in the [issues](https://www.github.com/LS6-Events/gengo/issues) section of this repository. 

Please try to be as specific as possible, 

We will try to fix them as soon as possible.

## Contributing
We aim to welcome contributions to this project, will provide further information in due course.
This project uses [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) for commit messages.

Developed and maintained by the [LS6 Events Team](https://www.ls6.events)