module github.com/pastelnetwork/gonode/integration

go 1.15

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/integration/fakes/common => ./fakes/common
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/raptorq => ../raptorq
)

require (
	github.com/DataDog/zstd v1.5.2
	github.com/btcsuite/btcutil v1.0.2
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.19.0
	github.com/pastelnetwork/gonode/common v0.0.0-20220615180506-00d063d0abf2
	github.com/pastelnetwork/gonode/integration/fakes/common v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/pastel v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/raptorq v0.0.0-20230125124533-72c853cea77b
	github.com/testcontainers/testcontainers-go v0.13.0
)
