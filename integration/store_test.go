package main_test

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	it "github.com/pastelnetwork/gonode/integration"
	"github.com/pastelnetwork/gonode/integration/helper"
)

var _ = Describe("StoreRetrieve", func() {
	var (
		itHelper   = helper.NewItHelper()
		storeReq   *helper.StoreRequest
		storeReply *helper.StoreReply
		getReply   *helper.RetrieveResponse
	)

	BeforeEach(func() {
		storeReq = &helper.StoreRequest{
			Value: []byte("test-data"),
		}
		storeReply = &helper.StoreReply{}
		getReply = &helper.RetrieveResponse{}
	})

	Context("when retrieving data from same node", func() {
		It("should exist", func() {
			resp, status, err := itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN1BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

			getResp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetRetrieveURI(it.SN1BaseURI, storeReply.Key), nil)
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
			resp, status, err := itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN1BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

			getResp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetRetrieveURI(it.SN2BaseURI, storeReply.Key), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(getResp, getReply)).To(Succeed())

			Expect(getReply.Key).To(Equal(storeReply.Key))
			Expect(storeReq.Value).To(Equal(getReply.Value))
		})
	})

	Context("when storing the data in sn1 & retrieving from sn3", func() {
		It("should exist", func() {
			storeReq.Value = []byte("test-data-3")
			resp, status, err := itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN1BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

			getResp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetRetrieveURI(it.SN3BaseURI, storeReply.Key), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(getResp, getReply)).To(Succeed())

			Expect(getReply.Key).To(Equal(storeReply.Key))
			Expect(storeReq.Value).To(Equal(getReply.Value))
		})
	})

	Context("when storing the data in sn3 & retrieving from sn1", func() {
		It("should exist", func() {
			storeReq.Value = []byte("test-data-4")
			resp, status, err := itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN3BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

			getResp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetRetrieveURI(it.SN1BaseURI, storeReply.Key), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(getResp, getReply)).To(Succeed())

			Expect(getReply.Key).To(Equal(storeReply.Key))
			Expect(storeReq.Value).To(Equal(getReply.Value))
		})
	})

	Context("when storing the data in sn2 & retrieving from sn3", func() {
		It("should exist", func() {
			storeReq.Value = []byte("test-data-5")
			resp, status, err := itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN2BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

			getResp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetRetrieveURI(it.SN3BaseURI, storeReply.Key), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(getResp, getReply)).To(Succeed())

			Expect(getReply.Key).To(Equal(storeReply.Key))
			Expect(storeReq.Value).To(Equal(getReply.Value))
		})
	})

	AfterEach(func() {
		storeReq.Value = []byte{}
	})
})
