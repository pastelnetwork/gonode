module github.com/pastelnetwork/supernode

go 1.16

require (
	github.com/nats-io/nats-server/v2 v2.2.1
	github.com/pastelnetwork/go-commons v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/go-pastel v0.0.0-00010101000000-000000000000
)

replace (
	github.com/pastelnetwork/go-commons => ../go-commons
	github.com/pastelnetwork/go-pastel => ../go-pastel
)
