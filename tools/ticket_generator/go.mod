module github.com/pastelnetwork/gonode/tools/ticket_generator

go 1.16

require (
	github.com/go-errors/errors v1.4.0
	github.com/gorilla/websocket v1.4.2
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/pastel v0.0.0
	github.com/pastelnetwork/gonode/walletnode v0.0.0-20210723172801-5d493665cdd7
	github.com/stretchr/testify v1.7.0 // indirect
	goa.design/goa/v3 v3.4.3
)

replace (
	github.com/pastelnetwork/gonode/common => ../../common
	github.com/pastelnetwork/gonode/pastel => ../../pastel
	github.com/pastelnetwork/gonode/proto => ../../proto
	github.com/pastelnetwork/gonode/raptorq => ../../raptorq
	github.com/pastelnetwork/gonode/walletnode => ../../walletnode
)
