package selfhealing

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

const (
	watchlistThreshold  = 6
	defaultClosestNodes = 6
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
	DataHash   string
	TicketType TicketType
}

// SelfHealingWorker checks the ping info and decide self-healing files
func (task *SHTask) SelfHealingWorker(ctx context.Context) error {
	log.WithContext(ctx).Infoln("Self Healing Worker has been invoked")

	watchlistPingInfos, err := task.retrieveWatchlistPingInfo(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("retrieveWatchlistPingInfo")
		return errors.Errorf("error retrieving watchlist ping info")
	}
	log.WithContext(ctx).Info("watchlist ping history has been retrieved")

	if len(watchlistPingInfos) < watchlistThreshold {
		log.WithContext(ctx).WithField("no_of_nodes_on_watchlist", len(watchlistPingInfos)).Info("not enough nodes on the watchlist, skipping further processing")
		return nil
	}
	log.WithContext(ctx).Info("watchlist threshold has been reached, proceeding forward")

	keys, symbolFileKeyMap, err := task.ListSymbolFileKeysFromNFTAndActionTickets(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving symbol file keys from NFT & action tickets")
		return errors.Errorf("error retrieving symbol file keys")
	}
	log.WithContext(ctx).Info("all the keys from NFT and action tickets have been listed")

	mapOfClosestNodesAgainstKeys := task.createClosestNodesMapAgainstKeys(ctx, keys)
	if len(mapOfClosestNodesAgainstKeys) == 0 {
		log.WithContext(ctx).Error("unable to create map of closest nodes against keys")
	}
	log.WithContext(ctx).Info("map of closest nodes against keys have been created")

	selfHealingTicketsMap := task.identifySelfHealingTickets(ctx, watchlistPingInfos, mapOfClosestNodesAgainstKeys, symbolFileKeyMap)
	log.WithContext(ctx).WithField("self_healing_tickets", selfHealingTicketsMap).Info("self-healing tickets have been identified")

	challengeRecipientMap, err := task.identifyChallengeRecipients(ctx, selfHealingTicketsMap)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error identifying challenge recipients")
	}
	log.WithContext(ctx).WithField("challenge_recipients", challengeRecipientMap).Info("challenge recipients have been identified")

	if err := task.prepareAndSendSelfHealingMessage(ctx, challengeRecipientMap); err != nil {
		log.WithContext(ctx).WithError(err).Error("error sending self-healing messages")
	}

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

// ListSymbolFileKeysFromNFTAndActionTickets : Get an NFT and Action Ticket's associated raptor q ticket file id's.
func (service *SHService) ListSymbolFileKeysFromNFTAndActionTickets(ctx context.Context) ([]string, map[string]SymbolFileKeyDetails, error) {
	var keys = make([]string, 0)
	symbolFileKeyMap := make(map[string]SymbolFileKeyDetails)

	regTickets, err := service.SuperNodeService.PastelClient.RegTickets(ctx)
	if err != nil {
		return keys, symbolFileKeyMap, err
	}
	if len(regTickets) == 0 {
		log.WithContext(ctx).WithField("count", len(regTickets)).Info("no reg tickets retrieved")
		return keys, symbolFileKeyMap, nil
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
func (task *SHTask) createClosestNodesMapAgainstKeys(ctx context.Context, keys []string) map[string][]string {
	mapOfClosestNodesAgainstKeys := make(map[string][]string, len(keys))

	for _, key := range keys {
		closestNodes := task.identifyClosestNodes(ctx, key)

		mapOfClosestNodesAgainstKeys[key] = append(mapOfClosestNodesAgainstKeys[key], closestNodes...)
	}

	return mapOfClosestNodesAgainstKeys
}

// identifyClosestNodes find closest nodes against the given key
func (task *SHTask) identifyClosestNodes(ctx context.Context, key string) []string {
	logger := log.WithContext(ctx).WithField("key", key)
	logger.Info("identifying closest nodes against the key")

	closestNodes := task.GetNClosestSupernodesToAGivenFileUsingKademlia(ctx, defaultClosestNodes, key, "")
	if len(closestNodes) < 1 {
		log.WithContext(ctx).WithField("file_hash", key).Info("no closest nodes have found against the file")
		return nil
	}

	logger.Info("closest nodes against the key has been retrieved")
	return closestNodes
}

func (task *SHTask) identifySelfHealingTickets(ctx context.Context, watchlistPingInfos types.PingInfos, keyClosestNodesMap map[string][]string, symbolFileKeyMap map[string]SymbolFileKeyDetails) map[string]SymbolFileKeyDetails {
	log.WithContext(ctx).Info("identifying tickets for self healing")
	selfHealingTicketsMap := make(map[string]SymbolFileKeyDetails)

	for key, closestNodes := range keyClosestNodesMap {
		if task.requireSelfHealing(closestNodes, watchlistPingInfos) {
			ticketDetails := symbolFileKeyMap[key]
			selfHealingTicketsMap[ticketDetails.TicketTxID] = SymbolFileKeyDetails{
				TicketType: ticketDetails.TicketType,
			}

			selfHealingTicketDetails := selfHealingTicketsMap[ticketDetails.TicketTxID]
			selfHealingTicketDetails.Keys = append(selfHealingTicketDetails.Keys, key)
			selfHealingTicketsMap[ticketDetails.TicketTxID] = selfHealingTicketDetails

			log.WithContext(ctx).WithField("ticket_tx_id", ticketDetails.TicketTxID).Info("ticket added for self healing")
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

func (task *SHTask) identifyChallengeRecipients(ctx context.Context, selfHealingTicketsMap map[string]SymbolFileKeyDetails) (map[string][]SymbolFileKeyDetails, error) {
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

		challengeRecipients := task.SHService.GetNClosestSupernodeIDsToComparisonString(ctx, 1, string(dataHash), nil)
		if len(challengeRecipients) < 1 {
			log.WithContext(ctx).WithField("file_hash", dataHash).Info("no closest nodes have found against the file")
			continue
		}
		recipient := challengeRecipients[0]

		challengeRecipientMap[recipient] = append(challengeRecipientMap[recipient], SymbolFileKeyDetails{
			TicketTxID: ticketDetails.TicketTxID,
			TicketType: ticketDetails.TicketType,
			Keys:       ticketDetails.Keys,
			DataHash:   string(dataHash),
		})
	}

	return challengeRecipientMap, nil
}

func (task *SHTask) getTicket(ctx context.Context, txID string, ticketType TicketType) (*pastel.NFTTicket, *pastel.APICascadeTicket, *pastel.APISenseTicket, error) {
	logger := log.WithContext(ctx).WithField("TxID", txID)

	switch ticketType {
	case nftTicketType:
		nftTicket, err := task.getNFTTicket(ctx, txID)
		if err != nil {
			logger.WithError(err).Error("error retrieving nft ticket")
			return nil, nil, nil, err
		}

		return nftTicket, nil, nil, nil
	case cascadeTicketType:
		cascadeTicket, err := task.getCascadeTicket(ctx, txID)
		if err != nil {
			logger.WithError(err).Error("error retrieving cascade ticket")
			return nil, nil, nil, err
		}

		return nil, cascadeTicket, nil, nil
	case senseTicketType:
		senseTicket, err := task.getSenseTicket(ctx, txID)
		if err != nil {
			logger.WithError(err).Error("error retrieving sense ticket")
			return nil, nil, nil, err
		}

		return nil, nil, senseTicket, nil
	}

	return nil, nil, nil, errors.Errorf("not a valid NFT or action ticket")
}

func (task *SHTask) getSenseTicket(ctx context.Context, txID string) (*pastel.APISenseTicket, error) {
	regTicket, err := task.pastelHandler.PastelClient.ActionRegTicket(ctx, txID)
	if err != nil {
		return nil, fmt.Errorf("not valid TxID for sense reg ticket")
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

func (task *SHTask) getCascadeTicket(ctx context.Context, txID string) (*pastel.APICascadeTicket, error) {
	regTicket, err := task.pastelHandler.PastelClient.ActionRegTicket(ctx, txID)
	if err != nil {
		return nil, fmt.Errorf("not valid TxID for cascade reg ticket")
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

func (task *SHTask) getNFTTicket(ctx context.Context, txID string) (*pastel.NFTTicket, error) {
	regTicket, err := task.pastelHandler.PastelClient.RegTicket(ctx, txID)
	if err != nil {
		return nil, fmt.Errorf("not valid TxID for nft reg ticket")
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

func (task *SHTask) prepareAndSendSelfHealingMessage(ctx context.Context, challengeRecipientMap map[string][]SymbolFileKeyDetails) error {
	log.WithContext(ctx).WithField("method", "prepareAndSendSelfHealingMessage").Info("method has been invoked")

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
		msgData := types.SelfHealingMessageData{
			ChallengerID: task.nodeID,
			RecipientID:  challengeRecipient,
			Challenge: types.SelfHealingChallengeData{
				Block:      currentBlockCount,
				Merkelroot: merkleroot,
				Timestamp:  time.Now().UTC(),
				Tickets:    getTicketsForSelfHealingChallengeMessage(ticketsDetails),
			},
		}

		shMsg := types.SelfHealingMessage{
			ChallengeID:            uuid.NewString(),
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
			log.WithContext(ctx).WithField("challenge_id", shMsg.ChallengeID).WithError(err).Error("unable to send challenge msg")
			continue
		}
	}
	log.WithContext(ctx).Info("self-healing messages have been sent")

	return nil
}

func getTicketsForSelfHealingChallengeMessage(ticketDetails []SymbolFileKeyDetails) []types.Ticket {
	var challengeTickets []types.Ticket
	for _, detail := range ticketDetails {
		challengeTickets = append(challengeTickets, types.Ticket{
			TxID:        detail.TicketTxID,
			TicketType:  types.TicketType(detail.TicketType),
			DataHash:    detail.DataHash,
			MissingKeys: detail.Keys,
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
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeID).Info("Sending self-healing challenge to processing supernode address: " + processingSupernodeAddr)

	//Connect over grpc
	nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
	if err != nil {
		err = fmt.Errorf("Could not connect to: " + processingSupernodeAddr)
		log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeID).WithField("method", "sendProcessStorageChallenge").Warn(err.Error())
		return err
	}
	defer nodeClientConn.Close()

	sig, data, err := task.SignMessage(ctx, challengeMessage.SelfHealingMessageData)
	if err != nil {
		log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeID).WithError(err).
			Error("error signing self-healing challenge msg")
	}
	challengeMessage.SenderSignature = sig

	msg := &pb.SelfHealingMessage{
		ChallengeId:     challengeMessage.ChallengeID,
		MessageType:     pb.SelfHealingMessageMessageType(challengeMessage.MessageType),
		SenderId:        challengeMessage.SenderID,
		SenderSignature: challengeMessage.SenderSignature,
		Data:            data,
	}
	selfHealingIF := nodeClientConn.SelfHealingChallenge()

	return selfHealingIF.ProcessSelfHealingChallenge(ctx, msg)
}
