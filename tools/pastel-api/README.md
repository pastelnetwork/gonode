[![Maintained by PastelNetwork](https://img.shields.io/badge/maintained%20by-pastel.network-%235849a6.svg)](https://pastel.network)

# Pastel API (imitator)

This repo contains `pastel-api` server that imitates the real Pastel API.

## Quick Start

``` shell
go run . --log-level debug
```

## Supported API commands

* `masternode top`
Returns the result from the file `api/services/fake/data/masternode_top.json`.

* `storagefee getnetworkfee`
Returns the result from the file `api/services/fake/data/storagefee_getnetworkfee.json`.

* `tickets list id mine`
Returns the result from the file `api/services/fake/data/tickets_list_id_mine.json`.

* `masternode list-conf`
Filter the result from the file `api/services/fake/data/tickets_list_id_mine.json` with `extPort` using the username from AuthBasic as the value.

* `masternode status`
Filter the result from the file `api/services/fake/data/tickets_list_id_mine.json` with `extPort` using the username from AuthBasic as the value.
