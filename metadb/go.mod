module github.com/pastelnetwork/gonode/metadb

go 1.17

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/proto => ../proto
)

require (
	github.com/hashicorp/go-hclog v1.2.1
	github.com/hashicorp/raft v1.3.9
	github.com/hashicorp/raft-boltdb/v2 v2.2.2
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pastelnetwork/gonode/common v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/pastel v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/proto v0.0.0-00010101000000-000000000000
	github.com/rqlite/go-sqlite3 v1.25.0
	github.com/stretchr/testify v1.7.2
	github.com/tj/assert v0.0.3
	golang.org/x/net v0.0.0-20220617184016-355a448f1bc9
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/DataDog/zstd v1.5.2 // indirect
	github.com/armon/go-metrics v0.3.10 // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/hashicorp/go-uuid v1.0.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.6.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd // indirect
	google.golang.org/grpc v1.47.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
