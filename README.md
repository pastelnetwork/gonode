[![Maintained by PastelNetwork](https://img.shields.io/badge/maintained%20by-pastel.network-%235849a6.svg)](https://pastel.network)

# PastelNetwork Go Commons

This repo contains common libraries and helpers used by PastelNetwork for building Apps in Go.

## Packages

This repo contains the following packages:

* [collection](#collection)
* [cli](#cli)
* [configer](#configer)
* [errors](#errors)
* [log](#log)
* [sys](#sys)
* [version](#version)

Each of these packages is described below.

### collection

This package contains useful helper methods for working with collections such as lists, as Go has very few of these built-in.

### cli

Most CLI apps built by PastelNetwork should use this package to run their app, as it takes care of common tasks such as setting the text template for app help topic. Under the hood, the `cli` package is using [urfave/cli](https://github.com/urfave/cli), which is currently our library of choice for CLI apps. It mostly works transparently, but it makes it possible to use another package without any changes in the applications themselves in the future.

Here is the typical usage pattern in `cli/app.go`:

```go
package cli

import (
        "github.com/pastelnetwork/go-commons/cli"
        "github.com/pastelnetwork/go-commons/version"
)

func NewApp() {
    config := config.New()

    app := cli.NewApp("app-name")
    app.SetUsage("App Description")
    app.SetVersion(version.Version())
    app.AddFlags(
        cli.NewFlag("log-level", &config.LogLevel).SetUsage("Set the log `level`.").SetValue(config.LogLevel),
        cli.NewFlag("log-file", &config.LogFile).SetUsage("The log `file` to write to."),
    )
    app.SetActionFunc(func(args []string) error {
        return run(config)
    })

    return app
}
```

### configer

The package contains utilites for load config from the file. Under the hood, the `configer` package is using the [spf13/viper](https://github.com/spf13/viper). Since the spf13/viper supports different configuration formats, you can use any of them: *json*, *toml*, *yaml*, *hcl*, *ini*.

To load a configuration from a file:

``` go
err := configer.ParseFile("./config.yml", &config)
```

If `configFile` you consist of the filename without a path, for example `config.yml`, then it will try to find the `configFile` in the paths defined by default: `.`, `$HOME/.pastel`. To change them, you can use the `SetDefaultConfigPaths` function:

``` go
configer.SetDefaultConfigPaths("/etc/app-config.yml", "/usr/local/etc/app-config.json")
```

### errors

In our apps, we should try to ensure that:

1. Every error has a stacktrace. This makes debugging easier.
1. Every error generated by our own code (as opposed to errors from Go built-in functions or errors from 3rd party libraries) has a custom type. This makes error handling more precise, as we can decide to handle different types of errors differently.

To accomplish these two goals, we have created an `errors` package that has several helper methods, such as `errors.New(err error)`, which wraps the given `error` in an Error object that contains a stacktrace. Under the hood, the `errors` package is using the [go-errors](https://github.com/go-errors/errors) library, but this may change in the future, so the rest of the code should not depend on `go-errors` directly.

Here is how the `errors` package should be used:

1. Any time you want to create your own error, create a custom type for it, and when instantiating that type, wrap it with a call to `errros.New` or `errros.Errorf`. That way, any time you call a method defined in our own code, you know the error it returns already has a stacktrace and you don't have to wrap it yourself.
1. Any time you get back an error object from a function built into Go or a 3rd party library, immediately wrap it with `errros.New` or `errros.Errorf`. This gives us a stacktrace as close to the source as possible.
1. If you need to get back the underlying error, you can use the `errors.Unwrap` function.


Wrap the error object from a function built into Go:

``` go
import "github.com/pastelnetwork/go-commons/errors"

if err := json.Unmarshal(data, &config); err != nil {
    return errors.Errorf("unable to decode into struct, %v", err)
}
```

To check the error and return the corresponding exit code, call the `CheckErrorAndExit` function:

``` go
errors.CheckErrorAndExit(error)

// The above function also displays stacktracke, if the `log.SetDebugMode` function was called before:
log.SetDebugMode(true)
```

This package also contains useful utilities such as taking control of a panicking goroutine.

``` go
defer errors.Recover(errors.CheckErrorAndExit)
```

### log

This package contains utilities for logging from our apps. Instead of using Go's built-in logging library, we are using [logrus](github.com/sirupsen/logrus), as it supports log levels (INFO, WARN, DEBUG, etc), structured logging (making key=value pairs easier to parse), log formatting (including text and JSON), hooks to connect logging to a variety of external systems (e.g. syslog, airbrake, papertrail), and even hooks for automated testing.

```go
import "github.com/pastelnetwork/go-commons/log"

log.Info("Some info")
log.Debugf("Debug message, %s", msg)
```

To change the logging level globally, call the `SetLevelName` function:

```go
log.SetLevelName("info") // debug, warn, error, fatal
```

To print the file name and line number where the log was called, call the `SetDebugMode` function:

```go
log.SetDebugMode(true)
```

### sys

This package contains useful helper methods for working with system environment and signals.


### version

To inject a version of the app at build time, we use Go linker (go tool link) to set the value string variable.

``` shell
go build -o my-app -ldflags "-X github.com/pastelnetwork/go-commons/version.version=v0.0.1"
```

To get the version, call the `Version` function:

``` go
version.Version()
```


## Running tests

```
go test -v ./...
```

## License

This code is released under the MIT License. See [LICENSE.txt](LICENSE.txt).
