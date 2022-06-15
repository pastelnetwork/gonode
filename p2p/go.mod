module github.com/pastelnetwork/gonode/p2p

go 1.17

require (
	github.com/anacrolix/envpprof v1.1.1 // indirect
	github.com/anacrolix/sync v0.2.0 // indirect
	github.com/anacrolix/utp v0.1.0
	github.com/bradfitz/iter v0.0.0-20191230175014-e8f45d346db8 // indirect
	github.com/btcsuite/btcutil v1.0.2
	github.com/jmoiron/sqlx v1.3.5
	github.com/otrv4/ed448 v0.0.0-20210127123821-203e597250c3
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/pastel v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
	github.com/tj/assert v0.0.3
	go.uber.org/ratelimit v0.2.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	google.golang.org/grpc v1.43.0

)

require github.com/rqlite/go-sqlite3 v1.25.0

require (
	github.com/DataDog/zstd v1.5.0 // indirect
	github.com/anacrolix/missinggo v1.2.1 // indirect
	github.com/anacrolix/missinggo/perf v1.0.0 // indirect
	github.com/andres-erbsen/clock v0.0.0-20160526145045-9e14626cd129 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-errors/errors v1.4.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/huandu/xstrings v1.2.0 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/pastel => ../pastel
)
