module github.com/pastelnetwork/gonode/tools/download_thumbnail

go 1.17

require (
	github.com/gorilla/websocket v1.4.2
	github.com/pastelnetwork/gonode/walletnode v0.0.0-20210723172801-5d493665cdd7
	github.com/pkg/errors v0.9.1
	goa.design/goa/v3 v3.4.3
)

replace (
	github.com/pastelnetwork/gonode/common => ../../common
	github.com/pastelnetwork/gonode/pastel => ../../pastel
	github.com/pastelnetwork/gonode/proto => ../../proto
	github.com/pastelnetwork/gonode/raptorq => ../../raptorq
	github.com/pastelnetwork/gonode/walletnode => ../../walletnode
)
