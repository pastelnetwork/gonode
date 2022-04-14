module github.com/pastelnetwork/gonode/pastel

go 1.16

replace github.com/pastelnetwork/gonode/common => ../common

require (
	github.com/DataDog/zstd v1.5.0
	github.com/btcsuite/btcutil v1.0.2
	github.com/onsi/gomega v1.18.1
	github.com/pastelnetwork/gonode/common v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
)
