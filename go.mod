module github.com/pastelnetwork/supernode

go 1.16

require (
	github.com/nats-io/nats-server/v2 v2.2.1
	github.com/pastelnetwork/go-commons v0.0.1
	github.com/pastelnetwork/go-pastel v0.0.1
	github.com/pastelnetwork/supernode-proto v0.0.0-20210422111752-b8a67357d93e // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	google.golang.org/grpc v1.37.0 // indirect
)

replace github.com/pastelnetwork/go-commons => ../go-commons

replace github.com/pastelnetwork/supernode-proto => ../supernode-proto

replace github.com/pastelnetwork/go-pastel => ../go-pastel
