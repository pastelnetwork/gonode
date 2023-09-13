package main_test

import (
	"net/http"
	"strings"

	json "github.com/json-iterator/go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	it "github.com/pastelnetwork/gonode/integration"
	helper "github.com/pastelnetwork/gonode/integration/helper"
	"github.com/pastelnetwork/gonode/integration/mock"
)

var _ = Describe("NFTSearch", func() {
	var (
		itHelper   = helper.NewItHelper()
		mocker     *mock.Mocker
		storeReq   *helper.StoreRequest
		storeReply *helper.StoreReply
	)

	BeforeEach(func() {
		mocker = mock.New(it.PasteldServers, it.DDServers, it.SNServers, itHelper)
		Expect(mocker.MockAllRegExpections()).To(Succeed())

		storeReq = &helper.StoreRequest{}
		storeReply = &helper.StoreReply{}
		storeReq.Value = []byte("thumbnail-one")

		resp, status, err := itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN1BaseURI), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))
		Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

		storeReq.Value = []byte("thumbnail-two")

		resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN1BaseURI), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))
		Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

		storeReq = &helper.StoreRequest{}
		storeReply = &helper.StoreReply{}
		storeReq.Value = []byte("thumbnail-one")

		resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN2BaseURI), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))
		Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

		storeReq.Value = []byte("thumbnail-two")

		resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN2BaseURI), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))
		Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

		storeReq = &helper.StoreRequest{}
		storeReply = &helper.StoreReply{}
		storeReq.Value = []byte("thumbnail-one")

		resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN3BaseURI), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))
		Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

		storeReq.Value = []byte("thumbnail-two")

		resp, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN3BaseURI), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))
		Expect(json.Unmarshal(resp, storeReply)).To(Succeed())

	})

	Context("when searching an existing nft with matching creator name", func() {
		It("should be found & returned", func() {
			searchParams := make(map[string]string)
			searchParams["query"] = "mona"
			searchParams["creator_name"] = "true"
			searchParams["art_title"] = "true"
			searchParams["series"] = "true"

			Expect(helper.DoNFTSearchWSReq(strings.TrimPrefix(it.WNBaseURI, "http://"), helper.GetNFTSearchURI(searchParams),
				[]string{"nft.txid", "nft.thumbnail_1", "nft.thumbnail_2", "nft.creator_name"},
				[]string{"reg_b", "dGh1bWJuYWlsLW9uZQ==", "dGh1bWJuYWlsLXR3bw==", "Mona Lisa"})).To(Succeed())
		})
	})

	Context("when searching an existing nft with matching art-title", func() {
		It("should be found & returned", func() {
			searchParams := make(map[string]string)
			searchParams["query"] = "art-b"
			searchParams["creator_name"] = "false"
			searchParams["art_title"] = "true"
			searchParams["series"] = "false"

			Expect(helper.DoNFTSearchWSReq(strings.TrimPrefix(it.WNBaseURI, "http://"), helper.GetNFTSearchURI(searchParams),
				[]string{"nft.txid", "nft.thumbnail_1", "nft.thumbnail_2", "nft.creator_name"},
				[]string{"reg_b", "dGh1bWJuYWlsLW9uZQ==", "dGh1bWJuYWlsLXR3bw==", "Mona Lisa"})).To(Succeed())
		})
	})

	Context("when searching an existing nft with matching artist-written statement", func() {
		It("should be found & returned", func() {
			searchParams := make(map[string]string)
			searchParams["query"] = "Lakes"
			searchParams["creator_name"] = "false"
			searchParams["descr"] = "true"
			searchParams["series"] = "false"

			Expect(helper.DoNFTSearchWSReq(strings.TrimPrefix(it.WNBaseURI, "http://"), helper.GetNFTSearchURI(searchParams),
				[]string{"nft.txid", "nft.thumbnail_1", "nft.thumbnail_2", "nft.creator_name"},
				[]string{"reg_a", "dGh1bWJuYWlsLW9uZQ==", "dGh1bWJuYWlsLXR3bw==", "Alan Majchrowicz"})).To(Succeed())
		})
	})

	Context("when searching an existing nft with matching string in series", func() {
		It("should be found & returned", func() {
			searchParams := make(map[string]string)
			searchParams["query"] = "nft-series"
			searchParams["series"] = "true"

			Expect(helper.DoNFTSearchWSReq(strings.TrimPrefix(it.WNBaseURI, "http://"), helper.GetNFTSearchURI(searchParams),
				[]string{"nft.series_name", "nft.thumbnail_1", "nft.thumbnail_2"},
				[]string{"nft-series", "dGh1bWJuYWlsLW9uZQ==", "dGh1bWJuYWlsLXR3bw=="})).To(Succeed())
		})
	})

	AfterEach(func() {
		Expect(mocker.CleanupAll()).To(Succeed())
	})
})
