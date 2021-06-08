module github.com/pastelnetwork/gonode/tools/pastel-utility

go 1.16

require (
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/p2p v0.0.0
	github.com/pastelnetwork/gonode/pastel v0.0.0
	github.com/pastelnetwork/gonode/probe v0.0.0
	github.com/pastelnetwork/gonode/proto v0.0.0
	github.com/pastelnetwork/gonode/supernode v0.0.0
	github.com/pastelnetwork/gonode/walletnode v0.0.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.37.0
)

replace (
	github.com/pastelnetwork/gonode/common => ../../common
	github.com/pastelnetwork/gonode/p2p => ../../p2p
	github.com/pastelnetwork/gonode/pastel => ../../pastel
	github.com/pastelnetwork/gonode/probe => ../../probe
	github.com/pastelnetwork/gonode/proto => ../../proto
	github.com/pastelnetwork/gonode/supernode => ../../supernode
	github.com/pastelnetwork/gonode/walletnode => ../../walletnode
)
