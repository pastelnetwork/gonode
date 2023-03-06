package hermes

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pastelnetwork/gonode/hermes/service/hermes/domain"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/store"

	"github.com/pastelnetwork/gonode/hermes/service/hermes/scorer"

	_ "github.com/mattn/go-sqlite3" //go-sqlite3
	"github.com/pastelnetwork/gonode/common/errgroup"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	svc "github.com/pastelnetwork/gonode/hermes/service"
	"github.com/pastelnetwork/gonode/hermes/service/node"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	synchronizationIntervalSec = 5
	synchronizationTimeoutSec  = 60
	runTaskInterval            = 10 * time.Minute
	pastelIDRestartInterval    = 15 * time.Minute
	masterNodeSuccessfulStatus = "Masternode successfully started"
)

type service struct {
	config       *Config
	pastelClient pastel.Client
	p2p          node.HermesP2PInterface
	sn           node.SNClientInterface
	store        store.DDStore

	scorer             *scorer.Scorer
	isMasterNodeSynced bool
	latestBlockHeight  int

	currentNFTBlock    int
	currentActionBlock int
	currentBlockCount  int32
}

func toFloat64Array(data []float32) []float64 {
	ret := make([]float64, len(data))
	for idx, value := range data {
		ret[idx] = float64(value)
	}

	return ret
}

// Run starts task
func (s *service) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(5 * time.Second):
			if err := s.run(ctx); err != nil {
				if utils.IsContextErr(err) {
					log.WithContext(ctx).WithError(err).Error("closing hermes due to context err")
					return err
				}

				log.WithContext(ctx).WithError(err).Error("failed to run hermes, retrying.")
			} else {
				return nil
			}
		}
	}
}

func (s *service) run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, "hermes")
	if _, err := os.Stat(s.config.DataFile); os.IsNotExist(err) {
		return errors.Errorf("dataFile dd service not found: %w", err)
	}

	if err := s.waitSynchronization(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to initial wait synchronization")
	} else {
		s.isMasterNodeSynced = true
	}

	snAddr := fmt.Sprintf("%s:%d", s.config.SNHost, s.config.SNPort)
	log.WithContext(ctx).WithField("sn-addr", snAddr).Info("connecting with SN-Service")

	conn, err := s.sn.Connect(ctx, snAddr)
	if err != nil {
		return errors.Errorf("unable to connect with SN service: %w", err)
	}
	s.p2p = conn.HermesP2P()
	log.WithContext(ctx).Info("connection established with SN-Service")

	group, gctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return s.runTask(gctx)
	})

	/*group.Go(func() error {
		return s.CleanupInactiveTickets(gctx)
	})*/

	/*group.Go(func() error {
		return s.scorer.Start(gctx)
	})*/

	group.Go(func() error {
		return s.restartPastelID(gctx)
	})

	if err := group.Wait(); err != nil {
		log.WithContext(gctx).WithError(err).Errorf("First runTask() failed")
	}

	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(runTaskInterval):
			// Check if node is synchronized or not
			if !s.isMasterNodeSynced {
				if err := s.checkSynchronized(ctx); err != nil {
					log.WithContext(ctx).WithError(err).Debug("Failed to check synced status from master node")
					continue
				}

				log.WithContext(ctx).Debug("Done for waiting synchronization status")
				s.isMasterNodeSynced = true
			}

			group, gctx := errgroup.WithContext(ctx)
			group.Go(func() error {
				return s.runTask(gctx)
			})

			/*group.Go(func() error {
				return s.CleanupInactiveTickets(gctx)
			})*/

			/*group.Go(func() error {
				return s.scorer.Start(gctx)
			})*/

			if err := group.Wait(); err != nil {
				log.WithContext(gctx).WithError(err).Errorf("run task failed")
			}
		}
	}
}

func (s *service) checkSynchronized(ctx context.Context) error {
	st, err := s.pastelClient.MasterNodeStatus(ctx)
	if err != nil {
		return errors.Errorf("getMasterNodeStatus: %w", err)
	}

	if st == nil {
		return errors.New("empty status")
	}

	if st.Status == masterNodeSuccessfulStatus {
		return nil
	}

	return errors.Errorf("node not synced, status is %s", st.Status)
}

func (s *service) restartPastelID(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.WithContext(ctx).Info("context done being called in restartPastelID hermes service")
		case <-time.After(pastelIDRestartInterval):
			if isAvailable := s.checkNextBlockAvailable(ctx); !isAvailable {
				log.WithContext(ctx).Info("next block not available after 15 minutes, should restart pastelid")

				pastelD, err := getPastelDPath(ctx)
				if err != nil {
					log.WithContext(ctx).WithError(err).Error("Error retrieving pasteld path")
					continue
				}

				pastelCli, err := getPastelCliPath(ctx)
				if err != nil {
					log.WithContext(ctx).WithError(err).Error("Error retrieving pastel-cli path")
					continue
				}

				homeDir, err := os.UserHomeDir()
				if err != nil {
					log.WithContext(ctx).WithError(err).Error("Error getting home dir")
					continue

				}

				cmd := exec.Command(pastelCli, "stop")
				cmd.Dir = homeDir
				if err := cmd.Run(); err != nil {
					log.WithContext(ctx).WithError(err).Error("Error stopping pastel-cli")
					continue
				}

				log.WithContext(ctx).Info("pastel-cli has been stopped")
				time.Sleep(30 * time.Second)

				var extIP string
				if extIP, err = utils.GetExternalIPAddress(); err != nil {
					log.WithContext(ctx).WithError(err).Error("Could not get external IP address")
					continue
				}

				dataDir, err := getDataDir(ctx, filepath.Join(homeDir, ".pastel/supernode.yml"))
				if err != nil {
					log.WithContext(ctx).WithError(err).Error("Error retrieving data-dir from supernode.yml")
					continue
				}

				mnPrivKey, err := getMasternodePrivKey(ctx, filepath.Join(dataDir, "testnet3/masternode.conf"), extIP)
				if err != nil {
					log.WithContext(ctx).WithError(err).Error("Error retrieving masternode private key")
					continue
				}

				var pasteldArgs []string
				pasteldArgs = append(pasteldArgs,
					fmt.Sprintf("--datadir=%s", dataDir),
					fmt.Sprintf("--externalip=%s", extIP),
					fmt.Sprintf("--masternodeprivkey=%s", mnPrivKey),
					"--txindex=1",
					"--reindex",
					"--masternode",
					"--daemon",
				)

				log.WithContext(ctx).Infof("Starting -> %s %s", pastelD, strings.Join(pasteldArgs, " "))
				cmd = exec.Command(pastelD, pasteldArgs...)
				cmd.Dir = homeDir

				if err := cmd.Run(); err != nil {
					log.WithContext(ctx).WithError(err).Error("Error starting pastelid")
				}

				log.WithContext(ctx).Info("pasteld has been restarted")
			} else {
				log.WithContext(ctx).Info("block count has been updated")
			}
		}
	}
}

func getDataDir(ctx context.Context, confFilePath string) (dataDir string, err error) {
	log.WithContext(ctx).Info("opening the supernode.yml")
	file, err := os.OpenFile(confFilePath, os.O_RDWR, 0644)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Could not open supernode.yml - %s", confFilePath)
		return "", err
	}
	log.WithContext(ctx).Info("File opened")
	defer file.Close()

	log.WithContext(ctx).Info("Parsing data-dir from supernode.yml")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "work-dir:") {
			dataDir = line
			if err != nil {
				log.WithContext(ctx).Error(fmt.Sprintf("Could not parse work-dir: from file: - %s", confFilePath))
				return "", err
			}

			dataDirPath := strings.Split(dataDir, ": ")
			return dataDirPath[1], nil
		}
	}

	return dataDir, nil
}

func loadMasternodeConfFile(ctx context.Context, masternodeConfFilePath string) (map[string]domain.MasterNodeConf, error) {
	// Read ConfData from masternode.conf
	confFile, err := os.ReadFile(masternodeConfFilePath)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(fmt.Sprintf("Failed to read existing masternode.conf file - %s", masternodeConfFilePath))
		return nil, err
	}

	var conf map[string]domain.MasterNodeConf
	err = json.Unmarshal(confFile, &conf)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(fmt.Sprintf("Invalid existing masternode.conf file - %s", masternodeConfFilePath))
		return nil, err
	}

	return conf, nil
}
func getMasternodePrivKey(ctx context.Context, masterNodeConffilePath, extIP string) (string, error) {
	conf, err := loadMasternodeConfFile(ctx, masterNodeConffilePath)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to load existing masternode.conf file")
		return "", err
	}

	log.WithContext(ctx).Info("Attempting to load existing masternode.conf file using external IP address...")
	for mnName, mnConf := range conf {
		extAddrPort := strings.Split(mnConf.MnAddress, ":")
		extAddr := extAddrPort[0] // get Ext IP and Port
		if extAddr == extIP {
			log.WithContext(ctx).Info(fmt.Sprintf("Loading masternode.conf file using %s conf", mnName))
			return mnConf.MnPrivKey, nil
		}
	}

	return "", errors.New("not able to found private key")
}

// checkNextBlockAvailable calls pasteld and checks if a new block is available
func (s *service) checkNextBlockAvailable(ctx context.Context) bool {
	blockCount, err := s.pastelClient.GetBlockCount(ctx)
	if err != nil {
		return false
	}
	if blockCount > s.currentBlockCount {
		atomic.StoreInt32(&s.currentBlockCount, blockCount)
		return true
	}

	return false
}

func (s *service) waitSynchronization(ctx context.Context) error {
	checkTimeout := func(checked chan<- struct{}) {
		time.Sleep(synchronizationTimeoutSec * time.Second)
		close(checked)
	}

	timeoutCh := make(chan struct{})
	go checkTimeout(timeoutCh)

	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(synchronizationIntervalSec * time.Second):
			err := s.checkSynchronized(ctx)
			if err != nil {
				log.WithContext(ctx).WithError(err).Debug("Failed to check synced status from master node")
			} else {
				log.WithContext(ctx).Info("Done for waiting synchronization status")
				return nil
			}
		case <-timeoutCh:
			return errors.New("timeout expired")
		}
	}
}

// Utility function to get dd and fp file from an id hash, where the file should be stored
func (s *service) tryToGetFingerprintFileFromHash(ctx context.Context, hash string) (*pastel.DDAndFingerprints, error) {
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

func (s *service) runTask(ctx context.Context) error {
	log.WithContext(ctx).Info("getting Activation tickets")
	actTickets, err := s.pastelClient.ActTickets(ctx, pastel.ActTicketAll, s.latestBlockHeight)
	if err != nil {
		log.WithError(err).Error("unable to get act tickets - exit runtask now")
		return nil
	}
	if len(actTickets) == 0 {
		return nil
	}
	log.WithContext(ctx).WithField("count", len(actTickets)).Info("Act tickets retrieved")

	//track latest block height, but don't set it until we check all the nft reg tickets and the sense tickets.
	latestBlockHeight := s.latestBlockHeight
	lastKnownGoodHeight := s.latestBlockHeight

	//loop through nft tickets and store newly found nft reg tickets
	for i := 0; i < len(actTickets); i++ {
		if actTickets[i].Height <= s.latestBlockHeight {
			continue
		}

		regTicket, err := s.pastelClient.RegTicket(ctx, actTickets[i].ActTicketData.RegTXID)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("regTxid", actTickets[i].ActTicketData.RegTXID).
				WithField("act-Txid", actTickets[i].TXID).Error("unable to get reg ticket")
			continue
		}

		log.WithContext(ctx).WithField("txid", actTickets[i].ActTicketData.RegTXID).Info("Found new NFT ticket")

		decTicket, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		regTicket.RegTicketData.NFTTicketData = *decTicket

		ddFPIDs := regTicket.RegTicketData.NFTTicketData.AppTicketData.DDAndFingerprintsIDs

		ddAndFpFromTicket := &pastel.DDAndFingerprints{}
		//Get the dd and fp file from the ticket
		for _, id := range ddFPIDs {
			ddAndFpFromTicket, err = s.tryToGetFingerprintFileFromHash(ctx, id)
			if err != nil {
				//probably too verbose even for debug.
				log.WithContext(ctx).WithField("error", err).WithField("txid", actTickets[i].ActTicketData.RegTXID).WithField("id", id).Debug("Could not get the fingerprint for this file hash")
				continue
			}
			break
		}

		if ddAndFpFromTicket == nil {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).
				WithField("NFTRegTicket", regTicket.RegTicketData.NFTTicketData).Debugf("None of the dd and fp id files for this nft reg ticket could be properly unmarshalled")
			continue
		}
		if ddAndFpFromTicket.HashOfCandidateImageFile == "" {
			log.WithContext(ctx).WithField("NFTRegTicket", regTicket.RegTicketData.NFTTicketData).Debugf("This NFT Reg ticket's DDAndFp struct has no HashOfCandidateImageFile, perhaps it's an older version.")
			continue
		}

		existsInDatabase, err := s.store.IfFingerprintExists(ctx, ddAndFpFromTicket.HashOfCandidateImageFile)
		if existsInDatabase {
			log.WithContext(ctx).WithField("txid", regTicket.TXID).Info("fingerprints already exist in database, skipping")
			//can't directly update latest block height from here - if there's another ticket in this block we don't want to skip
			if actTickets[i].Height > lastKnownGoodHeight {
				lastKnownGoodHeight = regTicket.RegTicketData.NFTTicketData.BlockNum
			}
			continue
		}
		if err != nil {
			log.WithContext(ctx).Error("Could not properly query the dd database for this hash")
			continue
		}

		//make sure ImageFingerprintOfCnadidateImageFile exists.
		// this could fail if the ticket is an older version of the DDAndFingerprints struct, so we will continue to next fingerprint
		if ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile == nil {
			log.WithContext(ctx).WithField("NFTRegTicket", regTicket.RegTicketData.NFTTicketData).Debugf("This NFT Reg ticket's DDAndFp struct has no ImageFingerprintOfCandidateImageFile, perhaps it's an older version.")
			continue
		}
		if len(ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile) < 1 {
			log.WithContext(ctx).WithField("NFTRegTicket", regTicket.RegTicketData.NFTTicketData).Debugf("This NFT Reg ticket's DDAndFp struct's ImageFingerprintOfCandidateImageFile is zero length, perhaps it's an older version.")
			continue
		}

		// thumbnailHash := regTicket.RegTicketData.NFTTicketData.AppTicketData.Thumbnail1Hash
		// thumbnail, err := s.p2pClient.Retrieve(ctx, string(thumbnailHash))
		// if err != nil {
		// 	log.WithContext(ctx).WithField("thumbnailHash", regTicket.RegTicketData.NFTTicketData.AppTicketData.Thumbnail1Hash).Warnf("Could not get the thumbnail with this hash for nftticketdata")
		// 	continue
		// }

		collection := "PASTEL"
		if regTicket.RegTicketData.NFTTicketData.AppTicketData.NFTSeriesName != "" {
			collection = regTicket.RegTicketData.NFTTicketData.AppTicketData.NFTSeriesName
		}

		if err := s.store.StoreFingerprint(ctx, &domain.DDFingerprints{
			Sha256HashOfArtImageFile:           ddAndFpFromTicket.HashOfCandidateImageFile,
			ImageFingerprintVector:             toFloat64Array(ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile),
			DatetimeFingerprintAddedToDatabase: time.Now().Format("2006-01-02 15:04:05"),
			PathToArtImageFile:                 "",
			ImageThumbnailAsBase64:             "",
			RequestType:                        typeMapper(regTicket.RegTicketData.Type),
			IDString:                           collection,
		}); err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to store fingerprint")
		}
		if regTicket.RegTicketData.NFTTicketData.BlockNum > latestBlockHeight {
			latestBlockHeight = regTicket.RegTicketData.NFTTicketData.BlockNum
		}
	}
	//loop through action tickets and store newly found nft reg tickets

	senseRegTickets, err := s.pastelClient.ActionTickets(ctx)
	if err != nil {
		log.WithError(err).Errorf("get registered ticket - exit runTask")
		return nil
	}

	if len(senseRegTickets) == 0 {
		return nil
	}
	for i := 0; i < len(senseRegTickets); i++ {
		if senseRegTickets[i].ActionTicketData.CalledAt <= s.latestBlockHeight {
			continue
		}
		log.WithContext(ctx).WithField("txid", senseRegTickets[i].TXID).Info("Found sense ticket, checking activation...")

		idt, err := s.pastelClient.FindActionActByActionRegTxid(ctx, senseRegTickets[i].TXID)
		if err != nil {
			log.WithContext(ctx).WithField("reg-txid", senseRegTickets[i].TXID).WithError(err).
				Error("Failed to find action act by action reg txid")
			continue
		}

		if idt == nil || idt.TXID == "" {
			continue
		}

		log.WithContext(ctx).WithField("txid", senseRegTickets[i].TXID).Info("Found activated sense ticket")

		decTicket, err := pastel.DecodeActionTicket(senseRegTickets[i].ActionTicketData.ActionTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		senseRegTickets[i].ActionTicketData.ActionTicketData = *decTicket

		if senseRegTickets[i].ActionTicketData.ActionTicketData.ActionType != "sense" {
			continue
		}

		senseTicket, err := senseRegTickets[i].ActionTicketData.ActionTicketData.APISenseTicket()
		if err != nil {
			log.WithContext(ctx).WithField("senseRegTickets.ActionTicketData", senseRegTickets[i].ActionTicketData).Debugf("Could not get sense ticket for action ticket data")
			continue
		}

		ddAndFpFromTicket := &pastel.DDAndFingerprints{}
		//Get the dd and fp file from the ticket
		for _, id := range senseTicket.DDAndFingerprintsIDs {
			ddAndFpFromTicket, err = s.tryToGetFingerprintFileFromHash(ctx, id)
			if err != nil {
				//probably too verbose even for debug.
				log.WithContext(ctx).WithField("error", err).WithField("sense-txid", senseRegTickets[i].TXID).Debug("Could not get the fingerprint for this file hash")
				continue
			}
			break
		}
		if ddAndFpFromTicket == nil {
			log.WithContext(ctx).WithField("senseTicket", senseTicket).Debugf("None of the dd and fp id files for this sense reg ticket could be properly unmarshalled")
			continue
		}
		if ddAndFpFromTicket.HashOfCandidateImageFile == "" {
			log.WithContext(ctx).WithField("senseTicket", senseTicket).Debugf("This NFT sense ticket's DDAndFp struct has no HashOfCandidateImageFile, perhaps it's an older version.")
			continue
		}

		existsInDatabase, err := s.store.IfFingerprintExists(ctx, ddAndFpFromTicket.HashOfCandidateImageFile)
		if existsInDatabase {
			log.WithContext(ctx).WithField("txid", senseRegTickets[i].TXID).Info("Fingerprint exists in database, skipping...")
			//can't directly update latest block height from here - if there's another ticket in this block we don't want to skip
			if senseRegTickets[i].ActionTicketData.ActionTicketData.BlockNum > lastKnownGoodHeight {
				lastKnownGoodHeight = senseRegTickets[i].ActionTicketData.ActionTicketData.BlockNum
			}
			continue
		}
		if err != nil {
			log.WithContext(ctx).WithField("hashOfCandidateImageFile", ddAndFpFromTicket.HashOfCandidateImageFile).Error("Could not properly query the dd database for this hash")
			continue
		}

		// this could fail if the ticket is an older version of the DDAndFingerprints struct, so we will continue to next fingerprint
		if ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile == nil {
			log.WithContext(ctx).WithField("senseTicket", senseRegTickets[i].ActionTicketData).Debugf("This sense ticket's DDAndFp struct has no ImageFingerprintOfCandidateImageFile, perhaps it's an older version.")
			continue
		}
		if len(ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile) < 1 {
			log.WithContext(ctx).WithField("senseTicket", senseRegTickets[i].ActionTicketData).Debugf("This sense reg ticket's DDAndFp struct's ImageFingerprintOfCandidateImageFile is zero length, perhaps it's an older version.")
			continue
		}

		collection := "PASTEL"
		if ddAndFpFromTicket.OpenAPIGroupIDString != "" && !strings.EqualFold(ddAndFpFromTicket.OpenAPIGroupIDString, "NA") {
			collection = ddAndFpFromTicket.OpenAPISubsetIDString
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
			RequestType:                                typeMapper(senseRegTickets[i].ActionTicketData.Type),
			IDString:                                   collection,
			OpenAPIGroupIDString:                       groupID,
			CollectionNameString:                       ddAndFpFromTicket.CollectionNameString,
			DoesNotImpactTheFollowingCollectionsString: ddAndFpFromTicket.DoesNotImpactTheFollowingCollectionStrings,
		}); err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to store fingerprint")
		}
		if senseRegTickets[i].ActionTicketData.ActionTicketData.BlockNum > latestBlockHeight {
			latestBlockHeight = senseRegTickets[i].ActionTicketData.ActionTicketData.BlockNum
		}
	}

	if lastKnownGoodHeight > latestBlockHeight {
		s.latestBlockHeight = lastKnownGoodHeight
	} else {
		s.latestBlockHeight = latestBlockHeight
	}

	log.WithContext(ctx).WithField("latest blockheight", s.latestBlockHeight).Debugf("hermes successfully scanned to latest block height")

	return nil
}

// Stats return status of dupe detection
func (s *service) Stats(ctx context.Context) (map[string]interface{}, error) {
	stats := map[string]interface{}{}

	// Get last inserted item
	lastItem, err := s.store.GetLatestFingerprints(ctx)
	if err != nil {
		return nil, errors.Errorf("getLatestFingerprint: %w", err)
	}
	stats["last_insert_time"] = lastItem.DatetimeFingerprintAddedToDatabase

	// Get total of records
	recordCount, err := s.store.GetFingerprintsCount(ctx)
	if err != nil {
		return nil, errors.Errorf("getRecordCount: %w", err)
	}
	stats["record_count"] = recordCount

	fi, err := os.Stat(s.config.DataFile)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get dd db size")
	} else {
		stats["db_size"] = utils.BytesToMB(uint64(fi.Size()))
	}

	return stats, nil
}

func getPastelDPath(ctx context.Context) (path string, err error) {
	//create command
	findCmd := exec.Command("find", ".", "-print")
	grepCmd := exec.Command("grep", "-x", "./pastel/pasteld")

	findCmd.Dir, err = os.UserHomeDir()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error getting home dir")
		return "", err
	}

	//make a pipe and set the input and output to reader and writer
	reader, writer := io.Pipe()
	var buf bytes.Buffer

	findCmd.Stdout = writer
	grepCmd.Stdin = reader

	//cache the output of "grep" to memory
	grepCmd.Stdout = &buf

	//starting the commands
	findCmd.Start()
	grepCmd.Start()

	//waiting for commands to complete and close the reader & writer
	findCmd.Wait()
	writer.Close()

	grepCmd.Wait()
	reader.Close()

	pathWithEscapeCharacter := buf.String()
	return strings.Replace(pathWithEscapeCharacter, "\n", "", 1), nil
}

func getPastelCliPath(ctx context.Context) (path string, err error) {
	//create command
	findCmd := exec.Command("find", ".", "-print")
	grepCmd := exec.Command("grep", "-x", "./pastel/pastel-cli")

	findCmd.Dir, err = os.UserHomeDir()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error getting home dir")
		return "", err
	}

	//make a pipe and set the input and output to reader and writer
	reader, writer := io.Pipe()
	var buf bytes.Buffer

	findCmd.Stdout = writer
	grepCmd.Stdin = reader

	//cache the output of "grep" to memory
	grepCmd.Stdout = &buf

	//starting the commands
	findCmd.Start()
	grepCmd.Start()

	//waiting for commands to complete and close the reader & writer
	findCmd.Wait()
	writer.Close()

	grepCmd.Wait()
	reader.Close()

	pathWithEscapeCharacter := buf.String()
	return strings.Replace(pathWithEscapeCharacter, "\n", "", 1), nil
}

// NewService returns a new ddscan service
func NewService(config *Config, pastelClient pastel.Client, sn node.SNClientInterface) (svc.SvcInterface, error) {
	store, err := store.NewSQLiteStore(config.DataFile)
	if err != nil {
		return nil, fmt.Errorf("unable to initialise database: %w", err)
	}

	scorerConfig := &scorer.Config{
		NumberSuperNodes:          config.NumberSuperNodes,
		ConnectToNodeTimeout:      config.ConnectToNodeTimeout,
		ConnectToNextNodeDelay:    config.ConnectToNextNodeDelay,
		AcceptNodesTimeout:        config.AcceptNodesTimeout,
		CreatorPastelID:           config.CreatorPastelID,
		CreatorPastelIDPassphrase: config.CreatorPastelIDPassphrase,
	}

	return &service{
		config:             config,
		pastelClient:       pastelClient,
		store:              store,
		sn:                 sn,
		currentNFTBlock:    1,
		currentActionBlock: 1,
		scorer:             scorer.New(scorerConfig, pastelClient, sn, store),
	}, nil
}

func typeMapper(val string) string {
	if val == "action-reg" {
		return "SENSE"
	}
	if val == "nft-reg" {
		return "NFT"
	}

	return val
}
