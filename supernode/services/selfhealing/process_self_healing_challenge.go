package selfhealing

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
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
	rqConnection, err = task.RQClient.Connect(ctx, task.config.RaptorQServiceAddress)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error establishing RQ connection")
	}
	defer rqConnection.Done()

	rqNodeConfig := &rqnode.Config{
		RqFilesDir: task.config.RqFilesDir,
	}
	rqService := rqConnection.RaptorQ(rqNodeConfig)
	log.WithContext(ctx).Info("connection established with rq service")

	challengeFileHash := challengeMessage.ChallengeFile.FileHashToChallenge

	var (
		regTicket     pastel.RegTicket
		cascadeTicket *pastel.APICascadeTicket
		actionTicket  pastel.ActionRegTicket
	)

	log.WithContext(ctx).Info("retrieving the reg ticket based on file hash")
	if _, ok := regTicketKeys[challengeFileHash]; ok {
		regTicket, err = task.getRegTicket(ctx, regTicketKeys[challengeFileHash])
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("unable to retrieve reg ticket")
		}
		log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("reg ticket has been retrieved")

	} else if _, ok = actionTicketKeys[challengeFileHash]; ok {
		actionTicket, err = task.getActionTicket(ctx, actionTicketKeys[challengeFileHash])
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("unable to retrieve cascade ticket")
			return err
		}

		cascadeTicket, err = actionTicket.ActionTicketData.ActionTicketData.APICascadeTicket()
		if err != nil {
			log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTicket).
				Warnf("Could not get cascade ticket for action ticket data")
			return err
		}

		log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("cascade reg ticket has been retrieved")
	} else {
		log.WithContext(ctx).WithError(err).Error("Failed to find reg ticket")
	}

	//Checking Process
	//1. false, nil, err    - should not update the challenge to completed, so that it can be retried again
	//2. false, nil, nil    - reconstruction not required
	//3. true, symbols, nil - reconstruction required
	isReconstructionReq, availableSymbols, err := task.checkingProcess(ctx, challengeFileHash)
	if err != nil && !isReconstructionReq {
		log.WithContext(ctx).WithError(err).WithField("failed_challenge_id", challengeMessage.ChallengeId).Error("Error in checking process")
		return err
	}

	if !isReconstructionReq && availableSymbols == nil {
		log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info(fmt.Sprintf("Reconstruction is not required for file: %s", challengeFileHash))

		store, err := local.OpenHistoryDB()
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		}

		if store != nil {
			defer store.CloseHistoryDB(ctx)

			log.WithContext(ctx).Println("Storing failed challenge to DB for self healing inspection")
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

	log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("File has been reconstructed")

	log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("sending reconstructed file hash for verification")
	msg := &pb.SelfHealingData{
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
	if err := task.sendSelfHealingVerificationMessage(ctx, msg); err != nil {
		log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).WithError(err)
		return err
	}

	if err := task.StorageHandler.StoreRaptorQSymbolsIntoP2P(ctx, file, ""); err != nil {
		log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).WithError(err).Error("Error storing symbols to P2P")
		return err
	}

	log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("Self-healing completed")
	return nil
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

func (task *SHTask) selfHealing(ctx context.Context, rqService rqnode.RaptorQ, msg *pb.SelfHealingData, regTicket pastel.RegTicket, cascadeTicket *pastel.APICascadeTicket, symbols map[string][]byte) (file, reconstructedFileHash []byte, err error) {
	log.WithContext(ctx).WithField("failed_challenge_id", msg.ChallengeId).Info("Self-healing initiated")

	var encodeInfo rqnode.Encode
	if regTicket.TXID != "" {
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

			responseMessages = append(responseMessages, *res)
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
