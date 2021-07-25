module github.com/pastelnetwork/gonode/tools/ui-cli

go 1.16

require (
	github.com/gorilla/websocket v1.4.2
	github.com/pastelnetwork/gonode/walletnode v0.0.0-20210723172801-5d493665cdd7 // indirect
	github.com/pkg/errors v0.9.1
	goa.design/goa/v3 v3.4.3
)

replace (
	github.com/pastelnetwork/gonode/common => /home/nvnguyen/ccwork/pastel/gonode/common
	github.com/pastelnetwork/gonode/pastel => /home/nvnguyen/ccwork/pastel/gonode/pastel
	github.com/pastelnetwork/gonode/proto => /home/nvnguyen/ccwork/pastel/gonode/proto
	github.com/pastelnetwork/gonode/raptorq => /home/nvnguyen/ccwork/pastel/gonode/raptorq
	github.com/pastelnetwork/gonode/walletnode => /home/nvnguyen/ccwork/pastel/gonode/walletnode
)
