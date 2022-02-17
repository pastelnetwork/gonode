package storagechallenge

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	"golang.org/x/crypto/sha3"

	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

func (task *StorageChallengeTask) ProcessStorageChallenge(ctx context.Context, incomingChallengeMessage *pb.StorageChallengeData) error {
	log.WithContext(ctx).WithField("method", "ProcessStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug("Start processing storage challenge")

	// incomming challenge message validation
	if err := task.validateProcessingStorageChallengeIncomingData(incomingChallengeMessage); err != nil {
		return err
	}

	/* ----------------------------------------------- */
	/* ----- Main logic implementation goes here ----- */
	challengeFileData, err := task.repository.GetSymbolFileByKey(ctx, incomingChallengeMessage.ChallengeFile.FileHashToChallenge, true)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeId).Error("could not read file data in to memory")
		return err
	}
	challengeResponseHash := task.computeHashOfFileSlice(challengeFileData, incomingChallengeMessage.ChallengeFile.ChallengeSliceStartIndex, incomingChallengeMessage.ChallengeFile.ChallengeSliceEndIndex)
	challengeStatus := pb.StorageChallengeData_Status_RESPONDED
	messageType := pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE
	messageIDInputData := incomingChallengeMessage.ChallengingMasternodeId + incomingChallengeMessage.RespondingMasternodeId + incomingChallengeMessage.ChallengeFile.FileHashToChallenge + challengeStatus.String() + messageType.String() + incomingChallengeMessage.MerklerootWhenChallengeSent
	messageID := utils.GetHashFromString(messageIDInputData)
	blockNumChallengeRespondedTo, err := task.pclient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeId).Error("could not get current block count")
		return err
	}

	outgoingChallengeMessage := &pb.StorageChallengeData{
		MessageId:                    messageID,
		MessageType:                  messageType,
		ChallengeStatus:              challengeStatus,
		BlockNumChallengeSent:        incomingChallengeMessage.BlockNumChallengeSent,
		BlockNumChallengeRespondedTo: blockNumChallengeRespondedTo,
		BlockNumChallengeVerified:    0,
		MerklerootWhenChallengeSent:  incomingChallengeMessage.MerklerootWhenChallengeSent,
		ChallengingMasternodeId:      incomingChallengeMessage.ChallengingMasternodeId,
		//Change to my own
		RespondingMasternodeId: incomingChallengeMessage.RespondingMasternodeId,
		ChallengeFile: &pb.StorageChallengeDataChallengeFile{
			FileHashToChallenge:      incomingChallengeMessage.ChallengeFile.FileHashToChallenge,
			ChallengeSliceStartIndex: int64(incomingChallengeMessage.ChallengeFile.ChallengeSliceStartIndex),
			ChallengeSliceEndIndex:   int64(incomingChallengeMessage.ChallengeFile.ChallengeSliceEndIndex),
		},
		ChallengeSliceCorrectHash: "",
		ChallengeResponseHash:     challengeResponseHash,
		ChallengeId:               incomingChallengeMessage.ChallengeId,
	}

	blocksToRespondToStorageChallenge := outgoingChallengeMessage.BlockNumChallengeRespondedTo - incomingChallengeMessage.BlockNumChallengeSent
	log.WithContext(ctx).WithField("method", "ProcessStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug(fmt.Sprintf("Supernode %s responded to storage challenge for file hash %s in %v nano second!", outgoingChallengeMessage.RespondingMasternodeId, outgoingChallengeMessage.ChallengeFile.FileHashToChallenge, blocksToRespondToStorageChallenge))

	// send to verifying Supernode to validate challenge response hash
	if err = task.sendVerifyStorageChallenge(ctx, outgoingChallengeMessage); err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeId).Error("could not send processed challenge message to verifying node")
		return err
	}

	task.repository.SaveChallengMessageState(
		ctx,
		"respond",
		outgoingChallengeMessage.ChallengeId,
		outgoingChallengeMessage.ChallengingMasternodeId,
		outgoingChallengeMessage.BlockNumChallengeSent,
	)

	return err
}

func (task *StorageChallengeTask) validateProcessingStorageChallengeIncomingData(incomingChallengeMessage *pb.StorageChallengeData) error {
	if incomingChallengeMessage.ChallengeStatus != pb.StorageChallengeData_Status_PENDING {
		return fmt.Errorf("incorrect status to processing storage challenge")
	}
	if incomingChallengeMessage.MessageType != pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE {
		return fmt.Errorf("incorrect message type to processing storage challenge")
	}
	return nil
}

func (task *StorageChallengeTask) computeHashOfFileSlice(fileData []byte, challengeSliceStartIndex, challengeSliceEndIndex int64) string {
	challengeDataSlice := fileData[challengeSliceStartIndex:challengeSliceEndIndex]
	algorithm := sha3.New256()
	algorithm.Write(challengeDataSlice)
	return hex.EncodeToString(algorithm.Sum(nil))
}

func (task *StorageChallengeTask) sendVerifyStorageChallenge(ctx context.Context, challengeMessage *pb.StorageChallengeData) error {
	Supernodes, err := task.pclient.MasterNodesExtra(ctx)
	if err != nil {
		return err
	}

	mapSupernodes := make(map[string]pastel.MasterNode)
	for _, mn := range Supernodes {
		mapSupernodes[mn.ExtKey] = mn
	}

	verifierSupernodesAddr := []string{}
	var mn pastel.MasterNode
	var ok bool
	if mn, ok = mapSupernodes[challengeMessage.RespondingMasternodeId]; !ok {
		return fmt.Errorf("cannot get Supernode info of Supernode id %v", challengeMessage.RespondingMasternodeId)
	}
	verifierSupernodesAddr = append(verifierSupernodesAddr, mn.ExtAddress)

	err = nil
	for _, nodeToConnectTo := range verifierSupernodesAddr {
		nodeClientConn, err := task.nodeClient.Connect(ctx, nodeToConnectTo)
		if err != nil {
			err = fmt.Errorf("Could not use nodeclient to connect to: " + nodeToConnectTo)
			log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("method", "sendprocessStorageChallenge").Warn(err.Error())
			return err
		}
		storageChallengeIF := nodeClientConn.StorageChallenge()
		//Sends the process storage challenge message to the connected processing supernode
		err = storageChallengeIF.VerifyStorageChallenge(ctx, challengeMessage)
		if err != nil {
			log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("verifierSuperNodeAddress", nodeToConnectTo).Warn("Storage challenge verification failure: ", err.Error())
		}
	}
	if err == nil {
		log.WithContext(ctx).Println("After calling storage process on " + challengeMessage.ChallengeId + " no nodes returned an error code in verification")
	}
	return err
	//return s.actor.Send(ctx, s.domainActorID, newSendVerifyStorageChallengeMsg(ctx, verifierSupernodesAddr, challengeMessage))
}
