module github.com/pastelnetwork/gonode/integration

go 1.15

require (
	github.com/DataDog/zstd v1.5.0
	github.com/btcsuite/btcutil v1.0.2
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.18.1
	github.com/pastelnetwork/gonode/common v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/integration/fakes/common v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/pastel v0.0.0-20220301004244-b3d3f466e5d2
	github.com/testcontainers/testcontainers-go v0.12.0
)

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/integration/fakes/common => ./fakes/common
	github.com/pastelnetwork/gonode/pastel => ../pastel
)
