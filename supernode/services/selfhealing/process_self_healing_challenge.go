package selfhealing

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/crypto/sha3"
	"os"
	"sync"
	"time"

	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	common "github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	verifyThreshold = 2
)

// ProcessSelfHealingChallenge is called from grpc server, which processes the self-healing challenge,
// and will execute the reconstruction work.
func (task *SHTask) ProcessSelfHealingChallenge(ctx context.Context, incomingChallengeMessage types.SelfHealingMessage) error {
	logger := log.WithContext(ctx).WithField("trigger_id", incomingChallengeMessage.TriggerID)

	logger.Info("ProcessSelfHealingChallenge has been invoked")

	// wait if test env so tests can mock the p2p
	if os.Getenv("INTEGRATION_TEST_ENV") == "true" {
		time.Sleep(10 * time.Second)
	}

	selfHealingMetrics := task.SelfHealingMetricsMap[incomingChallengeMessage.TriggerID]
	selfHealingMetrics.ChallengeID = incomingChallengeMessage.TriggerID

	// incoming challenge message validation
	if err := task.validateSelfHealingChallengeIncomingData(ctx, incomingChallengeMessage); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error validating self-healing challenge incoming data: ")
		return err
	}
	logger.Info("self healing challenge message has been validated")

	log.WithContext(ctx).Info("retrieving block no and verbose")
	currentBlockCount, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block count")
		return err
	}
	blkVerbose1, err := task.SuperNodeService.PastelClient.GetBlockVerbose1(ctx, currentBlockCount)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block verbose 1")
		return err
	}
	merkleroot := blkVerbose1.MerkleRoot
	logger.Info("block count & merkelroot has been retrieved")

	responseMsg := types.SelfHealingMessage{
		TriggerID:   incomingChallengeMessage.TriggerID,
		SenderID:    task.nodeID,
		MessageType: types.SelfHealingResponseMessage,
		SelfHealingMessageData: types.SelfHealingMessageData{
			ChallengerID: incomingChallengeMessage.SelfHealingMessageData.ChallengerID,
			Challenge: types.SelfHealingChallengeData{
				Block:            incomingChallengeMessage.SelfHealingMessageData.Challenge.Block,
				Merkelroot:       incomingChallengeMessage.SelfHealingMessageData.Challenge.Merkelroot,
				ChallengeTickets: incomingChallengeMessage.SelfHealingMessageData.Challenge.ChallengeTickets,
				Timestamp:        incomingChallengeMessage.SelfHealingMessageData.Challenge.Timestamp,
				NodesOnWatchlist: incomingChallengeMessage.SelfHealingMessageData.Challenge.NodesOnWatchlist,
			},
			Response: types.SelfHealingResponseData{
				Block:      currentBlockCount,
				Merkelroot: merkleroot,
			},
		},
	}

	var (
		nftTicket     *pastel.NFTTicket
		cascadeTicket *pastel.APICascadeTicket
		senseTicket   *pastel.APISenseTicket
	)

	challengeTickets := incomingChallengeMessage.SelfHealingMessageData.Challenge.ChallengeTickets
	if challengeTickets == nil {
		logger.Info("no tickets found to challenge")
		return nil
	}

	selfHealingMetrics.SentTicketsForSelfHealing = len(challengeTickets)

	for _, ticket := range challengeTickets {
		if task.SHService.ticketsMap[ticket.TxID] {
			log.WithContext(ctx).WithField("txid", ticket.TxID).Info("ticket processing is already in progress")
			selfHealingMetrics.TicketsInProgress++
			continue
		}
		task.SHService.ticketsMap[ticket.TxID] = true
		selfHealingMetrics.EstimatedMissingKeys = selfHealingMetrics.EstimatedMissingKeys + len(ticket.MissingKeys)

		logger.WithField("ticket_txid", ticket.TxID).Info("going to initiate checking process")
		nftTicket, cascadeTicket, senseTicket, err = task.getTicket(ctx, ticket.TxID, TicketType(ticket.TicketType))
		if err != nil {
			logger.WithError(err).Error("Error getRelevantTicketFromMsg")
		}
		logger.Info("reg ticket has been retrieved")

		challengeID := utils.GetHashFromString(ticket.TxID + responseMsg.TriggerID)
		responseMsg.SelfHealingMessageData.Response.ChallengeID = challengeID

		if nftTicket != nil || cascadeTicket != nil {
			newCtx := context.Background()
			if !task.isReconstructionRequired(newCtx, ticket.MissingKeys) {
				logger.WithField("txid", ticket.TxID).Info("reconstruction is not required for ticket")

				responseMsg.SelfHealingMessageData.Response.RespondedTicket = types.RespondedTicket{
					TxID:                     ticket.TxID,
					TicketType:               ticket.TicketType,
					MissingKeys:              ticket.MissingKeys,
					ReconstructedFileHash:    nil,
					RaptorQSymbols:           nil,
					IsReconstructionRequired: false,
				}

				newCtx = context.Background()

				logger.WithField("txid", ticket.TxID).Info("sending response for verification")
				verifications, err := task.getVerifications(newCtx, responseMsg)
				if err != nil {
					logger.WithField("txid", ticket.TxID).Info("unable to get verifications")
				}
				logger.WithField("txid", ticket.TxID).WithField("verifications", len(verifications)).
					Info("verifications have been evaluated")

				if len(verifications) >= verifyThreshold {
					selfHealingMetrics.SuccessfullyVerifiedTickets++

					logger.WithField("txid", ticket.TxID).Info("self-healing not required for this ticket, " +
						"verified by verifiers")
				}

				continue
			}
			selfHealingMetrics.TicketsRequiredSelfHealing++
			logger.WithField("ticket_txid", ticket.TxID).Info("going to initiate self-healing")

			newCtx = context.Background()
			raptorQSymbols, reconstructedFileHash, err := task.selfHealing(newCtx, ticket.TxID, nftTicket, cascadeTicket)
			if err != nil {
				logger.WithError(err).Error("Error self-healing the file")
				continue
			}
			logger.WithField("ticket_txid", ticket.TxID).Info("file has been reconstructed")

			responseMsg.SelfHealingMessageData.Response.RespondedTicket = types.RespondedTicket{
				TxID:                     ticket.TxID,
				TicketType:               ticket.TicketType,
				MissingKeys:              ticket.MissingKeys,
				ReconstructedFileHash:    reconstructedFileHash,
				RaptorQSymbols:           task.RaptorQSymbols,
				IsReconstructionRequired: true,
			}
			selfHealingMetrics.SuccessfullySelfHealedTickets++

			newCtx = context.Background()
			logger.WithField("txid", ticket.TxID).Info("sending response for verification")
			verifications, err := task.getVerifications(newCtx, responseMsg)
			if err != nil {
				logger.WithField("txid", ticket.TxID).Info("unable to get verifications")
			}
			logger.WithField("txid", ticket.TxID).WithField("verifications", len(verifications)).
				Info("verifications have been evaluated")

			if len(verifications) >= verifyThreshold {
				selfHealingMetrics.SuccessfullyVerifiedTickets++

				newCtx = context.Background()
				if err := task.StorageHandler.StoreRaptorQSymbolsIntoP2P(newCtx, raptorQSymbols, ""); err != nil {
					logger.WithError(err).Error("Error storing symbols to P2P")
					continue
				}
				logger.WithField("ticket_txid", ticket.TxID).Info("raptor q symbols have been stored")

				sig, _, err := task.SignMessage(context.Background(), responseMsg.SelfHealingMessageData)
				if err != nil {
					logger.WithError(err).Error("error storing completion execution metric")
					continue
				}

				completionMsg := []types.SelfHealingMessage{
					{
						TriggerID:       responseMsg.TriggerID,
						MessageType:     types.SelfHealingCompletionMessage,
						SenderID:        task.nodeID,
						SenderSignature: sig,
					},
				}
				completionMsgBytes, err := json.Marshal(completionMsg)
				if err != nil {
					logger.WithError(err).Error("error converting the completion msg to bytes for metrics")
				}

				if err := task.StoreSelfHealingExecutionMetrics(ctx, types.SelfHealingExecutionMetric{
					TriggerID:       responseMsg.TriggerID,
					ChallengeID:     responseMsg.SelfHealingMessageData.ChallengerID,
					MessageType:     int(types.SelfHealingCompletionMessage),
					SenderID:        task.nodeID,
					Data:            completionMsgBytes,
					SenderSignature: sig,
				}); err != nil {
					logger.WithError(err).Error("error storing self-healing completion execution metric")
				}

				continue
			}

			logger.WithField("ticket_txid", ticket.TxID).Info("verify threshold does not meet, not saving symbols")
		} else if senseTicket != nil {
			newCtx := context.Background()
			reqSelfHealing, mostCommonFile, sigs := task.senseCheckingProcess(newCtx, senseTicket.DDAndFingerprintsIDs)
			if err != nil {
				logger.WithField("ticket_txid", ticket.TxID).WithError(err).Error("Error in checking process for sense action ticket")
				continue
			}

			if !reqSelfHealing {
				logger.WithField("ticket_txid", ticket.TxID).WithError(err).Error("self-healing not required for sense ticket")
				responseMsg.SelfHealingMessageData.Response.RespondedTicket = types.RespondedTicket{
					TxID:                     ticket.TxID,
					TicketType:               ticket.TicketType,
					MissingKeys:              ticket.MissingKeys,
					ReconstructedFileHash:    nil,
					RaptorQSymbols:           nil,
					FileIDs:                  nil,
					IsReconstructionRequired: false,
				}

				newCtx = context.Background()

				logger.WithField("txid", ticket.TxID).Info("sending response for verification")
				verifications, err := task.getVerifications(newCtx, responseMsg)
				if err != nil {
					logger.WithField("txid", ticket.TxID).Info("unable to get verifications")
				}
				logger.WithField("txid", ticket.TxID).WithField("verifications", len(verifications)).
					Info("verifications have been evaluated")

				if len(verifications) >= verifyThreshold {
					selfHealingMetrics.SuccessfullyVerifiedTickets++

					logger.WithField("txid", ticket.TxID).Info("self-healing not required for this sense ticket, " +
						"verified by verifiers")
				}

				continue
			}

			selfHealingMetrics.TicketsRequiredSelfHealing++
			logger.WithField("txid", ticket.TxID).Info("going to initiate sense self-healing")

			ids, idFiles, err := task.senseSelfHealing(ctx, senseTicket, mostCommonFile, sigs)
			if err != nil {
				logger.WithField("ticket_txid", ticket.TxID).WithError(err).Error("self-healing not required for sense action ticket")
				continue
			}
			task.IDFiles = idFiles
			logger.WithField("ticket_txid", ticket.TxID).Info("sense ticket has been reconstructed")

			fileBytes, err := json.Marshal(mostCommonFile)
			if err != nil {
				logger.WithError(err).Error(err)
			}

			fileHash, err := utils.Sha3256hash(fileBytes)
			if err != nil {
				logger.WithError(err).Error(err)
			}

			responseMsg.SelfHealingMessageData.Response.RespondedTicket = types.RespondedTicket{
				TxID:                     ticket.TxID,
				TicketType:               ticket.TicketType,
				MissingKeys:              ticket.MissingKeys,
				ReconstructedFileHash:    fileHash,
				FileIDs:                  ids,
				IsReconstructionRequired: true,
			}
			selfHealingMetrics.SuccessfullySelfHealedTickets++

			newCtx = context.Background()
			logger.WithField("txid", ticket.TxID).Info("sending response for verification")
			verifications, err := task.getVerifications(newCtx, responseMsg)
			if err != nil {
				logger.WithField("txid", ticket.TxID).Info("unable to get verifications")
			}
			logger.WithField("txid", ticket.TxID).WithField("verifications", len(verifications)).
				Info("verifications have been evaluated")

			if len(verifications) >= verifyThreshold {
				selfHealingMetrics.SuccessfullyVerifiedTickets++

				newCtx = context.Background()
				if err := task.StorageHandler.StoreBatch(newCtx, idFiles, common.P2PDataDDMetadata); err != nil {
					logger.WithError(err).Error("Error storing symbols to P2P")
					continue
				}
				logger.WithField("ticket_txid", ticket.TxID).Info("dd & fp ids have been stored")

				sig, _, err := task.SignMessage(context.Background(), responseMsg.SelfHealingMessageData)
				if err != nil {
					logger.WithError(err).Error("error storing completion execution metric")
					continue
				}

				completionMsg := []types.SelfHealingMessage{
					{
						TriggerID:       responseMsg.TriggerID,
						MessageType:     types.SelfHealingCompletionMessage,
						SenderID:        task.nodeID,
						SenderSignature: sig,
					},
				}
				completionMsgBytes, err := json.Marshal(completionMsg)
				if err != nil {
					logger.WithError(err).Error("error converting the completion msg to bytes for metrics")
				}

				if err := task.StoreSelfHealingExecutionMetrics(ctx, types.SelfHealingExecutionMetric{
					TriggerID:       responseMsg.TriggerID,
					ChallengeID:     responseMsg.SelfHealingMessageData.ChallengerID,
					MessageType:     int(types.SelfHealingCompletionMessage),
					SenderID:        task.nodeID,
					Data:            completionMsgBytes,
					SenderSignature: sig,
				}); err != nil {
					logger.WithError(err).Error("error storing self-healing completion execution metric")
				}

				continue
			}
		}

		logger.WithField("self_healing_metrics", selfHealingMetrics).Info("self-healing metrics have been updated")
	}

	task.SHService.SelfHealingMetricsMap[incomingChallengeMessage.TriggerID] = selfHealingMetrics
	logger.WithField("self_healing_metrics", task.SHService.SelfHealingMetricsMap[incomingChallengeMessage.TriggerID]).
		Info("self-healing completed")

	return nil
}

func (task *SHTask) validateSelfHealingChallengeIncomingData(ctx context.Context, incomingChallengeMessage types.SelfHealingMessage) error {
	if incomingChallengeMessage.MessageType != types.SelfHealingChallengeMessage {
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

// VerifyMessageSignature verifies the sender's signature on message
func (task *SHTask) VerifyMessageSignature(ctx context.Context, msg types.SelfHealingMessage) (bool, error) {
	data, err := json.Marshal(msg.SelfHealingMessageData)
	if err != nil {
		return false, errors.Errorf("unable to marshal message data")
	}

	isVerified, err := task.PastelClient.Verify(ctx, data, string(msg.SenderSignature), msg.SenderID, pastel.SignAlgorithmED448)
	if err != nil {
		return false, errors.Errorf("error verifying self-healing message: %w", err)
	}

	return isVerified, nil
}

func (task *SHTask) checkingProcess(ctx context.Context, fileHash string) (requiredReconstruction bool) {
	rqIDsData, err := task.P2PClient.Retrieve(ctx, fileHash, false)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("SymbolIDsFileId", fileHash).Info("Retrieve compressed symbol IDs file from P2P failed")
		return true
	}

	if len(rqIDsData) == 0 {
		log.WithContext(ctx).WithField("SymbolIDsFileId", fileHash).Info("Retrieve compressed symbol IDs file from P2P is empty")
		return true
	}

	return false
}

func (task *SHTask) selfHealing(ctx context.Context, txid string, regTicket *pastel.NFTTicket, cascadeTicket *pastel.APICascadeTicket) (file, reconstructedFileHash []byte, err error) {
	log.WithContext(ctx).WithField("txid", txid).Info("Self-healing initiated")
	if regTicket != nil {
		log.WithContext(ctx).WithField("txid", txid).WithField("reg_ticket", regTicket).Debug("going to self heal reg ticket")

		file, err = task.downloadTask.RestoreFile(ctx, regTicket.AppTicketData.RQIDs, regTicket.AppTicketData.RQOti, regTicket.AppTicketData.DataHash, txid)
		if err != nil {
			log.WithContext(ctx).WithField("txid", txid).WithError(err).Error("Unable to restore file")
			return nil, nil, errors.Errorf("unable to restore file")
		}

		fileHash := sha3.Sum256(file)
		log.WithContext(ctx).WithField("txid", txid).Info("file has been restored")
		reconstructedFileHash = fileHash[:]
	} else if cascadeTicket != nil {
		log.WithContext(ctx).WithField("txid", txid).WithField("cascade_ticket", cascadeTicket).Debug("going to self heal cascade ticket")

		file, err = task.downloadTask.RestoreFile(ctx, cascadeTicket.RQIDs, cascadeTicket.RQOti, cascadeTicket.DataHash, txid)
		if err != nil {
			log.WithContext(ctx).WithField("txid", txid).WithError(err).Error("Unable to restore file")
			return nil, nil, errors.Errorf("unable to restore file")
		}

		fileHash := sha3.Sum256(file)
		log.WithContext(ctx).WithField("txid", txid).Info("file has been restored")
		reconstructedFileHash = fileHash[:]
	} else {
		log.WithContext(ctx).WithField("txid", txid).Info("not a valid ticket")

		return file, reconstructedFileHash, errors.Errorf("not a valid cascade or NFT ticket")
	}

	return file, reconstructedFileHash, nil
}

func (task *SHTask) isReconstructionRequired(ctx context.Context, fileHashes []string) bool {
	for _, fileHash := range fileHashes {
		if task.checkingProcess(ctx, fileHash) {
			return true
		}
	}

	return false
}

func (task *SHTask) senseCheckingProcess(ctx context.Context, ddFPIDs []string) (requiredReconstruction bool, mostCommonFile *pastel.DDAndFingerprints, signatures signatures) {
	//download all the DD and FPIds
	availableDDFPFiles, err := task.DownloadDDAndFingerprints(ctx, ddFPIDs)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Downloading DD and Fingerprints")
	}

	//if not able to retrieve all should self-heal
	if len(ddFPIDs) == len(availableDDFPFiles) {
		log.WithContext(ctx).WithField("ddfp_ids_count", len(ddFPIDs)).WithField("file_count", len(availableDDFPFiles)).
			Info("no of ids and retrieved files from the network are equal, self-healing not required for this sense ticket")
		return false, nil, signatures
	}
	log.WithContext(ctx).WithField("ddfp_ids_count", len(ddFPIDs)).WithField("file_count", len(availableDDFPFiles)).
		Info("no of ids and retrieved files from the network are not equal, self-healing required")

	//Find most common file (source of truth) for self-healing
	hashTally := make(map[string]int)
	fileHashMap := make(map[string]*pastel.DDAndFingerprints)
	fileHashToDDFPKeyMap := make(map[string]string)

	for key, DDFPIdFile := range availableDDFPFiles {
		bytes, err := json.Marshal(DDFPIdFile)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to marshal DDFPIdFile")
			continue
		}

		fileHash := utils.GetHashStringFromBytes(bytes)

		//Generating hashes, maintaining hash tally and fileHash map to find the file with the most common hash
		hashTally[fileHash]++
		fileHashMap[fileHash] = &DDFPIdFile
		fileHashToDDFPKeyMap[fileHash] = key
	}

	var mostCommonHash string
	var greatestTally int
	for hash, tally := range hashTally {
		if tally > greatestTally {
			mostCommonHash = hash
			greatestTally = tally
		}
	}

	return true, fileHashMap[mostCommonHash], task.KeySignaturesMap[fileHashToDDFPKeyMap[mostCommonHash]]
}

func (task *SHTask) senseSelfHealing(ctx context.Context, senseTicket *pastel.APISenseTicket, ddFPID *pastel.DDAndFingerprints, sig signatures) (ids []string, idFiles [][]byte, err error) {
	if len(sig.SNSignature1) == 0 || len(sig.SNSignature2) == 0 || len(sig.SNSignature3) == 0 {
		return nil, nil, errors.Errorf("empty signature")
	}

	ddDataJSON, err := json.Marshal(ddFPID)
	if err != nil {
		return nil, nil, errors.Errorf("failed to marshal dd-data: %w", err)
	}

	ddEncoded := utils.B64Encode(ddDataJSON)

	var buffer bytes.Buffer
	buffer.Write(ddEncoded)
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(sig.SNSignature1)
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(sig.SNSignature2)
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(sig.SNSignature3)
	ddFpFile := buffer.Bytes()

	ids, idFiles, err = pastel.GetIDFiles(ctx, ddFpFile, senseTicket.DDAndFingerprintsIc, senseTicket.DDAndFingerprintsMax)
	if err != nil {
		return nil, nil, errors.Errorf("get ID Files: %w", err)
	}

	return ids, idFiles, nil
}

func (task *SHTask) getVerifications(ctx context.Context, msg types.SelfHealingMessage) (map[string]types.SelfHealingMessage, error) {
	nodesToConnect, err := task.GetNodesAddressesToConnect(ctx, &msg)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error finding nodes to connect for self-healing verification")
	}

	if nodesToConnect == nil {
		return nil, errors.Errorf("no nodes found to connect for getting affirmations")
	}

	responseMsg, err := task.prepareResponseMessage(ctx, msg)
	if err != nil {
		return nil, errors.Errorf("error preparing response message")
	}

	successfulVerifications, allVerifications := task.processSelfHealingVerifications(ctx, nodesToConnect, &responseMsg)

	verificationExecMetrics, err := task.getVerificationExecMetrics(msg, allVerifications)
	if err != nil {
		return successfulVerifications, nil
	}

	if verificationExecMetrics != nil {
		err := task.StoreSelfHealingExecutionMetrics(ctx, *verificationExecMetrics)
		if err != nil {
			log.WithContext(ctx).WithField("trigger_id", msg.TriggerID).
				WithField("txid", msg.SelfHealingMessageData.Response.RespondedTicket.TxID).
				WithField("challenge_id", msg.SelfHealingMessageData.Response.ChallengeID).
				Error("error storing verification execution metrics")
		}
	}
	log.WithContext(ctx).WithField("trigger_id", msg.TriggerID).
		WithField("txid", msg.SelfHealingMessageData.Response.RespondedTicket.TxID).
		WithField("challenge_id", msg.SelfHealingMessageData.Response.ChallengeID).
		Info("verification execution metrics have been stored")

	return successfulVerifications, nil
}

// GetNodesAddressesToConnect basically retrieves the masternode address against the pastel-id from the list and return that
func (task *SHTask) GetNodesAddressesToConnect(ctx context.Context, challengeMessage *types.SelfHealingMessage) ([]pastel.MasterNode, error) {
	var nodesToConnect []pastel.MasterNode

	var (
		bc         int32
		merkleroot string
	)

	blockCount, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block verbose 1")
		bc = challengeMessage.SelfHealingMessageData.Response.Block
	} else {
		bc = blockCount
	}

	blkVerbose1, err := task.SuperNodeService.PastelClient.GetBlockVerbose1(ctx, bc)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block verbose 1")
		merkleroot = challengeMessage.SelfHealingMessageData.Response.Merkelroot
	} else {
		merkleroot = blkVerbose1.MerkleRoot
	}

	supernodes, err := task.SuperNodeService.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		log.WithContext(ctx).WithField("challengeID", challengeMessage.TriggerID).WithField("method", "GetNodesAddressesToConnect").WithError(err).Warn("could not get Supernode extra: ", err.Error())
		return nil, err
	}

	mapSupernodes := make(map[string]pastel.MasterNode)
	for _, mn := range supernodes {
		if mn.ExtAddress == "" || mn.ExtKey == "" {
			log.WithContext(ctx).WithField("challengeID", challengeMessage.TriggerID).WithField("method", "GetNodesAddressesToConnect").
				WithField("node_id", mn.ExtKey).Warn("node address or node id is empty")

			continue
		}

		mapSupernodes[mn.ExtKey] = mn
	}

	logger := log.WithContext(ctx).WithField("challenge_id", challengeMessage.TriggerID).
		WithField("message_type", challengeMessage.MessageType).WithField("node_id", task.nodeID)

	listOfSupernodes, err := task.getListOfOnlineSupernodes(ctx)
	if err != nil {
		logger.WithError(err).Error("error getting list of online supernodes")
	}

	switch challengeMessage.MessageType {
	case types.SelfHealingResponseMessage:

		//Finding 5 closest nodes to previous block hash
		sliceOfSupernodesClosestToPreviousBlockHash := task.GetNClosestSupernodeIDsToComparisonString(ctx, 5, merkleroot, listOfSupernodes, task.nodeID)
		challengeMessage.SelfHealingMessageData.Response.Verifiers = sliceOfSupernodesClosestToPreviousBlockHash
		log.WithContext(ctx).Info(fmt.Sprintf("sliceOfSupernodesClosestToPreviousBlockHash:%s", sliceOfSupernodesClosestToPreviousBlockHash))

		for _, verifier := range sliceOfSupernodesClosestToPreviousBlockHash {
			nodesToConnect = append(nodesToConnect, mapSupernodes[verifier])
		}

		logger.WithField("nodes_to_connect", nodesToConnect).Info("nodes to send self-healing response msg have been selected")
		return nodesToConnect, nil
	default:
		return nil, errors.Errorf("no nodes found to send message")
	}

	return nil, err
}

// prepareEvaluationMessage prepares the evaluation message by gathering the data required for it
func (task *SHTask) prepareResponseMessage(ctx context.Context, responseMessage types.SelfHealingMessage) (pb.SelfHealingMessage, error) {
	logger := log.WithContext(ctx).WithField("trigger_id", responseMessage.TriggerID).
		WithField("ticket_txid", responseMessage.SelfHealingMessageData.Response.RespondedTicket.TxID).
		WithField("challenge_id", responseMessage.SelfHealingMessageData.Response.ChallengeID)

	responseMessage.SelfHealingMessageData.Response.Timestamp = time.Now().UTC()

	signature, data, err := task.SignMessage(ctx, responseMessage.SelfHealingMessageData)
	if err != nil {
		logger.WithError(err).Error("error signing the response message")
		return pb.SelfHealingMessage{}, err
	}
	responseMessage.SenderSignature = signature

	responseExecutionMetric := []types.SelfHealingMessage{
		{
			TriggerID:   responseMessage.TriggerID,
			MessageType: responseMessage.MessageType,
			SelfHealingMessageData: types.SelfHealingMessageData{
				Response: responseMessage.SelfHealingMessageData.Response,
			},
			SenderID: task.nodeID,
		},
	}

	responseMetricBytes, err := json.Marshal(responseExecutionMetric)
	if err != nil {
		return pb.SelfHealingMessage{}, err
	}

	if err := task.StoreSelfHealingExecutionMetrics(ctx, types.SelfHealingExecutionMetric{
		TriggerID:       responseMessage.TriggerID,
		ChallengeID:     responseMessage.SelfHealingMessageData.Response.ChallengeID,
		MessageType:     int(responseMessage.MessageType),
		Data:            responseMetricBytes,
		SenderID:        task.nodeID,
		SenderSignature: signature,
	}); err != nil {
		logger.WithError(err).Error("error storing response execution metric")
	}
	logger.Info("response execution metric has been stored")

	return pb.SelfHealingMessage{
		MessageType:     pb.SelfHealingMessageMessageType(responseMessage.MessageType),
		TriggerId:       responseMessage.TriggerID,
		Data:            data,
		SenderId:        responseMessage.SenderID,
		SenderSignature: responseMessage.SenderSignature,
	}, nil
}

// processSelfHealingVerifications simply sends a response msg, to receive verifications
func (task *SHTask) processSelfHealingVerifications(ctx context.Context, nodesToConnect []pastel.MasterNode, responseMsg *pb.SelfHealingMessage) (successfulVerifications map[string]types.SelfHealingMessage, allVerifications map[string]types.SelfHealingMessage) {
	successfulVerifications = make(map[string]types.SelfHealingMessage)
	allVerifications = make(map[string]types.SelfHealingMessage)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, node := range nodesToConnect {
		node := node
		wg.Add(1)

		go func() {
			defer wg.Done()

			verificationMsg, err := task.sendSelfHealingResponseMessage(ctx, responseMsg, node.ExtAddress)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("error sending response message for verification")
				return
			}

			if verificationMsg != nil && verificationMsg.TriggerID != "" &&
				verificationMsg.MessageType == types.SelfHealingVerificationMessage {

				if verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket.IsVerified {
					log.WithContext(ctx).WithField("is_success", true).
						WithField("node_address", node.ExtAddress).
						Info("affirmation message has been received verified")

					mu.Lock()
					successfulVerifications[verificationMsg.SenderID] = *verificationMsg
					mu.Unlock()
				}

				mu.Lock()
				allVerifications[verificationMsg.SenderID] = *verificationMsg
				mu.Unlock()

				log.WithContext(ctx).WithField("node_address", node.ExtAddress).
					Info("affirmation message has been stored")
			}
		}()
	}

	wg.Wait()

	return successfulVerifications, allVerifications
}

// sendSelfHealingResponseMessage establish a connection with the processingSupernodeAddr and sends response msg to
// receive verification upon it
func (task *SHTask) sendSelfHealingResponseMessage(ctx context.Context, challengeMessage *pb.SelfHealingMessage, processingSupernodeAddr string) (*types.SelfHealingMessage, error) {
	log.WithContext(ctx).WithField("trigger_id", challengeMessage.TriggerId).
		Info("Sending self-healing response message to supernode address: " + processingSupernodeAddr)

	//Connect over grpc
	nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
	if err != nil {
		err = fmt.Errorf("Could not use node client to connect to: " + processingSupernodeAddr)
		log.WithContext(ctx).
			WithField("trigger_id", challengeMessage.TriggerId).
			WithField("method", "sendSelfHealingResponseMsg").
			WithField("node_address", processingSupernodeAddr).
			Warn(err.Error())
		return nil, err
	}
	defer nodeClientConn.Close()

	selfHealingChallengeIF := nodeClientConn.SelfHealingChallenge()

	affirmationResult, err := selfHealingChallengeIF.VerifySelfHealingChallenge(ctx, challengeMessage)
	if err != nil {
		log.WithContext(ctx).
			WithField("trigger_id", challengeMessage.TriggerId).
			WithField("method", "sendSelfHealingResponseMsg").
			WithField("node_address", processingSupernodeAddr).
			Warn(err.Error())
		return nil, err
	}

	log.WithContext(ctx).WithField("verifying_node:", processingSupernodeAddr).Info("affirmation result received")
	return &affirmationResult, nil
}

// DownloadDDAndFingerprints gets dd and fp file from ticket based on id and returns the file.
func (task *SHTask) DownloadDDAndFingerprints(ctx context.Context, DDAndFingerprintsIDs []string) (availableDDFPFiles map[string]pastel.DDAndFingerprints, err error) {
	keyToDDFpFileMap := make(map[string]pastel.DDAndFingerprints)
	log.WithContext(ctx).WithField("ids", DDAndFingerprintsIDs).Debug("going to retrieve file from given ids in the ticket")

	for i := 0; i < len(DDAndFingerprintsIDs); i++ {
		key := DDAndFingerprintsIDs[i]

		file, err := task.P2PClient.Retrieve(ctx, DDAndFingerprintsIDs[i], false)
		if err != nil {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails tried to get this file and failed. ")
			continue
		}
		log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Debug("file has been retrieved from the network")

		decompressedData, err := utils.Decompress(file)
		if err != nil {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails self healing - failed to decompress this file. ")
			continue
		}
		log.WithContext(ctx).Debug("file has been decompressed")

		if len(decompressedData) == 0 {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("decompressed data is empty")
			continue
		}

		splits := bytes.Split(decompressedData, []byte{pastel.SeparatorByte})
		if len(splits) < 4 {
			log.WithContext(ctx).WithField("hash", DDAndFingerprintsIDs[i]).
				Info("decompressed data has incorrect number of separator bytes")
			continue
		}

		SN1Sig := splits[1]
		log.WithContext(ctx).WithField("sn_sig_1", SN1Sig).Debug("first signature has been retrieved")

		SN2Sig := splits[2]
		log.WithContext(ctx).WithField("sn_sig_2", SN2Sig).Debug("second signature has been retrieved")

		SN3Sig := splits[3]
		log.WithContext(ctx).WithField("sn_sig_3", SN3Sig).Debug("third signature has been retrieved")

		dataToJSONDecode, err := utils.B64Decode(splits[0])
		if err != nil {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails could not base64 decode. ")
			continue
		}

		if len(dataToJSONDecode) == 0 {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails file contents are empty")
			continue
		}

		ddAndFingerprintFile := &pastel.DDAndFingerprints{}
		err = json.Unmarshal(dataToJSONDecode, ddAndFingerprintFile)
		if err != nil {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails could not JSON unmarshal. ")
			continue
		}

		task.KeySignaturesMap[key] = signatures{
			SNSignature1: SN1Sig,
			SNSignature2: SN2Sig,
			SNSignature3: SN3Sig,
		}

		keyToDDFpFileMap[key] = *ddAndFingerprintFile

		log.WithContext(ctx).WithField("key", key).Info("file has been decoded successfully with signatures")
	}

	return keyToDDFpFileMap, nil
}

func (task *SHTask) getVerificationExecMetrics(responseMsg types.SelfHealingMessage, verifications map[string]types.SelfHealingMessage) (*types.SelfHealingExecutionMetric, error) {
	var vmsgs []types.SelfHealingMessage

	if len(verifications) > 0 {
		for _, verificationMsg := range verifications {
			verificationMsg.SelfHealingMessageData.Challenge = types.SelfHealingChallengeData{}
			verificationMsg.SelfHealingMessageData.Response = types.SelfHealingResponseData{}

			vmsgs = append(vmsgs, verificationMsg)
		}

		vBytes, err := json.Marshal(vmsgs)
		if err != nil {
			return nil, errors.Errorf("error marshaling verifcation msgs")
		}

		sig, err := task.SignBroadcastMessage(context.Background(), vBytes)
		if err != nil {
			return nil, errors.Errorf("Error signing message")
		}

		return &types.SelfHealingExecutionMetric{
			TriggerID:       responseMsg.TriggerID,
			ChallengeID:     responseMsg.SelfHealingMessageData.Response.ChallengeID,
			MessageType:     int(types.SelfHealingVerificationMessage),
			SenderID:        task.nodeID,
			SenderSignature: sig,
			Data:            vBytes,
		}, nil
	}

	verificationMsg := types.SelfHealingMessage{
		TriggerID:   responseMsg.TriggerID,
		MessageType: types.SelfHealingVerificationMessage,
		SelfHealingMessageData: types.SelfHealingMessageData{
			ChallengerID: responseMsg.SelfHealingMessageData.ChallengerID,
			RecipientID:  responseMsg.SelfHealingMessageData.RecipientID,
			Verification: types.SelfHealingVerificationData{
				ChallengeID: responseMsg.SelfHealingMessageData.Response.ChallengeID,
				Block:       responseMsg.SelfHealingMessageData.Response.Block,
				Merkelroot:  responseMsg.SelfHealingMessageData.Response.Merkelroot,
				Timestamp:   responseMsg.SelfHealingMessageData.Response.Timestamp,
				VerifiedTicket: types.VerifiedTicket{
					IsVerified: false,
					Message:    "no successful verifications from verifiers",
				},
			},
		},
		SenderID:        responseMsg.SenderID,
		SenderSignature: responseMsg.SenderSignature,
	}

	vmsgs = append(vmsgs, verificationMsg)

	vBytes, err := json.Marshal(vmsgs)
	if err != nil {
		return nil, errors.Errorf("error marshaling verification msg")
	}

	sig, err := task.SignBroadcastMessage(context.Background(), vBytes)
	if err != nil {
		return nil, errors.Errorf("Error signing message")
	}

	return &types.SelfHealingExecutionMetric{
		TriggerID:       verificationMsg.TriggerID,
		ChallengeID:     verificationMsg.SelfHealingMessageData.Verification.ChallengeID,
		MessageType:     int(verificationMsg.MessageType),
		SenderID:        task.nodeID,
		SenderSignature: sig,
		Data:            vBytes,
	}, nil
}
