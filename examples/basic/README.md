# Basic Example
This is a basic example of how to use the Astra API. It utilises some basic features of the API, and some basic types for CRUD operations for a blog system. The results and data are faked, but the important thing to note here is that the types are carried through the entire system, and Astra can be used to generate a full CRUD system with a few lines of code.

## Running the example

To run the example, you need to have a working Go installation. You can run the example by running:

```bash
go run .
```

## Important files

The important files in this example are:
* `main.go` - This is the main file that runs the example.
* `openapi.generated.yaml` - This is the OpenAPI file that is generated by Astra.