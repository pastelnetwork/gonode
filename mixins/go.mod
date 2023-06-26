module github.com/pastelnetwork/gonode/mixins

go 1.17

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/dupedetection => ../dupedetection
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/raptorq => ../raptorq
)

require (
	github.com/DataDog/zstd v1.5.2
	github.com/btcsuite/btcutil v1.0.2
	github.com/pastelnetwork/gonode/common v0.0.0-20220615180506-00d063d0abf2
	github.com/pastelnetwork/gonode/pastel v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/gonode/raptorq v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/uuid v1.3.0
	github.com/tj/assert v0.0.3
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/makiuchi-d/gozxing v0.1.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	golang.org/x/crypto v0.10.0 // indirect
	golang.org/x/image v0.0.0-20220302094943-723b81ca9867 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/term v0.9.0 // indirect
	golang.org/x/text v0.10.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
