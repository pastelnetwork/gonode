package selfhealing

import (
	"bytes"
	"context"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"golang.org/x/crypto/sha3"
)

// ProcessSelfHealingChallenge is called from grpc server, which processes the self-healing challenge,
// and will execute the reconstruction work.
func (task *SHTask) ProcessSelfHealingChallenge(ctx context.Context, data *pb.SelfHealingData) error {
	log.WithContext(ctx).WithField("challenge_id", data.ChallengeId).Info("File Healing worker has been invoked invoked")

	mapOfFileHashes, err := task.MapSymbolFileKeysFromNFTTickets(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not generate the map of symbol file keys")
		return err
	}
	log.WithContext(ctx).Info("symbol file keys with their corresponding registered nft tickets ID have been retrieved")

	log.WithContext(ctx).Info("retrieving failed storage challenges in created status from DB")
	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).Error("Failed to open DB")
		return err
	}

	failedChallenges, err := store.QueryFailedStorageChallenges()
	if err != nil {
		log.WithContext(ctx).Error("Unable to retrieve failed storage challenges in created status from DB")

	}
	log.WithContext(ctx).Info(fmt.Sprintf("Total challenges retrieved from the DB are:%d", len(failedChallenges)))

	if len(failedChallenges) == 0 {
		return nil
	}

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

	for _, fch := range failedChallenges {
		log.WithContext(ctx).Info("retrieving reg ticket")
		regTicket, err := task.PastelClient.RegTicket(ctx, mapOfFileHashes[fch.FileHash])
		if err != nil {
			log.WithContext(ctx).Error("Error retrieving regTicket")
			continue
		}
		log.WithContext(ctx).WithField("reg_ticket_tx_id", regTicket.TXID).Info("reg ticket has been retrieved")

		//Checking Process
		//1. false, nil, err    - should not update the challenge to completed, so that it can be retried again
		//2. false, nil, nil    - reconstruction not required
		//3. true, symbols, nil - reconstruction required
		isReconstructionReq, availableSymbols, err := task.checkingProcess(ctx, fch.FileHash)
		if err != nil && !isReconstructionReq {
			log.WithContext(ctx).WithError(err).WithField("failed_challenge_id", fch.ChallengeID).Error("Error in checking process")
			continue
		}

		if !isReconstructionReq && availableSymbols == nil {
			log.WithContext(ctx).WithField("failed_challenge_id", fch.ChallengeID).Info(fmt.Sprintf("Reconstruction is not required for file: %s", fch.FileHash))
			fch.Status = types.CompletedSelfHealingStatus.String()
			if err := store.UpdateFailedStorageChallenge(fch); err != nil {
				log.WithContext(ctx).WithError(err).Error("Error updating failed storage challenge status")
			}

			continue
		}

		fch.Status = types.InProgressSelfHealingStatus.String()
		if err := store.UpdateFailedStorageChallenge(fch); err != nil {
			log.WithContext(ctx).WithError(err).Error("Error updating failed storage challenge status")
		}

		if err := task.selfHealing(ctx, rqService, fch, regTicket, availableSymbols); err != nil {
			log.WithContext(ctx).WithError(err).Error("Error self-healing the file")
			continue
		}
		log.WithContext(ctx).WithField("failed_challenge_id", fch.ChallengeID).Info("Self-healing completed")

		fch.FileReconstructingNode = task.nodeID
		fch.Status = types.CompletedSelfHealingStatus.String()
		if err := store.UpdateFailedStorageChallenge(fch); err != nil {
			log.WithContext(ctx).WithError(err).Error("Error updating failed storage challenge status")
		}
		log.WithContext(ctx).WithField("failed_challenge_id", fch.ChallengeID).Info("Self-healing status updated")

	}

	return nil
}

func (task *SHTask) checkingProcess(ctx context.Context, fileHash string) (requiredReconstruction bool, symbols map[string][]byte, err error) {
	rqIDsData, err := task.P2PClient.Retrieve(ctx, fileHash)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("SymbolIDsFileId", fileHash).Warn("Retrieve compressed symbol IDs file from P2P failed")
		err = errors.Errorf("retrieve compressed symbol IDs file: %w", err)
		return requiredReconstruction, nil, err
	}

	if len(rqIDsData) == 0 {
		log.WithContext(ctx).WithField("SymbolIDsFileId", fileHash).Warn("Retrieve compressed symbol IDs file from P2P is empty")
		err = errors.New("retrieve compressed symbol IDs file empty")
		return requiredReconstruction, nil, err
	}

	var rqIDs []string
	rqIDs, err = task.getRQSymbolIDs(ctx, fileHash, rqIDsData)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("SymbolIDsFileId", err).Warn("Parse symbol IDs failed")
		err = errors.Errorf("parse symbol IDs: %w", err)
		return requiredReconstruction, nil, err
	}

	log.WithContext(ctx).Debugf("Symbol IDs: %v", rqIDs)
	symbols = make(map[string][]byte)
	for _, id := range rqIDs {
		var symbol []byte
		symbol, err = task.P2PClient.Retrieve(ctx, id)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("SymbolID", id).Error("Could not retrieve symbol")
			requiredReconstruction = true
			continue
		}

		if len(symbol) == 0 {
			log.WithContext(ctx).WithField("symbolID", id).Error("symbol received from symbolid is empty")
			requiredReconstruction = true
			continue
		}

		// Validate that the hash of each "symbol/chunk" matches its id
		h := sha3.Sum256(symbol)

		storedID := base58.Encode(h[:])
		if storedID != id {
			log.WithContext(ctx).Warnf("Symbol ID mismatched, expect %v, got %v", id, storedID)
			requiredReconstruction = true
			continue
		}

		if len(symbol) != 0 {
			symbols[id] = symbol
		}
	}

	if len(symbols) != len(rqIDs) || requiredReconstruction {
		log.WithContext(ctx).WithField("SymbolIDsFileId", fileHash).Warn("Could not retrieve all symbols")
		err = errors.New("could not retrieve all symbols from Kademlia")
		return true, symbols, nil
	}

	return false, nil, nil
}

func (task *SHTask) selfHealing(ctx context.Context, rqService rqnode.RaptorQ, fch types.FailedStorageChallenge, regTicket pastel.RegTicket, symbols map[string][]byte) error {
	log.WithContext(ctx).WithField("failed_challenge_id", fch.ChallengeID).Info("Self-healing initiated")

	var decodeInfo *rqnode.Decode
	encodeInfo := rqnode.Encode{
		Symbols: symbols,
		EncoderParam: rqnode.EncoderParameters{
			Oti: regTicket.RegTicketData.NFTTicketData.AppTicketData.RQOti,
		},
	}

	decodeInfo, err := rqService.Decode(ctx, &encodeInfo)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Restore file with rqserivce")
		return err
	}

	fileHash := sha3.Sum256(decodeInfo.File)
	if !bytes.Equal(fileHash[:], regTicket.RegTicketData.NFTTicketData.AppTicketData.DataHash) {
		log.WithContext(ctx).Error("hash file mismatched")
		return err
	}

	file := decodeInfo.File
	// - check format of block hash - should be base58 or not
	return task.StoreRaptorQSymbolsIntoP2P(ctx, file, regTicket.TXID)
}
