package storagechallenge

import (
	"context"
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"strings"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

// StoreChallengeMessage stores the challenge message to db for further verification
func (task SCTask) StoreChallengeMessage(ctx context.Context, msg types.Message) error {
	store, err := queries.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	data, err := json.Marshal(msg.Data)
	if err != nil {
		return err
	}

	if store != nil {
		log.WithContext(ctx).Debug("store")
		storageChallengeLog := types.StorageChallengeLogMessage{
			ChallengeID:     msg.ChallengeID,
			MessageType:     int(msg.MessageType),
			Data:            data,
			Sender:          msg.Sender,
			SenderSignature: msg.SenderSignature,
		}

		err = store.InsertStorageChallengeMessage(storageChallengeLog)
		if err != nil {
			if strings.Contains(err.Error(), ErrUniqueConstraint.Error()) {
				log.WithContext(ctx).WithField("challenge_id", msg.ChallengeID).
					WithField("message_type", msg.MessageType).
					WithField("sender_id", msg.Sender).
					Debug("message already exists, not updating")

				return nil
			}

			log.WithContext(ctx).WithError(err).Error("Error storing challenge message to DB")
			return err
		}
	}

	return nil
}

// StoreStorageChallengeMetric stores the challenge message to db for further verification
func (task SCTask) StoreStorageChallengeMetric(ctx context.Context, msg types.Message) error {
	store, err := queries.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	data, err := json.Marshal(msg.Data)
	if err != nil {
		return err
	}

	if store != nil {
		log.WithContext(ctx).Debug("store")
		storageChallengeLog := types.StorageChallengeMetric{
			ChallengeID: msg.ChallengeID,
			MessageType: int(msg.MessageType),
			Data:        data,
			SenderID:    msg.Sender,
		}

		err = store.InsertStorageChallengeMetric(storageChallengeLog)
		if err != nil {
			if strings.Contains(err.Error(), ErrUniqueConstraint.Error()) {
				log.WithContext(ctx).WithField("challenge_id", msg.ChallengeID).
					WithField("message_type", msg.MessageType).
					WithField("sender_id", msg.Sender).
					Debug("message already exists, not updating")

				return nil
			}

			return err
		}
	}

	return nil
}

// SignMessage signs the message using sender's pastelID and passphrase
func (task SCTask) SignMessage(ctx context.Context, data types.MessageData) (sig []byte, dat []byte, err error) {
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

// VerifyMessageSignature verifies the sender's signature on message
func (task SCTask) VerifyMessageSignature(ctx context.Context, msg types.Message) (bool, error) {
	data, err := json.Marshal(msg.Data)
	if err != nil {
		return false, errors.Errorf("unable to marshal message data")
	}

	isVerified, err := task.PastelClient.Verify(ctx, data, string(msg.SenderSignature), msg.Sender, pastel.SignAlgorithmED448)
	if err != nil {
		return false, errors.Errorf("error verifying storage challenge message: %w", err)
	}

	return isVerified, nil
}

// SendMessage establish a connection with the processingSupernodeAddr and sends the given message to it.
func (task SCTask) SendMessage(ctx context.Context, challengeMessage *pb.StorageChallengeMessage, processingSupernodeAddr string) error {
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Debug("Sending storage challenge to processing supernode address: " + processingSupernodeAddr)

	//Connect over grpc
	nodeClientConn, err := task.nodeClient.ConnectSN(ctx, processingSupernodeAddr)
	if err != nil {
		logError(ctx, "StorageChallenge", err)
		return fmt.Errorf("Could not connect to: " + processingSupernodeAddr)
	}
	defer nodeClientConn.Close()

	storageChallengeIF := nodeClientConn.StorageChallenge()

	switch challengeMessage.MessageType {
	case pb.StorageChallengeMessage_MessageType_STORAGE_CHALLENGE_CHALLENGE_MESSAGE:
		return storageChallengeIF.ProcessStorageChallenge(ctx, challengeMessage)
	case pb.StorageChallengeMessage_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE:
		return storageChallengeIF.VerifyStorageChallenge(ctx, challengeMessage)

	default:
		log.WithContext(ctx).Debug("message type not supported by any Process & Verify worker")
	}

	return nil
}

// GetNodesAddressesToConnect basically retrieves the masternode address against the pastel-id from the list and return that
func (task SCTask) GetNodesAddressesToConnect(ctx context.Context, challengeMessage types.Message) ([]pastel.MasterNode, error) {
	var nodesToConnect []pastel.MasterNode
	supernodes, err := task.SuperNodeService.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeID).WithField("method", "sendProcessStorageChallenge").WithError(err).Warn("could not get Supernode extra: ", err.Error())
		return nil, err
	}

	mapSupernodes := make(map[string]pastel.MasterNode)
	for _, mn := range supernodes {
		if mn.ExtAddress == "" || mn.ExtKey == "" {
			log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeID).WithField("method", "sendProcessStorageChallenge").
				WithField("node_id", mn.ExtKey).Debug("node address or node id is empty")

			continue
		}

		mapSupernodes[mn.ExtKey] = mn
	}

	logger := log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeID).
		WithField("message_type", challengeMessage.MessageType).WithField("node_id", task.nodeID)

	switch challengeMessage.MessageType {
	case types.ChallengeMessageType:
		//send message to partial observers;
		for _, po := range challengeMessage.Data.Observers {
			if value, ok := mapSupernodes[po]; ok {
				nodesToConnect = append(nodesToConnect, value)
			} else {
				logger.WithField("observer_id", po).Debug("observer not found in masternode list")
			}
		}

		// and challenge recipient
		nodesToConnect = append(nodesToConnect, mapSupernodes[challengeMessage.Data.RecipientID])

		logger.WithField("nodes_to_connect", nodesToConnect).Debug("nodes to send msg have been selected")
		return nodesToConnect, nil
	case types.ResponseMessageType:
		//should process and send by recipient only
		if challengeMessage.Data.RecipientID == task.nodeID {
			//send message to partial observers for verification
			for _, po := range challengeMessage.Data.Observers {
				if value, ok := mapSupernodes[po]; ok {
					nodesToConnect = append(nodesToConnect, value)
				} else {
					logger.WithField("observer_id", po).Debug("observer not found in masternode list")
				}
			}

			//and challenger
			nodesToConnect = append(nodesToConnect, mapSupernodes[challengeMessage.Data.ChallengerID])

			logger.WithField("nodes_to_connect", nodesToConnect).Debug("nodes to send msg have been selected")
			return nodesToConnect, nil
		}

		return nil, errors.Errorf("no recipients found to send this message")
	case types.EvaluationMessageType:
		//should process and send by challenger only
		if challengeMessage.Data.ChallengerID == task.nodeID {
			//send message to challenge recipient
			nodesToConnect = append(nodesToConnect, mapSupernodes[challengeMessage.Data.RecipientID])

			//and observers
			for _, po := range challengeMessage.Data.Observers {
				if value, ok := mapSupernodes[po]; ok {
					nodesToConnect = append(nodesToConnect, value)
				} else {
					logger.WithField("observer_id", po).Debug("observer not found in masternode list")
				}
			}

			logger.WithField("nodes_to_connect", nodesToConnect).Debug("nodes to send msg have been selected")
			return nodesToConnect, nil
		}

		return nil, errors.Errorf("no recipients found to send this message")
	case types.AffirmationMessageType:
	// not required

	case types.BroadcastMessageType:
		return supernodes, nil

	default:
		return nil, errors.Errorf("no nodes found to send message")
	}

	return nil, err
}
