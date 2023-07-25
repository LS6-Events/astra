# Logging

In GenGo we use [ZeroLog](https://www.github.com/rs/zerolog) for logging, which is a fast and lightweight logging library. By default we have `info` level logging configured, but to specify `debug`, you can add a configuration option to the `New` function

```go
gen := gengo.New(...(your options)..., gengo.WithLogLevel("debug"))
```

## Custom Logger

If you want to use your own logger, you can do so by using the `WithCustomLogger` option:

```go
import "github.com/rs/zerolog/log"

logger := log.With().Caller().Logger()
gengo.New(...(previous configuration)..., gengo.WithCustomLogger(logger))
```