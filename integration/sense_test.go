package main_test

import (
	"net/http"
	"path/filepath"
	"strings"

	json "github.com/json-iterator/go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	it "github.com/pastelnetwork/gonode/integration"
	"github.com/pastelnetwork/gonode/integration/fakes/common/testconst"
	helper "github.com/pastelnetwork/gonode/integration/helper"
	"github.com/pastelnetwork/gonode/integration/mock"
)

var _ = Describe("Sense", func() {
	var (
		itHelper = helper.NewItHelper()

		uploadImageReq   *helper.UploadImageReq
		uploadImageReply *helper.UploadImageResp

		startTaskReq *helper.SenseCascadeStartTaskReq
		regReply     *helper.RegistrationResp

		mocker *mock.Mocker
	)

	BeforeEach(func() {
		mocker = mock.New(it.PasteldServers, it.DDServers, it.SNServers, itHelper)

		uploadImageReq = &helper.UploadImageReq{
			Filename: filepath.Join(filepath.Dir("."), "testdata", "test.jpg"),
		}

		uploadImageReply = &helper.UploadImageResp{}
		startTaskReq = &helper.SenseCascadeStartTaskReq{
			AppPastelID: testconst.ArtistPastelID,
			BurnTXID:    "896950d860eaf408e76a1a153deff80a7cda9e76291e5085060634e30b145c6a",
		}
		regReply = &helper.RegistrationResp{}
	})

	Context("when using sense feature", func() {
		It("should be registered successfully", func() {
			// Register Mocks
			Expect(mocker.MockAllRegExpections()).To(Succeed())

			// Upload Image
			resp, err := itHelper.HTTPCurlUploadFile(helper.HttpPost, helper.GetSenseUploadImageURI(it.WNBaseURI), uploadImageReq.Filename, "test-img")
			Expect(err).NotTo(HaveOccurred())
			Expect(json.Unmarshal(resp, uploadImageReply)).To(Succeed())
			Expect(uploadImageReply.ImageID).NotTo(BeEmpty())
			Expect(uploadImageReply.EstimatedFee).NotTo(BeZero())

			// Start Task
			startResp, status, err := itHelper.Request(helper.HttpPost, startTaskReq, helper.GetSenseStartTaskURI(it.WNBaseURI, uploadImageReply.ImageID), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))
			Expect(json.Unmarshal(startResp, regReply)).To(Succeed())

			// Check status for Task Completed
			Expect(helper.DoWebSocketReq(strings.TrimPrefix(it.WNBaseURI, "http://"), helper.GetSenseTaskStateURI(regReply.TaskID),
				"status", "Task Completed")).To(Succeed())
		})
	})

	AfterEach(func() {
		Expect(mocker.CleanupAll()).To(Succeed())
	})
})
