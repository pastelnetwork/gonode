package fingerprint

import (
	"bytes"
	"context"
	"encoding/json"
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

		log.WithContext(ctx).WithField("txid", senseActTickets[i].TXID).WithField("height", senseActTickets[i].Height).
			Info("Found sense activation ticket, checking reg ticket...")

		regTicket, err := s.pastelClient.ActionRegTicket(ctx, senseActTickets[i].ActTicketData.RegTXID)
		if err != nil {
			log.WithContext(ctx).WithField("reg-txid", senseActTickets[i].ActTicketData.RegTXID).WithError(err).
				Error("Failed to find action act by action reg txid")
			continue
		}

		if regTicket.ActionTicketData.ActionType != pastel.ActionTypeSense {
			continue
		}

		log.WithContext(ctx).WithField("txid", regTicket.TXID).WithField("act-txid", senseActTickets[i].TXID).Info("Found activated sense ticket")

		decTicket, err := pastel.DecodeActionTicket(regTicket.ActionTicketData.ActionTicket)
		if err != nil {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).WithField("act-txid", senseActTickets[i].TXID).
				WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		regTicket.ActionTicketData.ActionTicketData = *decTicket

		senseTicket, err := regTicket.ActionTicketData.ActionTicketData.APISenseTicket()
		if err != nil {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).WithField("act-txid", senseActTickets[i].TXID).Error("Could not get sense ticket for action ticket data")
			continue
		}

		p2pServiceRunning := true
		ddAndFpFromTicket := &pastel.DDAndFingerprints{}
		//Get the dd and fp file from the ticket
		for _, id := range senseTicket.DDAndFingerprintsIDs {
			ddAndFpFromTicket, err = s.tryToGetFingerprintFileFromHash(ctx, id)
			if err != nil {
				if strings.Contains(err.Error(), "p2p service is not running") {
					p2pServiceRunning = false
					break
				}

				//probably too verbose even for debug.
				log.WithContext(ctx).WithField("error", err).WithField("txid", regTicket.TXID).WithField("act-txid", senseActTickets[i].TXID).Error("Could not get the fingerprint for this file hash")
				continue
			}
			break
		}

		if !p2pServiceRunning {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).WithField("act-txid", senseActTickets[i].TXID).
				Info("P2P service is not running, so we can't get the fingerprint for this file hash, stopping this run of the task")

			break
		}

		if ddAndFpFromTicket == nil {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).WithField("act-txid", senseActTickets[i].TXID).Error("None of the dd and fp id files for this sense reg ticket could be properly unmarshalled")
			continue
		}
		if ddAndFpFromTicket.HashOfCandidateImageFile == "" {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).WithField("act-txid", senseActTickets[i].TXID).Error("This NFT sense ticket's DDAndFp struct has no HashOfCandidateImageFile, perhaps it's an older version.")
			continue
		}
		log.WithContext(ctx).WithField("image_hash", ddAndFpFromTicket.HashOfCandidateImageFile).
			WithField("txid", regTicket.TXID).
			WithField("act-txid", senseActTickets[i].TXID).
			Info("image hash retrieved")

		existsInDatabase, err := s.store.IfFingerprintExists(ctx, ddAndFpFromTicket.HashOfCandidateImageFile)
		if existsInDatabase {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).WithField("act-txid", senseActTickets[i].TXID).
				Info("Fingerprint exists in database, skipping...")

			//can't directly update latest block height from here - if there's another ticket in this block we don't want to skip
			if senseActTickets[i].Height > lastKnownGoodHeight {
				lastKnownGoodHeight = senseActTickets[i].Height
			}
			continue
		}
		if err != nil {
			log.WithContext(ctx).WithField("hashOfCandidateImageFile", ddAndFpFromTicket.HashOfCandidateImageFile).
				WithField("txid", regTicket.TXID).WithField("act-txid", senseActTickets[i].TXID).Error("Could not properly query the dd database for this hash")
			continue
		}

		// this could fail if the ticket is an older version of the DDAndFingerprints struct, so we will continue to next fingerprint
		if ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile == nil {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).WithField("act-txid", senseActTickets[i].TXID).Info("This sense ticket's DDAndFp struct has no ImageFingerprintOfCandidateImageFile, perhaps it's an older version.")
			continue
		}
		if len(ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile) < 1 {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).WithField("act-txid", senseActTickets[i].TXID).Info("This sense reg ticket's DDAndFp struct's ImageFingerprintOfCandidateImageFile is zero length, perhaps it's an older version.")
			continue
		}

		collection := "PASTEL"
		if ddAndFpFromTicket.OpenAPISubsetID != "" && !strings.EqualFold(ddAndFpFromTicket.OpenAPISubsetID, "NA") {
			collection = ddAndFpFromTicket.OpenAPISubsetID
		}

		groupID := "PASTEL"
		if ddAndFpFromTicket.OpenAPIGroupIDString != "" && !strings.EqualFold(ddAndFpFromTicket.OpenAPIGroupIDString, "NA") {
			groupID = ddAndFpFromTicket.OpenAPIGroupIDString
		}

		if err := s.store.StoreFingerprint(ctx, &domain.DDFingerprints{
			Sha256HashOfArtImageFile:                   ddAndFpFromTicket.HashOfCandidateImageFile,
			ImageFingerprintVector:                     toFloat64Array(ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile),
			DatetimeFingerprintAddedToDatabase:         time.Now().Format("2006-01-02 15:04:05"),
			PathToArtImageFile:                         ddAndFpFromTicket.ImageFilePath,
			ImageThumbnailAsBase64:                     ddAndFpFromTicket.CandidateImageThumbnailWebpAsBase64String,
			RequestType:                                typeMapper(regTicket.ActionTicketData.Type),
			IDString:                                   collection,
			OpenAPIGroupIDString:                       groupID,
			CollectionNameString:                       ddAndFpFromTicket.CollectionNameString,
			DoesNotImpactTheFollowingCollectionsString: ddAndFpFromTicket.DoesNotImpactTheFollowingCollectionStrings,
		}); err != nil {
			log.WithContext(ctx).WithError(err).WithField("txid", regTicket.TXID).WithField("act-txid", senseActTickets[i].TXID).Error("Failed to store fingerprint")
			continue
		}

		if senseActTickets[i].Height > lastKnownGoodHeight {
			lastKnownGoodHeight = senseActTickets[i].Height
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

		log.WithContext(ctx).WithField("act-txid", actTickets[i].TXID).WithField("regTxid", actTickets[i].ActTicketData.RegTXID).
			Info("Found new NFT ticket")

		regTicket, err := s.pastelClient.RegTicket(ctx, actTickets[i].ActTicketData.RegTXID)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("regTxid", actTickets[i].ActTicketData.RegTXID).
				WithField("act-Txid", actTickets[i].TXID).Error("unable to get reg ticket")
			continue
		}

		log.WithContext(ctx).WithField("txid", actTickets[i].ActTicketData.RegTXID).Info("Found Reg Ticket for NFT-Act ticket")

		decTicket, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
		if err != nil {
			log.WithContext(ctx).WithField("txid", actTickets[i].ActTicketData.RegTXID).WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		regTicket.RegTicketData.NFTTicketData = *decTicket

		ddFPIDs := regTicket.RegTicketData.NFTTicketData.AppTicketData.DDAndFingerprintsIDs

		ddAndFpFromTicket := &pastel.DDAndFingerprints{}
		//Get the dd and fp file from the ticket

		p2pServiceRunning := true
		for _, id := range ddFPIDs {
			ddAndFpFromTicket, err = s.tryToGetFingerprintFileFromHash(ctx, id)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("Failed to get dd and fp file from ticket")
				if strings.Contains(err.Error(), "p2p service is not running") {
					p2pServiceRunning = false
					break
				}

				//probably too verbose even for debug.
				log.WithContext(ctx).WithField("error", err).WithField("txid", actTickets[i].ActTicketData.RegTXID).WithField("id", id).Error("Could not get the fingerprint for this file hash")
				continue
			}
			break
		}
		if !p2pServiceRunning {
			log.WithContext(ctx).WithField("txid", actTickets[i].ActTicketData.RegTXID).
				Info("P2P service is not running, so we can't get the fingerprint for this file hash, stopping this run of the task")

			break
		}

		if ddAndFpFromTicket == nil {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).
				WithField("txid", actTickets[i].ActTicketData.RegTXID).Info("None of the dd and fp id files for this nft reg ticket could be properly unmarshalled")
			continue
		}
		if ddAndFpFromTicket.HashOfCandidateImageFile == "" {
			log.WithContext(ctx).WithField("txid", actTickets[i].ActTicketData.RegTXID).Info("This NFT Reg ticket's DDAndFp struct has no HashOfCandidateImageFile, perhaps it's an older version.")
			continue
		}

		log.WithContext(ctx).WithField("image_hash", ddAndFpFromTicket.HashOfCandidateImageFile).
			WithField("txid", regTicket.TXID).
			WithField("act-txid", actTickets[i].TXID).
			Info("image hash retrieved")

		existsInDatabase, err := s.store.IfFingerprintExists(ctx, ddAndFpFromTicket.HashOfCandidateImageFile)
		if existsInDatabase {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).Info("fingerprint already exist in database, skipping")
			//can't directly update latest block height from here - if there's another ticket in this block we don't want to skip
			if actTickets[i].Height > lastKnownGoodHeight {
				lastKnownGoodHeight = actTickets[i].Height
			}
			continue
		}
		if err != nil {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).Error("Could not properly query the dd database for this hash")
			continue
		}

		//make sure ImageFingerprintOfCnadidateImageFile exists.
		// this could fail if the ticket is an older version of the DDAndFingerprints struct, so we will continue to next fingerprint
		if ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile == nil {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).Info("This NFT Reg ticket's DDAndFp struct has no ImageFingerprintOfCandidateImageFile, perhaps it's an older version.")
			continue
		}
		if len(ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile) < 1 {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).Info("This NFT Reg ticket's DDAndFp struct's ImageFingerprintOfCandidateImageFile is zero length, perhaps it's an older version.")
			continue
		}

		collection := "PASTEL"
		if regTicket.RegTicketData.NFTTicketData.AppTicketData.NFTSeriesName != "" {
			collection = regTicket.RegTicketData.NFTTicketData.AppTicketData.NFTSeriesName
		}

		groupID := "PASTEL"
		if ddAndFpFromTicket.OpenAPIGroupIDString != "" && !strings.EqualFold(ddAndFpFromTicket.OpenAPIGroupIDString, "NA") {
			groupID = ddAndFpFromTicket.OpenAPIGroupIDString
		}

		if err := s.store.StoreFingerprint(ctx, &domain.DDFingerprints{
			Sha256HashOfArtImageFile:                   ddAndFpFromTicket.HashOfCandidateImageFile,
			ImageFingerprintVector:                     toFloat64Array(ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile),
			DatetimeFingerprintAddedToDatabase:         time.Now().Format("2006-01-02 15:04:05"),
			PathToArtImageFile:                         ddAndFpFromTicket.ImageFilePath,
			ImageThumbnailAsBase64:                     ddAndFpFromTicket.CandidateImageThumbnailWebpAsBase64String,
			RequestType:                                typeMapper(regTicket.RegTicketData.Type),
			IDString:                                   collection,
			OpenAPIGroupIDString:                       groupID,
			CollectionNameString:                       ddAndFpFromTicket.CollectionNameString,
			DoesNotImpactTheFollowingCollectionsString: ddAndFpFromTicket.DoesNotImpactTheFollowingCollectionStrings,
		}); err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to store fingerprint")
			continue
		}
		if actTickets[i].Height > lastKnownGoodHeight {
			lastKnownGoodHeight = actTickets[i].Height
		}
	}
	//loop through action tickets and store newly found nft reg tickets

	if lastKnownGoodHeight > s.latestNFTBlockHeight {
		s.latestNFTBlockHeight = lastKnownGoodHeight
	}

	log.WithContext(ctx).WithField("latest NFT blockheight", s.latestNFTBlockHeight).Debugf("hermes successfully scanned to latest block height")

	return nil
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