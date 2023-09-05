module github.com/pastelnetwork/gonode/dupedetection

go 1.21

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/metadb => ../metadb
	github.com/pastelnetwork/gonode/p2p => ../p2p
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/proto => ../proto
)

require (
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.3.1
	github.com/pastelnetwork/gonode/common v0.0.0-20230829114806-701c926c0c44
	github.com/pastelnetwork/gonode/pastel v0.0.0-20230829114806-701c926c0c44
	github.com/stretchr/testify v1.8.4
	google.golang.org/grpc v1.57.0
)

require (
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stretchr/objx v0.5.1 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/term v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	google.golang.org/genproto v0.0.0-20230803162519-f966b187b2e5 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
