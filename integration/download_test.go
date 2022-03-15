package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pastelnetwork/gonode/integration/fakes/common/testconst"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	it "github.com/pastelnetwork/gonode/integration"
	helper "github.com/pastelnetwork/gonode/integration/helper"
	"github.com/pastelnetwork/gonode/integration/mock"
)

var _ = Describe("NFTDownload", func() {
	var (
		itHelper    = helper.NewItHelper()
		mocker      *mock.Mocker
		downloadRsp map[string]interface{}

		storeReq   *helper.StoreRequest
		storeReply *helper.StoreReply
		rqfile     *helper.RawSymbolIDFile
		symbol     []byte
	)

	BeforeEach(func() {
		storeReq = &helper.StoreRequest{}
		storeReply = &helper.StoreReply{}

		downloadRsp = make(map[string]interface{})
		rqfile = &helper.RawSymbolIDFile{
			ID:       "09f6c459-ec2a-4db1-a8fe-0648fd97b5cb",
			PastelID: "jXY1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gW",
		}

		mocker = mock.New(it.PasteldServers, it.DDServers, it.RQServers, it.SNServers, itHelper)
		symbol = []byte("test-symbol")
	})

	Context("when downloading an existing & owned valid nft", func() {
		When("storing rq symbols in p2p for the test", func() {
			It("should be stored successfully", func() {
				Expect(mocker.MockAllRegExpections()).To(Succeed())

				id, err := helper.GetP2PID(symbol)
				Expect(err).NotTo(HaveOccurred())
				storeReq.Value = symbol

				// TODO: Why keys are not propagating through p2p network? need to investigate
				resp, status, err := itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN1BaseURI), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
				Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

				resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN2BaseURI), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
				Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

				resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN3BaseURI), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
				Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

				storeReq = &helper.StoreRequest{}
				storeReply = &helper.StoreReply{}

				rqfile.SymbolIdentifiers = []string{id}
				idFile, err := rqfile.GetIDFile()
				Expect(err).NotTo(HaveOccurred())

				storeReq.Value = idFile

				// TODO: Why keys are not propagating through p2p network? need to investigate
				resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN1BaseURI), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
				Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

				resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN2BaseURI), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
				Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

				resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN3BaseURI), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
				Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

				downResp, status, err := itHelper.Request(helper.HttpGet, nil,
					helper.GetDownloadURI(it.WNBaseURI, testconst.ArtistPastelID,
						testconst.TestRegTXID), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))

				Expect(json.Unmarshal(downResp, &downloadRsp)).To(Succeed())
				Expect(len(fmt.Sprintf("%v", downloadRsp["file"]))).NotTo(BeZero())
			})
		})
	})

	AfterEach(func() {
		Expect(mocker.CleanupAll()).To(Succeed())
	})
})
