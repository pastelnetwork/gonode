package main_test

/*
import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/pastelnetwork/gonode/integration/fakes/common/testconst"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	it "github.com/pastelnetwork/gonode/integration"
	helper "github.com/pastelnetwork/gonode/integration/helper"
	"github.com/pastelnetwork/gonode/integration/mock"
)

var _ = Describe("NFTRegistration", func() {
	var (
		itHelper         = helper.NewItHelper()
		uploadImageReq   *helper.UploadImageReq
		uploadImageReply *helper.UploadImageResp
		regReq           *helper.RegistrationReq
		regReply         *helper.RegistrationResp
		mocker           *mock.Mocker
	)

	BeforeEach(func() {
		mocker = mock.New(it.PasteldServers, it.DDServers, it.SNServers, itHelper)

		uploadImageReq = &helper.UploadImageReq{
			Filename: filepath.Join(filepath.Dir("."), "testdata", "test.jpg"),
		}

		uploadImageReply = &helper.UploadImageResp{}
		regReq = &helper.RegistrationReq{
			CreatorPastelid:           testconst.ArtistPastelID,
			CreatorPastelidPassphrase: "passphrase",
			CreatorWebsiteURL:         "www.example.com",
			CreatorName:               "integration-test",
			Description:               "sample-description",
			MaximumFee:                500,
			IssuedCopies:              5,
			Keywords:                  "key,word",
			Name:                      "test-nft",
			SpendableAddress:          testconst.RegSpendableAddress,
			Royalty:                   10,
			ThumbnailCoordinate: helper.ThumbnailCoordinate{
				BottomRightX: 640,
				BottomRightY: 480,
				TopLeftX:     0,
				TopLeftY:     0,
			},
			YoutubeURL: "www.youtube.com",
		}
		regReply = &helper.RegistrationResp{}
	})

	Context("when registering a valid nft", func() {
		It("should registered successfully", func() {
			Expect(mocker.MockAllRegExpections()).To(Succeed())
			resp, err := itHelper.HTTPCurlUploadFile(helper.HttpPost, helper.GetUploadImageURI(it.WNBaseURI), uploadImageReq.Filename, "test-img")
			Expect(err).NotTo(HaveOccurred())

			Expect(json.Unmarshal(resp, uploadImageReply)).To(Succeed())
			Expect(uploadImageReply.ImageID).NotTo(BeEmpty())

			regReq.ImageID = uploadImageReply.ImageID
			regResp, status, err := itHelper.Request(helper.HttpPost, regReq, helper.GetRegistrationURI(it.WNBaseURI), nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(status).To(Equal(http.StatusCreated))
			Expect(json.Unmarshal(regResp, regReply)).To(Succeed())

			Expect(helper.DoWebSocketReq(strings.TrimPrefix(it.WNBaseURI, "http://"), helper.GetTaskStatePath(regReply.TaskID),
				"status", "Task Completed")).To(Succeed())
		})
	})

	AfterEach(func() {
		Expect(mocker.CleanupAll()).To(Succeed())
	})
})
*/
