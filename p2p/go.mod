module github.com/pastelnetwork/gonode/p2p

go 1.22.0

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/utils => ../utils
)

require (
	github.com/anacrolix/utp v0.2.0
	github.com/btcsuite/btcutil v1.0.2
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/json-iterator/go v1.1.12
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/otrv4/ed448 v0.0.0-20221017120334-a33859724cfd
	github.com/pastelnetwork/gonode/common v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/pastel v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.9.0
	github.com/tj/assert v0.0.3
	go.uber.org/ratelimit v0.3.1
	golang.org/x/sync v0.6.0
	google.golang.org/grpc v1.62.1
)

require (
	github.com/anacrolix/missinggo v1.3.0 // indirect
	github.com/anacrolix/missinggo/perf v1.0.0 // indirect
	github.com/anacrolix/missinggo/v2 v2.5.1 // indirect
	github.com/anacrolix/sync v0.4.0 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/huandu/xstrings v1.3.1 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/term v0.16.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
