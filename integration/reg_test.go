package integration_test

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	it "github.com/pastelnetwork/gonode/integration"
	"github.com/pastelnetwork/gonode/integration/nft"
)

var _ = Describe("NFTRegistration", func() {
	var (
		helper           = it.NewItHelper()
		uploadImageReq   *nft.UploadImageReq
		uploadImageReply *nft.UploadImageResp
		regReq           *nft.RegistrationReq
		regReply         *nft.RegistrationResp
	)

	BeforeEach(func() {
		uploadImageReq = &nft.UploadImageReq{
			Filename: "/nft/test.jpg",
		}

		uploadImageReply = &nft.UploadImageResp{}
		regReq = &nft.RegistrationReq{
			ArtistPastelid:           "jXY1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gW",
			ArtistPastelidPassphrase: "passphrase",
			ArtistWebsiteURL:         "www.example.com",
			ArtistName:               "integration-test",
			Description:              "sample-description",
			MaximumFee:               500,
			IssuesCopies:             5,
			Keywords:                 "key,word",
			Name:                     "test-nft",
			SpendableAddress:         "tPoXb35ABswxqMCw3gD8K7PvzRTcZLdpNEY",
			Royalty:                  10,
			ThumbnailCoordinate: nft.ThumbnailCoordinate{
				BottomRightX: 640,
				BottomRightY: 480,
				TopLeftX:     0,
				TopLeftY:     480,
			},
			YoutubeURL: "www.youtube.com",
		}
		regReply = &nft.RegistrationResp{}
	})

	Context("when registering a valid nft", func() {
		It("should registered successfully", func() {
			resp, status, err := helper.UploadFile(uploadImageReq.Filename, nft.GetUploadImageURI(it.WNBaseURI))
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))
			Expect(json.Unmarshal(resp, uploadImageReply)).To(Succeed())

			regReq.ImageID = uploadImageReply.ImageID
			regResp, status, err := helper.Request(it.HttpPost, regReq, nft.GetRegistrationURI(it.WNBaseURI))
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))
			Expect(json.Unmarshal(regResp, regReply)).To(Succeed())
		})
	})

	AfterEach(func() {})
})
