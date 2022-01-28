module github.com/pastelnetwork/gonode/mixins

go 1.17

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/dupedetection => ../dupedetection
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/raptorq => ../raptorq
)

require (
	github.com/DataDog/zstd v1.5.0
	github.com/btcsuite/btcutil v1.0.2
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/pastel v0.0.0-20210921141840-854e94914b9b
	github.com/pastelnetwork/gonode/raptorq v0.0.0-00010101000000-000000000000
)

require (
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/go-errors/errors v1.4.1 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/makiuchi-d/gozxing v0.1.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)
