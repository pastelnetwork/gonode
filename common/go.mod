module github.com/pastelnetwork/gonode/common

go 1.22.0

replace (
	github.com/kolesa-team/go-webp => ../go-webp
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/raptorq => ../raptorq
	golang.org/x/crypto => github.com/golang/crypto v0.13.0
)

require (
	github.com/cloudflare/circl v1.3.3
	github.com/disintegration/imaging v1.6.2
	github.com/fogleman/gg v1.3.0
	github.com/go-errors/errors v1.5.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629
	github.com/json-iterator/go v1.1.12
	github.com/klauspost/compress v1.17.7
	github.com/makiuchi-d/gozxing v0.1.1
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/otrv4/ed448 v0.0.0-20221017120334-a33859724cfd
	github.com/pastelnetwork/gonode/go-webp v0.0.0-20240307111701-c491864e6826
	github.com/pastelnetwork/gonode/pastel v0.0.0-00010101000000-000000000000
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/viper v1.16.0
	github.com/stretchr/testify v1.9.0
	github.com/urfave/cli/v2 v2.25.7
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	golang.org/x/crypto v0.18.0
	golang.org/x/sync v0.6.0
	golang.org/x/sys v0.18.0
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/term v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
