# WalletNode

`walletnode` is a client application (without a graphical interface), the main purpose if which is to provide a REST API for the GUI on the client side.


## Quick Start

1. Build WalletNode app:
    ``` shell
    cd ./wallletnode
    go build ./
    ```

2. Start WalletNode - there're some options:
    1. Without any parameters, `walletnode` tries to find and read all settings from the config file `walletnode.yml` in the [default dir](#default_dir):

    ``` shell
    ./walletnode
    ```

    2. Specified config files.

    ``` shell
    ./walletnode
        --pastel-config-file ./examples/configs/pastel.conf
        --config-file ./examples/configs/mainnet.yml
    ```

    3. The following set of parameters is suitable for debugging the process.

    ``` shell
    WALLETNODE_DEBUG=1 \
    ./walletnode --log-level debug
    ```


### CLI Options

`walletnode` supports the following CLI parameters:

##### --config-file

Specifies `walletnode` config file in yaml format. By default [default_dir](#default_dir)`/walletnode.yml`


##### --pastel-config-file

Specifies `pastel.conf` config file in env format. By default [default_dir](#default_dir)`/pastel.conf`


##### --temp-dir

Sets the directory for storing temp data. The default path is determined automatically and depends on the OS. On Unix systems, it returns `$TMPDIR` if non-empty, otherwise `/tmp`. On Windows, it returns the first non-empty value from `%TMP%`, `%TEMP%`, `%USERPROFILE%`, or the Windows directory. Use the `--help` flag to see the default path.


##### --log-level

Sets the log level. Can be set to `debug`, `info`, `warning`, `fatal`, e.g. `--log-level debug`. By default `info`.


##### --log-file

Sets the path to log file to write to. Disabled by default.


##### --quiet

Disallows log output to stdout. By default `false`.


##### --swagger

Activate the REST API docs that are available at the URL http://localhost:8080/swagger


### Enviroment variable

Used for integration testing and testing at the development stage.

* `WALLETNODE_DEBUG` displays advanced log messages and stack of the errors. Can be set to `0`, `1`, e.g. `WALLETNODE_DEBUG=1`. By default `0`.


### Default settings

##### default_dir

The path depends on the OS:

* MacOS `~/Library/Application Support/Pastel`
* Linux `~/.pastel`
* Windows (>= Vista) `C:\Users\Username\AppData\Roaming\Pastel`
* Windows (< Vista) `C:\Documents and Settings\Username\Application Data\Pastel`

## Troubleshooting

##### `go: finding module for package .... mocks`

We do not commit mock files, this is only needed for development and their lack, does not affect to build the package. But if you need to run unit tests or `go mod tidy`, you probably get an error. To fix this issue you need to [generate mock files](#generate-mock-files).

## Running unit tests

##### Generate mock files

if you run unit tests first time, you need to generate mock files:

* Install [mockery](https://github.com/vektra/mockery)

``` shell
go get github.com/vektra/mockery/v2/.../
```

* In root directory `./gonode` run generating mock files for all modules:

``` shell
for d in ./*/ ; do (cd "$d" && [[ -f go.mod ]] && go mod tidy && go generate ./...); done
```

##### Run tests

To run unit tests:

``` shell
go test -v -race -timeout 45m -short ./...
```
