module github.com/pastelnetwork/gonode/dupedetection

go 1.17

require (
	github.com/DataDog/zstd v1.5.0
	github.com/google/uuid v1.3.0
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/metadb v0.0.0-20210803173228-bd776aecea6a
	github.com/pastelnetwork/gonode/p2p v0.0.0-20210803173228-bd776aecea6a
	github.com/pastelnetwork/gonode/pastel v0.0.0
	github.com/sbinet/npyio v0.5.2
	github.com/stretchr/testify v1.7.0
	gonum.org/v1/gonum v0.8.2
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/anacrolix/log v0.9.0 // indirect
	github.com/anacrolix/missinggo v1.3.0 // indirect
	github.com/anacrolix/missinggo/perf v1.0.0 // indirect
	github.com/anacrolix/missinggo/v2 v2.5.2 // indirect
	github.com/anacrolix/sync v0.4.0 // indirect
	github.com/anacrolix/utp v0.1.0 // indirect
	github.com/andres-erbsen/clock v0.0.0-20160526145045-9e14626cd129 // indirect
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/campoy/embedmd v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-errors/errors v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/jmoiron/sqlx v1.3.4 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/kr/text v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rqlite/go-sqlite3 v1.22.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	go.uber.org/ratelimit v0.2.0 // indirect
	golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871 // indirect
	golang.org/x/exp v0.0.0-20191030013958-a1ab85dbe136 // indirect
	golang.org/x/net v0.0.0-20211123203042-d83791d6bcd9 // indirect
	golang.org/x/sys v0.0.0-20211124211545-fe61309f8881 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/metadb => ../metadb
	github.com/pastelnetwork/gonode/p2p => ../p2p
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/proto => ../proto
)
