package storagechallenge

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"

	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

// VerifyStorageChallenge : Verifying the storage challenge will occur if we are one of the <default 10> closest node ID's to the file hash being challenged.
//
//	 On receipt of challenge message we:
//			Validate it
//	 	Get the file assuming we host it locally (if not, return)
//	 	Compute the hash of the data at the indicated byte range
//			If the hash is correct and within the given byte range, success is indicated otherwise failure is indicated via SaveChallengeMessageState
func (task *SCTask) VerifyStorageChallenge(ctx context.Context, incomingChallengeMessage *pb.StorageChallengeData) (*pb.StorageChallengeData, error) {
	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug("Start verifying storage challenge") // Incoming challenge message validation
	if err := task.validateVerifyingStorageChallengeIncomingData(incomingChallengeMessage); err != nil {
		return nil, err
	}

	log.WithContext(ctx).Info("getting the file from hash to verify challenge")
	//Get the file assuming we host it locally (if not, return)
	challengeFileData, err := task.GetSymbolFileByKey(ctx, incomingChallengeMessage.ChallengeFile.FileHashToChallenge, true)
	if err != nil {
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Error("could not read local file data in to memory, so not continuing with verification.", "file.ReadFileIntoMemory", err.Error())
		return nil, err
	}
	log.WithContext(ctx).Info("file has been retrieved for verification")

	//Compute the hash of the data at the indicated byte range
	log.WithContext(ctx).Info("generating hash for the data against given indices")
	challengeCorrectHash := task.computeHashOfFileSlice(challengeFileData, incomingChallengeMessage.ChallengeFile.ChallengeSliceStartIndex, incomingChallengeMessage.ChallengeFile.ChallengeSliceEndIndex)
	log.WithContext(ctx).Info("hash of the data has been generated against the given indices")

	messageType := pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_VERIFICATION_MESSAGE
	//Identify current block count to see if validation was performed within the mandated length of time (default one block)
	blockNumChallengeVerified, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeId).Error("could not get current block count")
		return nil, err
	}

	blocksVerifyStorageChallengeInBlocks := blockNumChallengeVerified - incomingChallengeMessage.BlockNumChallengeSent

	var challengeStatus pb.StorageChallengeDataStatus
	var saveStatus string

	log.WithContext(ctx).Info("determining the challenge outcome based on calculated and received hash")
	// determine success or failure
	if (incomingChallengeMessage.ChallengeResponseHash == challengeCorrectHash) && (blocksVerifyStorageChallengeInBlocks <= task.storageChallengeExpiredBlocks) {
		challengeStatus = pb.StorageChallengeData_Status_SUCCEEDED
		saveStatus = "succeeded"
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug(fmt.Sprintf("Supernode %s correctly responded in %d blocks to a storage challenge for file %s", incomingChallengeMessage.RespondingMasternodeId, blocksVerifyStorageChallengeInBlocks, incomingChallengeMessage.ChallengeFile.FileHashToChallenge))
		task.SaveChallengeMessageState(ctx, "succeeded", incomingChallengeMessage.ChallengeId, incomingChallengeMessage.ChallengingMasternodeId, incomingChallengeMessage.BlockNumChallengeSent)
	} else if incomingChallengeMessage.ChallengeResponseHash == challengeCorrectHash {

		challengeStatus = pb.StorageChallengeData_Status_FAILED_TIMEOUT
		saveStatus = "timeout"
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug(fmt.Sprintf("Supernode %s  correctly responded in %d blocks to a storage challenge for file %s, but was too slow so failed the challenge anyway!", incomingChallengeMessage.RespondingMasternodeId, blocksVerifyStorageChallengeInBlocks, incomingChallengeMessage.ChallengeFile.FileHashToChallenge))
		task.SaveChallengeMessageState(ctx, "timeout", incomingChallengeMessage.ChallengeId, incomingChallengeMessage.ChallengingMasternodeId, incomingChallengeMessage.BlockNumChallengeSent)
	} else {

		challengeStatus = pb.StorageChallengeData_Status_FAILED_INCORRECT_RESPONSE
		saveStatus = "failed"
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug(fmt.Sprintf("Supernode %s failed by incorrectly responding to a storage challenge for file %s", incomingChallengeMessage.RespondingMasternodeId, incomingChallengeMessage.ChallengeFile.FileHashToChallenge))
		task.SaveChallengeMessageState(ctx, "failed", incomingChallengeMessage.ChallengeId, incomingChallengeMessage.ChallengingMasternodeId, incomingChallengeMessage.BlockNumChallengeSent)
	}

	messageIDInputData := incomingChallengeMessage.ChallengingMasternodeId + incomingChallengeMessage.RespondingMasternodeId + incomingChallengeMessage.ChallengeFile.FileHashToChallenge + challengeStatus.String() + messageType.String() + incomingChallengeMessage.MerklerootWhenChallengeSent
	messageID := utils.GetHashFromString(messageIDInputData)

	// Create a complete message - currently this goes nowhere.
	var outgoingChallengeMessage = &pb.StorageChallengeData{
		MessageId:                    messageID,
		MessageType:                  messageType,
		ChallengeStatus:              challengeStatus,
		BlockNumChallengeSent:        incomingChallengeMessage.BlockNumChallengeSent,
		BlockNumChallengeRespondedTo: incomingChallengeMessage.BlockNumChallengeRespondedTo,
		BlockNumChallengeVerified:    blockNumChallengeVerified,
		MerklerootWhenChallengeSent:  incomingChallengeMessage.MerklerootWhenChallengeSent,
		ChallengingMasternodeId:      incomingChallengeMessage.ChallengingMasternodeId,
		RespondingMasternodeId:       incomingChallengeMessage.RespondingMasternodeId,
		ChallengeFile: &pb.StorageChallengeDataChallengeFile{
			FileHashToChallenge:      incomingChallengeMessage.ChallengeFile.FileHashToChallenge,
			ChallengeSliceStartIndex: int64(incomingChallengeMessage.ChallengeFile.ChallengeSliceStartIndex),
			ChallengeSliceEndIndex:   int64(incomingChallengeMessage.ChallengeFile.ChallengeSliceEndIndex),
		},
		ChallengeSliceCorrectHash: challengeCorrectHash,
		ChallengeResponseHash:     incomingChallengeMessage.ChallengeResponseHash,
		ChallengeId:               incomingChallengeMessage.ChallengeId,
	}

	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)

		log.WithContext(ctx).Println("Storing challenge logs to DB")
		storageChallengeLog := types.StorageChallenge{
			ChallengeID:     outgoingChallengeMessage.ChallengeId,
			FileHash:        outgoingChallengeMessage.ChallengeFile.FileHashToChallenge,
			ChallengingNode: incomingChallengeMessage.ChallengingMasternodeId,
			Status:          types.VerifiedStorageChallengeStatus,
			RespondingNode:  incomingChallengeMessage.RespondingMasternodeId,
			GeneratedHash:   challengeCorrectHash,
			StartingIndex:   int(outgoingChallengeMessage.ChallengeFile.ChallengeSliceStartIndex),
			EndingIndex:     int(outgoingChallengeMessage.ChallengeFile.ChallengeSliceEndIndex),
		}

		_, err = store.InsertStorageChallenge(storageChallengeLog)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error storing challenge log to DB")
		}
	}

	blocksToRespondToStorageChallenge := outgoingChallengeMessage.BlockNumChallengeRespondedTo - incomingChallengeMessage.BlockNumChallengeSent
	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug("Supernode " + outgoingChallengeMessage.RespondingMasternodeId + " responded to storage challenge for file hash " + outgoingChallengeMessage.ChallengeFile.FileHashToChallenge + " in " + fmt.Sprint(blocksToRespondToStorageChallenge) + " blocks!")

	task.SaveChallengeMessageState(
		ctx,
		saveStatus,
		outgoingChallengeMessage.ChallengeId,
		outgoingChallengeMessage.ChallengingMasternodeId,
		outgoingChallengeMessage.BlockNumChallengeSent,
	)

	return outgoingChallengeMessage, nil
}

func (task *SCTask) validateVerifyingStorageChallengeIncomingData(incomingChallengeMessage *pb.StorageChallengeData) error {
	if incomingChallengeMessage.ChallengeStatus != pb.StorageChallengeData_Status_RESPONDED {
		return fmt.Errorf("incorrect status to verify storage challenge")
	}
	if incomingChallengeMessage.MessageType != pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE {
		return fmt.Errorf("incorrect message type to verify storage challenge")
	}
	return nil
}
