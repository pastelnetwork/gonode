module github.com/pastelnetwork/gonode/p2p

go 1.17

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/pastel => ../pastel
)

require (
	github.com/anacrolix/utp v0.1.0
	github.com/btcsuite/btcutil v1.0.2
	github.com/dgraph-io/badger/v3 v3.2103.2
	github.com/otrv4/ed448 v0.0.0-20210127123821-203e597250c3
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/pastel v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
	github.com/tj/assert v0.0.3
	go.uber.org/ratelimit v0.2.0
	golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871
	google.golang.org/grpc v1.42.0
)

require (
	github.com/DataDog/zstd v1.4.8 // indirect
	github.com/anacrolix/missinggo v1.3.0 // indirect
	github.com/anacrolix/missinggo/perf v1.0.0 // indirect
	github.com/anacrolix/missinggo/v2 v2.5.1 // indirect
	github.com/anacrolix/sync v0.4.0 // indirect
	github.com/andres-erbsen/clock v0.0.0-20160526145045-9e14626cd129 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgraph-io/ristretto v0.1.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/go-errors/errors v1.1.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/flatbuffers v1.12.1 // indirect
	github.com/huandu/xstrings v1.3.1 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/klauspost/compress v1.12.3 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	go.opencensus.io v0.22.5 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
