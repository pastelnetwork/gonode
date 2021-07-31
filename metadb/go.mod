module github.com/pastelnetwork/gonode/metadb

go 1.13

require (
	github.com/armon/go-metrics v0.3.7 // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/hashicorp/go-hclog v0.16.0
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/hashicorp/go-msgpack v1.1.5 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/raft v1.3.0
	github.com/hashicorp/raft-boltdb v0.0.0-20210422161416-485fa74b0b01 // indirect
	github.com/hashicorp/raft-boltdb/v2 v2.0.0-20210422161416-485fa74b0b01
	github.com/kr/pretty v0.2.1 // indirect
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/rqlite/go-sqlite3 v1.20.2
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7 // indirect
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
)

replace github.com/pastelnetwork/gonode/common => ../common
