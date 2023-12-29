# OpenAPI Schema Formats
This tests the OpenAPI schema formats feature of Astra. It is a simple example of how to use OpenAPI schema formats to create a new format that can be used. This tests:
- Creating a custom `date-time` format for the `time.Time` type, and `uuid.UUID` type from google for `uuid` format.
- `int32` is the default format for `int` types
- `int64` is respected for `int64` types
- `float32` is the default format for `float32` types
- `float64` is respected for `float64` types
- Accounted for all other data types in `formatTypes.go`
