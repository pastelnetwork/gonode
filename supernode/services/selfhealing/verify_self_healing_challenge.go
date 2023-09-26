package selfhealing

import (
	"bytes"
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
)

// VerifySelfHealingChallenge verifies the self-healing challenge
func (task *SHTask) VerifySelfHealingChallenge(ctx context.Context, challengeMessage *pb.SelfHealingData) (*pb.SelfHealingData, error) {
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("VerifySelfHealingChallenge has been invoked")

	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("retrieving reg ticket")

	var (
		regTicket     pastel.RegTicket
		cascadeTicket *pastel.APICascadeTicket
		senseTicket   *pastel.APISenseTicket
		actionTicket  pastel.ActionRegTicket
		err           error
	)

	if challengeMessage.RegTicketId != "" {
		regTicket, err = task.getRegTicket(ctx, challengeMessage.RegTicketId)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("unable to retrieve reg ticket")
		}

		log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("reg ticket has been retrieved")
	} else if challengeMessage.ActionTicketId != "" {
		actionTicket, err = task.getActionTicket(ctx, challengeMessage.ActionTicketId)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("unable to retrieve cascade ticket")
			return nil, err
		}

		if challengeMessage.IsSenseTicket {
			senseTicket, err = actionTicket.ActionTicketData.ActionTicketData.APISenseTicket()
			if err != nil {
				log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTicket.TXID).
					Warnf("Could not get sense ticket for action ticket data verify self healing")
				return nil, err
			}
		} else {
			cascadeTicket, err = actionTicket.ActionTicketData.ActionTicketData.APICascadeTicket()
			if err != nil {
				log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTicket.TXID).
					Warnf("Could not get cascade ticket for action ticket data verify self healing")
				return nil, err
			}
		}

		log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("cascade reg ticket has been retrieved")
	} else {
		log.WithContext(ctx).Info("unable to find reg ticket")
		return nil, errors.New("reg ticket not found")
	}

	log.WithContext(ctx).Info("establishing connection with rq service")
	var rqConnection rqnode.Connection
	rqConnection, err = task.StorageHandler.RqClient.Connect(ctx, task.config.RaptorQServiceAddress)
	if err != nil {
		log.WithContext(ctx).Error("Error establishing RQ connection")
	}
	defer rqConnection.Close()

	rqNodeConfig := &rqnode.Config{
		RqFilesDir: task.config.RqFilesDir,
	}
	rqService := rqConnection.RaptorQ(rqNodeConfig)
	log.WithContext(ctx).Info("connection established with rq service")

	challengeFileHash := challengeMessage.ChallengeFile.FileHashToChallenge
	log.WithContext(ctx).Info("retrieving reg ticket")

	responseMessage := &pb.SelfHealingData{
		MessageId:                   challengeMessage.MessageId,
		MessageType:                 pb.SelfHealingData_MessageType_SELF_HEALING_RESPONSE_MESSAGE,
		MerklerootWhenChallengeSent: challengeMessage.MerklerootWhenChallengeSent,
		ChallengingMasternodeId:     task.nodeID,
		RespondingMasternodeId:      challengeMessage.ChallengingMasternodeId,
		ChallengeFile: &pb.SelfHealingDataChallengeFile{
			FileHashToChallenge: challengeMessage.ChallengeFile.FileHashToChallenge,
		},
		ChallengeId:    challengeMessage.ChallengeId,
		RegTicketId:    challengeMessage.RegTicketId,
		ActionTicketId: challengeMessage.ActionTicketId,
		SenseFileIds:   challengeMessage.SenseFileIds,
		IsSenseTicket:  challengeMessage.IsSenseTicket,
	}

	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	shChallenge := types.SelfHealingChallenge{
		ChallengeID:     responseMessage.ChallengeId,
		MerkleRoot:      responseMessage.MerklerootWhenChallengeSent,
		FileHash:        responseMessage.ChallengeFile.FileHashToChallenge,
		ChallengingNode: responseMessage.ChallengingMasternodeId,
		RespondingNode:  responseMessage.RespondingMasternodeId,
		VerifyingNode:   task.nodeID,
	}

	//Checking Process
	//1. false, nil, err    - should not update the challenge to completed, so that it can be retried again
	//2. false, nil, nil    - reconstruction not required
	//3. true, symbols, nil - reconstruction required

	var reconstructedFileHash []byte
	if cascadeTicket != nil || regTicket.TXID != "" {
		isReconstructionReq, availableSymbols, err := task.checkingProcess(ctx, challengeFileHash)
		if err != nil && !isReconstructionReq {
			log.WithContext(ctx).WithError(err).WithField("failed_challenge_id", challengeMessage.ChallengeId).Error("Error in checking process")
			responseMessage.ChallengeStatus = pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE
			shChallenge.Status = types.FailedSelfHealingStatus
			storeLogs(ctx, store, shChallenge)
			return responseMessage, err
		}

		if !isReconstructionReq && availableSymbols == nil {
			log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info(fmt.Sprintf("Reconstruction is not required for file: %s", challengeFileHash))
			responseMessage.ChallengeStatus = pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE
			shChallenge.Status = types.FailedSelfHealingStatus
			storeLogs(ctx, store, shChallenge)
			return responseMessage, nil
		}

		_, reconstructedFileHash, err = task.selfHealing(ctx, rqService, challengeMessage, &regTicket, cascadeTicket, availableSymbols)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error self-healing the file")
			return responseMessage, err
		}
		log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("File has been reconstructed")

		log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("Comparing hashes")
		if !bytes.Equal(reconstructedFileHash, challengeMessage.ReconstructedFileHash) {
			log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("reconstructed file hash does not match with the verifier reconstructed file")

			responseMessage.ChallengeStatus = pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE
			shChallenge.Status = types.FailedSelfHealingStatus
			shChallenge.ReconstructedFileHash = reconstructedFileHash
			storeLogs(ctx, store, shChallenge)
			return responseMessage, nil
		}
		log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("hashes have been matched of the responder and verifiers")
	} else if senseTicket != nil {
		reqSelfHealing, mostCommonFile := task.senseCheckingProcess(ctx, senseTicket.DDAndFingerprintsIDs)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error in checking process for sense action ticket")
		}

		if !reqSelfHealing {
			log.WithContext(ctx).WithError(err).Error("self-healing not required for sense action ticket")
			responseMessage.ChallengeStatus = pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE
			shChallenge.Status = types.FailedSelfHealingStatus
			shChallenge.ReconstructedFileHash = nil
			storeLogs(ctx, store, shChallenge)

			return responseMessage, nil
		}

		ids, _, err := task.senseSelfHealing(ctx, senseTicket, mostCommonFile)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error while self-healing sense action ticket")
			responseMessage.ChallengeStatus = pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE
			shChallenge.Status = types.FailedSelfHealingStatus
			shChallenge.ReconstructedFileHash = nil
			storeLogs(ctx, store, shChallenge)

			return responseMessage, err
		}

		if ok := compareFileIDs(ids, challengeMessage.SenseFileIds); !ok {
			responseMessage.ChallengeStatus = pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE
			shChallenge.Status = types.FailedSelfHealingStatus
			shChallenge.ReconstructedFileHash = nil
			storeLogs(ctx, store, shChallenge)

			return responseMessage, nil
		}
	}

	responseMessage.ChallengeStatus = pb.SelfHealingData_Status_SUCCEEDED
	shChallenge.Status = types.CompletedSelfHealingStatus
	shChallenge.ReconstructedFileHash = reconstructedFileHash
	storeLogs(ctx, store, shChallenge)
	return responseMessage, nil
}

func storeLogs(ctx context.Context, store storage.LocalStoreInterface, msg types.SelfHealingChallenge) {
	log.WithContext(ctx).Println("Storing challenge to DB for self healing inspection")

	if store != nil {
		_, err := store.InsertSelfHealingChallenge(msg)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error storing challenge to DB")
		}
	}
}

func compareFileIDs(verifyingFileIDs []string, processedFileIDs []string) bool {
	verifyingFileIDsMap := make(map[string]bool)
	for _, id := range verifyingFileIDs {
		verifyingFileIDsMap[id] = false
	}

	for _, id := range processedFileIDs {
		verifyingFileIDsMap[id] = true
	}

	for _, IsMatched := range verifyingFileIDsMap {
		if !IsMatched {
			return false
		}
	}

	return true //all file ids matched
}
