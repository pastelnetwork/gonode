module github.com/pastelnetwork/gonode/metadb

go 1.17

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/proto => ../proto
)

require (
	github.com/hashicorp/go-hclog v0.16.2
	github.com/hashicorp/raft v1.3.1
	github.com/hashicorp/raft-boltdb/v2 v2.0.0-20210422161416-485fa74b0b01
	github.com/kr/pretty v0.2.0 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/pastel v0.0.0-20210921141840-854e94914b9b
	github.com/pastelnetwork/gonode/proto v0.0.0
	github.com/rqlite/go-sqlite3 v1.20.4
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/grpc v1.39.0 // indirect
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/DataDog/zstd v1.4.8 // indirect
	github.com/armon/go-metrics v0.0.0-20190430140413-ec5e00d3c878 // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/go-errors/errors v1.1.1 // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.0.0 // indirect
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/hashicorp/go-uuid v1.0.1 // indirect
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace github.com/pastelnetwork/gonode/pastel => ../pastel
