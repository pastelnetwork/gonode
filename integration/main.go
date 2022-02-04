// package main only used for debugging conatiners
// in order to run integration tests, do
// `INTEGRATION_TEST_ENV=true go test -v --timeout=10m`
package main

import (
	"github.com/pastelnetwork/gonode/integration/helper"
	"github.com/pastelnetwork/gonode/integration/mock"
)

const (
	// SN1BaseURI of SN1 Server
	SN1BaseURI = "http://localhost:19090"
	// SN2BaseURI of SN2 Server
	SN2BaseURI = "http://localhost:19091"
	// SN3BaseURI of SN3 Server
	SN3BaseURI = "http://localhost:19092"
	// SN4BaseURI of SN4 Server
	SN4BaseURI = "http://localhost:19093"

	// WNBaseURI of WN
	WNBaseURI = "http://localhost:18089"

	// DD-Servers
	dd1BaseURI = "http://localhost:51052"
	dd2BaseURI = "http://localhost:52052"
	dd3BaseURI = "http://localhost:53052"
	dd4BaseURI = "http://localhost:54052"

	// RQ-Servers
	rq0BaseURI = "http://localhost:49051"
	rq1BaseURI = "http://localhost:51051"
	rq2BaseURI = "http://localhost:52051"
	rq3BaseURI = "http://localhost:53051"
	rq4BaseURI = "http://localhost:54051"

	// Pastel-Servers
	p0BaseURI = "http://localhost:29930"
	p1BaseURI = "http://localhost:29931"
	p2BaseURI = "http://localhost:29932"
	p3BaseURI = "http://localhost:29933"
	p4BaseURI = "http://localhost:29934"
)

var SNServers = []string{
	SN1BaseURI,
	SN2BaseURI,
	SN3BaseURI,
	SN4BaseURI,
}

var DDServers = []string{
	dd1BaseURI,
	dd2BaseURI,
	dd3BaseURI,
	dd4BaseURI,
}

var RQServers = []string{
	rq0BaseURI,
	rq1BaseURI,
	rq2BaseURI,
	rq3BaseURI,
	rq4BaseURI,
}

var PasteldServers = []string{
	p0BaseURI,
	p1BaseURI,
	p2BaseURI,
	p3BaseURI,
	p4BaseURI,
}

func main() {
	itHelper := helper.NewItHelper()

	mocker := mock.New(PasteldServers, DDServers, RQServers, itHelper)
	if err := mocker.MockAllRegExpections(); err != nil {
		panic(err)
	}
}
