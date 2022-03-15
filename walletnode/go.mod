module github.com/pastelnetwork/gonode/walletnode

go 1.17

require (
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/mixins v0.0.0-20220304222854-decb94f5cc25
	github.com/pastelnetwork/gonode/pastel v0.0.0
	github.com/pastelnetwork/gonode/proto v0.0.0
	github.com/pastelnetwork/gonode/raptorq v0.0.0-20220304222854-decb94f5cc25
	github.com/sahilm/fuzzy v0.1.0
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/stretchr/testify v1.7.0
	goa.design/goa/v3 v3.6.2
	goa.design/plugins/v3 v3.6.1
	golang.org/x/crypto v0.0.0-20220307211146-efcb8507fb70
	golang.org/x/image v0.0.0-20220302094943-723b81ca9867 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	google.golang.org/grpc v1.45.0
)

require github.com/DataDog/zstd v1.5.0

require (
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dimfeld/httppath v0.0.0-20170720192232-ee938bf73598 // indirect
	github.com/dimfeld/httptreemux/v5 v5.4.0 // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/makiuchi-d/gozxing v0.1.1 // indirect
	github.com/manveru/faker v0.0.0-20171103152722-9fbc68a78c4d // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/otrv4/ed448 v0.0.0-20210127123821-203e597250c3 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.8.1 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.10.1 // indirect
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/urfave/cli/v2 v2.3.0 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	github.com/zach-klippenstein/goregen v0.0.0-20160303162051-795b5e3961ea // indirect
	golang.org/x/mod v0.5.1 // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220310020820-b874c991c1a5 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.9 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20220310185008-1973136f34c6 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/metadb => ../metadb
	github.com/pastelnetwork/gonode/mixins => ../mixins
	github.com/pastelnetwork/gonode/p2p => ../p2p
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/proto => ../proto
	github.com/pastelnetwork/gonode/raptorq => ../raptorq
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20210921155107-089bfa567519
)
