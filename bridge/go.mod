module github.com/pastelnetwork/gonode/bridge

go 1.17

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/metadb => ../metadb
	github.com/pastelnetwork/gonode/mixins => ../mixins
	github.com/pastelnetwork/gonode/p2p => ../p2p
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/proto => ../proto
	github.com/pastelnetwork/gonode/raptorq => ../raptorq
)

require (
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/mixins v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/pastel v0.0.0-20220527201150-f7f5e361d5b7
	github.com/pastelnetwork/gonode/proto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.43.0
)

require (
	github.com/DataDog/zstd v1.5.0 // indirect
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0-20190314233015-f79a8a8ca69d // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-errors/errors v1.4.1 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/makiuchi-d/gozxing v0.1.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/otrv4/ed448 v0.0.0-20210127123821-203e597250c3 // indirect
	github.com/pastelnetwork/gonode/raptorq v0.0.0-00010101000000-000000000000 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.10.1 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/urfave/cli/v2 v2.3.0 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
