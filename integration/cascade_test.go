package main_test

/*
import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	it "github.com/pastelnetwork/gonode/integration"
	"github.com/pastelnetwork/gonode/integration/fakes/common/testconst"
	helper "github.com/pastelnetwork/gonode/integration/helper"
	"github.com/pastelnetwork/gonode/integration/mock"
)

var _ = Describe("Cascade", func() {
	var (
		itHelper = helper.NewItHelper()

		uploadAssetReq   *helper.UploadImageReq
		uploadAssetReply *helper.UploadAssetResp

		startTaskReq *helper.SenseCascadeStartTaskReq
		regReply     *helper.RegistrationResp

		mocker *mock.Mocker
	)

	BeforeEach(func() {
		mocker = mock.New(it.PasteldServers, it.DDServers, it.SNServers, itHelper)

		uploadAssetReq = &helper.UploadImageReq{
			Filename: filepath.Join(filepath.Dir("."), "testdata", "test.jpg"),
		}

		uploadAssetReply = &helper.UploadAssetResp{}
		startTaskReq = &helper.SenseCascadeStartTaskReq{
			AppPastelID: testconst.ArtistPastelID,
			BurnTXID:    "896950d860eaf408e76a1a153deff80a7cda9e76291e5085060634e30b145c6a",
		}
		regReply = &helper.RegistrationResp{}
	})

	Context("when using cascade feature", func() {
		It("should get registered successfully", func() {
			// Register Mocks
			Expect(mocker.MockAllRegExpections()).To(Succeed())

			// Upload asset
			resp, err := itHelper.HTTPCurlUploadFile(helper.HttpPost, helper.GetCascadeUploadImageURI(it.WNBaseURI), uploadAssetReq.Filename, "test-img")
			Expect(err).NotTo(HaveOccurred())
			Expect(json.Unmarshal(resp, uploadAssetReply)).To(Succeed())
			Expect(uploadAssetReply.FileID).NotTo(BeEmpty())
			Expect(uploadAssetReply.EstimatedFee).NotTo(BeZero())

			// Start Task
			startResp, status, err := itHelper.Request(helper.HttpPost, startTaskReq, helper.GetCascadeStartTaskURI(it.WNBaseURI, uploadAssetReply.FileID), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))
			Expect(json.Unmarshal(startResp, regReply)).To(Succeed())

			// Check status for Task Completed
			Expect(helper.DoWebSocketReq(strings.TrimPrefix(it.WNBaseURI, "http://"), helper.GetCascadeTaskStateURI(regReply.TaskID),
				"status", "Task Completed")).To(Succeed())
		})
	})

	AfterEach(func() {
		Expect(mocker.CleanupAll()).To(Succeed())
	})
})
*/
