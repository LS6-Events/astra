# Path Blacklisting
By default, GenGo will traverse every path in the OpenAPI file, and generate the types for each path. However, sometimes you may want to blacklist certain paths, for example, if you have a path that is used for internal purposes, and you don't want to generate the types for it.

Therefore we have put together a simple path blacklisting system, which allows you to blacklist certain paths from being traversed by GenGo. This is done by adding a comment to the path in the OpenAPI file, which will be picked up by GenGo and will prevent it from traversing that path.

## Blacklisting a path
To blacklist a path, you need to add an option to the `New` function whereever you setup the service.

```go
gen := gengo.New(...(your options)..., gengo.WithPathBlacklist("path/to/blacklist"))
```

This will blacklist the _exact_ path `path/to/blacklist` from being traversed by GenGo.

There is also a regex version of this function, which allows you to blacklist multiple paths at once.

```go
gen := gengo.New(...(your options)..., gengo.WithPathBlacklistRegex(regex.MustCompile("^path.*$")))
```

This will blacklist any path that starts with `path`.

Finally there is a higher level function that allows you to pass in a function as a parameter to perform your own custom path blacklisting.

```go
gen := gengo.New(...(your options)..., gengo.WithPathBlacklistFunc(func(path string) bool {
    return path == "path/to/blacklist" || path == "path/to/another/blacklist"
}))
```

**Note:** The path is blacklisting, so if you return true to this function or if the other conditions are met, the path will _not_ be included in the generated types.