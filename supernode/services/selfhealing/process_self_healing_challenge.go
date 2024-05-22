package selfhealing

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/btcsuite/btcutil/base58"
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
	verifyThreshold = 3
)

// ProcessSelfHealingChallenge is called from service.Run to process the stored events.
func (task *SHTask) ProcessSelfHealingChallenge(ctx context.Context, event types.SelfHealingChallengeEvent) error {
	logger := log.WithContext(ctx).WithField("trigger_id", event.TriggerID).WithField("ticket_txid", event.TicketID)
	logger.Info("ProcessSelfHealingChallenge has been invoked")

	log.WithContext(ctx).Debug("retrieving block no and verbose")
	currentBlockCount, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		logger.WithError(err).Error("could not get current block count")
		return errors.Errorf("could not get current block count")
	}
	blkVerbose1, err := task.SuperNodeService.PastelClient.GetBlockVerbose1(ctx, currentBlockCount)
	if err != nil {
		logger.WithError(err).Error("could not get current block verbose")
		return errors.Errorf("could not get current block verbose")
	}
	merkleroot := blkVerbose1.MerkleRoot
	logger.Debug("block count & merkelroot has been retrieved")

	responseMsg := types.SelfHealingMessage{
		TriggerID:   event.TriggerID,
		SenderID:    task.nodeID,
		MessageType: types.SelfHealingResponseMessage,
		SelfHealingMessageData: types.SelfHealingMessageData{
			Response: types.SelfHealingResponseData{
				ChallengeID: event.ChallengeID,
				Block:       currentBlockCount,
				Merkelroot:  merkleroot,
			},
		},
	}

	var ticket types.ChallengeTicket
	if err := json.Unmarshal(event.Data, &ticket); err != nil {
		logger.WithError(err).Error("error un-marshaling the challenge event data")

		return errors.Errorf("error un-marshaling challenge event data")
	}

	var (
		nftTicket     *pastel.NFTTicket
		cascadeTicket *pastel.APICascadeTicket
		senseTicket   *pastel.APISenseTicket
	)

	logger.Debug("going to initiate checking process")
	nftTicket, cascadeTicket, senseTicket, err = task.getTicket(ctx, ticket.TxID, TicketType(ticket.TicketType))
	if err != nil {
		logger.WithError(err).Error("Error getRelevantTicketFromMsg")
		return errors.Errorf("error retrieving the ticket:%s", err.Error())
	}
	logger.Debug("ticket has been retrieved")

	if nftTicket != nil || cascadeTicket != nil {
		newCtx := context.Background()
		if !task.isReconstructionRequired(newCtx, ticket.MissingKeys) {
			logger.Debug("reconstruction is not required for ticket")

			responseMsg.SelfHealingMessageData.Response.RespondedTicket = types.RespondedTicket{
				TxID:                     ticket.TxID,
				TicketType:               ticket.TicketType,
				MissingKeys:              ticket.MissingKeys,
				ReconstructedFileHash:    nil,
				IsReconstructionRequired: false,
			}
			logger.Debug("sending response for verification")

			newCtx = context.Background()
			verifications, err := task.getVerifications(newCtx, responseMsg)
			if err != nil {
				logger.WithError(err).Debug("unable to get verifications")
				return err
			}
			logger.WithField("verifications", len(verifications)).
				Debug("verifications have been evaluated")

			if len(verifications) >= verifyThreshold {
				logger.Info("self-healing not required for this ticket, verified by verifiers")
			}

			return nil
		}

		logger.Debug("going to initiate self-healing")

		newCtx = context.Background()
		raptorQSymbols, reconstructedFileHash, err := task.selfHealing(newCtx, ticket.TxID, nftTicket, cascadeTicket)
		if err != nil {
			logger.WithError(err).Error("Error self-healing the file")
			return err
		}
		logger.Debug("file has been reconstructed")
		task.StorageHandler.TxID = ticket.TxID
		responseMsg.SelfHealingMessageData.Response.RespondedTicket = types.RespondedTicket{
			TxID:                     ticket.TxID,
			TicketType:               ticket.TicketType,
			MissingKeys:              ticket.MissingKeys,
			ReconstructedFileHash:    reconstructedFileHash,
			IsReconstructionRequired: true,
		}
		logger.Debug("sending response for verification")

		newCtx = context.Background()
		verifications, err := task.getVerifications(newCtx, responseMsg)
		if err != nil {
			logger.WithError(err).Debug("unable to get verifications")
			return err
		}
		logger.WithField("verifications", len(verifications)).Debug("verifications have been evaluated")

		if len(verifications) >= verifyThreshold {
			newCtx = context.Background()
			if err := task.StorageHandler.StoreRaptorQSymbolsIntoP2P(newCtx, raptorQSymbols, ""); err != nil {
				logger.WithError(err).Error("Error storing symbols to P2P")
				return err
			}
			logger.Debug("raptor q symbols have been stored")

			sig, _, err := task.SignMessage(context.Background(), responseMsg.SelfHealingMessageData)
			if err != nil {
				logger.WithError(err).Error("error signing self-healing completion metric")
				return err
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
				return err
			}

			if err := task.StoreSelfHealingExecutionMetrics(ctx, types.SelfHealingExecutionMetric{
				TriggerID:       responseMsg.TriggerID,
				ChallengeID:     responseMsg.SelfHealingMessageData.Response.ChallengeID,
				MessageType:     int(types.SelfHealingCompletionMessage),
				SenderID:        task.nodeID,
				Data:            completionMsgBytes,
				SenderSignature: sig,
			}); err != nil {
				logger.WithError(err).Error("error storing self-healing completion execution metric")
				return err
			}

			return nil
		}

		logger.Debug("verify threshold does not meet, not saving symbols")
	} else if senseTicket != nil {
		newCtx := context.Background()
		reqSelfHealing, sortedFiles := task.senseCheckingProcess(newCtx, senseTicket.DDAndFingerprintsIDs)
		if !reqSelfHealing {
			logger.Debug("self-healing not required for sense ticket")
			responseMsg.SelfHealingMessageData.Response.RespondedTicket = types.RespondedTicket{
				TxID:                     ticket.TxID,
				TicketType:               ticket.TicketType,
				MissingKeys:              ticket.MissingKeys,
				ReconstructedFileHash:    nil,
				IsReconstructionRequired: false,
			}
			logger.Debug("sending response for verification")

			newCtx = context.Background()
			verifications, err := task.getVerifications(newCtx, responseMsg)
			if err != nil {
				logger.WithError(err).Debug("unable to get verifications")
				return err
			}
			logger.WithField("verifications", len(verifications)).
				Debug("verifications have been evaluated")

			if len(verifications) >= verifyThreshold {
				logger.Info("self-healing not required for this sense ticket, verified by verifiers")
			}

			return nil
		}
		logger.Debug("going to initiate sense self-healing")

		_, idFiles, fileHash, err := task.senseSelfHealing(ctx, senseTicket, sortedFiles)
		if err != nil {
			logger.WithError(err).Error("self-healing failed for sense ticket")
			return err
		}
		logger.Debug("sense ticket has been reconstructed")

		responseMsg.SelfHealingMessageData.Response.RespondedTicket = types.RespondedTicket{
			TxID:                     ticket.TxID,
			TicketType:               ticket.TicketType,
			MissingKeys:              ticket.MissingKeys,
			ReconstructedFileHash:    fileHash,
			IsReconstructionRequired: true,
		}
		logger.Debug("sending response for verification")

		newCtx = context.Background()
		verifications, err := task.getVerifications(newCtx, responseMsg)
		if err != nil {
			logger.WithError(err).Debug("unable to get verifications")
			return err
		}
		logger.WithField("verifications", len(verifications)).Debug("verifications have been evaluated")

		if len(verifications) >= verifyThreshold {
			newCtx = context.Background()
			if err := task.StorageHandler.StoreBatch(newCtx, idFiles, common.P2PDataDDMetadata); err != nil {
				logger.WithError(err).Error("Error storing symbols to P2P")
				return err
			}
			logger.Debug("dd & fp ids have been stored")

			sig, _, err := task.SignMessage(context.Background(), responseMsg.SelfHealingMessageData)
			if err != nil {
				logger.WithError(err).Error("error storing completion execution metric")
				return err
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
				return err
			}

			if err := task.StoreSelfHealingExecutionMetrics(ctx, types.SelfHealingExecutionMetric{
				TriggerID:       responseMsg.TriggerID,
				ChallengeID:     responseMsg.SelfHealingMessageData.Response.ChallengeID,
				MessageType:     int(types.SelfHealingCompletionMessage),
				SenderID:        task.nodeID,
				Data:            completionMsgBytes,
				SenderSignature: sig,
			}); err != nil {
				logger.WithError(err).Error("error storing self-healing completion execution metric")
				return err
			}
		}
	}

	logger.Info("self-healing processing completed for the given ticket")

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
	logger := log.WithContext(ctx).WithField("SymbolIDsFileId", fileHash)

	decoded := base58.Decode(fileHash)
	if len(decoded) != 256/8 {
		logger.Error("error decoding file hash to db key")
		return false
	}

	dbKey := hex.EncodeToString(decoded)

	logger.WithField("db_key", dbKey).Debug("checking file in p2p db")
	rqIDsData, err := task.P2PClient.Retrieve(ctx, fileHash, false)
	if err != nil {
		logger.WithField("db_key", dbKey).
			WithError(err).Error("Retrieve compressed symbol IDs file from P2P failed")
		return true
	}

	if rqIDsData == nil {
		logger.WithField("db_key", dbKey).Debug("Retrieved rq-data is empty")
		return true
	}

	if len(rqIDsData) == 0 {
		logger.WithField("db_key", dbKey).Debug("Retrieved rq-data is empty")
		return true
	}

	return false
}

func (task *SHTask) selfHealing(ctx context.Context, txid string, regTicket *pastel.NFTTicket, cascadeTicket *pastel.APICascadeTicket) (file, reconstructedFileHash []byte, err error) {
	log.WithContext(ctx).WithField("txid", txid).Debug("Self-healing initiated")
	if regTicket != nil {
		log.WithContext(ctx).WithField("txid", txid).WithField("reg_ticket", regTicket).Debug("going to self heal reg ticket")

		file, err = task.downloadTask.RestoreFile(ctx, regTicket.AppTicketData.RQIDs, regTicket.AppTicketData.RQOti, regTicket.AppTicketData.DataHash, txid)
		if err != nil {
			log.WithContext(ctx).WithField("txid", txid).WithError(err).Error("Unable to restore file")
			return nil, nil, errors.Errorf("unable to restore file")
		}

		fileHash := sha3.Sum256(file)
		log.WithContext(ctx).WithField("txid", txid).Debug("file has been restored")
		reconstructedFileHash = fileHash[:]
	} else if cascadeTicket != nil {
		log.WithContext(ctx).WithField("txid", txid).WithField("cascade_ticket", cascadeTicket).Debug("going to self heal cascade ticket")

		file, err = task.downloadTask.RestoreFile(ctx, cascadeTicket.RQIDs, cascadeTicket.RQOti, cascadeTicket.DataHash, txid)
		if err != nil {
			log.WithContext(ctx).WithField("txid", txid).WithError(err).Error("Unable to restore file")
			return nil, nil, errors.Errorf("unable to restore file")
		}

		fileHash := sha3.Sum256(file)
		log.WithContext(ctx).WithField("txid", txid).Debug("file has been restored")
		reconstructedFileHash = fileHash[:]
	} else {
		log.WithContext(ctx).WithField("txid", txid).Debug("not a valid ticket")

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

func (task *SHTask) senseCheckingProcess(ctx context.Context, ddFPIDs []string) (requiredReconstruction bool, filesSortedOnHash [][]byte) {
	availableDDFPFiles := task.DownloadDDAndFingerprints(ctx, ddFPIDs)
	if len(ddFPIDs) == len(availableDDFPFiles) {
		log.WithContext(ctx).WithField("ddfp_ids_count", len(ddFPIDs)).WithField("file_count", len(availableDDFPFiles)).
			Debug("no of ids and retrieved files from the network are equal, self-healing not required for this sense ticket")

		return false, nil
	}

	log.WithContext(ctx).WithField("ddfp_ids_count", len(ddFPIDs)).WithField("file_count", len(availableDDFPFiles)).
		Debug("no of ids and retrieved files from the network are not equal, self-healing required")

	hashTally := make(map[string]int)
	type hashGroup struct {
		Hash  string
		Files [][]byte
	}

	hashGroupsMap := make(map[string][][]byte)
	for _, DDFPIdFile := range availableDDFPFiles {
		fileHash := utils.GetHashStringFromBytes(DDFPIdFile)

		//Generating hashes, maintaining hash tally and fileHash map to find the file with the most common hash
		hashTally[fileHash]++
		hashGroupsMap[fileHash] = append(hashGroupsMap[fileHash], DDFPIdFile)
	}

	var hashGroups []hashGroup
	for hash, files := range hashGroupsMap {
		hashGroups = append(hashGroups, hashGroup{Hash: hash, Files: files})
	}

	sort.Slice(hashGroups, func(i, j int) bool {
		return hashTally[hashGroups[i].Hash] > hashTally[hashGroups[j].Hash]
	})

	var sortedFiles [][]byte
	for _, group := range hashGroups {
		for _, file := range group.Files {
			sortedFiles = append(sortedFiles, file)
		}
	}

	return true, sortedFiles
}

func (task *SHTask) senseSelfHealing(ctx context.Context, senseTicket *pastel.APISenseTicket, sortedFiles [][]byte) (ids []string, idFiles [][]byte, fileHash []byte, err error) {
	if sortedFiles == nil {
		return nil, nil, nil, errors.Errorf("empty list of sorted files")
	}

	if len(sortedFiles) == 0 {
		return nil, nil, nil, errors.Errorf("empty list of sorted files")
	}

	var (
		ddAndFingerprintFile pastel.DDAndFingerprints
		sig                  signatures
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

		SN1Sig := splits[1]
		log.WithContext(ctx).WithField("sn_sig_1", SN1Sig).Debug("first signature has been retrieved")

		SN2Sig := splits[2]
		log.WithContext(ctx).WithField("sn_sig_2", SN2Sig).Debug("second signature has been retrieved")

		SN3Sig := splits[3]
		log.WithContext(ctx).WithField("sn_sig_3", SN3Sig).Debug("third signature has been retrieved")

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

		sig = signatures{
			SNSignature1: SN1Sig,
			SNSignature2: SN2Sig,
			SNSignature3: SN3Sig,
		}
		log.WithContext(ctx).Debug("signatures retrieved from the most common file")

		data, err = json.Marshal(ddAndFingerprintFile)
		if err != nil {
			continue
		}

		fileHash, err = utils.Sha3256hash(data)
		if err != nil {
			continue
		}

		log.WithContext(ctx).Debug("sense file has been reconstructed")
		break
	}

	ddEncoded := utils.B64Encode(data)

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
		return nil, nil, nil, errors.Errorf("get ID Files: %w", err)
	}

	return ids, idFiles, fileHash, nil
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
		return successfulVerifications, err
	}

	if verificationExecMetrics != nil {
		err := task.StoreSelfHealingExecutionMetrics(ctx, *verificationExecMetrics)
		if err != nil {
			log.WithContext(ctx).WithField("trigger_id", msg.TriggerID).
				WithField("txid", msg.SelfHealingMessageData.Response.RespondedTicket.TxID).
				WithField("challenge_id", msg.SelfHealingMessageData.Response.ChallengeID).
				Error("error storing verification execution metrics")

			return nil, err
		}
	}
	log.WithContext(ctx).WithField("trigger_id", msg.TriggerID).
		WithField("txid", msg.SelfHealingMessageData.Response.RespondedTicket.TxID).
		WithField("challenge_id", msg.SelfHealingMessageData.Response.ChallengeID).
		Debug("verification execution metrics have been stored")

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
		log.WithContext(ctx).Debug(fmt.Sprintf("sliceOfSupernodesClosestToPreviousBlockHash:%s", sliceOfSupernodesClosestToPreviousBlockHash))

		for _, verifier := range sliceOfSupernodesClosestToPreviousBlockHash {
			nodesToConnect = append(nodesToConnect, mapSupernodes[verifier])
		}

		logger.WithField("nodes_to_connect", nodesToConnect).Debug("nodes to send self-healing response msg have been selected")
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
		return pb.SelfHealingMessage{}, err
	}
	logger.Debug("response execution metric has been stored")

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
				log.WithError(err).Debug("error sending response message for verification")
				return
			}

			if verificationMsg != nil && verificationMsg.TriggerID != "" &&
				verificationMsg.MessageType == types.SelfHealingVerificationMessage {

				if verificationMsg.SelfHealingMessageData.Verification.VerifiedTicket.IsVerified {
					log.WithContext(ctx).WithField("is_success", true).
						WithField("node_address", node.ExtAddress).
						Debug("affirmation message has been received verified")

					mu.Lock()
					successfulVerifications[verificationMsg.SenderID] = *verificationMsg
					mu.Unlock()
				}

				mu.Lock()
				allVerifications[verificationMsg.SenderID] = *verificationMsg
				mu.Unlock()

				log.WithContext(ctx).WithField("node_address", node.ExtAddress).
					Debug("affirmation message has been stored")
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
		Debug("Sending self-healing response message to supernode address: " + processingSupernodeAddr)

	//Connect over grpc
	nodeClientConn, err := task.nodeClient.ConnectSN(ctx, processingSupernodeAddr)
	if err != nil {
		logError(ctx, "SelfHealingResponse", err)
		return nil, fmt.Errorf("Could not use node client to connect to: " + processingSupernodeAddr)
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

	log.WithContext(ctx).WithField("verifying_node:", processingSupernodeAddr).Debug("affirmation result received")
	return &affirmationResult, nil
}

// DownloadDDAndFingerprints gets dd and fp file from ticket based on id and returns the file.
func (task *SHTask) DownloadDDAndFingerprints(ctx context.Context, DDAndFingerprintsIDs []string) (availableDDFPFiles map[string][]byte) {
	keyToDDFpFileMap := make(map[string][]byte)
	log.WithContext(ctx).WithField("ids", DDAndFingerprintsIDs).Debug("going to retrieve file from given ids in the ticket")

	for i := 0; i < len(DDAndFingerprintsIDs); i++ {
		key := DDAndFingerprintsIDs[i]

		compressedFile, err := task.P2PClient.Retrieve(ctx, DDAndFingerprintsIDs[i], false)
		if err != nil {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails tried to get this file and failed. ")
			continue
		}
		log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Debug("file has been retrieved from the network")

		if len(compressedFile) == 0 {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Debug("compressed data is empty")
			continue
		}

		keyToDDFpFileMap[key] = compressedFile

		log.WithContext(ctx).WithField("key", key).Debug("file has been retrieved successfully")
	}

	return keyToDDFpFileMap
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
