package selfhealing

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

const (
	minNodesOnWatchlistThreshold = 6
	minTimeForWatchlistNodes     = 30 //in minutes
	defaultClosestNodes          = 6
)

// TicketType will identify the type of ticket
type TicketType int

const (
	cascadeTicketType TicketType = iota + 1
	senseTicketType
	nftTicketType
)

// SymbolFileKeyDetails is a struct represents with all the symbol file key details
type SymbolFileKeyDetails struct {
	TicketTxID string
	Keys       []string
	DataHash   []byte
	TicketType TicketType
	Recipient  string
}

// GenerateSelfHealingChallenge worker checks the ping info and identify self-healing tickets and their recipients
func (task *SHTask) GenerateSelfHealingChallenge(ctx context.Context) error {
	log.WithContext(ctx).Infoln("Self Healing Worker has been invoked")

	watchlistPingInfos, err := task.retrieveWatchlistPingInfo(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("retrieveWatchlistPingInfo")
		return errors.Errorf("error retrieving watchlist ping info")
	}
	log.WithContext(ctx).Info("watchlist ping history has been retrieved")

	shouldTrigger, watchlistPingInfos := task.shouldTriggerSelfHealing(watchlistPingInfos)
	if !shouldTrigger {
		log.WithContext(ctx).WithField("no_of_nodes_on_watchlist", len(watchlistPingInfos)).Info("not enough nodes on the watchlist, skipping further processing")
		return nil
	}
	log.WithContext(ctx).Info("self-healing has been triggered, proceeding with the identification of files & recipients")

	nodesOnWatchlist := getNodesOnWatchList(watchlistPingInfos)

	keys, symbolFileKeyMap, err := task.ListSymbolFileKeysFromNFTAndActionTickets(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving symbol file keys from NFT & action tickets")
		return errors.Errorf("error retrieving symbol file keys")
	}
	log.WithContext(ctx).WithField("total_keys", len(keys)).Info("all the keys from NFT and action tickets have been listed")

	mapOfClosestNodesAgainstKeys := task.createClosestNodesMapAgainstKeys(ctx, keys, watchlistPingInfos)
	if len(mapOfClosestNodesAgainstKeys) == 0 {
		log.WithContext(ctx).Error("unable to create map of closest nodes against keys")
		return nil
	}
	log.WithContext(ctx).Info("map of closest nodes against keys have been created")

	selfHealingTicketsMap := task.identifySelfHealingTickets(ctx, watchlistPingInfos, mapOfClosestNodesAgainstKeys, symbolFileKeyMap)
	if len(selfHealingTicketsMap) == 0 {
		log.WithContext(ctx).Info("no tickets required self-healing")

		if err := task.updateWatchlist(ctx, watchlistPingInfos); err != nil {
			log.WithContext(ctx).WithError(err).Error("error updating watchlist ping info")
		}

		return nil
	}
	log.WithContext(ctx).WithField("self_healing_tickets", len(selfHealingTicketsMap)).Info("self-healing tickets have been identified")

	challengeRecipientMap, err := task.identifyChallengeRecipients(ctx, selfHealingTicketsMap, watchlistPingInfos)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error identifying challenge recipients")
	}
	log.WithContext(ctx).WithField("challenge_recipients", len(challengeRecipientMap)).Info("challenge recipients have been identified")

	blockNum, err := task.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving block count")
		return err
	}

	blockVerbose, err := task.PastelClient.GetBlockVerbose1(ctx, blockNum)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving block verbose")
		return err
	}

	triggerID := getTriggerID(watchlistPingInfos, blockNum, blockVerbose.MerkleRoot)
	if err := task.prepareAndSendSelfHealingMessage(ctx, challengeRecipientMap, triggerID, nodesOnWatchlist); err != nil {
		log.WithContext(ctx).WithError(err).Error("error sending self-healing messages")
	}

	generationMetrics, err := task.getSelfHealingGenerationMetrics(ctx, triggerID, nodesOnWatchlist, challengeRecipientMap)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting generation metrics")
	}

	if generationMetrics != nil {
		task.StoreSelfHealingGenerationMetrics(ctx, *generationMetrics)
	}

	if err := task.updateWatchlist(ctx, watchlistPingInfos); err != nil {
		log.WithContext(ctx).WithError(err).Error("error updating watchlist ping info")
	}
	log.WithContext(ctx).Info("watchlist has been adjusted")

	return nil
}

// retrieveWatchlistPingInfo retrieves all the nodes on watchlist
func (task *SHTask) retrieveWatchlistPingInfo(ctx context.Context) (types.PingInfos, error) {
	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return nil, err
	}

	var infos types.PingInfos
	if store != nil {
		defer store.CloseHistoryDB(ctx)

		infos, err = store.GetWatchlistPingInfo()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return infos, nil
			}

			log.WithContext(ctx).
				Error("error retrieving watchlist ping info")

			return nil, err
		}
	}

	return infos, nil
}

func (task *SHTask) shouldTriggerSelfHealing(infos types.PingInfos) (bool, types.PingInfos) {
	var filteredPings types.PingInfos
	currentTime := time.Now().UTC()

	for _, ping := range infos {
		if ping.LastSeen.Valid {
			// Calculate the difference in minutes between the current UTC time and the LastSeen UTC time
			if currentTime.Sub(ping.LastSeen.Time).Minutes() <= minTimeForWatchlistNodes {
				filteredPings = append(filteredPings, ping)
			}
		}
	}

	if len(filteredPings) < minNodesOnWatchlistThreshold {
		return false, filteredPings
	}

	return true, filteredPings
}

// ListSymbolFileKeysFromNFTAndActionTickets : Get an NFT and Action Ticket's associated raptor q ticket file id's.
func (service *SHService) ListSymbolFileKeysFromNFTAndActionTickets(ctx context.Context) ([]string, map[string]SymbolFileKeyDetails, error) {
	var keys = make([]string, 0)
	symbolFileKeyMap := make(map[string]SymbolFileKeyDetails)

	regTickets, err := service.SuperNodeService.PastelClient.RegTickets(ctx)
	if err != nil {
		return keys, symbolFileKeyMap, err
	}
	log.WithContext(ctx).WithField("count", len(regTickets)).Info("Reg tickets retrieved")

	for i := 0; i < len(regTickets); i++ {
		decTicket, err := pastel.DecodeNFTTicket(regTickets[i].RegTicketData.NFTTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}

		regTickets[i].RegTicketData.NFTTicketData = *decTicket
		RQIDs := regTickets[i].RegTicketData.NFTTicketData.AppTicketData.RQIDs

		for j := 0; j < len(RQIDs); j++ {
			symbolFileKeyMap[RQIDs[j]] = SymbolFileKeyDetails{TicketTxID: regTickets[i].TXID, TicketType: nftTicketType}
			keys = append(keys, RQIDs[j])
		}

	}

	actionTickets, err := service.SuperNodeService.PastelClient.ActionTickets(ctx)
	if err != nil {
		return keys, symbolFileKeyMap, err
	}
	if len(actionTickets) == 0 {
		log.WithContext(ctx).WithField("count", len(actionTickets)).Info("no action tickets retrieved")
		return keys, symbolFileKeyMap, nil
	}
	log.WithContext(ctx).WithField("count", len(actionTickets)).Info("Action tickets retrieved")

	for i := 0; i < len(actionTickets); i++ {
		decTicket, err := pastel.DecodeActionTicket(actionTickets[i].ActionTicketData.ActionTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		actionTickets[i].ActionTicketData.ActionTicketData = *decTicket

		switch actionTickets[i].ActionTicketData.ActionType {
		case pastel.ActionTypeCascade:
			cascadeTicket, err := actionTickets[i].ActionTicketData.ActionTicketData.APICascadeTicket()
			if err != nil {
				log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTickets[i]).
					Warnf("Could not get cascade ticket for action ticket data")
				continue
			}

			for j := 0; j < len(cascadeTicket.RQIDs); j++ {
				symbolFileKeyMap[cascadeTicket.RQIDs[j]] = SymbolFileKeyDetails{TicketTxID: actionTickets[i].TXID, TicketType: cascadeTicketType}
				keys = append(keys, cascadeTicket.RQIDs[j])
			}
		case pastel.ActionTypeSense:
			senseTicket, err := actionTickets[i].ActionTicketData.ActionTicketData.APISenseTicket()
			if err != nil {
				log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTickets[i]).
					Warnf("Could not get sense ticket for action ticket data")
				continue
			}

			for j := 0; j < len(senseTicket.DDAndFingerprintsIDs); j++ {
				symbolFileKeyMap[senseTicket.DDAndFingerprintsIDs[j]] = SymbolFileKeyDetails{TicketTxID: actionTickets[i].TXID, TicketType: senseTicketType}
				keys = append(keys, senseTicket.DDAndFingerprintsIDs[j])
			}
		}
	}

	return keys, symbolFileKeyMap, nil
}

// FindClosestNodesAgainstKeys
func (task *SHTask) createClosestNodesMapAgainstKeys(ctx context.Context, keys []string, infos types.PingInfos) map[string][]string {
	mapOfClosestNodesAgainstKeys := make(map[string][]string, len(keys))
	var sliceOfNodesOnWatchlist []string

	for _, info := range infos {
		sliceOfNodesOnWatchlist = append(sliceOfNodesOnWatchlist, info.SupernodeID)
	}

	for _, key := range keys {
		closestNodes := task.identifyClosestNodes(ctx, key, sliceOfNodesOnWatchlist)

		mapOfClosestNodesAgainstKeys[key] = append(mapOfClosestNodesAgainstKeys[key], closestNodes...)
	}

	return mapOfClosestNodesAgainstKeys
}

// identifyClosestNodes find closest nodes against the given key
func (task *SHTask) identifyClosestNodes(ctx context.Context, key string, nodesOnWatchlist []string) []string {
	logger := log.WithContext(ctx).WithField("key", key)
	logger.Debug("identifying closest nodes against the key")

	closestNodes := task.GetNClosestSupernodesToAGivenFileUsingKademlia(ctx, defaultClosestNodes, key, []string{}, nodesOnWatchlist)

	logger.Debug("closest nodes against the key has been retrieved")
	return closestNodes
}

func (task *SHTask) identifySelfHealingTickets(ctx context.Context, watchlistPingInfos types.PingInfos, keyClosestNodesMap map[string][]string, symbolFileKeyMap map[string]SymbolFileKeyDetails) map[string]SymbolFileKeyDetails {
	log.WithContext(ctx).Info("identifying tickets for self healing")
	selfHealingTicketsMap := make(map[string]SymbolFileKeyDetails)

	for key, closestNodes := range keyClosestNodesMap {
		if task.requireSelfHealing(closestNodes, watchlistPingInfos) {
			ticketDetails := symbolFileKeyMap[key]
			selfHealingTicketDetails := selfHealingTicketsMap[ticketDetails.TicketTxID]

			selfHealingTicketDetails.TicketTxID = ticketDetails.TicketTxID
			selfHealingTicketDetails.TicketType = ticketDetails.TicketType
			selfHealingTicketDetails.Keys = append(selfHealingTicketDetails.Keys, key)

			selfHealingTicketsMap[ticketDetails.TicketTxID] = selfHealingTicketDetails

			log.WithContext(ctx).WithField("ticket_tx_id", ticketDetails.TicketTxID).Debug("ticket added for self healing")
		}
	}

	return selfHealingTicketsMap
}

func (task *SHTask) requireSelfHealing(closestNodes []string, watchlistPingInfos types.PingInfos) bool {
	for _, nodeID := range closestNodes {
		if !task.isOnWatchlist(nodeID, watchlistPingInfos) {
			return false
		}
	}

	return true
}

func (task *SHTask) isOnWatchlist(nodeID string, watchlistPingInfos types.PingInfos) bool {
	for _, info := range watchlistPingInfos {
		if nodeID == info.SupernodeID {
			return true
		}
	}

	return false
}

func (task *SHTask) identifyChallengeRecipients(ctx context.Context, selfHealingTicketsMap map[string]SymbolFileKeyDetails, watchlist types.PingInfos) (map[string][]SymbolFileKeyDetails, error) {
	challengeRecipientMap := make(map[string][]SymbolFileKeyDetails)
	for txID, ticketDetails := range selfHealingTicketsMap {
		logger := log.WithContext(ctx).WithField("TxID", txID)

		nftTicket, cascadeTicket, senseTicket, err := task.getTicket(ctx, txID, ticketDetails.TicketType)
		if err != nil {
			logger.WithError(err).Error("unable to retrieve NFT or action ticket from ticket details")
			continue
		}

		dataHash, err := task.getDataHash(nftTicket, cascadeTicket, senseTicket, ticketDetails.TicketType)
		if err != nil {
			logger.WithError(err).Error("unable to retrieve data hash from ticket details")
			continue
		}

		listOfSupernodes, err := task.getListOfOnlineSupernodes(ctx)
		if err != nil {
			logger.WithError(err).Error("unable to retrieve list of supernodes")
		}

		challengeRecipients := task.SHService.GetNClosestSupernodeIDsToComparisonString(ctx, 1, string(dataHash), task.filterWatchlistAndCurrentNode(watchlist, listOfSupernodes))
		if len(challengeRecipients) < 1 {
			log.WithContext(ctx).WithField("file_hash", dataHash).Info("no closest nodes have found against the file")
			continue
		}
		recipient := challengeRecipients[0]

		challengeRecipientMap[recipient] = append(challengeRecipientMap[recipient], SymbolFileKeyDetails{
			TicketTxID: ticketDetails.TicketTxID,
			TicketType: ticketDetails.TicketType,
			Keys:       ticketDetails.Keys,
			DataHash:   dataHash,
		})
	}

	return challengeRecipientMap, nil
}

func (task *SHTask) getTicket(ctx context.Context, txID string, ticketType TicketType) (*pastel.NFTTicket, *pastel.APICascadeTicket, *pastel.APISenseTicket, error) {
	logger := log.WithContext(ctx).WithField("TxID", txID)

	switch ticketType {
	case nftTicketType:
		nftTicket, err := task.getNFTTicket(txID)
		if err != nil {
			logger.WithError(err).Error("error retrieving nft ticket")
			return nil, nil, nil, err
		}

		return nftTicket, nil, nil, nil
	case cascadeTicketType:
		cascadeTicket, err := task.getCascadeTicket(txID)
		if err != nil {
			logger.WithError(err).Error("error retrieving cascade ticket")
			return nil, nil, nil, err
		}

		return nil, cascadeTicket, nil, nil
	case senseTicketType:
		senseTicket, err := task.getSenseTicket(txID)
		if err != nil {
			logger.WithError(err).Error("error retrieving sense ticket")
			return nil, nil, nil, err
		}

		return nil, nil, senseTicket, nil
	}

	return nil, nil, nil, errors.Errorf("not a valid NFT or action ticket")
}

func (task *SHTask) getSenseTicket(txID string) (*pastel.APISenseTicket, error) {
	regTicket, err := task.pastelHandler.PastelClient.ActionRegTicket(context.Background(), txID)
	if err != nil {
		return nil, fmt.Errorf("unable to get sense reg ticket %s:%s", txID, err.Error())
	}

	if regTicket.ActionTicketData.ActionType != pastel.ActionTypeSense {
		return nil, fmt.Errorf("not valid sense ticket")
	}

	decTicket, err := pastel.DecodeActionTicket(regTicket.ActionTicketData.ActionTicket)
	if err != nil {
		return nil, fmt.Errorf("failed to decode reg sense ticket")
	}

	regTicket.ActionTicketData.ActionTicketData = *decTicket

	senseTicket, err := regTicket.ActionTicketData.ActionTicketData.APISenseTicket()
	if err != nil {
		return nil, fmt.Errorf("failed to typecast sense ticket: %w", err)
	}

	return senseTicket, nil
}

func (task *SHTask) getCascadeTicket(txID string) (*pastel.APICascadeTicket, error) {
	regTicket, err := task.pastelHandler.PastelClient.ActionRegTicket(context.Background(), txID)
	if err != nil {
		return nil, fmt.Errorf("unable to get cascade reg ticket %s: %s", txID, err.Error())
	}

	if regTicket.ActionTicketData.ActionType != pastel.ActionTypeCascade {
		return nil, fmt.Errorf("not valid cascade ticket")
	}

	decTicket, err := pastel.DecodeActionTicket(regTicket.ActionTicketData.ActionTicket)
	if err != nil {
		return nil, fmt.Errorf("failed to decode reg cascade ticket")
	}

	regTicket.ActionTicketData.ActionTicketData = *decTicket

	cascadeTicket, err := regTicket.ActionTicketData.ActionTicketData.APICascadeTicket()
	if err != nil {
		return nil, fmt.Errorf("failed to typecast cascade ticket: %w", err)
	}

	return cascadeTicket, nil
}

func (task *SHTask) getNFTTicket(txID string) (*pastel.NFTTicket, error) {
	regTicket, err := task.pastelHandler.PastelClient.RegTicket(context.Background(), txID)
	if err != nil {
		return nil, fmt.Errorf("unable to get nft reg ticket by txid %s:%s", txID, err.Error())
	}

	decTicket, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
	if err != nil {
		return nil, fmt.Errorf("failed to decode reg nft ticket")
	}
	regTicket.RegTicketData.NFTTicketData = *decTicket

	return &regTicket.RegTicketData.NFTTicketData, nil
}

func (task *SHTask) getDataHash(nftTicket *pastel.NFTTicket, cascadeTicket *pastel.APICascadeTicket, senseTicket *pastel.APISenseTicket, ticketType TicketType) ([]byte, error) {
	switch ticketType {
	case nftTicketType:
		if nftTicket != nil {
			return nftTicket.AppTicketData.DataHash, nil
		}
	case cascadeTicketType:
		if cascadeTicket != nil {
			return cascadeTicket.DataHash, nil
		}
	case senseTicketType:
		if senseTicket != nil {
			return senseTicket.DataHash, nil
		}
	}

	return nil, errors.Errorf("data hash cannot be found because of missing ticket details")
}

func (task *SHTask) prepareAndSendSelfHealingMessage(ctx context.Context, challengeRecipientMap map[string][]SymbolFileKeyDetails, triggerID string, nodesOnWatchlist string) error {
	log.WithContext(ctx).WithField("method", "prepareAndSendSelfHealingMessage").Info("method has been invoked")

	var err error
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

	for challengeRecipient, ticketsDetails := range challengeRecipientMap {
		challengeTickets := getTicketsForSelfHealingChallengeMessage(ticketsDetails, challengeRecipient)

		log.WithContext(ctx).WithField("recipient_id", challengeRecipient).
			WithField("total_tickets", len(challengeTickets)).Info("sending for self-healing")

		msgData := types.SelfHealingMessageData{
			ChallengerID: task.nodeID,
			RecipientID:  challengeRecipient,
			Challenge: types.SelfHealingChallengeData{
				Block:            currentBlockCount,
				Merkelroot:       merkleroot,
				Timestamp:        time.Now().UTC(),
				ChallengeTickets: challengeTickets,
				NodesOnWatchlist: nodesOnWatchlist,
			},
		}

		shMsg := types.SelfHealingMessage{
			TriggerID:              triggerID,
			MessageType:            types.SelfHealingChallengeMessage,
			SelfHealingMessageData: msgData,
			SenderID:               task.nodeID,
		}

		node, err := task.GetNodeToConnect(ctx, shMsg.SelfHealingMessageData.RecipientID)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error getting node to connect")
			continue
		}

		if err := task.SendMessage(ctx, shMsg, node.ExtAddress); err != nil {
			log.WithContext(ctx).WithField("trigger_id", shMsg.TriggerID).WithError(err).Error("unable to send sh challenge msg")
			continue
		}
	}
	log.WithContext(ctx).Info("self-healing messages have been sent")

	return nil
}

func getTicketsForSelfHealingChallengeMessage(ticketDetails []SymbolFileKeyDetails, recipientID string) []types.ChallengeTicket {
	var challengeTickets []types.ChallengeTicket
	for _, detail := range ticketDetails {
		challengeTickets = append(challengeTickets, types.ChallengeTicket{
			TxID:        detail.TicketTxID,
			TicketType:  types.TicketType(detail.TicketType),
			DataHash:    detail.DataHash,
			MissingKeys: detail.Keys,
			Recipient:   recipientID,
		})
	}

	return challengeTickets
}

// SignMessage signs the message using sender's pastelID and passphrase
func (task *SHTask) SignMessage(ctx context.Context, data types.SelfHealingMessageData) (sig []byte, dat []byte, err error) {
	d, err := json.Marshal(data)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error marshaling the data")
	}

	signature, err := task.PastelClient.Sign(ctx, d, task.config.PastelID, task.config.PassPhrase, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, nil, errors.Errorf("error signing storage challenge message: %w", err)
	}

	return signature, d, nil
}

// GetNodeToConnect gets the node detail from the nodeID
func (task *SHTask) GetNodeToConnect(ctx context.Context, nodeID string) (*pastel.MasterNode, error) {
	supernodes, err := task.SuperNodeService.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		log.WithContext(ctx).WithField("method", "FetchAndMaintainPingInfo").WithError(err).Warn("could not get Supernode extra: ", err.Error())
		return nil, err
	}

	mapSupernodes := make(map[string]pastel.MasterNode)
	for _, mn := range supernodes {
		if mn.ExtAddress == "" || mn.ExtKey == "" {
			log.WithContext(ctx).WithField("method", "FetchAndMaintainPingInfo").
				WithField("node_id", mn.ExtKey).Warn("node address or node id is empty")

			continue
		}

		mapSupernodes[mn.ExtKey] = mn
	}

	if masternode, ok := mapSupernodes[nodeID]; ok {
		return &masternode, nil
	}

	log.WithContext(ctx).WithField("node_id", nodeID).WithError(err).Error("node address does not found")
	return nil, errors.Errorf("node address not found")
}

// SendMessage establish a connection with the processingSupernodeAddr and sends the given message to it.
func (task *SHTask) SendMessage(ctx context.Context, challengeMessage types.SelfHealingMessage, processingSupernodeAddr string) error {
	logger := log.WithContext(ctx).WithField("trigger_id", challengeMessage.TriggerID)

	logger.Info("sending self-healing challenge to processing supernode address: " + processingSupernodeAddr)

	ctx, cancel := context.WithTimeout(ctx, timeoutDuration)
	defer cancel()

	//Connect over grpc
	nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
	if err != nil {
		err = fmt.Errorf("Could not connect to: " + processingSupernodeAddr)
		log.WithContext(ctx).WithField("trigger_id", challengeMessage.TriggerID).WithField("method", "sendProcessStorageChallenge").Warn(err.Error())
		return err
	}
	defer nodeClientConn.Close()

	sig, data, err := task.SignMessage(ctx, challengeMessage.SelfHealingMessageData)
	if err != nil {
		log.WithContext(ctx).WithField("trigger_id", challengeMessage.TriggerID).WithError(err).
			Error("error signing self-healing challenge msg")
	}
	challengeMessage.SenderSignature = sig

	msg := &pb.SelfHealingMessage{
		TriggerId:       challengeMessage.TriggerID,
		MessageType:     pb.SelfHealingMessageMessageType(challengeMessage.MessageType),
		SenderId:        challengeMessage.SenderID,
		SenderSignature: challengeMessage.SenderSignature,
		Data:            data,
	}
	selfHealingIF := nodeClientConn.SelfHealingChallenge()

	return selfHealingIF.ProcessSelfHealingChallenge(ctx, msg)
}

func (task *SHTask) updateWatchlist(ctx context.Context, watchlistPingInfos types.PingInfos) error {
	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return nil
	}

	if store != nil {
		defer store.CloseHistoryDB(ctx)

		for _, info := range watchlistPingInfos {
			err = store.UpdatePingInfo(info.SupernodeID, false, true)
			if err != nil {
				log.WithContext(ctx).WithField("supernode_id", info.SupernodeID).
					Error("error updating watchlist ping info")

				continue
			}
		}
	}

	return nil
}

func (task *SHTask) getListOfOnlineSupernodes(ctx context.Context) ([]string, error) {
	pingInfos, err := task.GetAllPingInfo(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Unable to retrieve ping infos")
	}

	var ret []string
	for _, info := range pingInfos {

		if info.IsOnline {
			ret = append(ret, info.SupernodeID)
		}
	}

	return ret, nil
}

func (task *SHTask) filterWatchlistAndCurrentNode(watchList types.PingInfos, listOfSupernodes []string) []string {
	var filteredList []string

	for _, nodeID := range listOfSupernodes {
		if task.isOnWatchlist(nodeID, watchList) {
			continue
		}

		if task.nodeID == nodeID {
			continue
		}

		filteredList = append(filteredList, nodeID)
	}

	return filteredList
}

// StoreSelfHealingGenerationMetrics stores the self-healing generation metrics to db for further verification
func (task *SHTask) StoreSelfHealingGenerationMetrics(ctx context.Context, metricLog types.SelfHealingGenerationMetric) error {
	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	if store != nil {
		err = store.InsertSelfHealingGenerationMetrics(metricLog)
		if err != nil {
			if strings.Contains(err.Error(), ErrUniqueConstraint.Error()) {
				log.WithContext(ctx).WithField("trigger_id", metricLog.TriggerID).
					WithField("message_type", metricLog.MessageType).
					WithField("sender_id", metricLog.SenderID).
					Debug("message already exists, not storing")

				return nil
			}

			log.WithContext(ctx).WithError(err).Error("error storing self-healing generation metrics to DB")
			return err
		}
	}

	return nil
}

// StoreSelfHealingExecutionMetrics stores the self-healing execution metrics to db for further verification
func (task *SHTask) StoreSelfHealingExecutionMetrics(ctx context.Context, executionMetricsLog types.SelfHealingExecutionMetric) error {
	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	if store != nil {
		err = store.InsertSelfHealingExecutionMetrics(executionMetricsLog)
		if err != nil {
			if strings.Contains(err.Error(), ErrUniqueConstraint.Error()) {
				log.WithContext(ctx).WithField("trigger_id", executionMetricsLog.TriggerID).
					WithField("message_type", executionMetricsLog.MessageType).
					WithField("sender_id", executionMetricsLog.SenderID).
					Debug("message already exists, not storing")

				return nil
			}

			log.WithContext(ctx).WithError(err).Error("error storing self-healing execution metric to DB")
			return err
		}
	}

	return nil
}

func getTriggerID(infos types.PingInfos, blockNum int32, merkelroot string) string {
	var triggerID string
	var nodeIDs []string

	for _, info := range infos {
		nodeIDs = append(nodeIDs, info.SupernodeID)
	}

	sort.Strings(nodeIDs)
	triggerID = strings.Join(nodeIDs, ":")

	triggerID = triggerID + "-" + string(blockNum) + "-" + merkelroot

	return utils.GetHashFromString(triggerID)
}

func getNodesOnWatchList(infos types.PingInfos) string {
	var nodeIDs []string
	for _, info := range infos {
		nodeIDs = append(nodeIDs, info.SupernodeID+":"+info.IPAddress)
	}

	return strings.Join(nodeIDs, "-")
}

func (task *SHTask) getSelfHealingGenerationMetrics(ctx context.Context, triggerID, nodesOnWatchlist string, challengeRecipientMap map[string][]SymbolFileKeyDetails) (*types.SelfHealingGenerationMetric, error) {
	var challengeTickets []types.ChallengeTicket

	currentBlockCount, err := task.SuperNodeService.PastelClient.GetBlockCount(context.Background())
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block count")
		return nil, err
	}
	blkVerbose1, err := task.SuperNodeService.PastelClient.GetBlockVerbose1(context.Background(), currentBlockCount)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block verbose 1")
		return nil, err
	}
	merkleroot := blkVerbose1.MerkleRoot

	for challengeRecipient, ticketDetails := range challengeRecipientMap {
		challengeTicketsForRecipient := getTicketsForSelfHealingChallengeMessage(ticketDetails, challengeRecipient)
		challengeTickets = append(challengeTickets, challengeTicketsForRecipient...)
	}

	data := types.SelfHealingMessageData{
		ChallengerID: task.nodeID,
		Challenge: types.SelfHealingChallengeData{
			Block:            currentBlockCount,
			Merkelroot:       merkleroot,
			Timestamp:        time.Now(),
			NodesOnWatchlist: nodesOnWatchlist,
			ChallengeTickets: challengeTickets,
		},
	}

	sig, _, err := task.SignMessage(context.Background(), data)
	if err != nil {
		return nil, err
	}

	metricMsg := types.SelfHealingMessage{
		TriggerID:              triggerID,
		MessageType:            types.SelfHealingChallengeMessage,
		SelfHealingMessageData: data,
		SenderID:               task.nodeID,
		SenderSignature:        sig,
	}

	generationMetricMessage := []types.SelfHealingMessage{
		metricMsg,
	}

	generationMetricMsgBytes, err := json.Marshal(generationMetricMessage)
	if err != nil {
		return nil, errors.Errorf("error converting to bytes")
	}

	m := &types.SelfHealingGenerationMetric{
		TriggerID:       triggerID,
		MessageType:     int(types.SelfHealingChallengeMessage),
		Data:            generationMetricMsgBytes,
		SenderID:        task.nodeID,
		SenderSignature: sig,
	}

	return m, nil
}
