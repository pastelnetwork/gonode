module github.com/pastelnetwork/gonode/supernode

go 1.16

require (
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/p2p v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/pastel v0.0.0
	github.com/pastelnetwork/gonode/proto v0.0.0
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/testify v1.7.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/grpc v1.37.0
)

replace github.com/pastelnetwork/gonode/common => ../common

replace github.com/pastelnetwork/gonode/proto => ../proto

replace github.com/pastelnetwork/gonode/pastel => ../pastel

replace github.com/pastelnetwork/gonode/p2p => ../p2p

replace github.com/pastelnetwork/gonode/p2p/kademlia => ../p2p/kademlia

replace github.com/pastelnetwork/gonode/p2p/kademlia/dao => ../p2p/kademlia/dao

replace github.com/pastelnetwork/gonode/p2p/kademlia/crypto => ../p2p/kademlia/crypto
