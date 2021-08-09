module github.com/pastelnetwork/gonode/metadb

go 1.16

replace (
	github.com/kolesa-team/go-webp => ../go-webp
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/dupedetection => ../dupedetection
	github.com/pastelnetwork/gonode/metadb => ../metadb
	github.com/pastelnetwork/gonode/p2p => ../p2p
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/probe => ../probe
	github.com/pastelnetwork/gonode/proto => ../proto
	github.com/pastelnetwork/gonode/raptorq => ../raptorq
	github.com/pastelnetwork/gonode/supernode => ../supernode
)

require (
	github.com/hashicorp/go-hclog v0.16.2
	github.com/hashicorp/raft v1.3.1
	github.com/hashicorp/raft-boltdb/v2 v2.0.0-20210422161416-485fa74b0b01
	github.com/kr/pretty v0.2.0 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/proto v0.0.0
	github.com/rqlite/go-sqlite3 v1.20.4
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/grpc v1.39.0 // indirect
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)
