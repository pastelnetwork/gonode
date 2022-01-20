# Integration Tests
This package contains integration tests currently covering the following flows
    - P2P (kademila) Store & Retrieve
    - Register NFT through Walletnode API
    
To run the tests
```
INTEGRATION_TEST_ENV=true go test -v --timeout=10m 
```
## Debugging Tests
As we use containerized approach in which it could become harder to debug, we can simply do `cd infra && docker-compose up --build` and use `main.go` to register mocks and use curl (or any tool of choice) to call APIs that need debugging.
`main.go` is only supposed to be used for debugging purposes.
to debug a test,
## Structure
Tests are structured as follows

### Fakes
`/fakes` contains fake RaptorQ, Duplication detection & Pastel servers. These servers can be mocked through `/register` POST endpoint. 
This endpoint accepts JSON payload that is supposed to be the response and 3 query parameters

    - method: method name for example, `masternode` in case of pasteld
    - params: comma separated list of params, for example, "list,extra"
    - count: number of times this particular response is to be returned for these method & params. 
    
### Infra
`/infra` contains docker files for the fake servers as well as `docker-compose.yml` to use those dockerfiles to spin up all the containers.

### Mock
`/mock` contains mock helper as well as expected responses that the servers are supposed to return.

## Contribution
We use [testcontainers](https://github.com/testcontainers/testcontainers-go) to spin up containers & [Ginkgo](https://onsi.github.io/ginkgo/) & [GoMega](https://onsi.github.io/gomega/) for writing tests. Please go through the documentation to follow the coding style.

