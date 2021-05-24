module github.com/pastelnetwork/gonode/p2p/kademlia/dao

go 1.16

require (
	github.com/mattn/go-sqlite3 v1.14.7
	github.com/pastelnetwork/gonode/common v0.0.0-20210523194944-207a3e9a74d9
	github.com/pastelnetwork/gonode/p2p/kademlia/crypto v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
)

replace github.com/pastelnetwork/gonode/common => ../../../common

replace github.com/pastelnetwork/gonode/p2p/kademlia/crypto => ../crypto
