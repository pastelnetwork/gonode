package selfhealing

import (
	"bytes"
	"context"
	"fmt"
	"time"

	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
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
	logger.Debug("self healing response message has been validated")

	var (
		nftTicket     *pastel.NFTTicket
		cascadeTicket *pastel.APICascadeTicket
		senseTicket   *pastel.APISenseTicket
	)

	log.WithContext(ctx).Debug("retrieving block no and verbose")
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
	logger.Debug("block count & merkelroot has been retrieved")

	verificationMsg := &types.SelfHealingMessage{
		TriggerID:   incomingResponseMessage.TriggerID,
		SenderID:    task.nodeID,
		MessageType: types.SelfHealingVerificationMessage,
		SelfHealingMessageData: types.SelfHealingMessageData{
			Response: types.SelfHealingResponseData{
				ChallengeID:     incomingResponseMessage.SelfHealingMessageData.Response.ChallengeID,
				Block:           incomingResponseMessage.SelfHealingMessageData.Response.Block,
				Merkelroot:      incomingResponseMessage.SelfHealingMessageData.Response.Merkelroot,
				Timestamp:       incomingResponseMessage.SelfHealingMessageData.Response.Timestamp,
				RespondedTicket: incomingResponseMessage.SelfHealingMessageData.Response.RespondedTicket,
			},
			Verification: types.SelfHealingVerificationData{
				NodeID:      task.nodeID,
				NodeAddress: task.nodeAddress,
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

	logger.WithField("ticket_txid", ticket.TxID).Debug("starting self-healing verification for the ticket")
	nftTicket, cascadeTicket, senseTicket, err = task.getTicket(ctx, ticket.TxID, TicketType(ticket.TicketType))
	if err != nil {
		logger.WithError(err).Error("Error getRelevantTicketFromMsg")
	}
	logger.Debug("reg ticket has been retrieved for verification")

	if cascadeTicket != nil || nftTicket != nil {
		isReconstructionReq := task.isReconstructionRequired(ctx, ticket.MissingKeys)
		if !isReconstructionReq {
			logger.WithField("ticket_txid", ticket.TxID).Debug("reconstruction is not required for ticket")

			if !ticket.IsReconstructionRequired {
				logger.WithField("ticket_txid", ticket.TxID).
					Info("is reconstruction required set to false by both recipient and verifier")

				verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
					TxID:                             ticket.TxID,
					TicketType:                       ticket.TicketType,
					MissingKeys:                      ticket.MissingKeys,
					ReconstructedFileHash:            nil,
					IsReconstructionRequired:         false,
					IsReconstructionRequiredByHealer: ticket.IsReconstructionRequired,
					IsVerified:                       true,
					Message:                          "is reconstruction required set to false by both recipient and verifier",
				}
			} else {
				logger.WithField("ticket_txid", ticket.TxID).
					Info("is reconstruction required set to true by recipient but false by verifier")

				verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
					TxID:                             ticket.TxID,
					TicketType:                       ticket.TicketType,
					MissingKeys:                      ticket.MissingKeys,
					ReconstructedFileHash:            nil,
					IsReconstructionRequired:         false,
					IsReconstructionRequiredByHealer: ticket.IsReconstructionRequired,
					IsVerified:                       false,
					Message:                          "is reconstruction required set to true by recipient but false by verifier",
				}
			}

			return task.prepareAndSendVerificationMsg(ctx, *verificationMsg)
		}

		_, reconstructedFileHash, err := task.selfHealing(ctx, ticket.TxID, nftTicket, cascadeTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error self-healing the file")
		}

		logger.Debug("comparing data hashes")
		if !bytes.Equal(reconstructedFileHash, ticket.ReconstructedFileHash) {
			logger.WithContext(ctx).WithField("txid", ticket.TxID).Info("reconstructed file hash does not match with the verifier reconstructed file")
			verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
				TxID:                             ticket.TxID,
				TicketType:                       ticket.TicketType,
				MissingKeys:                      ticket.MissingKeys,
				ReconstructedFileHash:            reconstructedFileHash,
				IsReconstructionRequiredByHealer: ticket.IsReconstructionRequired,
				IsReconstructionRequired:         true,
				IsVerified:                       false,
				Message:                          "reconstructed file hash mismatched",
			}
		} else {
			logger.WithContext(ctx).WithField("txid", ticket.TxID).Info("reconstructed file hash matched")
			verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
				TxID:                             ticket.TxID,
				TicketType:                       ticket.TicketType,
				MissingKeys:                      ticket.MissingKeys,
				ReconstructedFileHash:            reconstructedFileHash,
				IsReconstructionRequiredByHealer: ticket.IsReconstructionRequired,
				IsReconstructionRequired:         true,
				IsVerified:                       true,
				Message:                          "reconstructed file hash matched",
			}

		}

		return task.prepareAndSendVerificationMsg(ctx, *verificationMsg)
	} else if senseTicket != nil {
		reqSelfHealing, sortedFiles := task.senseCheckingProcess(ctx, senseTicket.DDAndFingerprintsIDs)
		if !reqSelfHealing {
			if !ticket.IsReconstructionRequired {
				logger.WithField("ticket_txid", ticket.TxID).
					Info("is reconstruction required set to false by both recipient and verifier for sense ticket")

				verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
					TxID:                             ticket.TxID,
					TicketType:                       ticket.TicketType,
					MissingKeys:                      ticket.MissingKeys,
					IsReconstructionRequiredByHealer: ticket.IsReconstructionRequired,
					ReconstructedFileHash:            nil,
					IsReconstructionRequired:         false,
					IsVerified:                       true,
					Message:                          "is reconstruction required set to false by both recipient and verifier",
				}
			} else {
				logger.WithField("ticket_txid", ticket.TxID).
					Info("is reconstruction required set to true by recipient but false by verifier")

				verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
					TxID:                             ticket.TxID,
					TicketType:                       ticket.TicketType,
					MissingKeys:                      ticket.MissingKeys,
					IsReconstructionRequiredByHealer: ticket.IsReconstructionRequired,
					ReconstructedFileHash:            nil,
					IsReconstructionRequired:         false,
					IsVerified:                       false,
					Message:                          "is reconstruction required set to true by recipient but false by verifier",
				}
			}

			return task.prepareAndSendVerificationMsg(ctx, *verificationMsg)
		}

		var fileHash []byte
		fileHash, err = task.getReconstructedFileHashForVerification(ctx, sortedFiles)
		if err != nil {
			logger.WithField("ticket_txid", ticket.TxID).Debug("error getting reconstructed sense file hash")
		}

		if bytes.Equal(fileHash, incomingResponseMessage.SelfHealingMessageData.Response.RespondedTicket.ReconstructedFileHash) {
			logger.WithField("ticket_txid", ticket.TxID).Info("reconstructed sense file hash matched")
			verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
				TxID:                             ticket.TxID,
				TicketType:                       ticket.TicketType,
				MissingKeys:                      ticket.MissingKeys,
				IsReconstructionRequiredByHealer: ticket.IsReconstructionRequired,
				ReconstructedFileHash:            fileHash,
				IsReconstructionRequired:         true,
				IsVerified:                       true,
				Message:                          "reconstructed sense file hash matched",
			}
		} else {
			logger.WithField("ticket_txid", ticket.TxID).Info("reconstructed sense file hash mismatched")
			verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket = types.VerifiedTicket{
				TxID:                             ticket.TxID,
				TicketType:                       ticket.TicketType,
				MissingKeys:                      ticket.MissingKeys,
				IsReconstructionRequiredByHealer: ticket.IsReconstructionRequired,
				ReconstructedFileHash:            fileHash,
				IsReconstructionRequired:         true,
				IsVerified:                       false,
				Message:                          "reconstructed sense file hash mismatched",
			}
		}

		logger.Debug("sending verification back to the recipient")
		return task.prepareAndSendVerificationMsg(ctx, *verificationMsg)
	}

	return task.prepareAndSendVerificationMsg(ctx, *verificationMsg)
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
	logger.Debug("verification msg has been signed")

	msg := &pb.SelfHealingMessage{
		TriggerId:       verificationMsg.TriggerID,
		MessageType:     pb.SelfHealingMessageMessageType(verificationMsg.MessageType),
		SenderId:        verificationMsg.SenderID,
		SenderSignature: signature,
		Data:            data,
	}

	return msg, nil
}

func (task *SHTask) getReconstructedFileHashForVerification(ctx context.Context, sortedFiles [][]byte) ([]byte, error) {
	if sortedFiles == nil {
		return []byte{}, errors.Errorf("empty list of sorted files")
	}

	if len(sortedFiles) == 0 {
		return []byte{}, errors.Errorf("empty list of sorted files")
	}

	var (
		ddAndFingerprintFile pastel.DDAndFingerprints
		data                 []byte
	)

	for i := 0; i < len(sortedFiles); i++ {
		file := sortedFiles[i]

		decompressedData, err := utils.Decompress(file)
		if err != nil {
			continue
		}
		log.WithContext(ctx).Debug("file has been decompressed")

		if len(decompressedData) == 0 {
			continue
		}

		splits := bytes.Split(decompressedData, []byte{pastel.SeparatorByte})
		if len(splits) < 4 {
			continue
		}

		dataToJSONDecode, err := utils.B64Decode(splits[0])
		if err != nil {
			continue
		}

		if len(dataToJSONDecode) == 0 {
			continue
		}

		err = json.Unmarshal(dataToJSONDecode, &ddAndFingerprintFile)
		if err != nil {
			continue
		}

		data, err = json.Marshal(ddAndFingerprintFile)
		if err != nil {
			continue
		}

		fileHash, err := utils.Sha3256hash(data)
		if err != nil {
			continue
		}

		log.WithContext(ctx).Debug("sense reconstructed file hash has been generated for verification")

		return fileHash, nil
	}

	return []byte{}, nil
}
