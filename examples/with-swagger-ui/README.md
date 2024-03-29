# With Swagger UI Example
This is an example for how to use Swagger UI with Astra. Here we can run both the web server and the Astra parser separately to generate the types, while they both have access to the router object. Swagger UI can be used to view the OpenAPI file that is generated by Astra, and investigate the endpoints and their responses.

It is important to note that the Swagger UI code is included in this example as a separate distributed web application, and is not included in the Astra API. This is because the Swagger UI code is not required for the API to function, and is only used for documentation purposes. The path is blacklisted for Astra to prevent it from being included in the OpenAPI file.

## Running the example

To run the example, you need to have a working Go installation. You can then run the web server by running:

```bash
go run .
```

Then open the Swagger UI page at [http://localhost:8080/swaggerui/](http://localhost:8080/swaggerui/), and you should see your types there.

If you have any suggestions for better alternatives to hosting Swagger UI, please open an issue.

## Important files

The important files in this example are:
* `main.go` - This is the main file that runs the example.
* `swaggerui/` - This is the Swagger UI web application.
* `openapi.generated.yaml` - This is the OpenAPI file that is generated by Astra.