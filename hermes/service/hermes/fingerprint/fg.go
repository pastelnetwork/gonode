package fingerprint

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/hermes/domain"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	runTaskInterval = 2 * time.Minute
)

// Run stores the latest block hash and height to DB if not stored already
func (s *fingerprintService) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(runTaskInterval):
			// Check if node is synchronized or not
			log.WithContext(ctx).Info("fingerprint service run() has been invoked")
			if err := s.sync.WaitSynchronization(ctx); err != nil {
				log.WithContext(ctx).WithError(err).Error("error syncing master-node")
				continue
			}

			group, gctx := errgroup.WithContext(ctx)
			group.Go(func() error {
				return s.run(gctx)
			})

			if err := group.Wait(); err != nil {
				log.WithContext(gctx).WithError(err).Errorf("run task failed")
			}
		}
	}

}

func (s *fingerprintService) run(ctx context.Context) error {
	log.WithContext(ctx).Info("fingerprint service run() has been invoked")

	log.WithContext(ctx).Info("getting Activation tickets, checking non seed records.")
	nonseed, err := s.store.CheckNonSeedRecord(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to get nonseed record")
	} else if !nonseed {
		log.WithContext(ctx).Info("No NonSeed Record, set latestBlockHeight to 0")
		s.latestNFTBlockHeight = 0
		s.latestSenseBlockHeight = 0
	}

	if err := s.parseSenseTickets(ctx); err != nil {
		return err
	}

	return s.parseNFTTickets(ctx)
}

func (s *fingerprintService) parseSenseTickets(ctx context.Context) error {
	lastKnownGoodHeight := s.latestSenseBlockHeight

	senseActTickets, err := s.pastelClient.ActionActivationTicketsFromBlockHeight(ctx, uint64(lastKnownGoodHeight))
	if err != nil {
		log.WithError(err).Errorf("get registered ticket - exit runTask")
		return nil
	}

	sort.Slice(senseActTickets, func(i, j int) bool {
		return senseActTickets[i].Height < senseActTickets[j].Height
	})

	if len(senseActTickets) == 0 {
		log.WithContext(ctx).WithField("block height", s.latestSenseBlockHeight).Info("No sense tickets found")
		return nil
	}

	for i := 0; i < len(senseActTickets); i++ {
		if senseActTickets[i].Height < s.latestSenseBlockHeight {
			continue
		}

		regTicket, senseTicket, err := s.getSenseTicket(ctx, senseActTickets[i].TXID)
		if err != nil {
			if !strings.Contains(err.Error(), "not a sense ticket") {
				log.WithContext(ctx).WithField("txid", senseActTickets[i].TXID).WithField("height", senseActTickets[i].Height).
					WithError(err).Error("unable to get sense ticket")
			}
			continue
		}

		log.WithContext(ctx).WithField("act-txid", senseActTickets[i].TXID).WithField("reg-txid", regTicket.TXID).WithField("height", senseActTickets[i].Height).
			Info("Found sense activation ticket")

		stored, p2pRunning := s.fetchDDFpFileAndStoreFingerprints(ctx, regTicket.ActionTicketData.ActionType, "", senseActTickets[i].TXID,
			regTicket.TXID, senseTicket.DDAndFingerprintsIDs)
		if stored {
			if senseActTickets[i].Height > lastKnownGoodHeight {
				lastKnownGoodHeight = senseActTickets[i].Height
			}
		} else if !p2pRunning {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).
				Info("P2P service is not running, so we can't get the fingerprint for this file hash, stopping this run of the task")

			break
		}
	}

	if lastKnownGoodHeight > s.latestSenseBlockHeight {
		s.latestSenseBlockHeight = lastKnownGoodHeight
	}

	log.WithContext(ctx).WithField("latest Sense blockheight", s.latestSenseBlockHeight).Info("hermes successfully scanned to latest block height")

	return nil
}

func (s *fingerprintService) parseNFTTickets(ctx context.Context) error {
	actTickets, err := s.pastelClient.ActTickets(ctx, pastel.ActTicketAll, s.latestNFTBlockHeight)
	if err != nil {
		log.WithError(err).Error("unable to get act tickets - exit runtask now")
		return nil
	}
	if len(actTickets) == 0 {
		return nil
	}

	sort.Slice(actTickets, func(i, j int) bool {
		return actTickets[i].Height < actTickets[j].Height
	})

	log.WithContext(ctx).WithField("count", len(actTickets)).Info("Act tickets retrieved")

	//track latest block height, but don't set it until we check all the nft reg tickets and the sense tickets.
	lastKnownGoodHeight := s.latestNFTBlockHeight

	//loop through nft tickets and store newly found nft reg tickets
	for i := 0; i < len(actTickets); i++ {
		if actTickets[i].Height < s.latestNFTBlockHeight {
			continue
		}

		log.WithContext(ctx).WithField("act-txid", actTickets[i].TXID).WithField("reg-txid", actTickets[i].ActTicketData.RegTXID).
			Info("Found new NFT ticket")

		regTicket, err := s.pastelClient.RegTicket(ctx, actTickets[i].ActTicketData.RegTXID)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("reg-txid", actTickets[i].ActTicketData.RegTXID).
				WithField("act-txid", actTickets[i].TXID).Error("unable to get reg ticket")
			continue
		}

		log.WithContext(ctx).WithField("act-txid", actTickets[i].TXID).WithField("reg-txid", actTickets[i].ActTicketData.RegTXID).Info("Found Reg Ticket for NFT-Act ticket")

		decTicket, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
		if err != nil {
			log.WithContext(ctx).WithField("act-txid", actTickets[i].TXID).WithField("reg-txid", actTickets[i].ActTicketData.RegTXID).WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		regTicket.RegTicketData.NFTTicketData = *decTicket

		stored, p2pRunning := s.fetchDDFpFileAndStoreFingerprints(ctx, regTicket.RegTicketData.Type, regTicket.RegTicketData.NFTTicketData.AppTicketData.NFTSeriesName,
			actTickets[i].TXID, regTicket.TXID, regTicket.RegTicketData.NFTTicketData.AppTicketData.DDAndFingerprintsIDs)
		if stored {
			if actTickets[i].Height > lastKnownGoodHeight {
				lastKnownGoodHeight = actTickets[i].Height
			}
		} else if !p2pRunning {
			log.WithContext(ctx).WithField("act-txid", actTickets[i].TXID).WithField("reg-txid", actTickets[i].ActTicketData.RegTXID).
				Info("P2P service is not running, so we can't get the fingerprint for this file hash, stopping this run of the task")

			break
		}
	}
	//loop through action tickets and store newly found nft reg tickets

	if lastKnownGoodHeight > s.latestNFTBlockHeight {
		s.latestNFTBlockHeight = lastKnownGoodHeight
	}

	log.WithContext(ctx).WithField("latest NFT blockheight", s.latestNFTBlockHeight).Debugf("hermes successfully scanned to latest block height")

	return nil
}

func (s *fingerprintService) fetchDDFpFileAndStoreFingerprints(ctx context.Context, tType string, series string, actTXID string, regTXID string, ddFpIDs []string) (stored bool, p2pRunning bool) {
	ddAndFpFromTicket := &pastel.DDAndFingerprints{}
	var err error
	//Get the dd and fp file from the ticket

	logMsg := log.WithContext(ctx).WithField("reg-txid", regTXID).WithField("act-txid", actTXID).WithField("type", typeMapper(tType))
	if exist, err := s.store.IfFingerprintExistsByRegTxid(ctx, regTXID); err == nil && exist {
		logMsg.Info("ticket already exists")
		return true, true
	}

	p2pServiceRunning := true
	for _, id := range ddFpIDs {
		ddAndFpFromTicket, err = s.tryToGetFingerprintFileFromHash(ctx, id)
		if err == nil && ddAndFpFromTicket != nil {
			break
		}

		logMsg.WithField("id", id).Error("Failed to get dd and fp file from ticket")
		if strings.Contains(err.Error(), "p2p service is not running") {
			p2pServiceRunning = false
			break
		}

		//probably too verbose even for debug.
		logMsg.WithField("id", id).
			Error("Could not get the fingerprint for this file hash")
	}

	if !p2pServiceRunning {
		return false, false
	}

	if ddAndFpFromTicket == nil {
		logMsg.Info("None of the dd and fp id files for this ticket could be properly unmarshalled")
		return false, true
	}
	if ddAndFpFromTicket.HashOfCandidateImageFile == "" {
		logMsg.Info("This ticket's DDAndFp struct has no HashOfCandidateImageFile, perhaps it's an older version.")
		return false, true
	}

	logMsg.WithField("image_hash", ddAndFpFromTicket.HashOfCandidateImageFile).Info("image hash retrieved")

	existsInDatabase, err := s.store.IfFingerprintExists(ctx, ddAndFpFromTicket.HashOfCandidateImageFile)
	if existsInDatabase {
		logMsg.Info("fingerprint already exist in database, skipping")
		return true, true
	}
	if err != nil {
		logMsg.Error("Could not properly query the dd database for this hash")
		return false, true
	}

	//make sure ImageFingerprintOfCnadidateImageFile exists.
	// this could fail if the ticket is an older version of the DDAndFingerprints struct, so we will continue to next fingerprint
	if ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile == nil {
		logMsg.Info("This ticket's DDAndFp struct has no ImageFingerprintOfCandidateImageFile, perhaps it's an older version.")
		return false, true
	}
	if len(ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile) < 1 {
		logMsg.Info("This ticket's DDAndFp struct's ImageFingerprintOfCandidateImageFile is zero length, perhaps it's an older version.")
		return false, true
	}

	groupID, subsetID := getGroupIDAndSubsetID(series, ddAndFpFromTicket.OpenAPIGroupIDString)

	if err := s.store.StoreFingerprint(ctx, &domain.DDFingerprints{
		Sha256HashOfArtImageFile:                   ddAndFpFromTicket.HashOfCandidateImageFile,
		ImageFingerprintVector:                     toFloat64Array(ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile),
		DatetimeFingerprintAddedToDatabase:         time.Now().Format("2006-01-02 15:04:05"),
		PathToArtImageFile:                         ddAndFpFromTicket.ImageFilePath,
		ImageThumbnailAsBase64:                     ddAndFpFromTicket.CandidateImageThumbnailWebpAsBase64String,
		RequestType:                                typeMapper(tType),
		IDString:                                   subsetID,
		OpenAPIGroupIDString:                       groupID,
		CollectionNameString:                       ddAndFpFromTicket.CollectionNameString,
		DoesNotImpactTheFollowingCollectionsString: ddAndFpFromTicket.DoesNotImpactTheFollowingCollectionStrings,
		RegTXID: regTXID,
	}); err != nil {
		logMsg.WithError(err).Error("Failed to store fingerprint")
		return false, true
	}

	return true, true
}

// Utility function to get dd and fp file from an id hash, where the file should be stored
func (s *fingerprintService) tryToGetFingerprintFileFromHash(ctx context.Context, hash string) (*pastel.DDAndFingerprints, error) {
	rawFile, err := s.p2p.Retrieve(ctx, hash)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("retrieve err")
		return nil, errors.Errorf("Error finding dd and fp file: %w", err)
	}

	decData, err := zstd.Decompress(nil, rawFile)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("decompress err")
		return nil, errors.Errorf("decompress: %w", err)
	}

	splits := bytes.Split(decData, []byte{pastel.SeparatorByte})
	if (len(splits)) < 2 {
		log.WithContext(ctx).WithError(err).Error("incorrecrt split err")
		return nil, errors.Errorf("error separating file by separator bytes, separator not found")
	}
	file, err := utils.B64Decode(splits[0])
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("b64 decode err")
		return nil, errors.Errorf("decode file: %w", err)
	}

	ddFingerprint := &pastel.DDAndFingerprints{}
	if err := json.Unmarshal(file, ddFingerprint); err != nil {
		log.WithContext(ctx).WithError(err).Error("unmarshal err")
		return nil, errors.Errorf("unmarshal json: %w", err)
	}

	return ddFingerprint, nil
}

func (s *fingerprintService) Stats(_ context.Context) (map[string]interface{}, error) {
	//chain-reorg stats can be implemented here
	return nil, nil
}

func (s *fingerprintService) getSenseTicket(ctx context.Context, regTXID string) (*pastel.ActionRegTicket, *pastel.APISenseTicket, error) {
	regTicket, err := s.pastelClient.ActionRegTicket(ctx, regTXID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find action act by action reg txid: %w", err)
	}

	if regTicket.ActionTicketData.ActionType != pastel.ActionTypeSense {
		return nil, nil, errors.New("not a sense ticket")
	}

	decTicket, err := pastel.DecodeActionTicket(regTicket.ActionTicketData.ActionTicket)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode reg ticket: %w", err)
	}
	regTicket.ActionTicketData.ActionTicketData = *decTicket

	senseTicket, err := regTicket.ActionTicketData.ActionTicketData.APISenseTicket()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to typecast sense ticket: %w", err)
	}

	return &regTicket, senseTicket, nil
}

// getGroupIDAndSubsetID returns the groupID and subsetID for a given Sense or NFT ticket.
// ddGroupID is open_api_group_id_string from the output of dd-service.
// ddSubsetID is open_api_subset_id_string from the output of dd-service.
func getGroupIDAndSubsetID(ddGroupID, ddSubsetID string) (groupID, subsetID string) {
	subsetID = "PASTEL"
	if ddSubsetID != "" && !strings.EqualFold(ddSubsetID, "NA") {
		subsetID = ddSubsetID
	}

	groupID = "PASTEL"
	if ddGroupID != "" && !strings.EqualFold(ddGroupID, "NA") {
		groupID = ddGroupID
	}

	return groupID, subsetID
}
