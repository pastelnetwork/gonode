package storagechallenge

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"

	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

func (task *StorageChallengeTask) VerifyStorageChallenge(ctx context.Context, incomingChallengeMessage *pb.StorageChallengeData) error {
	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug(incomingChallengeMessage.MessageType)
	// Incoming challenge message validation
	if err := task.validateVerifyingStorageChallengeIncomingData(incomingChallengeMessage); err != nil {
		return err
	}

	/* ----------------------------------------------- */
	/* ----- Main logic implementation goes here ----- */
	challengeFileData, err := task.GetSymbolFileByKey(ctx, incomingChallengeMessage.ChallengeFile.FileHashToChallenge, true)
	if err != nil {
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Error("could not read file data in to memory", "file.ReadFileIntoMemory", err.Error())
		return err
	}

	challengeCorrectHash := task.computeHashOfFileSlice(challengeFileData, incomingChallengeMessage.ChallengeFile.ChallengeSliceStartIndex, incomingChallengeMessage.ChallengeFile.ChallengeSliceEndIndex)
	messageType := pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_VERIFICATION_MESSAGE
	blockNumChallengeVerified, err := task.pclient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeId).Error("could not get current block count")
		return err
	}

	blocksVerifyStorageChallengeInBlocks := blockNumChallengeVerified - incomingChallengeMessage.BlockNumChallengeSent

	var challengeStatus pb.StorageChallengeDataStatus
	var saveStatus string

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

	blocksToRespondToStorageChallenge := outgoingChallengeMessage.BlockNumChallengeRespondedTo - incomingChallengeMessage.BlockNumChallengeSent
	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug("Supernode " + outgoingChallengeMessage.RespondingMasternodeId + " responded to storage challenge for file hash " + outgoingChallengeMessage.ChallengeFile.FileHashToChallenge + " in " + fmt.Sprint(blocksToRespondToStorageChallenge) + " blocks!")

	task.SaveChallengeMessageState(
		ctx,
		saveStatus,
		outgoingChallengeMessage.ChallengeId,
		outgoingChallengeMessage.ChallengingMasternodeId,
		outgoingChallengeMessage.BlockNumChallengeSent,
	)

	return nil
}

func (task *StorageChallengeTask) validateVerifyingStorageChallengeIncomingData(incomingChallengeMessage *pb.StorageChallengeData) error {
	if incomingChallengeMessage.ChallengeStatus != pb.StorageChallengeData_Status_RESPONDED {
		return fmt.Errorf("incorrect status to verify storage challenge")
	}
	if incomingChallengeMessage.MessageType != pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE {
		return fmt.Errorf("incorrect message type to verify storage challenge")
	}
	return nil
}
