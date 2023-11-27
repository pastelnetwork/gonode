package selfhealing

import (
	"context"
	"database/sql"
	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
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

// SymbolFileKeyType is a struct represents the type of ticket and its txid
type SymbolFileKeyType struct {
	TicketTxID string
	TicketType
}

// SelfHealingWorker checks the ping info and decide self-healing files
func (task *SHTask) SelfHealingWorker(ctx context.Context) error {
	log.WithContext(ctx).Infoln("Self Healing Worker has been invoked")

	watchlistPingInfos, err := task.retrieveWatchlistPingInfo(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("retrieveWatchlistPingInfo")
		return errors.Errorf("error retrieving watchlist ping info")
	}

	if len(watchlistPingInfos) < watchlistThreshold {
		log.WithContext(ctx).WithField("no_of_nodes_on_watchlist", len(watchlistPingInfos)).Info("not enough nodes on the watchlist, skipping further processing")
		return errors.Errorf("no of nodes on the watchlist are not sufficient for further processing")
	}

	keys, symbolFileKeyMap, err := task.ListSymbolFileKeysFromNFTAndActionTickets(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving symbol file keys from NFT & action tickets")
		return errors.Errorf("error retrieving symbol file keys")
	}

	mapOfClosestNodesAgainstKeys := task.createClosestNodesMapAgainstKeys(ctx, keys)
	if len(mapOfClosestNodesAgainstKeys) == 0 {
		log.WithContext(ctx).Error("unable to create map of closest nodes against keys")
	}

	selfHealingTicketsMap := task.identifySelfHealingTickets(ctx, watchlistPingInfos, mapOfClosestNodesAgainstKeys, symbolFileKeyMap)

	log.WithContext(ctx).WithField("self_healing_tickets", selfHealingTicketsMap).Info("self-healing tickets have been identified")
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
func (service *SHService) ListSymbolFileKeysFromNFTAndActionTickets(ctx context.Context) ([]string, map[string]SymbolFileKeyType, error) {
	var keys = make([]string, 0)
	symbolFileKeyMap := make(map[string]SymbolFileKeyType)

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
			symbolFileKeyMap[RQIDs[j]] = SymbolFileKeyType{TicketTxID: regTickets[i].TXID, TicketType: nftTicketType}
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
				symbolFileKeyMap[cascadeTicket.RQIDs[j]] = SymbolFileKeyType{TicketTxID: actionTickets[i].TXID, TicketType: cascadeTicketType}
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
				symbolFileKeyMap[senseTicket.DDAndFingerprintsIDs[j]] = SymbolFileKeyType{TicketTxID: actionTickets[i].TXID, TicketType: senseTicketType}
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

func (task *SHTask) identifySelfHealingTickets(ctx context.Context, watchlistPingInfos types.PingInfos, keyClosestNodesMap map[string][]string, symbolFileKeyMap map[string]SymbolFileKeyType) map[string]SymbolFileKeyType {
	log.WithContext(ctx).Info("identifying tickets for self healing")
	selfHealingTicketsMap := make(map[string]SymbolFileKeyType)

	for key, closestNodes := range keyClosestNodesMap {
		if task.requireSelfHealing(closestNodes, watchlistPingInfos) {
			ticketDetails := symbolFileKeyMap[key]
			selfHealingTicketsMap[ticketDetails.TicketTxID] = SymbolFileKeyType{TicketType: ticketDetails.TicketType}
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
