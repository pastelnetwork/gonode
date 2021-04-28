module github.com/pastelnetwork/supernode

go 1.16

require (
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/nats-io/nats-server/v2 v2.2.1
	github.com/pastelnetwork/go-commons v0.0.5
	github.com/pastelnetwork/go-pastel v0.0.4
	github.com/pastelnetwork/supernode-proto v0.0.0-20210422111752-b8a67357d93e // indirect
	github.com/ybbus/jsonrpc/v2 v2.1.6 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	google.golang.org/grpc v1.37.0 // indirect
)

replace github.com/pastelnetwork/supernode-proto => ../supernode-proto
