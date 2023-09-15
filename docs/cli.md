# CLI

Astra has also a command line interface (CLI) that can be used to generate code. It revolves around the idea that the setup function in the code is quick and easy to run, creates a file (by default `.astra/cache.json`) that contains the information about the code, and then the CLI can be used to generate code from that file.

## Setup

First, you need to install the CLI:

```bash
go install github.com/ls6-events/astra/cli/astra
```

Next, you must import the CLI package into your code:

```go
import "github.com/ls6-events/astra/cli"
```

And you must include the `cli.WithCLI()` option in your `astra.New` function:

```go
gen := astra.New(...<YOUR OPTIONS>..., cli.WithCLI())
````

Then, you need to utilise the `SetupParse` function instead of `Parse` in your code. It should look similar to the following:

```go
err := gen.SetupParse()
```

Finally, you must run the program once to generate the cache file. You can then use the CLI to generate code from that file. See the [usage](#usage) section for more information.

## Usage

To use the CLI, you need to run the following command:

```bash
astra generate <options>
```

It has the following options:
- `-c` or `--cache`: The path to the generated cache file. Defaults to `.astra/cache.json`.
- `-d` or `--dir`: The working directory where the code was generated (i.e where the `main.go` is located). Defaults to the current working directory.

**Note:** The CLI at this current point in time cannot accept any additional configuration options. If you wish to configure the service differently (i.e. different output file location), you must do so in the code.

## How it works

It leverages the new caching mechanism devised for the Astra library. The `SetupParse` function will only use the router and it's routes to locate the handler function utilising the `reflect.ValueOf` function. It will then create a cache file (by default `.astra/cache.json`) that contains the file name, line number and function name referencing these handlers. The CLI can then be used to generate code from that file.

The file paths are calculated _relative to the current working directory_ so that it can be used with any project structure. The CLI will also use the current working directory to generate the code.

## Example

The following is an example of how to use the CLI. It looks similar to the example you've seen passed around, but with a few key differences - mainly including the CLI.

```go
package main

import (
	"github.com/ls6-events/astra"
	astraGin "github.com/ls6-events/astra/inputs/gin"
	"github.com/ls6-events/astra/outputs/openapi"
	"github.com/gin-gonic/gin"
	"github.com/ls6-events/astra/cli" // Note the new import!
)

func main() {
    r := gin.Default()
    
    r.GET("/hello", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello world!",
        })
    })
    
    // Same as before
    // Note the new WithCLI() option!
    gen := astra.New(astraGin.WithGinInput(r), openapi.WithOpenAPIOutput("openapi.yaml"), cli.WithCLI()) 
    
    // The appropriate configuration for your API
    config := astra.Config{
        Title:   "Example API",
        Version: "1.0.0",
        Host:    "localhost",
        Port:    8000,
    }
    
    gen.SetConfig(config)
    
    // You need to use SetupParse() instead of Parse()
    err := gen.SetupParse()
    if err != nil {
        panic(err)
    }
    
    err = r.Run(":8000")
    if err != nil {
        panic(err)
    }
}
```

## CI/CD

If you wish to use the CLI in a CI/CD pipeline, you must configure the options a bit more first. The reason behind this is the entire `.astra` directory is git ignored by default, and we believe it wouldn't be ideal to have to commit the cache to the repository in this fashion. Therefore we recommend utilising a custom cache file with a custom output like so:
```go
gen := astra.New(...<YOUR OPTIONS>..., cache.WithCustomCachePath("./cache.json"))
```
This will enable the caching mechanism to work, and it will specify the output to which ever file you choose (`.json` by default, but it does support `.yaml`/`.yml`). You can then use the CLI to generate code from that file as such.

```bash
astra generate -c ./cache.json
```

If you include these commands and the installation of the CLI in your CI/CD pipeline (such as GitHub Actions) then you can utilise the CLI to generate code as part of your pipeline.

## Commands/Options

The CLI module has the following options:
- `cli.WithCLI()` - This option will enable the CLI module. It will also enable the caching mechanism. It will not enable or disable any of the steps outlined in the [how it works](./how-it-works.md) section. Instead you must make sure to use the correct functions (i.e. `SetupParse` instead of `Parse`).
- `cli.WithCLIBuilder()` - This option is used by the actual CLI options to enable the rest of the process to be completed providing a cache file is present. It will enable the caching mechanism to maintain state between runs, and set the mode in the service to the correct one.
## Limitations

The CLI is still in it's early stages and has a few limitations. These are:
- It cannot accept any additional configuration options. If you wish to configure the service differently (i.e. different output file location), you must do so in the code.
- It still _must_ utilise the router objects for the inputs (i.e. `gin.Engine` for Gin). This is because the CLI will use the router to locate the handler functions.
- At this current point in time, there is only the `generate` command, any other suggestions please let us know!