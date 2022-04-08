module github.com/pastelnetwork/gonode/dupedetection

go 1.17

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/metadb => ../metadb
	github.com/pastelnetwork/gonode/p2p => ../p2p
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/proto => ../proto
)

require (
	github.com/DataDog/zstd v1.5.0
	github.com/google/uuid v1.1.2
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/metadb v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/p2p v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/pastel v0.0.0-20210921141840-854e94914b9b
	github.com/sbinet/npyio v0.6.0
	github.com/stretchr/testify v1.7.0
	gonum.org/v1/gonum v0.9.3
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
)

require golang.org/x/image v0.0.0-20211028202545-6944b10bf410 // indirect

require (
	github.com/anacrolix/log v0.9.0 // indirect
	github.com/anacrolix/missinggo v1.2.1 // indirect
	github.com/anacrolix/missinggo/perf v1.0.0 // indirect
	github.com/anacrolix/sync v0.2.0 // indirect
	github.com/anacrolix/utp v0.1.0 // indirect
	github.com/andres-erbsen/clock v0.0.0-20160526145045-9e14626cd129 // indirect
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/campoy/embedmd v1.0.0 // indirect
	github.com/chai2010/webp v1.1.1
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/disintegration/imaging v1.6.2
	github.com/go-errors/errors v1.4.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/huandu/xstrings v1.3.0 // indirect
	github.com/jmoiron/sqlx v1.3.4 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rqlite/go-sqlite3 v1.22.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	go.uber.org/ratelimit v0.2.0 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
