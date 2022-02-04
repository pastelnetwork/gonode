module github.com/pastelnetwork/gonode/integration

go 1.15

require (
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.15.0
	github.com/pastelnetwork/gonode/integration/fakes/common v0.0.0-00010101000000-000000000000 // indirect
	github.com/testcontainers/testcontainers-go v0.12.0
)

replace github.com/pastelnetwork/gonode/integration/fakes/common => ./fakes/common
