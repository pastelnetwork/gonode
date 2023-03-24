package main_test

/*
import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pastelnetwork/gonode/common/types"
	it "github.com/pastelnetwork/gonode/integration"
	helper "github.com/pastelnetwork/gonode/integration/helper"
	"github.com/pastelnetwork/gonode/integration/mock"
)

var _ = Describe("StorageChallenge", func() {
	var (
		itHelper               = helper.NewItHelper()
		mocker                 *mock.Mocker
		scResp                 []types.StorageChallenge
		storeReq               *helper.StoreRequest
		localStoreReq          *helper.LocalStoreRequest
		delReq                 *helper.DeleteRequest
		ChallengeID            string
		trueChallengerPastelID string
		trueProcessorPastelID  string
		p2pID                  string
	)

	BeforeEach(func() {
		mocker = mock.New(it.SCPasteldServers, it.SCDDServers, it.SCSNServers, itHelper)
		scResp = []types.StorageChallenge{}
		storeReq = &helper.StoreRequest{}
		delReq = &helper.DeleteRequest{}
		trueChallengerPastelID = "jXXaQm7VuD4p8vG32XyZPq3Z8WBZWdpaRPM92nfZZvcW5pvjVf7RsEnB55qQ214WzuuXr7LoNJsZZKth3ArKAU" // SN3
		trueProcessorPastelID = "jXZX9sH6Fp5EUhqKQJteiivbNyGhsjPQKYXqVUVHRDTpH9UzGeEuG92c2S5tMTzUz1CLARVohpzjsWLUyZXSGN"  /// SN1

		localStoreReq = &helper.LocalStoreRequest{}

	})

	Context("when all nodes have valid files", func() {
		It("should be successful", func() {
			// generate RQSymbols
			_, file, oti, err := itHelper.GetRaptorQSymbols()
			Expect(err).NotTo(HaveOccurred())
			// set RQ Input parameters for mocker
			mocker.SetOti(oti)
			// mock storage challenge
			Expect(mocker.MockSCRegExpections()).To(Succeed())

			// set file ID so that it ends up with same hash
			file.ID = "48d745e2-e1d3-4ba8-ab01-b636af38619a"

			idFile, err := file.GetIDFile()
			Expect(err).NotTo(HaveOccurred())

			// get p2p ID for the file
			p2pID, err = helper.GetP2PID(idFile)
			Expect(err).NotTo(HaveOccurred())

			storeReq.Value = idFile
			// store the correct file in node 1,3 and 7
			_, status, err := itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN1BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			_, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN5BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			_, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN6BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			_, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN7BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			// Initiate Storage Challenge
			_, status, err = itHelper.Request(helper.HttpGet, nil, helper.GetSCInitiateURI(it.SN3BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			b := backoff.NewExponentialBackOff()
			b.MaxElapsedTime = 30 * time.Second
			b.MaxInterval = 5 * time.Second

			// Verify that trueChallenger (SN3) is challenger Node
			Expect(backoff.Retry(backoff.Operation(func() error {
				resp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetStorageChallengeURI(helper.SN3BaseURI), nil)
				if err != nil {
					return fmt.Errorf("err reaching server: - err: %w", err)
				}
				if status != http.StatusOK {
					return errors.New("status is not OK")
				}

				if err := json.Unmarshal(resp, &scResp); err != nil {
					return fmt.Errorf("err unmarshalling response: - err: %w", err)
				}

				for _, c := range scResp {
					if c.ChallengingNode == trueChallengerPastelID && c.RespondingNode == trueProcessorPastelID {
						ChallengeID = c.ChallengeID
					}
				}

				if ChallengeID == "" {
					return errors.New("challenge ID is empty")
				}

				return nil
			}), b)).To(Succeed())

			// Verify that trueProcessor (SN1) is responding Node
			Expect(backoff.Retry(backoff.Operation(func() error {
				resp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetStorageChallengeURI(helper.SN1BaseURI), nil)
				if err != nil {
					return fmt.Errorf("err reaching server: - err: %w", err)
				}
				if status != http.StatusOK {
					return fmt.Errorf("status not OK - status: %d", status)
				}

				if err := json.Unmarshal(resp, &scResp); err != nil {
					return fmt.Errorf("err unmarshalling response: - err: %w", err)
				}

				for _, c := range scResp {
					if c.ChallengeID == ChallengeID {
						if c.RespondingNode != trueProcessorPastelID {
							return fmt.Errorf("responding node is not true processor: %s", c.RespondingNode)
						}
					}
				}

				return nil
			}), b)).To(Succeed())

			Expect(backoff.Retry(backoff.Operation(func() error {

				verifications := 0
				for _, url := range []string{helper.SN2BaseURI, helper.SN3BaseURI, helper.SN4BaseURI, helper.SN5BaseURI, helper.SN6BaseURI, helper.SN7BaseURI} {

					resp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetStorageChallengeURI(url), nil)
					if err != nil {
						return fmt.Errorf("err reaching server: - err: %w", err)
					}
					if status != http.StatusOK {
						return fmt.Errorf("status not OK - status: %d", status)
					}

					if err := json.Unmarshal(resp, &scResp); err != nil {
						return fmt.Errorf("err unmarshalling response: - err: %w", err)
					}

					for _, c := range scResp {
						if c.ChallengeID == ChallengeID {
							if c.ChallengingNode != trueChallengerPastelID || c.RespondingNode != trueProcessorPastelID {
								return fmt.Errorf("challenging node is not true challenger: %s", c.ChallengingNode)
							}

							if c.Status == types.VerifiedStorageChallengeStatus {
								verifications++
							}
						}
					}
				}

				if verifications >= 5 {
					return nil
				}

				return fmt.Errorf(" unexpected verifications: %d", verifications)
			}), b)).To(Succeed())

		})
	})

	Context("when all nodes don't have valid files", func() {
		It("storage challenge should fail", func() {
			// generate RQSymbols
			_, file, oti, err := itHelper.GetRaptorQSymbols()
			Expect(err).NotTo(HaveOccurred())
			// set RQ Input parameters for mocker
			mocker.SetOti(oti)
			// mock storage challenge
			Expect(mocker.MockSCRegExpections()).To(Succeed())

			// set file ID so that it ends up with same hash
			file.ID = "48d745e2-e1d3-4ba8-ab01-b636af38619a"

			idFile, err := file.GetIDFile()
			Expect(err).NotTo(HaveOccurred())

			// get p2p ID for the file
			p2pID, err = helper.GetP2PID(idFile)
			Expect(err).NotTo(HaveOccurred())

			localStoreReq.Value = idFile
			localStoreReq.Key = p2pID

			// store the correct file in node 1,3 and 7
			_, status, err := itHelper.Request(helper.HttpPost, localStoreReq, helper.GetLocalStoreURI(it.SN1BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			_, status, err = itHelper.Request(helper.HttpPost, localStoreReq, helper.GetLocalStoreURI(it.SN3BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			_, status, err = itHelper.Request(helper.HttpPost, localStoreReq, helper.GetLocalStoreURI(it.SN7BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			// Initiate Storage Challenge
			_, status, err = itHelper.Request(helper.HttpGet, nil, helper.GetSCInitiateURI(it.SN3BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			b := backoff.NewExponentialBackOff()
			b.MaxElapsedTime = 30 * time.Second
			b.MaxInterval = 5 * time.Second

			// store wrong data in node 2,4,5 and 6 so that storage challenge fails.
			localStoreReq.Value = []byte("wrong data")
			for _, url := range []string{helper.SN2BaseURI, helper.SN4BaseURI, helper.SN5BaseURI, helper.SN6BaseURI} {

				_, status, err = itHelper.Request(helper.HttpPost, localStoreReq, helper.GetLocalStoreURI(url), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
			}

			// Verify that trueChallenger (SN3) is challenger Node
			Expect(backoff.Retry(backoff.Operation(func() error {
				resp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetStorageChallengeURI(helper.SN3BaseURI), nil)
				if err != nil {
					return fmt.Errorf("err reaching server: - err: %w", err)
				}
				if status != http.StatusOK {
					return errors.New("status is not OK")
				}

				if err := json.Unmarshal(resp, &scResp); err != nil {
					return fmt.Errorf("err unmarshalling response: - err: %w", err)
				}

				for _, c := range scResp {
					if c.ChallengingNode == trueChallengerPastelID && c.RespondingNode == trueProcessorPastelID {
						ChallengeID = c.ChallengeID
					}
				}

				if ChallengeID == "" {
					return errors.New("challenge ID is empty")
				}

				return nil
			}), b)).To(Succeed())

			// Verify that trueProcessor (SN1) is responding Node
			Expect(backoff.Retry(backoff.Operation(func() error {
				resp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetStorageChallengeURI(helper.SN1BaseURI), nil)
				if err != nil {
					return fmt.Errorf("err reaching server: - err: %w", err)
				}
				if status != http.StatusOK {
					return fmt.Errorf("status not OK - status: %d", status)
				}

				if err := json.Unmarshal(resp, &scResp); err != nil {
					return fmt.Errorf("err unmarshalling response: - err: %w", err)
				}

				for _, c := range scResp {
					if c.ChallengeID == ChallengeID {
						if c.RespondingNode != trueProcessorPastelID {
							return fmt.Errorf("responding node is not true processor: %s", c.RespondingNode)
						}
					}
				}

				return nil
			}), b)).To(Succeed())

			Expect(backoff.Retry(backoff.Operation(func() error {

				verifications := 0
				for _, url := range []string{helper.SN2BaseURI, helper.SN3BaseURI, helper.SN4BaseURI, helper.SN5BaseURI, helper.SN6BaseURI, helper.SN7BaseURI} {

					resp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetStorageChallengeURI(url), nil)
					if err != nil {
						return fmt.Errorf("err reaching server: - err: %w", err)
					}
					if status != http.StatusOK {
						return fmt.Errorf("status not OK - status: %d", status)
					}

					if err := json.Unmarshal(resp, &scResp); err != nil {
						return fmt.Errorf("err unmarshalling response: - err: %w", err)
					}

					for _, c := range scResp {
						if c.ChallengeID == ChallengeID {
							if c.ChallengingNode != trueChallengerPastelID || c.RespondingNode != trueProcessorPastelID {
								return fmt.Errorf("challenging node is not true challenger: %s", c.ChallengingNode)
							}

							if c.Status == types.VerifiedStorageChallengeStatus {
								verifications++
							}
						}
					}
				}

				if verifications < 3 {
					return nil
				}

				return fmt.Errorf(" unexpected verifications: %d", verifications)

			}), b)).To(Succeed())
		})
	})

	AfterEach(func() {
		Expect(mocker.CleanupAll()).To(Succeed())
		delReq.Key = p2pID
		for _, url := range []string{helper.SN1BaseURI, helper.SN2BaseURI, helper.SN3BaseURI, helper.SN4BaseURI, helper.SN5BaseURI, helper.SN6BaseURI, helper.SN7BaseURI} {

			_, status, err := itHelper.Request(helper.HttpPost, delReq, helper.GetRemoveURI(url), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			_, status, err = itHelper.Request(helper.HttpGet, nil, helper.GetDeleteStorageChallengeURI(url), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
		}

	})
})
*/
