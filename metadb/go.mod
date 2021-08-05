module github.com/pastelnetwork/gonode/metadb

go 1.16

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/supernode => ../supernode
)

require (
	github.com/hashicorp/go-hclog v0.16.2 // indirect
	github.com/hashicorp/raft v1.3.1 // indirect
	github.com/hashicorp/raft-boltdb/v2 v2.0.0-20210422161416-485fa74b0b01 // indirect
	github.com/pastelnetwork/gonode/common v0.0.0-20210803173228-bd776aecea6a // indirect
	github.com/rqlite/go-sqlite3 v1.20.4 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
