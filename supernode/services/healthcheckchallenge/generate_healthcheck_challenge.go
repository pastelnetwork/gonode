package healthcheckchallenge

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

// GenerateHealthCheckChallenges is called from service run, generate health-check challenges that will to check the connectivity of the network
func (task *HCTask) GenerateHealthCheckChallenges(ctx context.Context) error {
	logger := log.WithContext(ctx).WithField("Method", "GenerateHealthCheckChallenge")

	currentBlockCount, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		logger.WithError(err).Error("could not get current block count")
		return err
	}
	blkVerbose1, err := task.SuperNodeService.PastelClient.GetBlockVerbose1(ctx, currentBlockCount)
	if err != nil {
		logger.WithError(err).Error("could not get current block verbose 1")
		return err
	}

	challengingSupernodeID, recipient, sliceOfObservers, _, err := task.identifyChallengerRecipientsAndObserversAgainstMarkleRoot(ctx, blkVerbose1.MerkleRoot)
	if err != nil {
		logger.Error("unable to identify challenge recipients and partial observers")
		return err
	}

	//identify challenger supernode
	isMyNodeAChallenger := task.isMyNodeChallenger(ctx, challengingSupernodeID)
	if !isMyNodeAChallenger {
		return nil
	}

	//Create Challenge And Message IDs
	observersString := strings.Join(sliceOfObservers, ", ")
	challengeIDInputData := challengingSupernodeID + recipient + observersString + blkVerbose1.MerkleRoot + uuid.NewString()
	challengeID := utils.GetHashFromString(challengeIDInputData)

	// sending the message
	healthCheckChallengeMessage := types.HealthCheckMessage{
		MessageType: types.HealthCheckChallengeMessageType,
		ChallengeID: challengeID,
		Data: types.HealthCheckMessageData{
			ChallengerID: challengingSupernodeID,
			Observers:    sliceOfObservers,
			RecipientID:  recipient,
			Challenge: types.HealthCheckChallengeData{
				Block:      currentBlockCount,
				Merkelroot: blkVerbose1.MerkleRoot,
				Timestamp:  time.Now().UTC(),
			},
		},
		Sender:          challengingSupernodeID,
		SenderSignature: nil,
	}

	if err := task.StoreHealthCheckChallengeMetric(ctx, healthCheckChallengeMessage); err != nil {
		logger.WithField("hc_challenge_id", challengeID).
			WithField("message_type", healthCheckChallengeMessage.MessageType).Error(
			"error storing health check challenge metric")
	}

	if err = task.sendProcessHealthCheckChallenge(ctx, healthCheckChallengeMessage); err != nil {
		logger.WithError(err).WithField("hc_challenge_id", nil).
			Error("Error processing healthcheck challenge: ")
	}

	logger.Debug("health-check challenge has been sent for processing")
	return nil
}

func (task *HCTask) identifyChallengerRecipientsAndObserversAgainstMarkleRoot(ctx context.Context, merkleRoot string) (string, string, []string, map[string]pastel.MasterNode, error) {
	logger := log.WithContext(ctx).WithField("Method", "GenerateHealthCheckChallenge")

	logger.Debug("identifying challenge recipients and partial observers")
	sliceOfObservers := make([]string, 5)

	supernodes, err := task.getListOfOnlineSupernodes(ctx)
	if err != nil {
		logger.WithField("method", "sendHealthCheckChallenge").WithError(err).Warn("could not get Supernode extra: ", err.Error())
		return "", "", nil, nil, err
	}

	supernodesIDs := make([]string, len(supernodes))
	supernodesDetails := make(map[string]pastel.MasterNode, len(supernodes))
	for _, supernode := range supernodes {
		supernodesIDs = append(supernodesIDs, supernode.ExtKey)
		supernodesDetails[supernode.ExtKey] = supernode
	}

	respondingSupernodeIDs := task.GetNClosestSupernodeIDsToComparisonString(ctx, 7, merkleRoot, supernodesIDs)
	if len(respondingSupernodeIDs) < 6 {
		logger.WithField("closest_nodes", len(respondingSupernodeIDs)).Warn("not enough closest nodes found")
		return "", "", nil, nil, errors.Errorf("not enough closest nodes found")
	}

	//challenge challenger
	challenger := respondingSupernodeIDs[0]

	//challenge recipient
	recipient := respondingSupernodeIDs[1]

	//challenge observers
	count := 0
	for i := 2; i < len(respondingSupernodeIDs); i++ {
		sliceOfObservers[count] = respondingSupernodeIDs[i]
		count++
	}

	return challenger, recipient, sliceOfObservers, supernodesDetails, nil
}

// sendProcessHealthCheckChallenge will send the healthcheck challenge message to the responding node
func (task *HCTask) sendProcessHealthCheckChallenge(ctx context.Context, challengeMessage types.HealthCheckMessage) error {
	logger := log.WithContext(ctx).WithField("Method", "GenerateHealthCheckChallenge")

	nodesToConnect, err := task.GetNodesAddressesToConnect(ctx, challengeMessage)
	if err != nil {
		logger.WithError(err).Error("unable to find nodes to connect for send process healthcheck challenge")
		return err
	}

	if nodesToConnect == nil {
		return errors.Errorf("no nodes found to connect for send process healthcheck challenge")
	}

	signature, data, err := task.SignMessage(ctx, challengeMessage.Data)
	if err != nil {
		logger.WithError(err).Error("error signing challenge message")
		return err
	}
	challengeMessage.SenderSignature = signature

	if err := task.StoreChallengeMessage(ctx, challengeMessage); err != nil {
		logger.WithError(err).Error("error storing healthcheck challenge message")
	}

	msg := pb.HealthCheckChallengeMessage{
		MessageType:     pb.HealthCheckChallengeMessageMessageType(challengeMessage.MessageType),
		ChallengeId:     challengeMessage.ChallengeID,
		Data:            data,
		SenderId:        challengeMessage.Sender,
		SenderSignature: challengeMessage.SenderSignature,
	}

	var wg sync.WaitGroup
	var recipientNode pastel.MasterNode
	for _, node := range nodesToConnect {
		node := node
		wg.Add(1)

		go func() {
			defer wg.Done()

			if challengeMessage.Data.RecipientID == node.ExtKey {
				recipientNode = node
				return
			}
			logger := logger.WithField("node_address", node.ExtAddress)

			if err := task.SendMessage(ctx, &msg, node.ExtAddress); err != nil {
				logger.WithError(err).Error("error sending healthcheck challenge message for processing")
				return
			}

			logger.Info("health-check challenge message has been sent")
		}()
	}
	wg.Wait()
	logger.WithField("hc_challenge_id", challengeMessage.ChallengeID).Debug("health-check challenge message has been sent to observers")

	if err := task.SendMessage(ctx, &msg, recipientNode.ExtAddress); err != nil {
		logger.WithField("node_address", recipientNode.ExtAddress).WithError(err).Error("error sending healthcheck challenge message to recipient for processing")
		return err
	}
	logger.WithField("hc_challenge_id", challengeMessage.ChallengeID).
		WithField("node_address", recipientNode.ExtAddress).Info("health-check challenge message has been sent to observers & recipient")

	return nil
}

// GetNodesAddressesToConnect basically retrieves the masternode address against the pastel-id from the list and return that
func (task *HCTask) GetNodesAddressesToConnect(ctx context.Context, challengeMessage types.HealthCheckMessage) ([]pastel.MasterNode, error) {
	logger := log.WithContext(ctx).WithField("Method", "HealthCheckChallenge").WithField("hc_challenge_id", challengeMessage.ChallengeID)

	var nodesToConnect []pastel.MasterNode
	supernodes, err := task.SuperNodeService.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		logger.WithField("hc_challenge_id", challengeMessage.ChallengeID).WithError(err).Warn("could not get Supernode extra: ", err.Error())
		return nil, err
	}

	mapSupernodes := make(map[string]pastel.MasterNode)
	for _, mn := range supernodes {
		if mn.ExtAddress == "" || mn.ExtKey == "" {
			logger.WithField("hc_challenge_id", challengeMessage.ChallengeID).
				WithField("node_id", mn.ExtKey).Debug("node address or node id is empty")

			continue
		}

		mapSupernodes[mn.ExtKey] = mn
	}

	switch challengeMessage.MessageType {
	case types.HealthCheckChallengeMessageType:
		//send message to partial observers;
		for _, po := range challengeMessage.Data.Observers {
			if value, ok := mapSupernodes[po]; ok {
				nodesToConnect = append(nodesToConnect, value)
			} else {
				logger.WithField("observer_id", po).Info("observer not found in masternode list")
			}
		}

		// and challenge recipient
		nodesToConnect = append(nodesToConnect, mapSupernodes[challengeMessage.Data.RecipientID])

		logger.WithField("nodes_to_connect", nodesToConnect).Debug("nodes to send msg have been selected")
		return nodesToConnect, nil
	case types.HealthCheckResponseMessageType:
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
	case types.HealthCheckEvaluationMessageType:
		//should process and send by challenger only
		if challengeMessage.Data.ChallengerID == task.nodeID {
			//send message to challenge recipient
			nodesToConnect = append(nodesToConnect, mapSupernodes[challengeMessage.Data.RecipientID])

			//and observers
			for _, po := range challengeMessage.Data.Observers {
				if value, ok := mapSupernodes[po]; ok {
					nodesToConnect = append(nodesToConnect, value)
				} else {
					logger.WithField("observer_id", po).Info("observer not found in masternode list")
				}
			}

			logger.WithField("nodes_to_connect", nodesToConnect).Debug("nodes to send msg have been selected")
			return nodesToConnect, nil
		}

		return nil, errors.Errorf("no recipients found to send this message")
	case types.HealthCheckAffirmationMessageType:
	// not required

	case types.HealthCheckBroadcastMessageType:
		return supernodes, nil

	default:
		return nil, errors.Errorf("no nodes found to send message")
	}

	return nil, err
}

// SendMessage establish a connection with the processingSupernodeAddr and sends the given message to it.
func (task *HCTask) SendMessage(ctx context.Context, challengeMessage *pb.HealthCheckChallengeMessage, processingSupernodeAddr string) error {
	logger := log.WithContext(ctx).WithField("Method", "HealthCheckChallenge").WithField("hc_challenge_id", challengeMessage.ChallengeId)

	logger.Info("Sending health-check challenge to processing supernode address: " + processingSupernodeAddr)

	//Connect over grpc
	nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
	if err != nil {
		err = fmt.Errorf("Could not connect to: " + processingSupernodeAddr)
		logger.WithField("hc_challenge_id", challengeMessage.ChallengeId).Warn(err.Error())
		return err
	}
	defer nodeClientConn.Close()

	healthcheckChallengeIF := nodeClientConn.HealthCheckChallenge()

	switch challengeMessage.MessageType {
	case pb.HealthCheckChallengeMessage_MessageType_HEALTHCHECK_CHALLENGE_CHALLENGE_MESSAGE:
		return healthcheckChallengeIF.ProcessHealthCheckChallenge(ctx, challengeMessage)
	case pb.HealthCheckChallengeMessage_MessageType_HEALTHCHECK_CHALLENGE_RESPONSE_MESSAGE:
		return healthcheckChallengeIF.VerifyHealthCheckChallenge(ctx, challengeMessage)

	default:
		logger.Info("message type not supported by any Process & Verify worker")
	}

	return nil
}

func (task *HCTask) isMyNodeChallenger(ctx context.Context, challengingSupernodeIDForBlock string) bool {
	logger := log.WithContext(ctx).WithField("Method", "GenerateHealthCheckChallenge")

	if challengingSupernodeIDForBlock == task.nodeID {
		return true
	}
	logger.Info("exit because this node is not a challenger")

	return false
}

func (task *HCTask) getChallengers(ctx context.Context, merkleRoot string) (sliceOfChallengers []string, challengersPerBlock, challengesPerBlock int, err error) {
	logger := log.WithContext(ctx).WithField("Method", "GenerateHealthCheckChallenge")

	logger.Info("retrieving list of all super nodes")
	listOfSupernodes, err := task.GetListOfSupernode(ctx)
	if err != nil {
		logger.WithError(err).Error("could not get list of Supernode using pastel client")
		return nil, 0, 0, err
	}
	logger.Info("list of supernodes have been retrieved")

	logger.Info("identifying no of challengers & challenges to issue for this block")
	numberOfSupernodesDividedByThree := int(math.Ceil(float64(len(listOfSupernodes)) / 3))
	// Calculate the number of supernodes to issue challenge per block
	challengersPerBlock = numberOfSupernodesDividedByThree
	// Calculate the challenges per challenger
	challengesPerBlock = numberOfSupernodesDividedByThree

	// Identify challengers for this block
	sliceOfChallengers = task.GetNClosestSupernodeIDsToComparisonString(ctx, challengersPerBlock, merkleRoot, listOfSupernodes)

	return sliceOfChallengers, challengersPerBlock, challengesPerBlock, nil
}

// SignMessage signs the message using sender's pastelID and passphrase
func (task *HCTask) SignMessage(ctx context.Context, data types.HealthCheckMessageData) (sig []byte, dat []byte, err error) {
	logger := log.WithContext(ctx).WithField("Method", "HealthCheckChallenge")

	d, err := json.Marshal(data)
	if err != nil {
		logger.WithError(err).Error("error marshaling the data")
	}

	signature, err := task.PastelClient.Sign(ctx, d, task.config.PastelID, task.config.PassPhrase, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, nil, errors.Errorf("error signing health check challenge message: %w", err)
	}

	return signature, d, nil
}

// StoreChallengeMessage stores the challenge message to db for further verification
func (task *HCTask) StoreChallengeMessage(ctx context.Context, msg types.HealthCheckMessage) error {
	logger := log.WithContext(ctx).WithField("Method", "HealthCheckChallenge")

	store, err := queries.OpenHistoryDB()
	if err != nil {
		logger.WithError(err).Error("Error Opening DB")
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
		logger.Info("store")
		healthcheckChallengeLog := types.HealthCheckChallengeLogMessage{
			ChallengeID:     msg.ChallengeID,
			MessageType:     int(msg.MessageType),
			Data:            data,
			Sender:          msg.Sender,
			SenderSignature: msg.SenderSignature,
		}

		err = store.InsertHealthCheckChallengeMessage(healthcheckChallengeLog)
		if err != nil {
			if strings.Contains(err.Error(), ErrUniqueConstraint.Error()) {
				logger.WithField("challenge_id", msg.ChallengeID).
					WithField("message_type", msg.MessageType).
					WithField("sender_id", msg.Sender).
					Info("message already exists, not updating")

				return nil
			}

			logger.WithError(err).Error("Error storing challenge message to DB")
			return err
		}
	}

	return nil
}

// StoreHealthCheckChallengeMetric stores the challenge message to db for further verification
func (task *HCTask) StoreHealthCheckChallengeMetric(ctx context.Context, msg types.HealthCheckMessage) error {
	logger := log.WithContext(ctx).WithField("Method", "HealthCheckChallenge")

	store, err := queries.OpenHistoryDB()
	if err != nil {
		logger.WithError(err).Error("Error Opening DB")
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
		logger.Info("store")
		healthCheckChallengeLog := types.HealthCheckChallengeMetric{
			ChallengeID: msg.ChallengeID,
			MessageType: int(msg.MessageType),
			Data:        data,
			SenderID:    msg.Sender,
		}

		err = store.InsertHealthCheckChallengeMetric(healthCheckChallengeLog)
		if err != nil {
			if strings.Contains(err.Error(), ErrUniqueConstraint.Error()) {
				logger.WithField("challenge_id", msg.ChallengeID).
					WithField("message_type", msg.MessageType).
					WithField("sender_id", msg.Sender).
					Info("message already exists, not updating")

				return nil
			}

			return err
		}
	}

	return nil
}

// VerifyMessageSignature verifies the sender's signature on message
func (task *HCTask) VerifyMessageSignature(ctx context.Context, msg types.HealthCheckMessage) (bool, error) {
	data, err := json.Marshal(msg.Data)
	if err != nil {
		return false, errors.Errorf("unable to marshal message data")
	}

	isVerified, err := task.PastelClient.Verify(ctx, data, string(msg.SenderSignature), msg.Sender, pastel.SignAlgorithmED448)
	if err != nil {
		return false, errors.Errorf("error verifying health check challenge message: %w", err)
	}

	return isVerified, nil
}

func (task *HCTask) getListOfOnlineSupernodes(ctx context.Context) ([]pastel.MasterNode, error) {
	pingInfos, err := task.historyDB.GetAllPingInfoForOnlineNodes()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Unable to retrieve ping infos")
	}

	supernodes, err := task.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		return nil, errors.Errorf("unable to retrieve list of supernodes")
	}

	var listOfSupernodes []pastel.MasterNode
	for _, sn := range supernodes {
		if IsOnline(pingInfos, sn.ExtKey) {
			listOfSupernodes = append(listOfSupernodes, sn)
		}
	}

	return listOfSupernodes, nil
}

func IsOnline(infos types.PingInfos, supernodeID string) bool {
	for _, inf := range infos {
		if inf.SupernodeID == supernodeID && inf.IsOnline {
			return true
		}
	}

	return false
}
