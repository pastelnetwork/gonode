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

// ProcessSelfHealingChallenge is called from grpc server, which processes the self-healing challenge,
// and will execute the reconstruction work.
func (task *SHTask) ProcessSelfHealingChallenge(ctx context.Context, incomingChallengeMessage types.SelfHealingMessage) error {
	logger := log.WithContext(ctx).WithField("challenge_id", incomingChallengeMessage.ChallengeID)

	logger.WithField("challenge", incomingChallengeMessage).Info("ProcessSelfHealingChallenge has been invoked")

	// wait if test env so tests can mock the p2p
	if os.Getenv("INTEGRATION_TEST_ENV") == "true" {
		time.Sleep(10 * time.Second)
	}

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
		ChallengeID: incomingChallengeMessage.ChallengeID,
		SenderID:    task.nodeID,
		MessageType: types.SelfHealingResponseMessage,
		SelfHealingMessageData: types.SelfHealingMessageData{
			ChallengerID: incomingChallengeMessage.SelfHealingMessageData.ChallengerID,
			Challenge: types.SelfHealingChallengeData{
				Block:            incomingChallengeMessage.SelfHealingMessageData.Challenge.Block,
				Merkelroot:       incomingChallengeMessage.SelfHealingMessageData.Challenge.Merkelroot,
				ChallengeTickets: incomingChallengeMessage.SelfHealingMessageData.Challenge.ChallengeTickets,
				Timestamp:        incomingChallengeMessage.SelfHealingMessageData.Challenge.Timestamp,
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

	for _, ticket := range challengeTickets {
		if task.SHService.ticketsMap[ticket.TxID] {
			log.WithContext(ctx).WithField("txid", ticket.TxID).Info("ticket processing is already in progress")
			continue
		}
		task.SHService.ticketsMap[ticket.TxID] = true

		logger.WithField("ticket_txid", ticket.TxID).Info("checking process going to initiate")

		nftTicket, cascadeTicket, senseTicket, err = task.getTicket(ctx, ticket.TxID, TicketType(ticket.TicketType))
		if err != nil {
			logger.WithError(err).Error("Error getRelevantTicketFromMsg")
		}
		logger.WithField("challenge_id", incomingChallengeMessage.ChallengeID).Info("reg ticket has been retrieved")

		if nftTicket != nil || cascadeTicket != nil {
			if !task.isReconstructionRequired(ctx, ticket.MissingKeys) {
				logger.WithField("txid", ticket.TxID).Info("Reconstruction is not required for ticket")

				responseMsg.SelfHealingMessageData.Response.RespondedTickets = append(responseMsg.SelfHealingMessageData.Response.RespondedTickets,
					types.RespondedTicket{
						TxID:                     ticket.TxID,
						TicketType:               ticket.TicketType,
						MissingKeys:              ticket.MissingKeys,
						ReconstructedFileHash:    nil,
						RaptorQSymbols:           nil,
						IsReconstructionRequired: false,
					})
				continue
			}
			logger.WithField("ticket_txid", ticket.TxID).Info("self-healing required for this ticket")

			file, reconstructedFileHash, err := task.selfHealing(ctx, ticket.TxID, nftTicket, cascadeTicket)
			if err != nil {
				logger.WithError(err).Error("Error self-healing the file")
				continue
			}
			task.RaptorQSymbols = file
			logger.WithField("ticket_txid", ticket.TxID).Info("file has been reconstructed")

			if err := task.StorageHandler.StoreRaptorQSymbolsIntoP2P(ctx, task.RaptorQSymbols, ""); err != nil {
				logger.WithError(err).Error("Error storing symbols to P2P")
				continue
			}
			logger.WithField("ticket_txid", ticket.TxID).Info("raptor q symbols have been stored")

			responseMsg.SelfHealingMessageData.Response.RespondedTickets = append(responseMsg.SelfHealingMessageData.Response.RespondedTickets,
				types.RespondedTicket{
					TxID:                     ticket.TxID,
					TicketType:               ticket.TicketType,
					MissingKeys:              ticket.MissingKeys,
					ReconstructedFileHash:    reconstructedFileHash,
					RaptorQSymbols:           task.RaptorQSymbols,
					IsReconstructionRequired: true,
				})
		} else if senseTicket != nil {
			reqSelfHealing, mostCommonFile := task.senseCheckingProcess(ctx, senseTicket.DDAndFingerprintsIDs)
			if err != nil {
				logger.WithField("ticket_txid", ticket.TxID).WithError(err).Error("Error in checking process for sense action ticket")
				continue
			}

			if !reqSelfHealing {
				logger.WithField("ticket_txid", ticket.TxID).WithError(err).Error("self-healing not required for sense action ticket")
				responseMsg.SelfHealingMessageData.Response.RespondedTickets = append(responseMsg.SelfHealingMessageData.Response.RespondedTickets,
					types.RespondedTicket{
						TxID:                     ticket.TxID,
						TicketType:               ticket.TicketType,
						MissingKeys:              ticket.MissingKeys,
						ReconstructedFileHash:    nil,
						RaptorQSymbols:           nil,
						FileIDs:                  nil,
						IsReconstructionRequired: false,
					})

				continue
			}

			logger.WithField("ticket_txid", ticket.TxID).Info("self-healing required for sense ticket")
			ids, idFiles, err := task.senseSelfHealing(ctx, senseTicket, mostCommonFile)
			if err != nil {
				logger.WithField("ticket_txid", ticket.TxID).WithError(err).Error("self-healing not required for sense action ticket")
				continue
			}
			task.IDFiles = idFiles
			logger.WithField("ticket_txid", ticket.TxID).Info("sense ticket has been healed")

			if err := task.StorageHandler.StoreBatch(ctx, task.IDFiles, common.P2PDataDDMetadata); err != nil {
				logger.WithError(err).Error("Error storing id files to P2P")
				continue
			}
			logger.WithField("ticket_txid", ticket.TxID).Info("DD & FP ids have been stored")

			fileBytes, err := json.Marshal(mostCommonFile)
			if err != nil {
				logger.WithError(err).Error(err)
			}

			fileHash, err := utils.Sha3256hash(fileBytes)
			if err != nil {
				logger.WithError(err).Error(err)
			}

			responseMsg.SelfHealingMessageData.Response.RespondedTickets = append(responseMsg.SelfHealingMessageData.Response.RespondedTickets,
				types.RespondedTicket{
					TxID:                     ticket.TxID,
					TicketType:               ticket.TicketType,
					MissingKeys:              ticket.MissingKeys,
					ReconstructedFileHash:    fileHash,
					FileIDs:                  ids,
					IsReconstructionRequired: true,
				})
		}
	}

	verifications, err := task.getVerifications(ctx, responseMsg)
	if err != nil {
		log.WithContext(ctx).WithField("challenge_id", incomingChallengeMessage.ChallengeID).WithError(err)
	}
	logger.WithField("verifications", len(verifications)).Info("verifications have been received")

	logger.Info("self-healing completed")

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

func (task *SHTask) senseCheckingProcess(ctx context.Context, ddFPIDs []string) (requiredReconstruction bool, mostCommonFile *pastel.DDAndFingerprints) {
	//download all the DD and FPIds
	availableDDFPFiles, err := task.DownloadDDAndFingerprints(ctx, ddFPIDs)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Downloading DD and Fingerprints")
	}

	//if not able to retrieve all should self-heal
	if len(ddFPIDs) == len(availableDDFPFiles) {
		return false, nil
	}

	//Find most common file (source of truth) for self-healing
	hashTally := make(map[string]int)
	fileHashMap := make(map[string]*pastel.DDAndFingerprints)
	for _, DDFPIdFile := range availableDDFPFiles {
		bytes, err := json.Marshal(DDFPIdFile)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to marshal DDFPIdFile")
			continue
		}

		fileHash := utils.GetHashStringFromBytes(bytes)

		//Generating hashes, maintaining hash tally and fileHash map to find the file with the most common hash
		hashTally[fileHash]++
		fileHashMap[fileHash] = &DDFPIdFile
	}

	var mostCommonHash string
	var greatestTally int
	for hash, tally := range hashTally {
		if tally > greatestTally {
			mostCommonHash = hash
			greatestTally = tally
		}
	}

	return true, fileHashMap[mostCommonHash]
}

func (task *SHTask) senseSelfHealing(ctx context.Context, senseTicket *pastel.APISenseTicket, ddFPID *pastel.DDAndFingerprints) (ids []string, idFiles [][]byte, err error) {
	if len(task.SNsSignatures) != 3 {
		return nil, nil, errors.Errorf("wrong number of signature for fingerprints - %d", len(task.SNsSignatures))
	}

	ddDataJSON, err := json.Marshal(ddFPID)
	if err != nil {
		return nil, nil, errors.Errorf("failed to marshal dd-data: %w", err)
	}

	ddEncoded := utils.B64Encode(ddDataJSON)

	var buffer bytes.Buffer
	buffer.Write(ddEncoded)
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(task.SNsSignatures[0])
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(task.SNsSignatures[1])
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(task.SNsSignatures[2])
	ddFpFile := buffer.Bytes()

	ids, idFiles, err = pastel.GetIDFiles(ctx, ddFpFile, senseTicket.DDAndFingerprintsIc, senseTicket.DDAndFingerprintsMax)
	if err != nil {
		return nil, nil, errors.Errorf("get ID Files: %w", err)
	}

	return ids, idFiles, nil
}

func (task *SHTask) getVerifications(ctx context.Context, msg types.SelfHealingMessage) (map[string]types.SelfHealingMessage, error) {
	nodesToConnect, err := task.GetNodesAddressesToConnect(ctx, msg)
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

	return task.processSelfHealingVerifications(ctx, nodesToConnect, &responseMsg), nil
}

// GetNodesAddressesToConnect basically retrieves the masternode address against the pastel-id from the list and return that
func (task *SHTask) GetNodesAddressesToConnect(ctx context.Context, challengeMessage types.SelfHealingMessage) ([]pastel.MasterNode, error) {
	var nodesToConnect []pastel.MasterNode
	supernodes, err := task.SuperNodeService.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeID).WithField("method", "GetNodesAddressesToConnect").WithError(err).Warn("could not get Supernode extra: ", err.Error())
		return nil, err
	}

	mapSupernodes := make(map[string]pastel.MasterNode)
	for _, mn := range supernodes {
		if mn.ExtAddress == "" || mn.ExtKey == "" {
			log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeID).WithField("method", "GetNodesAddressesToConnect").
				WithField("node_id", mn.ExtKey).Warn("node address or node id is empty")

			continue
		}

		mapSupernodes[mn.ExtKey] = mn
	}

	logger := log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeID).
		WithField("message_type", challengeMessage.MessageType).WithField("node_id", task.nodeID)

	listOfSupernodes, err := task.getListOfOnlineSupernodes(ctx)
	if err != nil {
		logger.WithError(err).Error("error getting list of online supernodes")
	}

	switch challengeMessage.MessageType {
	case types.SelfHealingResponseMessage:

		//Finding 5 closest nodes to previous block hash
		sliceOfSupernodesClosestToPreviousBlockHash := task.GetNClosestSupernodeIDsToComparisonString(ctx, 5, challengeMessage.SelfHealingMessageData.Challenge.Merkelroot, listOfSupernodes, task.nodeID)
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
	signature, data, err := task.SignMessage(ctx, responseMessage.SelfHealingMessageData)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error signing the response message")
		return pb.SelfHealingMessage{}, err
	}
	responseMessage.SenderSignature = signature

	if err := task.StoreSelfHealingMessage(ctx, responseMessage); err != nil {
		log.WithContext(ctx).WithError(err).Error("error storing evaluation report message")
		return pb.SelfHealingMessage{}, err
	}

	return pb.SelfHealingMessage{
		MessageType:     pb.SelfHealingMessageMessageType(responseMessage.MessageType),
		ChallengeId:     responseMessage.ChallengeID,
		Data:            data,
		SenderId:        responseMessage.SenderID,
		SenderSignature: responseMessage.SenderSignature,
	}, nil
}

// processSelfHealingVerifications simply sends a response msg, to receive verifications
func (task *SHTask) processSelfHealingVerifications(ctx context.Context, nodesToConnect []pastel.MasterNode, responseMsg *pb.SelfHealingMessage) (verifications map[string]types.SelfHealingMessage) {
	verifications = make(map[string]types.SelfHealingMessage)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, node := range nodesToConnect {
		node := node
		log.WithContext(ctx).WithField("node_id", node.ExtAddress).Info("processing sn address")
		wg.Add(1)

		go func() {
			defer wg.Done()

			verificationMsg, err := task.sendSelfHealingResponseMessage(ctx, responseMsg, node.ExtAddress)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("error sending evaluation message for processing")
				return
			}

			if verificationMsg != nil && verificationMsg.ChallengeID != "" &&
				verificationMsg.MessageType == types.SelfHealingVerificationMessage {
				//isSuccess, err := task.isVerificationSuccessful(ctx, verificationMsg)
				//if err != nil {
				//	log.WithContext(ctx).WithField("node_id", node.ExtKey).WithError(err).
				//		Error("unsuccessful affirmation by observer")
				//}
				//log.WithContext(ctx).WithField("is_success", isSuccess).Info("affirmation message has been verified")

				mu.Lock()
				verifications[verificationMsg.SenderID] = *verificationMsg
				mu.Unlock()

				//if err := task.StoreChallengeMessage(ctx, *affirmationResponse); err != nil {
				//	log.WithContext(ctx).WithField("node_id", node.ExtKey).WithError(err).
				//		Error("error storing affirmation response from observer")
				//}
			}
		}()
	}

	wg.Wait()

	return verifications
}

// sendSelfHealingResponseMessage establish a connection with the processingSupernodeAddr and sends response msg to
// receive verification upon it
func (task *SHTask) sendSelfHealingResponseMessage(ctx context.Context, challengeMessage *pb.SelfHealingMessage, processingSupernodeAddr string) (*types.SelfHealingMessage, error) {
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("Sending self-healing response message to supernode address: " + processingSupernodeAddr)

	//Connect over grpc
	nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
	if err != nil {
		err = fmt.Errorf("Could not use node client to connect to: " + processingSupernodeAddr)
		log.WithContext(ctx).
			WithField("challengeID", challengeMessage.ChallengeId).
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
			WithField("challengeID", challengeMessage.ChallengeId).
			WithField("method", "sendSelfHealingResponseMsg").
			WithField("node_address", processingSupernodeAddr).
			Warn(err.Error())
		return nil, err
	}

	log.WithContext(ctx).WithField("affirmation:", affirmationResult).Info("affirmation result received")
	return &affirmationResult, nil
}

// DownloadDDAndFingerprints gets dd and fp file from ticket based on id and returns the file.
func (task *SHTask) DownloadDDAndFingerprints(ctx context.Context, DDAndFingerprintsIDs []string) (availableDDFPFiles []pastel.DDAndFingerprints, err error) {
	for i := 0; i < len(DDAndFingerprintsIDs); i++ {
		file, err := task.P2PClient.Retrieve(ctx, DDAndFingerprintsIDs[i])
		if err != nil {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails tried to get this file and failed. ")
			continue
		}
		log.WithContext(ctx).WithField("file", file).Println("Got the file")
		decompressedData, err := utils.Decompress(file)
		if err != nil {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails self healing - failed to decompress this file. ")
			continue
		}
		log.WithContext(ctx).Println("Decompressed the file")
		//base64 dataset doesn't contain periods, so we just find the first index of period and chop it and everything else
		firstIndexOfSeparator := bytes.IndexByte(decompressedData, pastel.SeparatorByte)
		if firstIndexOfSeparator < 1 {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails got a bad separator index. ")
			continue
		}
		secondIndexOfSeparator := bytes.IndexByte(decompressedData[firstIndexOfSeparator+1:], pastel.SeparatorByte)
		task.SNsSignatures = append(task.SNsSignatures, decompressedData[firstIndexOfSeparator+1:secondIndexOfSeparator])

		thirdIndexOfSeparator := bytes.IndexByte(decompressedData[secondIndexOfSeparator+1:], pastel.SeparatorByte)
		task.SNsSignatures = append(task.SNsSignatures, decompressedData[secondIndexOfSeparator+1:thirdIndexOfSeparator])

		fourthIndexOfSeparator := bytes.IndexByte(decompressedData[thirdIndexOfSeparator+1:], pastel.SeparatorByte)
		task.SNsSignatures = append(task.SNsSignatures, decompressedData[thirdIndexOfSeparator+1:fourthIndexOfSeparator])

		dataToBase64Decode := decompressedData[:firstIndexOfSeparator]
		dataToJSONDecode, err := utils.B64Decode(dataToBase64Decode)
		if err != nil {
			log.WithContext(ctx).WithField("Hash", DDAndFingerprintsIDs[i]).Warn("DDAndFingerPrintDetails could not base64 decode. ")
			continue
		}
		log.WithContext(ctx).Println("base64 decoded the file")

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

		availableDDFPFiles = append(availableDDFPFiles, *ddAndFingerprintFile)
	}

	return availableDDFPFiles, nil
}
