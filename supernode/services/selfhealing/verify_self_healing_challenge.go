package selfhealing

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

// VerifySelfHealingChallenge verifies the self-healing challenge
func (task *SHTask) VerifySelfHealingChallenge(ctx context.Context, incomingResponseMessage types.SelfHealingMessage) (*pb.SelfHealingMessage, error) {
	logger := log.WithContext(ctx).WithField("trigger_id", incomingResponseMessage.TriggerID).
		WithField("challenge_id", incomingResponseMessage.SelfHealingMessageData.Response.ChallengeID)

	logger.Info("VerifySelfHealingChallenge has been invoked")

	// incoming challenge message validation
	if err := task.validateSelfHealingResponseIncomingData(ctx, incomingResponseMessage); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error validating self-healing challenge incoming data: ")
		return nil, err
	}
	logger.Info("self healing response message has been validated")

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
	logger.Info("block count & merkelroot has been retrieved")

	verificationMsg := &types.SelfHealingMessage{
		TriggerID:   incomingResponseMessage.TriggerID,
		SenderID:    task.nodeID,
		MessageType: types.SelfHealingVerificationMessage,
		SelfHealingMessageData: types.SelfHealingMessageData{
			ChallengerID: incomingResponseMessage.SelfHealingMessageData.ChallengerID,
			Challenge: types.SelfHealingChallengeData{
				Block:            incomingResponseMessage.SelfHealingMessageData.Challenge.Block,
				Merkelroot:       incomingResponseMessage.SelfHealingMessageData.Challenge.Merkelroot,
				ChallengeTickets: incomingResponseMessage.SelfHealingMessageData.Challenge.ChallengeTickets,
				Timestamp:        incomingResponseMessage.SelfHealingMessageData.Challenge.Timestamp,
				NodesOnWatchlist: incomingResponseMessage.SelfHealingMessageData.Challenge.NodesOnWatchlist,
			},
			Response: types.SelfHealingResponseData{
				ChallengeID:     incomingResponseMessage.SelfHealingMessageData.Response.ChallengeID,
				Block:           incomingResponseMessage.SelfHealingMessageData.Response.Block,
				Merkelroot:      incomingResponseMessage.SelfHealingMessageData.Response.Merkelroot,
				Timestamp:       incomingResponseMessage.SelfHealingMessageData.Response.Timestamp,
				RespondedTicket: incomingResponseMessage.SelfHealingMessageData.Response.RespondedTicket,
			},
			Verification: types.SelfHealingVerificationData{
				ChallengeID: incomingResponseMessage.SelfHealingMessageData.Response.ChallengeID,
				Block:       currentBlockCount,
				Merkelroot:  merkleroot,
			},
		},
	}
	ticket := incomingResponseMessage.SelfHealingMessageData.Response.RespondedTicket

	if ticket.TxID == "" {
		return nil, nil
	}

	logger.WithField("ticket_txid", ticket.TxID).Info("starting self-healing verification for the ticket")
	nftTicket, cascadeTicket, senseTicket, err = task.getTicket(ctx, ticket.TxID, TicketType(ticket.TicketType))
	if err != nil {
		logger.WithError(err).Error("Error getRelevantTicketFromMsg")
	}
	logger.Info("reg ticket has been retrieved for verification")

	if cascadeTicket != nil || nftTicket != nil {
		isReconstructionReq := task.isReconstructionRequired(ctx, ticket.MissingKeys)
		if !isReconstructionReq {
			logger.WithField("ticket_txid", ticket.TxID).Error("reconstruction is not required for ticket")

			if !ticket.IsReconstructionRequired {
				logger.WithField("ticket_txid", ticket.TxID).
					Info("is reconstruction required set to false by both recipient and verifier")

				verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
					TxID:                     ticket.TxID,
					TicketType:               ticket.TicketType,
					MissingKeys:              ticket.MissingKeys,
					ReconstructedFileHash:    nil,
					IsReconstructionRequired: false,
					IsVerified:               true,
					Message:                  "is reconstruction required set to false by both recipient and verifier",
				}
			} else {
				logger.WithField("ticket_txid", ticket.TxID).
					Info("is reconstruction required set to true by recipient but false by verifier")

				verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
					TxID:                     ticket.TxID,
					TicketType:               ticket.TicketType,
					MissingKeys:              ticket.MissingKeys,
					ReconstructedFileHash:    nil,
					IsReconstructionRequired: false,
					IsVerified:               false,
					Message:                  "is reconstruction required set to true by recipient but false by verifier",
				}
			}

			return task.prepareAndSendVerificationMsg(ctx, *verificationMsg)
		}

		_, reconstructedFileHash, err := task.selfHealing(ctx, ticket.TxID, nftTicket, cascadeTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error self-healing the file")
		}

		logger.Info("comparing data hashes")
		if !bytes.Equal(reconstructedFileHash, ticket.ReconstructedFileHash) {
			logger.WithContext(ctx).WithField("txid", ticket.TxID).Info("reconstructed file hash does not match with the verifier reconstructed file")
			verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
				TxID:                     ticket.TxID,
				TicketType:               ticket.TicketType,
				MissingKeys:              ticket.MissingKeys,
				ReconstructedFileHash:    nil,
				IsReconstructionRequired: false,
				IsVerified:               false,
				Message:                  "reconstructed file hash mismatched",
			}
		} else {
			logger.WithContext(ctx).WithField("txid", ticket.TxID).Info("reconstructed file hash matched")
			verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
				TxID:                     ticket.TxID,
				TicketType:               ticket.TicketType,
				MissingKeys:              ticket.MissingKeys,
				ReconstructedFileHash:    nil,
				IsReconstructionRequired: false,
				IsVerified:               true,
				Message:                  "reconstructed file hash matched",
			}

		}

		return task.prepareAndSendVerificationMsg(ctx, *verificationMsg)
	} else if senseTicket != nil {
		reqSelfHealing, mostCommonFile, sigs := task.senseCheckingProcess(ctx, senseTicket.DDAndFingerprintsIDs)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error in checking process for sense action ticket")
		}

		if !reqSelfHealing {
			if !ticket.IsReconstructionRequired {
				logger.WithField("ticket_txid", ticket.TxID).
					Info("is reconstruction required set to false by both recipient and verifier for sense ticket")

				verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
					TxID:                     ticket.TxID,
					TicketType:               ticket.TicketType,
					MissingKeys:              ticket.MissingKeys,
					ReconstructedFileHash:    nil,
					IsReconstructionRequired: false,
					IsVerified:               true,
					Message:                  "is reconstruction required set to false by both recipient and verifier",
				}
			} else {
				logger.WithField("ticket_txid", ticket.TxID).
					Info("is reconstruction required set to true by recipient but false by verifier")

				verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
					TxID:                     ticket.TxID,
					TicketType:               ticket.TicketType,
					MissingKeys:              ticket.MissingKeys,
					ReconstructedFileHash:    nil,
					IsReconstructionRequired: false,
					IsVerified:               false,
					Message:                  "is reconstruction required set to true by recipient but false by verifier",
				}
			}
			return task.prepareAndSendVerificationMsg(ctx, *verificationMsg)
		}

		ids, _, err := task.senseSelfHealing(ctx, senseTicket, mostCommonFile, sigs)
		if err != nil {
			logger.WithField("ticket_txid", ticket.TxID).Info("error self-healing sense ticket")
		}

		if ok := compareFileIDs(ids, ticket.FileIDs); !ok {
			logger.WithField("ticket_txid", ticket.TxID).Info("sense file hash mismatched")
			verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
				TxID:                     ticket.TxID,
				TicketType:               ticket.TicketType,
				MissingKeys:              ticket.MissingKeys,
				ReconstructedFileHash:    nil,
				IsReconstructionRequired: false,
				IsVerified:               false,
				Message:                  "reconstructed file hash mismatched",
			}
		} else {
			logger.WithField("ticket_txid", ticket.TxID).Info("sense file hash matched")
			verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
				TxID:                     ticket.TxID,
				TicketType:               ticket.TicketType,
				MissingKeys:              ticket.MissingKeys,
				ReconstructedFileHash:    nil,
				IsReconstructionRequired: false,
				IsVerified:               true,
				Message:                  "reconstructed file hash matched",
			}
		}

		logger.Info("sending verification back to the recipient")
		return task.prepareAndSendVerificationMsg(ctx, *verificationMsg)
	}

	return task.prepareAndSendVerificationMsg(ctx, *verificationMsg)
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

func (task *SHTask) prepareAndSendVerificationMsg(ctx context.Context, verificationMsg types.SelfHealingMessage) (*pb.SelfHealingMessage, error) {
	logger := log.WithContext(ctx).WithField("trigger_id", verificationMsg.TriggerID).
		WithField("challenge_id", verificationMsg.SelfHealingMessageData.Verification.ChallengeID)

	verificationMsg.SelfHealingMessageData.Verification.Timestamp = time.Now().UTC()

	signature, data, err := task.SignMessage(ctx, verificationMsg.SelfHealingMessageData)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error signing the self-healing verification msg")
	}
	verificationMsg.SenderSignature = signature
	logger.Info("verification msg has been signed")

	msg := &pb.SelfHealingMessage{
		TriggerId:       verificationMsg.TriggerID,
		MessageType:     pb.SelfHealingMessageMessageType(verificationMsg.MessageType),
		SenderId:        verificationMsg.SenderID,
		SenderSignature: signature,
		Data:            data,
	}

	return msg, nil
}
