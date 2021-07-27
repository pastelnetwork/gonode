module github.com/pastelnetwork/gonode/supernode

go 1.16

require (
	github.com/DataDog/zstd v1.4.8
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/metadb v0.0.0
	github.com/pastelnetwork/gonode/p2p v0.0.0
	github.com/pastelnetwork/gonode/pastel v0.0.0
	github.com/pastelnetwork/gonode/probe v0.0.0
	github.com/pastelnetwork/gonode/proto v0.0.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	google.golang.org/grpc v1.37.0
)

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/metadb => ../metadb
	github.com/pastelnetwork/gonode/p2p => ../p2p
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/probe => ../probe
	github.com/pastelnetwork/gonode/proto => ../proto
)
