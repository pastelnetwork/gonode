module github.com/pastelnetwork/gonode/pastel

go 1.16

require (
	github.com/onsi/gomega v1.11.0
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/probe v0.0.0-20210716145051-623a1d4cc3b7
	github.com/stretchr/testify v1.7.0
)

replace github.com/pastelnetwork/gonode/common => ../common
