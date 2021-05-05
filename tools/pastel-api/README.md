[![Maintained by PastelNetwork](https://img.shields.io/badge/maintained%20by-pastel.network-%235849a6.svg)](https://pastel.network)

# Pastel API (imitator)

This repo contains `pastel-api` server that imitates the real Pastel API.

## Quick Start

``` shell
go run . --log-level debug
```

## Supported API commands

* `masternode top`
Returns a prepared result from a file `api/services/static/data/masternode_top.json`.

* `storagefee getnetworkfee`
Returns a prepared result from a file `api/services/static/data/storagefee_getnetworkfee.json`.
