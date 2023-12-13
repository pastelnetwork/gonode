package selfhealing

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
)

// VerifySelfHealingChallenge verifies the self-healing challenge
func (task *SHTask) VerifySelfHealingChallenge(ctx context.Context, incomingResponseMessage types.SelfHealingMessage) (*types.SelfHealingMessage, error) {
	log.WithContext(ctx).WithField("challenge_id", incomingResponseMessage.ChallengeID).Info("VerifySelfHealingChallenge has been invoked")

	// incoming challenge message validation
	if err := task.validateSelfHealingResponseIncomingData(ctx, incomingResponseMessage); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error validating self-healing challenge incoming data: ")
		return nil, err
	}

	var (
		nftTicket     *pastel.NFTTicket
		cascadeTicket *pastel.APICascadeTicket
		senseTicket   *pastel.APISenseTicket
	)

	log.WithContext(ctx).Info("retrieving block no and verbose")
	currentBlockCount, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block count")
		return nil, err
	}
	blkVerbose1, err := task.SuperNodeService.PastelClient.GetBlockVerbose1(ctx, currentBlockCount)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block verbose 1")
		return nil, err
	}
	merkleroot := blkVerbose1.MerkleRoot

	verificationMsg := &types.SelfHealingMessage{
		ChallengeID: incomingResponseMessage.ChallengeID,
		SenderID:    task.nodeID,
		MessageType: types.SelfHealingResponseMessage,
		SelfHealingMessageData: types.SelfHealingMessageData{
			ChallengerID: incomingResponseMessage.SelfHealingMessageData.ChallengerID,
			Challenge: types.SelfHealingChallengeData{
				Block:            incomingResponseMessage.SelfHealingMessageData.Challenge.Block,
				Merkelroot:       incomingResponseMessage.SelfHealingMessageData.Challenge.Merkelroot,
				ChallengeTickets: incomingResponseMessage.SelfHealingMessageData.Challenge.ChallengeTickets,
				Timestamp:        incomingResponseMessage.SelfHealingMessageData.Challenge.Timestamp,
			},
			Response: types.SelfHealingResponseData{
				Block:            incomingResponseMessage.SelfHealingMessageData.Response.Block,
				Merkelroot:       incomingResponseMessage.SelfHealingMessageData.Response.Merkelroot,
				Timestamp:        incomingResponseMessage.SelfHealingMessageData.Response.Timestamp,
				RespondedTickets: incomingResponseMessage.SelfHealingMessageData.Response.RespondedTickets,
			},
			Verification: types.SelfHealingVerificationData{
				Block:      currentBlockCount,
				Merkelroot: merkleroot,
			},
		},
	}

	var reconstructedFileHash []byte
	respondedTickets := incomingResponseMessage.SelfHealingMessageData.Response.RespondedTickets

	if respondedTickets == nil {
		log.WithContext(ctx).Info("no tickets have been responded by the recipient")
		return nil, nil
	}

	for _, ticket := range respondedTickets {
		nftTicket, cascadeTicket, senseTicket, err = task.getTicket(ctx, ticket.TxID, TicketType(ticket.TicketType))
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error getRelevantTicketFromFileHash")
		}
		log.WithContext(ctx).WithField("challenge_id", incomingResponseMessage.ChallengeID).Info("reg ticket has been retrieved")

		if cascadeTicket != nil || nftTicket != nil {
			isReconstructionReq := task.isReconstructionRequired(ctx, ticket.MissingKeys)
			{
				log.WithContext(ctx).WithField("txid", ticket.TxID).Info("Reconstruction is not required for ticket")
				continue
			}

			if !isReconstructionReq {
				log.WithContext(ctx).WithError(err).WithField("failed_challenge_id", incomingResponseMessage.ChallengeID).Error("Error in checking process")
				//responseMessage.ChallengeStatus = pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE
				//shChallenge.Status = types.FailedSelfHealingStatus
				//storeLogs(ctx, store, shChallenge)
				//return responseMessage, err
			}

			_, reconstructedFileHash, err = task.selfHealing(ctx, ticket.TxID, nftTicket, cascadeTicket)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("Error self-healing the file")
			}

			log.WithContext(ctx).WithField("failed_challenge_id", incomingResponseMessage.ChallengeID).Info("Comparing hashes")
			if !bytes.Equal(reconstructedFileHash, ticket.ReconstructedFileHash) {
				log.WithContext(ctx).WithField("txid", ticket.TxID).Info("reconstructed file hash does not match with the verifier reconstructed file")

				//responseMessage.ChallengeStatus = pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE
				//shChallenge.Status = types.FailedSelfHealingStatus
				//shChallenge.ReconstructedFileHash = reconstructedFileHash
				//storeLogs(ctx, store, shChallenge)
				//return responseMessage, nil
			}
			log.WithContext(ctx).WithField("txid", ticket.TxID).Info("reconstructed file hash matched")

		} else if senseTicket != nil {
			reqSelfHealing, mostCommonFile := task.senseCheckingProcess(ctx, senseTicket.DDAndFingerprintsIDs)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("Error in checking process for sense action ticket")
			}

			if !reqSelfHealing {
				log.WithContext(ctx).WithError(err).Error("self-healing not required for sense action ticket")
				//responseMessage.ChallengeStatus = pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE
				//shChallenge.Status = types.FailedSelfHealingStatus
				//shChallenge.ReconstructedFileHash = nil
				//storeLogs(ctx, store, shChallenge)
				//
				//return responseMessage, nil
			}

			ids, _, err := task.senseSelfHealing(ctx, senseTicket, mostCommonFile)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("error while self-healing sense action ticket")
				//responseMessage.ChallengeStatus = pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE
				//shChallenge.Status = types.FailedSelfHealingStatus
				//shChallenge.ReconstructedFileHash = nil
				//storeLogs(ctx, store, shChallenge)
				//
				//return responseMessage, err
			}

			if ok := compareFileIDs(ids, ticket.FileIDs); !ok {
				log.WithContext(ctx).WithField("challenge_id", incomingResponseMessage.ChallengeID).Info("Failed")
				//responseMessage.ChallengeStatus = pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE
				//shChallenge.Status = types.FailedSelfHealingStatus
				//shChallenge.ReconstructedFileHash = nil
				//storeLogs(ctx, store, shChallenge)
				//
				//return responseMessage, nil
			}
		}
	}

	//responseMessage.ChallengeStatus = pb.SelfHealingData_Status_SUCCEEDED
	//shChallenge.Status = types.CompletedSelfHealingStatus
	//shChallenge.ReconstructedFileHash = reconstructedFileHash
	//storeLogs(ctx, store, shChallenge)
	return verificationMsg, nil
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

func (task *SHTask) validateSelfHealingResponseIncomingData(ctx context.Context, incomingChallengeMessage types.SelfHealingMessage) error {
	if incomingChallengeMessage.MessageType != types.SelfHealingResponseMessage {
		return fmt.Errorf("incorrect message type to processing self-healing challenge")
	}

	isVerified, err := task.VerifyMessageSignature(ctx, incomingChallengeMessage)
	if err != nil {
		return errors.Errorf("error verifying sender's signature: %w", err)
	}

	if !isVerified {
		return errors.Errorf("not able to verify message signature")
	}

	return nil
}
