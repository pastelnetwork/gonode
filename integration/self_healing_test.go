package main_test

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

var _ = Describe("SelfHealing", func() {
	var (
		itHelper               = helper.NewItHelper()
		mocker                 *mock.Mocker
		scResp                 []types.StorageChallenge
		shResp                 []types.SelfHealingChallenge
		storeReq               *helper.StoreRequest
		localStoreReq          *helper.LocalStoreRequest
		delReq                 *helper.DeleteRequest
		trueChallengerPastelID string
		trueProcessorPastelID  string

		p2pID    string
		symbols  map[string][]byte
		getReply *helper.RetrieveResponse
		idFile   []byte
	)

	BeforeEach(func() {
		mocker = mock.New(it.SCPasteldServers, it.SCDDServers, it.SCSNServers, itHelper)
		scResp = []types.StorageChallenge{}
		shResp = []types.SelfHealingChallenge{}
		storeReq = &helper.StoreRequest{}
		localStoreReq = &helper.LocalStoreRequest{}
		getReply = &helper.RetrieveResponse{}
		delReq = &helper.DeleteRequest{}

		trueChallengerPastelID = "jXXaQm7VuD4p8vG32XyZPq3Z8WBZWdpaRPM92nfZZvcW5pvjVf7RsEnB55qQ214WzuuXr7LoNJsZZKth3ArKAU" // SN3
		trueProcessorPastelID = "jXZX9sH6Fp5EUhqKQJteiivbNyGhsjPQKYXqVUVHRDTpH9UzGeEuG92c2S5tMTzUz1CLARVohpzjsWLUyZXSGN"  /// SN1

		storeReq = &helper.StoreRequest{}

		// generate RQSymbols
		sym, file, oti, err := itHelper.GetRaptorQSymbols()
		Expect(err).NotTo(HaveOccurred())
		symbols = sym

		// set RQ Input parameters for mocker
		mocker.SetOti(oti)

		// mock storage challenge
		Expect(mocker.MockSCRegExpections()).To(Succeed())

		idFile, err = file.GetIDFile()
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

		// store wrong data in node 2,4,5 and 6 so that storage challenge fails.
		localStoreReq.Value = []byte("wrong data")
		for _, url := range []string{helper.SN2BaseURI, helper.SN4BaseURI, helper.SN5BaseURI, helper.SN6BaseURI} {

			_, status, err = itHelper.Request(helper.HttpPost, localStoreReq, helper.GetLocalStoreURI(url), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
		}

	})

	Context("when storage challenge fails but all symbols are not available", func() {
		It("file should be reconstructed", func() {
			ChallengeID := []string{}
			// store some symbols in node 1,5 and 7
			i := 0
			for _, symbol := range symbols {
				storeReq.Value = symbol
				_, status, err := itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN1BaseURI), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))

				_, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN5BaseURI), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))

				_, status, err = itHelper.Request(helper.HttpPost, storeReq, helper.GetStoreURI(it.SN7BaseURI), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
				i++

				if i == 3 {
					break
				}
			}

			// Initiate Storage Challenge
			_, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetSCInitiateURI(it.SN3BaseURI), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			time.Sleep(2 * time.Second)

			// store the right data  back in node 2,4,5 and 6 so that self healing verification can be done
			localStoreReq.Key = p2pID
			localStoreReq.Value = idFile
			for _, url := range []string{helper.SN2BaseURI, helper.SN4BaseURI, helper.SN5BaseURI, helper.SN6BaseURI} {

				_, status, err = itHelper.Request(helper.HttpPost, localStoreReq, helper.GetLocalStoreURI(url), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
			}

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
						ChallengeID = append(ChallengeID, c.ChallengeID)
					}
				}

				if len(ChallengeID) == 0 {
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
					if in(c.ChallengeID, ChallengeID) {
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
						if in(c.ChallengeID, ChallengeID) {
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

			// Verify that SN7 is self-healing message checker and decides that self healing should be done
			Expect(backoff.Retry(backoff.Operation(func() error {
				resp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetSelfHealingChallengeURI(helper.SN7BaseURI), nil)
				if err != nil {
					return fmt.Errorf("err reaching server: - err: %w", err)
				}
				if status != http.StatusOK {
					return fmt.Errorf("status not OK - status: %d", status)
				}

				if err := json.Unmarshal(resp, &shResp); err != nil {
					return fmt.Errorf("err unmarshalling response: - err: %w", err)
				}

				success := false
				for _, c := range shResp {
					if in(c.ChallengeID, ChallengeID) && c.Status == types.InProgressSelfHealingStatus {
						success = true
					}
				}

				if !success {
					return fmt.Errorf("self healing challenge not found or not in correct status")
				}

				return nil
			}), b)).To(Succeed())

			// Verify that verifier nodes verify self healing message
			for _, url := range []string{helper.SN3BaseURI, helper.SN2BaseURI, helper.SN4BaseURI, helper.SN5BaseURI, helper.SN6BaseURI} {

				Expect(backoff.Retry(backoff.Operation(func() error {
					resp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetSelfHealingChallengeURI(url), nil)
					if err != nil {
						return fmt.Errorf("err reaching server: - err: %w", err)
					}
					if status != http.StatusOK {
						return fmt.Errorf("status not OK - status: %d", status)
					}

					if err := json.Unmarshal(resp, &shResp); err != nil {
						return fmt.Errorf("err unmarshalling response: - err: %w", err)
					}

					success := false
					for _, c := range shResp {
						if in(c.ChallengeID, ChallengeID) && c.Status == types.CompletedSelfHealingStatus {
							success = true
						}
					}

					if !success {
						return fmt.Errorf("self healing challenge not found or not in correct verification status")
					}

					return nil
				}), b)).To(Succeed())
			}

			for id, symbol := range symbols {
				getResp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetRetrieveURI(it.SN7BaseURI, id), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
				Expect(json.Unmarshal(getResp, getReply)).To(Succeed())
				Expect(getReply.Value).To(Equal(symbol))

			}
		})
	})

	Context("when storage challenge fails but all symbols are available", func() {
		It("file is not reconstructed", func() {
			ChallengeID := []string{}
			// store all symbols in node 1,5 and 7
			for _, symbol := range symbols {
				storeReq.Value = symbol
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
			}

			// Initiate Storage Challenge
			_, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetSCInitiateURI(it.SN3BaseURI), nil)
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
						ChallengeID = append(ChallengeID, c.ChallengeID)
					}
				}

				if len(ChallengeID) == 0 {
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
					if in(c.ChallengeID, ChallengeID) {
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
						if in(c.ChallengeID, ChallengeID) {
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

			// store the right data  back in node 2,4,5 and 6 so that self healing verification can be done
			localStoreReq.Key = p2pID
			localStoreReq.Value = idFile
			for _, url := range []string{helper.SN2BaseURI, helper.SN4BaseURI, helper.SN5BaseURI, helper.SN6BaseURI} {
				delReq.Key = p2pID
				_, status, err := itHelper.Request(helper.HttpPost, delReq, helper.GetRemoveURI(url), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))

				_, status, err = itHelper.Request(helper.HttpPost, localStoreReq, helper.GetLocalStoreURI(url), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
			}

			// Verify that SN7 is self-healing message checker and decides that reconstruction isn't required
			Expect(backoff.Retry(backoff.Operation(func() error {
				resp, status, err := itHelper.Request(helper.HttpGet, nil, helper.GetSelfHealingChallengeURI(helper.SN7BaseURI), nil)
				if err != nil {
					return fmt.Errorf("err reaching server: - err: %w", err)
				}
				if status != http.StatusOK {
					return fmt.Errorf("status not OK - status: %d", status)
				}

				if err := json.Unmarshal(resp, &shResp); err != nil {
					return fmt.Errorf("err unmarshalling response: - err: %w", err)
				}

				success := false
				for _, c := range shResp {
					if in(c.ChallengeID, ChallengeID) && c.Status == types.ReconstructionNotRequiredSelfHealingStatus {
						success = true
					}
				}

				if !success {
					return fmt.Errorf("self healing challenge not found or not in correct status")
				}

				return nil
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

			_, status, err = itHelper.Request(helper.HttpDelete, nil, helper.GetDeleteSelfHealingURI(url), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
		}

		for id := range symbols {
			delReq.Key = id
			for _, url := range []string{helper.SN1BaseURI, helper.SN2BaseURI, helper.SN3BaseURI, helper.SN4BaseURI, helper.SN5BaseURI, helper.SN6BaseURI, helper.SN7BaseURI} {

				_, status, err := itHelper.Request(helper.HttpPost, delReq, helper.GetRemoveURI(url), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
			}

		}
	})
})

func in(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}
