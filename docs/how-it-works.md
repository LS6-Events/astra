# How it works

GenGo utilises the inbuilt Go AST (Abstract Syntax Tree) to parse the source code of your project. It then extracts the types from the source code and outputs them in the desired format. The documentation for the AST can be found [here](https://golang.org/pkg/go/ast/). The AST is a tree representation of the source code, which can be traversed to extract the types. We can then utilise the [Go tools packages library](https://pkg.go.dev/golang.org/x/tools/go/packages) to extract the types from the source code of the packages wherever they are located on the dependency tree.

We have methods to extract types from the following:
* Functions in your code
* Functions in the dependency tree
* Functions in the `main` package
* Anywhere that utilises the `gin.Context` type
* Functions from inline functions (_but it has to be set inside the function where you specify your routes_)

## The process

The process of extracting the types from the source code is as follows, which we have categorised into 7 steps:
1. Setup
2. Create Routes
3. Parse Routes
4. Process
5. Clean
6. Generate
7. Teardown

Each of these steps have their own function referenced in `parse.go` which the `Parse` function will run all of them. It will run `SetupParse` and `CompleteParse`. `SetupParse` includes the 'Setup' and 'Create Routes' steps, whereas `CompleteParse` contains 'Parse Routes', 'Process', 'Clean', 'Generate' and 'Teardown'.

The reason why these processes were separated into `SetupParse` and `CompleteParse` were because they were to be used if the developer chose to utilise the CLI elements of this process - the `SetupParse` would be run after the creation of the router as it would be the first step in the process, and `CompleteParse` would be run where needed utilising the CLI. Documentation for the CLI can be found [here](./cli.md).

If caching is enabled, after each of these steps the caching process is run to maintain state between steps should the program have any issues. You can read more about caching [here](./caching.md).

### Setup

The setup process creates the temporary directory `.gengo` at the current working directory. This temporary directory is used as a temporary storage for a copy of the main directory to be used in package traversal and also holds the cache files. The temporary directory is deleted at the end of the process.

### Create Routes

Create routes utilises the router objects specified by the inputs to access the handler function references for each of the endpoints, and utilising the `reflect.ValueOf` (I know, but we haven't found any issues so far) we can locate the file, line number and function name of these handlers. We then store this information inside the service to be used at a later step. This process is very quick and is used if the CLI process is required and no parsing inferring is necessary. The supported inputs for this step are:
- Gin

### Parse Routes

Parse routes uses the previously stored router directional objects to setup the AST traversal. It will traverse the AST of the files specified by the router objects and extract the types from the source code. It will then store these types in the service to be used at a later step. It has the capabilities to extract types from the following:
- Your handler function
- Variables that acquire their values from separate functions same file/package
- Functions in the `main` package (we copy the `main` package to the temporary directory to allow for this, as the `main` keyword is reserved)
- Anywhere that utilises the `gin.Context` type
- 'Inline' handlers (handlers that are set inside the function where you specify your routes)
- Any functions that either return the type used by the sending functions or any function that utilises the context imported from any other package
- Status codes in the constant format (e.g. `200` or `http.StatusOK`)

### Process

The process step is where any of the types extracted from the previous step are processed. This includes:
- Removing any duplicate types
- Searching packages for these types (utilising the `golang.org/x/tools/go/packages` library) [This does include the `main` package if required as flagged in the previous step]
- Finding new types recursively from the packages found
- Processing structs, maps, slices and references to find their types
- Accounting for embedded structs

### Clean

The clean step is where the types are cleaned up. This includes:
- Converting the package name of the temporary `main` package created back to `main`
- Changing `time.Time` and any other special types to their string representation for JSON encoding

### Generate

The generate step is where the types are outputted to the desired format, utilising the previously specified output format. This includes:
- OpenAPI 3.0
- JSON (purely for debugging purposes)

### Teardown

This step is where the temporary directory is deleted (unless caching is specified), and the process is complete.