module github.com/pastelnetwork/gonode/tools/downloadnft

go 1.21

toolchain go1.21.0

replace (
	github.com/pastelnetwork/gonode/common => ../../common
	github.com/pastelnetwork/gonode/mixins => ../../mixins
	github.com/pastelnetwork/gonode/p2p => ../../p2p
	github.com/pastelnetwork/gonode/pastel => ../../pastel
	github.com/pastelnetwork/gonode/proto => ../../proto
	github.com/pastelnetwork/gonode/raptorq => ../../raptorq
	github.com/pastelnetwork/gonode/walletnode => ../../walletnode
)

require (
	github.com/gorilla/websocket v1.5.0
	github.com/pastelnetwork/gonode/walletnode v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	goa.design/goa/v3 v3.13.0
)
