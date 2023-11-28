module github.com/pastelnetwork/gonode/tools/download_thumbnail

go 1.21

toolchain go1.21.0

require (
	github.com/gorilla/websocket v1.5.1
	github.com/pastelnetwork/gonode/walletnode v0.0.0-20210723172801-5d493665cdd7
	github.com/pkg/errors v0.9.1
	goa.design/goa/v3 v3.14.0
)

require (
	github.com/dimfeld/httptreemux/v5 v5.5.0 // indirect
	github.com/google/uuid v1.4.0 // indirect
)

replace (
	github.com/pastelnetwork/gonode/common => ../../common
	github.com/pastelnetwork/gonode/mixins => ../../mixins
	github.com/pastelnetwork/gonode/pastel => ../../pastel
	github.com/pastelnetwork/gonode/proto => ../../proto
	github.com/pastelnetwork/gonode/raptorq => ../../raptorq
	github.com/pastelnetwork/gonode/walletnode => ../../walletnode
)
