module github.com/pastelnetwork/gonode/pastel

go 1.16

require (
	github.com/DataDog/zstd v1.4.8
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/onsi/gomega v1.11.0
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
)

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/probe => ../probe
)
