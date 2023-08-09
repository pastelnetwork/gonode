package storagechallenge

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"

	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

// VerifyStorageChallenge : Verifying the storage challenge will occur only on the challenger node
//
//	 On receipt of challenge message we:
//			Validate it
//	 	Get the file assuming we host it locally (if not, return)
//	 	Compute the hash of the data at the indicated byte range
//			If the hash is correct and within the given byte range, success is indicated otherwise failure is indicated via SaveChallengeMessageState
func (task *SCTask) VerifyStorageChallenge(ctx context.Context, incomingResponseMessage types.Message) (*pb.StorageChallengeData, error) {
	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingResponseMessage.ChallengeID).Debug("Start verifying storage challenge") // Incoming challenge message validation
	if err := task.validateVerifyingStorageChallengeIncomingData(incomingResponseMessage); err != nil {
		return nil, err
	}

	log.WithContext(ctx).Info("getting the file from hash to verify challenge")
	//Get the file assuming we host it locally (if not, return)
	challengeFileData, err := task.GetSymbolFileByKey(ctx, incomingResponseMessage.Data.Challenge.FileHash, true)
	if err != nil {
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingResponseMessage.ChallengeID).Error("could not read local file data in to memory, so not continuing with verification.", "file.ReadFileIntoMemory", err.Error())
		return nil, err
	}
	log.WithContext(ctx).Info("file has been retrieved for verification")

	//Compute the hash of the data at the indicated byte range
	log.WithContext(ctx).Info("generating hash for the data against given indices")
	challengeCorrectHash := task.computeHashOfFileSlice(challengeFileData, incomingResponseMessage.Data.Challenge.StartIndex, incomingResponseMessage.Data.Challenge.EndIndex)
	log.WithContext(ctx).Info("hash of the data has been generated against the given indices")

	messageType := types.EvaluationMessageType
	//Identify current block count to see if validation was performed within the mandated length of time (default one block)
	blockNumChallengeVerified, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingResponseMessage.ChallengeID).Error("could not get current block count")
		return nil, err
	}

	blkVerbose1, err := task.SuperNodeService.PastelClient.GetBlockVerbose1(ctx, blockNumChallengeVerified)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block verbose 1")
		return nil, err
	}

	blocksVerifyStorageChallengeInBlocks := blockNumChallengeVerified - incomingResponseMessage.Data.Challenge.Block

	var challengeStatus pb.StorageChallengeDataStatus
	var saveStatus string

	var isVerified bool
	// determine success or failure
	log.WithContext(ctx).Info("determining the challenge outcome based on calculated and received hash")
	if (incomingResponseMessage.Data.Response.Hash == challengeCorrectHash) && (blocksVerifyStorageChallengeInBlocks <= task.storageChallengeExpiredBlocks) {
		challengeStatus = pb.StorageChallengeData_Status_SUCCEEDED
		saveStatus = "succeeded"
		isVerified = true
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingResponseMessage.ChallengeID).Debug(fmt.Sprintf("Supernode %s correctly responded in %d blocks to a storage challenge for file %s", incomingResponseMessage.Data.RecipientID, blocksVerifyStorageChallengeInBlocks, incomingResponseMessage.Data.Challenge.FileHash))
		task.SaveChallengeMessageState(ctx, "succeeded", incomingResponseMessage.ChallengeID, incomingResponseMessage.Data.ChallengerID, incomingResponseMessage.Data.Challenge.Block)
	} else if incomingResponseMessage.Data.Response.Hash == challengeCorrectHash {

		challengeStatus = pb.StorageChallengeData_Status_FAILED_TIMEOUT
		saveStatus = "timeout"
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingResponseMessage.ChallengeID).Debug(fmt.Sprintf("Supernode %s  correctly responded in %d blocks to a storage challenge for file %s, but was too slow so failed the challenge anyway!", incomingResponseMessage.Data.RecipientID, blocksVerifyStorageChallengeInBlocks, incomingResponseMessage.Data.Challenge.FileHash))
		task.SaveChallengeMessageState(ctx, "timeout", incomingResponseMessage.ChallengeID, incomingResponseMessage.Data.ChallengerID, incomingResponseMessage.Data.Challenge.Block)
	} else {

		challengeStatus = pb.StorageChallengeData_Status_FAILED_INCORRECT_RESPONSE
		saveStatus = "failed"
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingResponseMessage.ChallengeID).Debug(fmt.Sprintf("Supernode %s failed by incorrectly responding to a storage challenge for file %s", incomingResponseMessage.Data.RecipientID, incomingResponseMessage.Data.Challenge.FileHash))
		task.SaveChallengeMessageState(ctx, "failed", incomingResponseMessage.ChallengeID, incomingResponseMessage.Data.ChallengerID, incomingResponseMessage.Data.Challenge.Block)
	}

	_ = types.Message{
		MessageType: messageType,
		ChallengeID: incomingResponseMessage.ChallengeID,
		Data: types.MessageData{
			ChallengerID: incomingResponseMessage.Data.ChallengerID,
			Challenge: types.ChallengeData{
				Block:      incomingResponseMessage.Data.Challenge.Block,
				Merkelroot: incomingResponseMessage.Data.Challenge.Merkelroot,
				Timestamp:  incomingResponseMessage.Data.Challenge.Timestamp,
				FileHash:   incomingResponseMessage.Data.Challenge.FileHash,
				StartIndex: incomingResponseMessage.Data.Challenge.StartIndex,
				EndIndex:   incomingResponseMessage.Data.Challenge.EndIndex,
			},
			Observers:   append([]string(nil), incomingResponseMessage.Data.Observers...),
			RecipientID: incomingResponseMessage.Data.RecipientID,
			Response: types.ResponseData{
				Block:      incomingResponseMessage.Data.Response.Block,
				Merkelroot: incomingResponseMessage.Data.Response.Merkelroot,
				Hash:       incomingResponseMessage.Data.Response.Hash,
				Timestamp:  time.Now(),
			},
			ChallengerEvaluation: types.EvaluationData{
				Block:      blockNumChallengeVerified,
				Merkelroot: blkVerbose1.MerkleRoot,
				Hash:       challengeCorrectHash,
				IsVerified: isVerified,
				Timestamp:  time.Now(),
			},
		},
	}

	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)

		log.WithContext(ctx).Println("Storing challenge logs to DB")
		storageChallengeLog := types.StorageChallenge{
			ChallengeID:     incomingResponseMessage.ChallengeID,
			FileHash:        incomingResponseMessage.Data.Challenge.FileHash,
			ChallengingNode: incomingResponseMessage.Data.ChallengerID,
			Status:          types.StorageChallengeStatus(challengeStatus),
			RespondingNode:  incomingResponseMessage.Data.RecipientID,
			GeneratedHash:   challengeCorrectHash,
			StartingIndex:   incomingResponseMessage.Data.Challenge.StartIndex,
			EndingIndex:     incomingResponseMessage.Data.Challenge.EndIndex,
		}

		_, err = store.InsertStorageChallenge(storageChallengeLog)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error storing challenge log to DB")
		}
	}
	//
	//blocksToRespondToStorageChallenge := outgoingChallengeMessage.BlockNumChallengeRespondedTo - incomingChallengeMessage.BlockNumChallengeSent
	//log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug("Supernode " + outgoingChallengeMessage.RespondingMasternodeId + " responded to storage challenge for file hash " + outgoingChallengeMessage.ChallengeFile.FileHashToChallenge + " in " + fmt.Sprint(blocksToRespondToStorageChallenge) + " blocks!")

	task.SaveChallengeMessageState(
		ctx,
		saveStatus,
		incomingResponseMessage.ChallengeID,
		incomingResponseMessage.Data.ChallengerID,
		incomingResponseMessage.Data.ChallengerEvaluation.Block,
	)

	return nil, nil
}

func (task *SCTask) validateVerifyingStorageChallengeIncomingData(incomingChallengeMessage types.Message) error {
	if incomingChallengeMessage.MessageType != types.ResponseMessageType {
		return fmt.Errorf("incorrect message type to verify storage challenge")
	}
	return nil
}
