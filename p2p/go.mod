module github.com/pastelnetwork/gonode/p2p

go 1.16

replace github.com/pastelnetwork/gonode/p2p/kademlia => ./kademlia

replace github.com/pastelnetwork/gonode/p2p/kademlia/dao => ./kademlia/dao

replace github.com/pastelnetwork/gonode/p2p/kademlia/crypto => ./kademlia/crypto

replace github.com/pastelnetwork/gonode/common => ../common

require (
	github.com/pastelnetwork/gonode/common v0.0.0-20210523194944-207a3e9a74d9
	github.com/pastelnetwork/gonode/p2p/kademlia v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/p2p/kademlia/dao v0.0.0-00010101000000-000000000000
)
