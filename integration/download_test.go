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

var _ = Describe("Download", func() {
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
		id, err := helper.GetP2PID(symbol)
		Expect(err).NotTo(HaveOccurred())
		storeReq.Value = symbol

		resp, status, err := itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN1BaseURI), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))
		Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

		storeReq = &helper.StoreRequest{}
		storeReply = &helper.StoreReply{}

		rqfile.SymbolIdentifiers = []string{id}
		idFile, err := rqfile.GetIDFile()
		Expect(err).NotTo(HaveOccurred())

		storeReq.Value = idFile
		resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN1BaseURI), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))
		Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

		// Mock Expectations
		Expect(mocker.MockAllRegExpections()).To(Succeed())

	})

	When("user tries to download an existing NFT ", func() {
		Context("user owns the NFT", func() {
			It("should download NFT successfully", func() {
				downResp, status, err := itHelper.Request(helper.HttpGet, nil,
					helper.GetDownloadURI(it.WNBaseURI, testconst.ArtistPastelID,
						testconst.TestRegTXID), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))

				Expect(json.Unmarshal(downResp, &downloadRsp)).To(Succeed())
				Expect(len(fmt.Sprintf("%v", downloadRsp["file"]))).NotTo(BeZero())
			})
		})

		Context("user does not own the NFT", func() {
			It("should not download nft", func() {
				downResp, status, err := itHelper.Request(helper.HttpGet, nil,
					// Passing a different PastelID
					helper.GetDownloadURI(it.WNBaseURI, "j2Y1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gY",
						testconst.TestRegTXID), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusBadRequest))

				Expect(json.Unmarshal(downResp, &downloadRsp)).To(Succeed())
				Expect(fmt.Sprintf("%v", downloadRsp["message"])).To(Equal("failed to verify ownership"))
			})
		})
	})

	When("user tries to download an existing Cascade artifact ", func() {
		Context("user owns the artifact", func() {
			It("should download artifact successfully", func() {
				downResp, status, err := itHelper.Request(helper.HttpGet, nil,
					helper.GetCascadeDownloadURI(it.WNBaseURI, testconst.ArtistPastelID,
						testconst.TestCascadeRegTXID), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))

				Expect(json.Unmarshal(downResp, &downloadRsp)).To(Succeed())
				Expect(len(fmt.Sprintf("%v", downloadRsp["file"]))).NotTo(BeZero())
			})
		})

		Context("user does not own the Artifact", func() {
			It("should not download artifact", func() {
				downResp, status, err := itHelper.Request(helper.HttpGet, nil,
					// Passing a different PastelID
					helper.GetCascadeDownloadURI(it.WNBaseURI, "j2Y1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gY",
						testconst.TestCascadeRegTXID), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusBadRequest))

				Expect(json.Unmarshal(downResp, &downloadRsp)).To(Succeed())
				Expect(fmt.Sprintf("%v", downloadRsp["message"])).To(Equal("failed to verify ownership"))
			})
		})
	})

	When("user tries to download an existing Sense dd and fp file ", func() {
		Context("user owns the artifact", func() {
			It("should download ddAndFp file successfully", func() {
				downResp, status, err := itHelper.Request(helper.HttpGet, nil,
					helper.GetSenseDownloadURI(it.WNBaseURI, testconst.ArtistPastelID,
						testconst.TestSenseRegTXID), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))

				Expect(json.Unmarshal(downResp, &downloadRsp)).To(Succeed())
				Expect(len(fmt.Sprintf("%v", downloadRsp["file"]))).NotTo(BeZero())
			})
		})

		Context("user does not own the Artifact", func() {
			It("should still download ddAndFp file", func() {
				downResp, status, err := itHelper.Request(helper.HttpGet, nil,
					// Passing a different PastelID
					helper.GetSenseDownloadURI(it.WNBaseURI, "j2Y1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gY",
						testconst.TestSenseRegTXID), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusBadRequest))

				Expect(json.Unmarshal(downResp, &downloadRsp)).To(Succeed())
				Expect(fmt.Sprintf("%v", downloadRsp["message"])).To(Equal("failed to verify ownership"))
			})
		})
	})

	AfterEach(func() {
		Expect(mocker.CleanupAll()).To(Succeed())
	})
})
