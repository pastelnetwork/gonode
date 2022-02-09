# SuperNode

`supernode` is a server application that handles `walletnode` requests. The app does most of the work of registering and searching NFTs.


## Quick Start

1. Build SuperNode app:
    1. Install libtensorflow
        - For Linux:
            - For new CPUs which supports avx/avx2 instructions
                ``` shell
                wget https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-2.4.0.tar.gz
                sudo tar -C /usr/local -xzf ./libtensorflow-cpu-linux-x86_64-2.4.0.tar.gz
                sudo /sbin/ldconfig -v
                ```
            - For older CPU, have to build libtensorflow by yourself:
                ``` shell
                # Install bazel
                wget https://github.com/bazelbuild/bazel/releases/download/3.1.0/bazel-3.1.0-installer-linux-x86_64.sh
                chmod u+x bazel-3.1.0-installer-linux-x86_64.sh
                sudo ./bazel-3.1.0-installer-linux-x86_64.sh

                # Build libtensorflow
                wget https://github.com/tensorflow/tensorflow/archive/refs/tags/v2.4.0.tar.gz
                tar xvf ./v2.4.0.tar.gz
                cd tensorflow-2.4.0
                ./configure #set options you want
                bazel build --config opt //tensorflow/tools/lib_package:libtensorflow
                sudo tar -C /usr/local -xzf bazel-bin/tensorflow/tools/lib_package/libtensorflow.tar.gz

                ```

        - For Window/Mac (TBD):

    2. Build:
    ``` shell
    cd ./supernode
    go build ./
    ```  

2. Start supernode - there're some options:
    1. Without any parameters, `supernode` tries to find and read all settings from the config file `supernode.yml` in the [default dir](#default_dir):

    ``` shell
    ./supernode
    ```

    2. Specified config files:

    ``` shell
    ./supernode
        --pastel-config-file ./examples/configs/pastel.conf
        --config-file ./examples/configs/mainnet.yml
    ```

    3. The following set of parameters is suitable for debugging the process:

    ``` shell
    SUPERNODE_DEBUG=1 \
    ./supernode --log-level debug
    ```

3. Do some quick debug on supernode's http debug service - default listen on port 9090
    1. Get stats of p2p
    ``` shell
    $ curl -s  -H "Content-Type: application/json" "http://localhost:9090/p2p/stats"
        # Output : 
        {"config":{"listen_address":"0.0.0.0","port":6000,"data_dir":"/home/nvnguyen/.pastel/supernode/p2p-localnet-6000"},"dht":{"database":{"dir_size":2056071,"log_size":2050592,"record_count":55},"peers":[],"peers_count":0,"self":{"id":"fc+/CbE/r9YgnCgj/5Iet6Gc6MI=","ip":"0.0.0.0","port":6000}},"disk-info":{"all":250438021120,"used":209287962624,"free":41150058496}}
    ```
    2. Store a value to p2p
    ``` shell
    # Store string "Hello" - base64 encode is "SGVsbG8=" to p2p
    $ curl -s -X POST -H "Content-Type: application/json" "http://localhost:9090/p2p" -d '{ "value": "SGVsbG8=" }'
        # Output : 
        {"key":"AU3GnD1W1zCxijnQPgNm895LkqwgQpYK8hYvF2WrPPcA"}
    ```
    3. Retrieve a key from p2p:
    ``` shell
    $curl -s  -H "Content-Type: application/json" "http://localhost:9090/p2p/AU3GnD1W1zCxijnQPgNm895LkqwgQpYK8hYvF2WrPPcA"
        # Output : 
        {"key":"AU3GnD1W1zCxijnQPgNm895LkqwgQpYK8hYvF2WrPPcA","value":"SGVsbG8="}
    ```
    4. Get list of keys:
    ``` shell
    curl -s  -H "Content-Type: application/json" "http://localhost:9090/p2p/get?offset=0&limit=3"
    # Output :
    {
        "keys":["BpY5xbiAu8NTM5voZRWat1is3HJ9zurzE8AFnoTZQbP","9WhodxiYS13tuEQhTEjHL38Y4rJzJHh2BKL4vZk8XTj","BpY5xbiAu8NTM5voZRWat1is3HJ9zurzE8AFnoTZQbP"],
        "len":3,
        "limit":3,
        "offset":0
    }
    ```

### CLI Options

`supernode` supports the following CLI parameters:

##### --config-file

Specifies `supernode` config file in yaml format. By default [default_dir](#default_dir)`/supernode.yml`


##### --pastel-config-file

Specifies `pastel.conf` config file in env format. By default [default_dir](#default_dir)`/pastel.conf`


##### --temp-dir

Sets the directory for storing temp data. The default path is determined automatically and depends on the OS. On Unix systems, it returns `$TMPDIR` if non-empty, else `/tmp`. On Windows, it returns the first non-empty value from `%TMP%`, `%TEMP%`, `%USERPROFILE%`, or the Windows directory. Use the `--help` flag to see the default path.


##### --work-dir

Sets the directory for storing working data, such as tensorflow models. The default path is [default_dir](#default_dir)`/supernode`


##### --log-level

Sets the log level. Can be set to `debug`, `info`, `warning`, `fatal`, e.g. `--log-level debug`. By default `info`.


##### --log-file

Sets the path to log file to write to. Disabled by default.


##### --quiet

Disallows log output to stdout. By default `false`.


### Enviroment variable

Used for integration testing and testing at the developer stage.
* `SUPERNODE_DEBUG` displays advanced log messages and stack of the errors. Can be set to `0`, `1`, e.g. `SUPERNODE_DEBUG=1`. By default `0`.
* `LOAD_TFMODELS` sets the number of tensorflow models to load; it can be used to reduce the number of loaded models; a value of `0` means not to load any models at all, e.g. `LOAD_TFMODELS=0`. By default, it is set to the maximum number of `7`.


### Default settings

##### default_dir

The path depends on the OS:
* MacOS `~/Library/Application Support/Pastel`
* Linux `~/.pastel`
* Windows (>= Vista) `C:\Users\Username\AppData\Roaming\Pastel`
* Windows (< Vista) `C:\Documents and Settings\Username\Application Data\Pastel`

## Troubleshooting

##### `failed to load tensor model`

This means that `supernode` did not find tensforflow models; the path where it tries to load them is [work-dir](#--work-dir)`/tfmodels`. By default [default_dir](#default_dir)`/supernode/tfmodels`

##### `go: finding module for package .... mocks`

We do not commit mock files, this is only needed for development and their lack, does not affect to build the package. But if you need to run unit tests or `go mod tidy`, you probably get an error. To fix this issue you need to [generate mock files](#generate-mock-files).

## Running unit tests

##### Generate mock files

if you run unit tests for the first time, you must generate the mock files:

* Install [mockery](https://github.com/vektra/mockery)

``` shell
go get github.com/vektra/mockery/v2/.../
```

* In the root directory, running `./gonode` will generate mock files for all modules

``` shell
for d in ./*/ ; do (cd "$d" && [[ -f go.mod ]] && go mod tidy && go generate ./...); done
```

##### Run tests

To run unit tests:

``` shell
go test -v -race -timeout 45m -short ./...
```
