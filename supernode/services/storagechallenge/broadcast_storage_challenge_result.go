package storagechallenge

import (
	"context"

	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
)

// BroadcastStorageChallengeResult receives and process the storage challenge result
func (task *SCTask) BroadcastStorageChallengeResult(ctx context.Context, incomingBroadcastMsg types.BroadcastMessage) (types.Message, error) {
	log.WithContext(ctx).WithField("method", "BroadcastStorageChallengeResult").
		WithField("challengeID", incomingBroadcastMsg.ChallengeID).
		Debug("Start processing broadcasting message") // Incoming challenge message validation

	if err := task.storeBroadcastChallengeMsg(ctx, incomingBroadcastMsg); err != nil {
		log.WithContext(ctx).WithError(err).Error("error storing broadcast message")
	}
	log.WithContext(ctx).WithField("challenge_id", incomingBroadcastMsg.ChallengeID).
		Debug("Broadcast message has been stored")

	return types.Message{}, nil
}

func (task *SCTask) storeBroadcastChallengeMsg(ctx context.Context, msg types.BroadcastMessage) error {
	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	challengerNodeIDAndAddress := task.getNodeIDAndAddress(ctx, msg.Challenger)

	recipientNodeIDAndAddress := task.getNodeIDAndAddress(ctx, msg.Recipient)

	observersNodeIDAndAddresses := task.getNodeIDAndAddress(ctx, msg.Observers)

	if challengerNodeIDAndAddress == "" || recipientNodeIDAndAddress == "" || observersNodeIDAndAddresses == "" {
		return errors.Errorf("unable to retrieve node ID and addresses")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if store != nil {
		log.WithContext(ctx).Debug("store")
		broadcastMsgLog := types.BroadcastLogMessage{
			ChallengeID: msg.ChallengeID,
			Challenger:  challengerNodeIDAndAddress,
			Recipient:   recipientNodeIDAndAddress,
			Observers:   observersNodeIDAndAddresses,
			Data:        data,
		}

		err = store.InsertBroadcastMessage(broadcastMsgLog)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error storing broadcast message to DB")
			return err
		}
	}

	return nil
}

func (task *SCTask) getNodeIDAndAddress(ctx context.Context, msg map[string][]byte) string {
	supernodes, err := task.SuperNodeService.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		for key := range msg {
			return key
		}
	}

	mapSupernodes := make(map[string]pastel.MasterNode)
	for _, mn := range supernodes {
		mapSupernodes[mn.ExtKey] = mn
	}

	var result string
	for key := range msg {
		if value, ok := mapSupernodes[key]; ok {
			result = result + value.ExtKey + ";" + value.ExtAddress + ";"
		} else {
			result = key
		}
	}

	return result
}
