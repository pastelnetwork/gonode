module github.com/pastelnetwork/gonode/p2p

go 1.16

replace github.com/pastelnetwork/gonode/p2p/kademlia => ./kademlia

replace github.com/pastelnetwork/gonode/common => ../common

require (
	github.com/pastelnetwork/gonode/common v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/p2p/kademlia v0.0.0-00010101000000-000000000000
)
