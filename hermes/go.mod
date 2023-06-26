module github.com/pastelnetwork/gonode/hermes

go 1.17

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/mixins => ../mixins
	github.com/pastelnetwork/gonode/p2p => ../p2p
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/proto => ../proto
	github.com/pastelnetwork/gonode/raptorq => ../raptorq
)

require (
	github.com/DataDog/zstd v1.5.2
	github.com/btcsuite/btcutil v1.0.2
	github.com/bxcodec/faker/v4 v4.0.0-beta.3
	github.com/gorilla/mux v1.8.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/mattn/go-sqlite3 v1.14.14
	github.com/pastelnetwork/gonode/common v0.0.0-20220615180506-00d063d0abf2
	github.com/pastelnetwork/gonode/mixins v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/pastel v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/proto v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.1
	golang.org/x/crypto v0.10.0
	google.golang.org/grpc v1.55.0
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/makiuchi-d/gozxing v0.1.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/otrv4/ed448 v0.0.0-20210127123821-203e597250c3 // indirect
	github.com/pastelnetwork/gonode/raptorq v0.0.0-00010101000000-000000000000 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.12.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/subosito/gotenv v1.3.0 // indirect
	github.com/urfave/cli/v2 v2.8.1 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/image v0.0.0-20220302094943-723b81ca9867 // indirect
	golang.org/x/net v0.11.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/term v0.9.0 // indirect
	golang.org/x/text v0.10.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
