[![Maintained by PastelNetwork](https://img.shields.io/badge/maintained%20by-pastel.network-%235849a6.svg)](https://pastel.network)

# Pastel API

This repo contains `pastel-api`, a server which simulates the real Pastel API. This is used for integration testing and testing at the development stage.


## Quick Start

##### Run API server itsef

``` shell
PASTEL_API_DEBUG=1 ./pastel-api --log-level debug
```

* `PASTEL_API_DEBUG=1` displays advanced log messages and stack of the errors.


##### Run `supernode`

We have to run each supernode on a different port. The range of ports must be within the range of *extAddress* in the json response [`masternode top`](api/services/fake/data/masternode_top.json). To achieve this, we use different config files, with different port settings.

``` shell
SUPERNODE_DEBUG=1 LOAD_TFMODELS=0 ./supernode --log-level debug -c ./examples/configs/localnet-4444.yml
SUPERNODE_DEBUG=1 LOAD_TFMODELS=0 ./supernode --log-level debug -c ./examples/configs/localnet-4445.yml
SUPERNODE_DEBUG=1 LOAD_TFMODELS=0 ./supernode --log-level debug -c ./examples/configs/localnet-4446.yml
```

* `SUPERNODE_DEBUG=1` displays advanced log messages and stack of the errors.
* `LOAD_TFMODELS=0` prevent loading tensorflow models.


##### Run `walletnod`

``` shell
WALLETNODE_DEBUG=1 ./walletnode --log-level debug -c ./examples/configs/localnet.yml --swagger
```

* `WALLETNODE_DEBUG=1` displays advanced log messages and the stack of the errors.
* `--swagger` enables REST API docs on the http://localhost:8080/swagger


## Supported API requests

* `masternode top`
Returns the result from the file `api/services/fake/data/masternode_top.json`.

* `storagefee getnetworkfee`
Returns the result from the file `api/services/fake/data/storagefee_getnetworkfee.json`.

* `tickets list id mine`
Returns the result from the file `api/services/fake/data/tickets_list_id_mine.json`.

* `masternode list-conf`
Filters the result from the file `api/services/fake/data/tickets_list_id_mine.json` with `extPort` using the username from AuthBasic as the value.

* `masternode status`
Filters the result from the file `api/services/fake/data/tickets_list_id_mine.json` with `extPort` using the username from AuthBasic as the value.

* `pastelid sign "text" "PastelID" "passphrase"`
Returns a hash of `PastelID` and `text` as signature.

* `pastelid verify "text" "signature" "PastelID"`
Returns true if a hash of `PastelID` and `text` matches the `signature`.
