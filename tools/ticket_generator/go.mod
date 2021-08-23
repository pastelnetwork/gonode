module github.com/pastelnetwork/gonode/tools/ticket_generator

go 1.17

require (
	github.com/go-errors/errors v1.4.0
	github.com/gorilla/websocket v1.4.2
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/pastel v0.0.0
	github.com/pastelnetwork/gonode/walletnode v0.0.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	goa.design/goa/v3 v3.4.3
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d
)

require (
	github.com/DataDog/zstd v1.4.8 // indirect
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dimfeld/httptreemux/v5 v5.3.0 // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/makiuchi-d/gozxing v0.0.0-20210610151345-232d6193415d // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/pastelnetwork/gonode/raptorq v0.0.0-20210731175226-8a39b3090588 // indirect
	github.com/pelletier/go-toml v1.7.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/image v0.0.0-20210216034530-4410531fe030 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/ini.v1 v1.51.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace (
	github.com/pastelnetwork/gonode/common => ../../common
	github.com/pastelnetwork/gonode/p2p => ../../p2p
	github.com/pastelnetwork/gonode/pastel => ../../pastel
	github.com/pastelnetwork/gonode/proto => ../../proto
	github.com/pastelnetwork/gonode/raptorq => ../../raptorq
	github.com/pastelnetwork/gonode/walletnode => ../../walletnode
)
