// package main only used for debugging conatiners
// in order to run integration tests, do
// `INTEGRATION_TEST_ENV=true go test -v --timeout=10m`
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	// SN5BaseURI of SN5 Server
	SN5BaseURI = "http://localhost:19094"
	// SN6BaseURI of SN6 Server
	SN6BaseURI = "http://localhost:19095"
	// SN7BaseURI of SN6 Server
	SN7BaseURI = "http://localhost:19096"

	// WNBaseURI of WN
	WNBaseURI = "http://localhost:18089"

	// DD-Servers
	dd1BaseURI = "http://localhost:51052"
	dd2BaseURI = "http://localhost:52052"
	dd3BaseURI = "http://localhost:53052"
	dd4BaseURI = "http://localhost:54052"
	dd5BaseURI = "http://localhost:55052"
	dd6BaseURI = "http://localhost:56052"

	// Pastel-Servers
	p0BaseURI = "http://localhost:29930"
	p1BaseURI = "http://localhost:29931"
	p2BaseURI = "http://localhost:29932"
	p3BaseURI = "http://localhost:29933"
	p4BaseURI = "http://localhost:29934"
	p5BaseURI = "http://localhost:29935"
	p6BaseURI = "http://localhost:29936"
	p7BaseURI = "http://localhost:29937"
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

var PasteldServers = []string{
	p0BaseURI,
	p1BaseURI,
	p2BaseURI,
	p3BaseURI,
	p4BaseURI,
}

var SCSNServers = []string{
	SN1BaseURI,
	SN2BaseURI,
	SN3BaseURI,
	SN4BaseURI,
	SN5BaseURI,
	SN6BaseURI,
	SN7BaseURI,
}

var SCDDServers = []string{
	dd1BaseURI,
	dd2BaseURI,
	dd3BaseURI,
	dd4BaseURI,
	dd5BaseURI,
	dd6BaseURI,
}

var SCPasteldServers = []string{
	p0BaseURI,
	p1BaseURI,
	p2BaseURI,
	p3BaseURI,
	p4BaseURI,
	p5BaseURI,
	p6BaseURI,
	p7BaseURI,
}

func main() {
	itHelper := helper.NewItHelper()

	mocker := mock.New(PasteldServers, DDServers, SNServers, itHelper)
	/*if err := mocker.MockAllRegExpections(); err != nil {
		panic(err)
	}*/

	if err := mocker.MockSCRegExpections(); err != nil {
		panic(err)
	}

	// insert thumbnails for search
	storeReq := &helper.StoreRequest{}
	storeReq.Value = []byte("thumbnail-one")

	resp, status, err := itHelper.Request(helper.HttpPost, storeReq, helper.GetLocalStoreURI(SN3BaseURI), nil)
	if err != nil {
		panic(fmt.Sprintf("unable to store thumbnail for search tests: %s", err.Error()))
	} else if status != http.StatusOK {
		panic(fmt.Sprintf("unable to store thumbnail for search tests: %s", string(resp)))
	}

	storeReq.Value = []byte("thumbnail-two")
	resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetLocalStoreURI(SN3BaseURI), nil)
	if err != nil {
		panic(fmt.Sprintf("unable to store thumbnail for search tests: %s", err.Error()))
	} else if status != http.StatusOK {
		panic(fmt.Sprintf("unable to store thumbnail for search tests: %s", string(resp)))
	}

	rqfile := &helper.RawSymbolIDFile{
		ID:       "09f6c459-ec2a-4db1-a8fe-0648fd97b5cb",
		PastelID: "jXY1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gW",
	}
	symbol := []byte("test-symbol")
	storeReq.Value = symbol
	resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetLocalStoreURI(SN3BaseURI), nil)
	if err != nil {
		panic(fmt.Sprintf("store thumbnail for search tests: %s", err.Error()))
	} else if status != http.StatusOK {
		panic(fmt.Sprintf("store thumbnail for search tests: %s", string(resp)))
	}

	storeReq = &helper.StoreRequest{}
	id, err := helper.GetP2PID(symbol)
	if err != nil {
		panic(fmt.Sprintf("generate p2p id of symbol for download tests: %s", err.Error()))
	}
	rqfile.SymbolIdentifiers = []string{id}
	idFile, err := rqfile.GetIDFile()
	if err != nil {
		panic(fmt.Sprintf("store generate id file for download tests: %s", err.Error()))
	}

	storeReq.Value = idFile
	fmt.Println("this is the correct file:", string(idFile))
	resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetLocalStoreURI(SN3BaseURI), nil)
	if err != nil {
		panic(fmt.Sprintf("store rq_id file for download tests: %s", err.Error()))
	} else if status != http.StatusOK {
		panic(fmt.Sprintf("store rq_id file for download tests: %s", string(resp)))
	}

	resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetLocalStoreURI(SN4BaseURI), nil)
	if err != nil {
		panic(fmt.Sprintf("store rq_id file for download tests: %s", err.Error()))
	} else if status != http.StatusOK {
		panic(fmt.Sprintf("store rq_id file for download tests: %s", string(resp)))
	}

	repl := &helper.StoreReply{}
	if err := json.Unmarshal(resp, repl); err != nil {
		panic(err)
	}

	fmt.Printf("rq id file hash: %s", repl.Key)
}
