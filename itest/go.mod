module github.com/pastelnetwork/gonode/itest

go 1.16

require (
	github.com/onsi/ginkgo v1.16.2
	github.com/onsi/gomega v1.11.0
	github.com/pastelnetwork/gonode/common v0.0.0-00010101000000-000000000000
)

replace github.com/pastelnetwork/gonode/common => ../common
