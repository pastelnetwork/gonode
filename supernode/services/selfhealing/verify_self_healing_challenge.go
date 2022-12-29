package selfhealing

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
)

// VerifySelfHealingChallenge verifies the self-healing challenge
func (task *SHTask) VerifySelfHealingChallenge(ctx context.Context, challengeMessage *pb.SelfHealingData) (*pb.SelfHealingData, error) {
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("VerifySelfHealingChallenge has been invoked")

	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("retrieving reg ticket")
	regTicket, err := task.SuperNodeService.PastelClient.RegTicket(ctx, challengeMessage.RegTicketId)
	if err != nil {
		return nil, err
	}

	decTicket, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
		return nil, err
	}

	regTicket.RegTicketData.NFTTicketData = *decTicket
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("reg ticket has been retrieved")

	log.WithContext(ctx).Info("establishing connection with rq service")
	var rqConnection rqnode.Connection
	rqConnection, err = task.RQClient.Connect(ctx, task.config.RaptorQServiceAddress)
	if err != nil {
		log.WithContext(ctx).Error("Error establishing RQ connection")
	}
	defer rqConnection.Done()

	rqNodeConfig := &rqnode.Config{
		RqFilesDir: task.config.RqFilesDir,
	}
	rqService := rqConnection.RaptorQ(rqNodeConfig)
	log.WithContext(ctx).Info("connection established with rq service")

	challengeFileHash := challengeMessage.ChallengeFile.FileHashToChallenge
	log.WithContext(ctx).Info("retrieving reg ticket")

	responseMessage := challengeMessage
	responseMessage.MessageType = pb.SelfHealingData_MessageType_SELF_HEALING_RESPONSE_MESSAGE
	responseMessage.ChallengingMasternodeId = task.nodeID
	responseMessage.RespondingMasternodeId = challengeMessage.ChallengingMasternodeId

	//Checking Process
	//1. false, nil, err    - should not update the challenge to completed, so that it can be retried again
	//2. false, nil, nil    - reconstruction not required
	//3. true, symbols, nil - reconstruction required
	isReconstructionReq, availableSymbols, err := task.checkingProcess(ctx, challengeFileHash)
	if err != nil && !isReconstructionReq {
		log.WithContext(ctx).WithError(err).WithField("failed_challenge_id", challengeMessage.ChallengeId).Error("Error in checking process")
		responseMessage.ChallengeStatus = pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE
		return responseMessage, err
	}

	if !isReconstructionReq && availableSymbols == nil {
		log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info(fmt.Sprintf("Reconstruction is not required for file: %s", challengeFileHash))
		responseMessage.ChallengeStatus = pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE
		return responseMessage, nil
	}

	_, reconstructedFileHash, err := task.selfHealing(ctx, rqService, challengeMessage, regTicket, availableSymbols)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error self-healing the file")
		return responseMessage, err
	}
	log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("File has been reconstructed")

	log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("Comparing hashes")
	if string(reconstructedFileHash) == challengeMessage.ReconstructedFileHash {
		log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("reconstructed file hash does not match with the verifier reconstructed file")

		responseMessage.ChallengeStatus = pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE
		return responseMessage, nil
	}
	log.WithContext(ctx).WithField("failed_challenge_id", challengeMessage.ChallengeId).Info("hashes have been matched of the responder and verifiers")

	responseMessage.ChallengeStatus = pb.SelfHealingData_Status_SUCCEEDED
	return responseMessage, nil
}
