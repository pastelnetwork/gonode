package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pastelnetwork/gonode/common/types"
	it "github.com/pastelnetwork/gonode/integration"
	helper "github.com/pastelnetwork/gonode/integration/helper"
	"github.com/pastelnetwork/gonode/integration/mock"
)

var _ = Describe("StorageChallenge", func() {
	var (
		itHelper    = helper.NewItHelper()
		mocker      *mock.Mocker
		scResp      []types.StorageChallenge
		storeReq    *helper.StoreRequest
		rqfile      *helper.RawSymbolIDFile
		symbol      []byte
		idFile      []byte
		ChallengeID string
		sn3PastelID string
		sn4PastelID string
	)

	BeforeEach(func() {
		mocker = mock.New(it.SCPasteldServers, it.SCDDServers, it.SCSNServers, itHelper)
		scResp = []types.StorageChallenge{}
		storeReq = &helper.StoreRequest{}
		rqfile = &helper.RawSymbolIDFile{
			ID:       "09f6c459-ec2a-4db1-a8fe-0648fd97b5cb",
			PastelID: "jXY1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gW",
		}

		sn3PastelID = "jXXaQm7VuD4p8vG32XyZPq3Z8WBZWdpaRPM92nfZZvcW5pvjVf7RsEnB55qQ214WzuuXr7LoNJsZZKth3ArKAU"
		sn4PastelID = "jABegY26dJ1auzTj5xy1LSUHy3vEn91pMjQrMwAb97nVzH1859xBzbdiU1RsAhJiQ6Z8VbJwQeH5pwGbFAAaHc"
		symbol = []byte("test-symbol")
		storeReq = &helper.StoreRequest{}

	})

	Context("when all nodes have the valid file", func() {
		It("storage challenge should be successful", func() {
			Expect(mocker.MockSCRegExpections()).To(Succeed())

			id, err := helper.GetP2PID(symbol)
			Expect(err).NotTo(HaveOccurred())

			rqfile.SymbolIdentifiers = []string{id}
			idFile, err = rqfile.GetIDFile()
			Expect(err).NotTo(HaveOccurred())

			storeReq.Value = idFile
			_, status, err := itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN1BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			_, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN5BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			_, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN6BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			time.Sleep(25 * time.Second)
			// Verify that SN3 is challenger Node
			resp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetStorageChallengeURI(helper.SN3BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(resp, &scResp)).To(Succeed())

			for _, c := range scResp {
				if c.ChallengingNode == sn3PastelID && c.RespondingNode == sn4PastelID {
					ChallengeID = c.ChallengeID
				}
			}

			Expect(ChallengeID).NotTo(BeEmpty())

			// Verify that SN4 is processor Node
			resp, status, err = itHelper.Request(helper.HttpGet, nil, helper.GetStorageChallengeURI(helper.SN4BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(json.Unmarshal(resp, &scResp)).To(Succeed())

			for _, c := range scResp {
				if c.ChallengeID == ChallengeID {
					Expect(c.RespondingNode).
						To(Equal(sn4PastelID))
				}
			}

			verifications := 0
			for _, url := range []string{helper.SN1BaseURI, helper.SN2BaseURI, helper.SN3BaseURI, helper.SN5BaseURI, helper.SN6BaseURI, helper.SN7BaseURI} {

				// Verify that SN4 is processor Node
				resp, status, err = itHelper.Request(helper.HttpGet, nil, helper.GetStorageChallengeURI(url), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
				Expect(json.Unmarshal(resp, &scResp)).To(Succeed())

				for _, c := range scResp {
					if c.ChallengeID == ChallengeID {
						Expect(c.ChallengingNode).
							To(Equal(sn3PastelID))
						Expect(c.RespondingNode).
							To(Equal(sn4PastelID))

						if c.Status == types.VerifiedStorageChallengeStatus {
							verifications++
						}
					}
				}
			}

			Expect(verifications >= 5).To(BeTrue(), fmt.Sprintf(" nodes verified: %d", verifications))
		})
	})

	AfterEach(func() {
		Expect(mocker.CleanupAll()).To(Succeed())
	})
})
