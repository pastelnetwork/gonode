package integration_test

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	it "github.com/pastelnetwork/gonode/integration"
)

var _ = Describe("StoreRetrieve", func() {
	var (
		helper     = it.NewItHelper()
		storeReq   *it.StoreRequest
		storeReply *it.StoreReply
		getReply   *it.RetrieveResponse
		//err        error
	)

	BeforeEach(func() {
		storeReq = &it.StoreRequest{
			Value: []byte("test-data"),
		}
		storeReply = &it.StoreReply{}
		getReply = &it.RetrieveResponse{}
	})

	Context("when retrieving data from same node", func() {
		It("should exist", func() {
			resp, status, err := helper.Request(it.HttpPost, storeReq, it.GetStoreURI(it.SN1BaseURI))
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

			getResp, status, err := helper.Request(it.HttpGet, nil, it.GetRetrieveURI(it.SN1BaseURI, storeReply.Key))
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(getResp, getReply)).To(Succeed())

			Expect(getReply.Key).To(Equal(storeReply.Key))
			Expect(getReply.Value).To(Equal(getReply.Value))
		})
	})

	Context("when storing the data in sn1 & retrieving from sn2", func() {
		It("should exist", func() {
			storeReq.Value = []byte("test-data-2")
			resp, status, err := helper.Request(it.HttpPost, storeReq, it.GetStoreURI(it.SN1BaseURI))
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

			getResp, status, err := helper.Request(it.HttpGet, nil, it.GetRetrieveURI(it.SN2BaseURI, storeReply.Key))
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(getResp, getReply)).To(Succeed())

			Expect(getReply.Key).To(Equal(storeReply.Key))
			Expect(getReply.Value).To(Equal(getReply.Value))
		})
	})

	Context("when storing the data in sn1 & retrieving from sn3", func() {
		It("should exist", func() {
			storeReq.Value = []byte("test-data-3")
			resp, status, err := helper.Request(it.HttpPost, storeReq, it.GetStoreURI(it.SN1BaseURI))
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

			getResp, status, err := helper.Request(it.HttpGet, nil, it.GetRetrieveURI(it.SN3BaseURI, storeReply.Key))
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(getResp, getReply)).To(Succeed())

			Expect(getReply.Key).To(Equal(storeReply.Key))
			Expect(getReply.Value).To(Equal(getReply.Value))
		})
	})

	Context("when storing the data in sn3 & retrieving from sn1", func() {
		It("should exist", func() {
			storeReq.Value = []byte("test-data-4")
			resp, status, err := helper.Request(it.HttpPost, storeReq, it.GetStoreURI(it.SN3BaseURI))
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

			getResp, status, err := helper.Request(it.HttpGet, nil, it.GetRetrieveURI(it.SN1BaseURI, storeReply.Key))
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(getResp, getReply)).To(Succeed())

			Expect(getReply.Key).To(Equal(storeReply.Key))
			Expect(getReply.Value).To(Equal(getReply.Value))
		})
	})

	Context("when storing the data in sn2 & retrieving from sn3", func() {
		It("should exist", func() {
			storeReq.Value = []byte("test-data-5")
			resp, status, err := helper.Request(it.HttpPost, storeReq, it.GetStoreURI(it.SN2BaseURI))
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

			getResp, status, err := helper.Request(it.HttpGet, nil, it.GetRetrieveURI(it.SN3BaseURI, storeReply.Key))
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(getResp, getReply)).To(Succeed())

			Expect(getReply.Key).To(Equal(storeReply.Key))
			Expect(getReply.Value).To(Equal(getReply.Value))
		})
	})

	AfterEach(func() {
		storeReq.Value = []byte{}
	})
})
