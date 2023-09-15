# Caching

The caching mechanism behind Astra allows a developer to maintain a cache of the service's state between the steps outlined in the [How it works](https://www.github.com/LS6-Events/docs/how-it-works.md) section. This allows for a developer to maintain the state of the service between steps should the program have any issues. The caching mechanism is not enabled by default, but can be enabled by setting the appropriate option.

## Options

The caching mechanism has two options that can be set:
* `cache.WithCache()` - This option enables the caching mechanism with the default cache file of `.astra/cache.json` in the temporary directory. This also disables cleanup of this directory entirely throughout the process.
* `cache.WithCustomCachePath("cache.json")` - This option enables the caching mechanism with a custom cache file of `cache.json` in the current working directory. It supports both JSON and YAML files. This also disables cleanup of this directory entirely throughout the process.

Both of these options are imported from `github.com/LS6-Events/astra/cache`.

## Plans for the future
We aim to utilise this caching mechanism to maintain state between attempts that this program runs, i.e. only iterate and generate files that have changed and append them to existing outputs, or only run if a function has changed. Any suggestions for both ideas or implementations are welcome.