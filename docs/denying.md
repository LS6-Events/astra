# Path Denying
By default, Astra will traverse every path in the OpenAPI file, and generate the types for each path. However, sometimes you may want to denylist certain paths, for example, if you have a path that is used for internal purposes, and you don't want to generate the types for it.

Therefore we have put together a simple path denying system, which allows you to denylist certain paths from being traversed by Astra. This is done by adding a comment to the path in the OpenAPI file, which will be picked up by Astra and will prevent it from traversing that path.

## Denying a path
To denylist a path, you need to add an option to the `New` function wherever you setup the service.

```go
gen := astra.New(...(your options)..., astra.WithPathDenyList("path/to/denylist"))
```

This will deny the _exact_ path `path/to/denylist` from being traversed by Astra.

There is also a regex version of this function, which allows you to denylist multiple paths at once.

```go
gen := astra.New(...(your options)..., astra.WithPathDenyListRegex(regex.MustCompile("^path.*$")))
```

This will deny any path that starts with `path`.

Finally there is a higher level function that allows you to pass in a function as a parameter to perform your own custom path denying.

```go
gen := astra.New(...(your options)..., astra.WithPathDenyListFunc(func(path string) bool {
    return path == "path/to/denylist" || path == "path/to/another/denylist"
}))
```

**Note:** The path is denying, so if you return true to this function or if the other conditions are met, the path will _not_ be included in the generated types.