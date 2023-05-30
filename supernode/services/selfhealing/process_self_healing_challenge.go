package selfhealing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/DataDog/zstd"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"golang.org/x/crypto/sha3"
)

// ProcessSelfHealingChallenge is called from grpc server, which processes the self-healing challenge,
// and will execute the reconstruction work.
func (task *SHTask) ProcessSelfHealingChallenge(ctx context.Context, challengeMessage *pb.SelfHealingData) error {
	// wait if test env so tests can mock the p2p
	if os.Getenv("INTEGRATION_TEST_ENV") == "true" {
		time.Sleep(10 * time.Second)
	}

	regTicketKeys, actionTicketKeys, err := task.MapSymbolFileKeysFromNFTAndActionTickets(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not generate the map of symbol file keys")
		return err
	}
	log.WithContext(ctx).Info("symbol file keys with their corresponding registered nft tickets ID have been retrieved")

	log.WithContext(ctx).Info("establishing connection with rq service")
	var rqConnection rqnode.Connection
	rqConnection, err = task.StorageHandler.RqClient.Connect(ctx, task.config.RaptorQServiceAddress)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error establishing RQ connection")
	}
	defer rqConnection.Close()

	rqNodeConfig := &rqnode.Config{
		RqFilesDir: task.config.RqFilesDir,
	}
	rqService := rqConnection.RaptorQ(rqNodeConfig)
	log.WithContext(ctx).Info("connection established with rq service")

	challengeFileHash := challengeMessage.ChallengeFile.FileHashToChallenge

	var (
		regTicket     *pastel.RegTicket
		cascadeTicket *pastel.APICascadeTicket
		senseTicket   *pastel.APISenseTicket
		actionTicket  *pastel.ActionRegTicket
	)

	regTicket, cascadeTicket, senseTicket, actionTicket, err = task.getTicketInfoFromFileHash(ctx, challengeFileHash, regTicketKeys, actionTicketKeys)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error getRelevantTicketFromFileHash")
	}
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("reg ticket has been retrieved")

	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	//Checking Process
	//1. false, nil, err    - should not update the challenge to completed, so that it can be retried again
	//2. false, nil, nil    - reconstruction not required
	//3. true, symbols, nil - reconstruction required

	var msg pb.SelfHealingData
	if regTicket != nil || cascadeTicket != nil {
		isReconstructionReq, availableSymbols, err := task.checkingProcess(ctx, challengeFileHash)
		if err != nil && !isReconstructionReq {
			log.WithContext(ctx).WithError(err).WithField("failed_challenge_id", challengeMessage.ChallengeId).Error("Error in checking process")
			return err
		}

		if !isReconstructionReq && availableSymbols == nil {
			log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info(fmt.Sprintf("Reconstruction is not required for file: %s", challengeFileHash))

			if store != nil {
				log.WithContext(ctx).Println("Storing self-healing audit log")
				shChallenge := types.SelfHealingChallenge{
					ChallengeID:           challengeMessage.ChallengeId,
					MerkleRoot:            challengeMessage.MerklerootWhenChallengeSent,
					FileHash:              challengeMessage.ChallengeFile.FileHashToChallenge,
					ChallengingNode:       challengeMessage.ChallengingMasternodeId,
					RespondingNode:        challengeMessage.RespondingMasternodeId,
					ReconstructedFileHash: []byte{},
					Status:                types.ReconstructionNotRequiredSelfHealingStatus,
				}

				_, err = store.InsertSelfHealingChallenge(shChallenge)
				if err != nil {
					log.WithContext(ctx).WithError(err).Error("Error storing failed challenge to DB")
				}
			}

			return nil
		}

		file, reconstructedFileHash, err := task.selfHealing(ctx, rqService, challengeMessage, regTicket, cascadeTicket, availableSymbols)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error self-healing the file")
			return err
		}
		task.RaptorQSymbols = file

		log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("File has been reconstructed")

		log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("sending reconstructed file hash for verification")
		msg = pb.SelfHealingData{
			MessageId:                   challengeMessage.MessageId,
			MessageType:                 pb.SelfHealingData_MessageType_SELF_HEALING_VERIFICATION_MESSAGE,
			ChallengeStatus:             pb.SelfHealingData_Status_RESPONDED,
			MerklerootWhenChallengeSent: challengeMessage.MerklerootWhenChallengeSent,
			ChallengingMasternodeId:     challengeMessage.ChallengingMasternodeId,
			RespondingMasternodeId:      challengeMessage.RespondingMasternodeId,
			ReconstructedFileHash:       reconstructedFileHash,
			ChallengeFile: &pb.SelfHealingDataChallengeFile{
				FileHashToChallenge: challengeFileHash,
			},
			ChallengeId: challengeMessage.ChallengeId,
			RegTicketId: regTicket.TXID,
		}
		if cascadeTicket != nil {
			msg.ActionTicketId = actionTicket.TXID
		}
	} else if senseTicket != nil {
		reqSelfHealing, mostCommonFile := task.senseCheckingProcess(ctx, senseTicket.DDAndFingerprintsIDs)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error in checking process for sense action ticket")
		}

		if !reqSelfHealing {
			log.WithContext(ctx).WithError(err).Error("self-healing not required for sense action ticket")

			if store != nil {
				log.WithContext(ctx).Println("Storing self-healing audit log")
				shChallenge := types.SelfHealingChallenge{
					ChallengeID:           challengeMessage.ChallengeId,
					MerkleRoot:            challengeMessage.MerklerootWhenChallengeSent,
					FileHash:              challengeMessage.ChallengeFile.FileHashToChallenge,
					ChallengingNode:       challengeMessage.ChallengingMasternodeId,
					RespondingNode:        challengeMessage.RespondingMasternodeId,
					ReconstructedFileHash: []byte{},
					Status:                types.ReconstructionNotRequiredSelfHealingStatus,
				}

				_, err = store.InsertSelfHealingChallenge(shChallenge)
				if err != nil {
					log.WithContext(ctx).WithError(err).Error("Error storing failed challenge to DB")
				}
			}

			return nil
		}

		ids, idFiles, err := task.senseSelfHealing(senseTicket, mostCommonFile)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("self-healing not required for sense action ticket")
			return err
		}
		task.IDFiles = idFiles

		msg = pb.SelfHealingData{
			MessageId:                   challengeMessage.MessageId,
			MessageType:                 pb.SelfHealingData_MessageType_SELF_HEALING_VERIFICATION_MESSAGE,
			ChallengeStatus:             pb.SelfHealingData_Status_RESPONDED,
			MerklerootWhenChallengeSent: challengeMessage.MerklerootWhenChallengeSent,
			ChallengingMasternodeId:     challengeMessage.ChallengingMasternodeId,
			RespondingMasternodeId:      challengeMessage.RespondingMasternodeId,
			ReconstructedFileHash:       nil,
			ChallengeFile: &pb.SelfHealingDataChallengeFile{
				FileHashToChallenge: challengeFileHash,
			},
			ChallengeId:    challengeMessage.ChallengeId,
			RegTicketId:    "",
			IsSenseTicket:  true,
			SenseFileIds:   ids,
			ActionTicketId: actionTicket.TXID,
		}

		fileBytes, err := json.Marshal(mostCommonFile)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error(err)
		}

		fileHash, err := utils.Sha3256hash(fileBytes)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error(err)
		}

		msg.ReconstructedFileHash = fileHash

		log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("Self-healing completed")
	}

	if err := task.sendSelfHealingVerificationMessage(ctx, &msg); err != nil {
		log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).WithError(err)
		return err
	}

	log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info(fmt.Sprintf("Reconstruction is not required for file: %s", challengeFileHash))

	if task.RaptorQSymbols != nil {
		if err := task.StorageHandler.StoreRaptorQSymbolsIntoP2P(ctx, task.RaptorQSymbols, ""); err != nil {
			log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).WithError(err).Error("Error storing symbols to P2P")
			return err
		}
	} else {
		if err := task.StorageHandler.StoreListOfBytesIntoP2P(ctx, task.IDFiles); err != nil {
			log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).WithError(err).Error("Error storing id files to P2P")
			return err
		}
	}

	log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("Self-healing completed")

	return nil
}

func (task *SHTask) getTicketInfoFromFileHash(ctx context.Context, challengeFileHash string, regTicketKeys, actionTicketKeys map[string]string) (*pastel.RegTicket, *pastel.APICascadeTicket, *pastel.APISenseTicket, *pastel.ActionRegTicket, error) {
	log.WithContext(ctx).Info("retrieving the reg ticket based on file hash")
	if _, ok := regTicketKeys[challengeFileHash]; ok {
		regTicket, err := task.getRegTicket(ctx, regTicketKeys[challengeFileHash])
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("unable to retrieve reg ticket")
			return nil, nil, nil, nil, err
		}

		return &regTicket, nil, nil, nil, nil
	} else if _, ok = actionTicketKeys[challengeFileHash]; ok {
		actionTicket, err := task.getActionTicket(ctx, actionTicketKeys[challengeFileHash])
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("unable to retrieve cascade ticket")
			return nil, nil, nil, nil, err
		}

		switch actionTicket.ActionTicketData.ActionType {
		case pastel.ActionTypeCascade:
			cascadeTicket, err := actionTicket.ActionTicketData.ActionTicketData.APICascadeTicket()
			if err != nil {
				log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTicket.TXID).
					Warnf("Could not get cascade ticket for action ticket data self-healing")
				return nil, nil, nil, nil, err
			}

			return nil, cascadeTicket, nil, &actionTicket, nil
		case pastel.ActionTypeSense:
			senseTicket, err := actionTicket.ActionTicketData.ActionTicketData.APISenseTicket()
			if err != nil {
				log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTicket.TXID).
					Warnf("Could not get sense ticket for action ticket data self-healing")
				return nil, nil, nil, nil, err
			}

			return nil, nil, senseTicket, &actionTicket, nil
		}

	}

	return nil, nil, nil, nil, errors.New("failed to find the ticket")
}

func (task *SHTask) getRegTicket(ctx context.Context, regTicketID string) (regTicket pastel.RegTicket, err error) {
	regTicket, err = task.PastelClient.RegTicket(ctx, regTicketID)
	if err != nil {
		log.WithContext(ctx).Error("Error retrieving regTicket")
		return regTicket, err
	}
	decTicket, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
		return regTicket, err
	}

	regTicket.RegTicketData.NFTTicketData = *decTicket

	return regTicket, nil
}

func (task *SHTask) getActionTicket(ctx context.Context, actionTicketID string) (actionTicket pastel.ActionRegTicket, err error) {
	actionTicket, err = task.PastelClient.ActionRegTicket(ctx, actionTicketID)
	if err != nil {
		log.WithContext(ctx).Error("Error retrieving actionRegTicket")
		return actionTicket, err
	}

	decTicket, err := pastel.DecodeActionTicket(actionTicket.ActionTicketData.ActionTicket)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to decode actionRegTicket")
		return actionTicket, err
	}
	actionTicket.ActionTicketData.ActionTicketData = *decTicket

	return actionTicket, nil
}

func (task *SHTask) checkingProcess(ctx context.Context, fileHash string) (requiredReconstruction bool, symbols map[string][]byte, err error) {
	rqIDsData, err := task.P2PClient.Retrieve(ctx, fileHash, true)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("SymbolIDsFileId", fileHash).Warn("Retrieve compressed symbol IDs file from P2P failed")
		err = errors.Errorf("retrieve compressed symbol IDs file: %w", err)
		return requiredReconstruction, nil, err
	}

	if len(rqIDsData) == 0 {
		log.WithContext(ctx).WithField("SymbolIDsFileId", fileHash).Warn("Retrieve compressed symbol IDs file from P2P is empty")
		err = errors.New("retrieve compressed symbol IDs file empty")
		return requiredReconstruction, nil, err
	}

	var rqIDs []string
	rqIDs, err = task.getRQSymbolIDs(ctx, fileHash, rqIDsData)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("SymbolIDsFileId", err).Warn("Parse symbol IDs failed")
		err = errors.Errorf("parse symbol IDs: %w", err)
		return requiredReconstruction, nil, err
	}

	log.WithContext(ctx).Debugf("Symbol IDs: %v", rqIDs)
	symbols = make(map[string][]byte)
	for _, id := range rqIDs {
		var symbol []byte
		symbol, err = task.P2PClient.Retrieve(ctx, id)

		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("SymbolID", id).Error("Could not retrieve symbol")
			requiredReconstruction = true
			continue
		}

		if len(symbol) == 0 {
			log.WithContext(ctx).WithField("symbolID", id).Error("symbol received from symbolid is empty")
			requiredReconstruction = true
			continue
		}

		// Validate that the hash of each "symbol/chunk" matches its id
		h := sha3.Sum256(symbol)

		storedID := base58.Encode(h[:])
		if storedID != id {
			log.WithContext(ctx).Warnf("Symbol ID mismatched, expect %v, got %v", id, storedID)
			requiredReconstruction = true
			continue
		}

		if len(symbol) != 0 {
			symbols[id] = symbol
		}
	}

	if len(symbols) != len(rqIDs) || requiredReconstruction {
		log.WithContext(ctx).WithField("SymbolIDsFileId", fileHash).Warn("Could not retrieve all symbols")
		err = errors.New("could not retrieve all symbols from Kademlia")
		return true, symbols, nil
	}

	return false, nil, nil
}

func (task *SHTask) selfHealing(ctx context.Context, rqService rqnode.RaptorQ, msg *pb.SelfHealingData, regTicket *pastel.RegTicket, cascadeTicket *pastel.APICascadeTicket, symbols map[string][]byte) (file, reconstructedFileHash []byte, err error) {
	log.WithContext(ctx).WithField("failed_challenge_id", msg.ChallengeId).Info("Self-healing initiated")

	var encodeInfo rqnode.Encode
	if regTicket != nil {
		encodeInfo = rqnode.Encode{
			Symbols: symbols,
			EncoderParam: rqnode.EncoderParameters{
				Oti: regTicket.RegTicketData.NFTTicketData.AppTicketData.RQOti,
			},
		}
	} else if cascadeTicket != nil {
		encodeInfo = rqnode.Encode{
			Symbols: symbols,
			EncoderParam: rqnode.EncoderParameters{
				Oti: cascadeTicket.RQOti,
			},
		}
	}

	decodeInfo, err := rqService.Decode(ctx, &encodeInfo)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Restore file with rqserivce")
		return nil, nil, err
	}

	fileHash := sha3.Sum256(decodeInfo.File)
	if !bytes.Equal(fileHash[:], regTicket.RegTicketData.NFTTicketData.AppTicketData.DataHash) {
		err = errors.New("hash file mismatched")
		log.WithContext(ctx).Error("hash file mismatched")
		return nil, nil, err
	}

	return decodeInfo.File, fileHash[:], nil
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

func (task *SHTask) senseSelfHealing(senseTicket *pastel.APISenseTicket, ddFPID *pastel.DDAndFingerprints) (ids []string, idFiles [][]byte, err error) {
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

	ids, idFiles, err = pastel.GetIDFiles(ddFpFile, senseTicket.DDAndFingerprintsIc, senseTicket.DDAndFingerprintsMax)
	if err != nil {
		return nil, nil, errors.Errorf("get ID Files: %w", err)
	}

	return ids, idFiles, nil
}

func (task *SHTask) sendSelfHealingVerificationMessage(ctx context.Context, msg *pb.SelfHealingData) error {
	listOfSupernodes, err := task.SuperNodeService.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		log.WithContext(ctx).WithField("challengeID", msg.ChallengeId).WithField("method", "sendProcessSelfHealingChallenge").WithError(err).Warn("could not get Supernode extra: ", err.Error())
		return err
	}
	log.WithContext(ctx).Info(fmt.Sprintf("list of supernodes have been retrieved for process self healing challenge:%s", listOfSupernodes))

	//list of supernode ext keys without current node to find the 5 closest nodes for verification
	mapSupernodesWithoutCurrentNode := make(map[string]pastel.MasterNode)
	var sliceOfSupernodeKeysExceptCurrentNode []string
	for _, mn := range listOfSupernodes {
		if mn.ExtKey != task.nodeID {
			mapSupernodesWithoutCurrentNode[mn.ExtKey] = mn
			sliceOfSupernodeKeysExceptCurrentNode = append(sliceOfSupernodeKeysExceptCurrentNode, mn.ExtKey)
		}
	}
	log.WithContext(ctx).Info(fmt.Sprintf("current node has been filtered out from the supernodes list:%s", sliceOfSupernodeKeysExceptCurrentNode))

	//Finding 5 closest nodes to previous block hash
	sliceOfSupernodesClosestToPreviousBlockHash := task.GetNClosestSupernodeIDsToComparisonString(ctx, 5, msg.MerklerootWhenChallengeSent, sliceOfSupernodeKeysExceptCurrentNode)
	log.WithContext(ctx).Info(fmt.Sprintf("sliceOfSupernodesClosestToPreviousBlockHash:%s", sliceOfSupernodesClosestToPreviousBlockHash))

	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)

		log.WithContext(ctx).Println("Storing failed challenge to DB for self healing inspection")
		shChallenge := types.SelfHealingChallenge{
			ChallengeID:           msg.ChallengeId,
			MerkleRoot:            msg.MerklerootWhenChallengeSent,
			FileHash:              msg.ChallengeFile.FileHashToChallenge,
			ChallengingNode:       msg.ChallengingMasternodeId,
			RespondingNode:        msg.RespondingMasternodeId,
			VerifyingNode:         strings.Join(sliceOfSupernodesClosestToPreviousBlockHash, ","),
			ReconstructedFileHash: msg.ReconstructedFileHash,
			Status:                types.InProgressSelfHealingStatus,
		}

		_, err = store.InsertSelfHealingChallenge(shChallenge)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error storing failed challenge to DB")
		}
	}

	var (
		countOfFailures  int
		wg               sync.WaitGroup
		responseMessages []pb.SelfHealingData
	)
	err = nil
	// iterate through supernodes, connecting and sending the message
	for _, nodeToConnectTo := range sliceOfSupernodesClosestToPreviousBlockHash {
		nodeToConnectTo := nodeToConnectTo
		log.WithContext(ctx).WithField("outgoing_message", msg.ChallengeId).Info("outgoing message from ProcessSelfHealingChallenge")

		var (
			mn pastel.MasterNode
			ok bool
		)
		if mn, ok = mapSupernodesWithoutCurrentNode[nodeToConnectTo]; !ok {
			log.WithContext(ctx).WithField("challengeID", msg.ChallengeId).WithField("method", "sendVerifySelfHealingChallenge").Warn(fmt.Sprintf("cannot get Supernode info of Supernode id %s", mapSupernodesWithoutCurrentNode))
			continue
		}
		//We use the ExtAddress of the supernode to connect
		processingSupernodeAddr := mn.ExtAddress
		log.WithContext(ctx).WithField("challenge_id", msg.ChallengeId).Info("Sending self-healing challenge for verification to processing supernode address: " + processingSupernodeAddr)

		log.WithContext(ctx).Info(fmt.Sprintf("establishing connection with node: %s", processingSupernodeAddr))
		nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error(fmt.Sprintf("Connection failed to establish with node: %s", processingSupernodeAddr))
			continue
		}
		defer nodeClientConn.Close()

		selfHealingChallengeIF := nodeClientConn.SelfHealingChallenge()
		log.WithContext(ctx).Info(fmt.Sprintf("connection established with node:%s", nodeToConnectTo))

		msg.RespondingMasternodeId = mn.ExtKey
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := selfHealingChallengeIF.VerifySelfHealingChallenge(ctx, msg)
			if err != nil {
				log.WithContext(ctx).WithField("challengeID", msg.ChallengeId).WithField("verifierSuperNodeAddress", nodeToConnectTo).Warn("Storage challenge verification failed or didn't process: ", err.Error())
				return
			}

			responseMessages = task.appendResponses(responseMessages, res)
			log.WithContext(ctx).WithField("challenge_id", res.ChallengeId).
				Info("response has been received from verifying node")
		}()
		wg.Wait()
	}
	if err == nil {
		log.WithContext(ctx).Println("After calling storage process on " + msg.ChallengeId + " no nodes returned an error code in verification")
	}

	//Counting the number of challenges being failed by verifying nodes.
	var responseMessage pb.SelfHealingData
	for _, responseMessage = range responseMessages {
		if responseMessage.ChallengeStatus == pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE {
			countOfFailures++
		}
	}

	if countOfFailures > 0 {
		log.WithContext(ctx).Info(fmt.Sprintf("self-healing verification failed by verifiers:%d, symbols should not be updated in P2P", countOfFailures))
		return errors.New("self-healing verification failed")
	}

	return nil
}

func (task *SHTask) appendResponses(responseMessages []pb.SelfHealingData, message *pb.SelfHealingData) []pb.SelfHealingData {
	task.responseMessageMu.Lock()
	defer task.responseMessageMu.Unlock()
	responseMessages = append(responseMessages, *message)

	return responseMessages
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
		decompressedData, err := zstd.Decompress(nil, file)
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
		log.WithContext(ctx).WithField("first index of separator", firstIndexOfSeparator).WithField("decompd data", string(decompressedData)).WithField("datatobase64decode", string(dataToBase64Decode)).Println("About to base64 decode file")
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
