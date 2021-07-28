module github.com/pastelnetwork/gonode/pastel

go 1.16

require (
	github.com/onsi/gomega v1.11.0
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.7.0
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887 // indirect
)

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/probe => ../probe
)
