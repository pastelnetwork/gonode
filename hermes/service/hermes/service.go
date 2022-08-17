package hermes

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //go-sqlite3
	"github.com/pastelnetwork/gonode/common/errgroup"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	svc "github.com/pastelnetwork/gonode/hermes/service"
	"github.com/pastelnetwork/gonode/hermes/service/node"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/sbinet/npyio"
	"gonum.org/v1/gonum/mat"
)

const (
	synchronizationIntervalSec       = 5
	synchronizationTimeoutSec        = 60
	runTaskInterval                  = 2 * time.Minute
	masterNodeSuccessfulStatus       = "Masternode successfully started"
	getLatestFingerprintStatement    = `SELECT * FROM image_hash_to_image_fingerprint_table ORDER BY datetime_fingerprint_added_to_database DESC LIMIT 1`
	getFingerprintFromHashStatement  = `SELECT * FROM image_hash_to_image_fingerprint_table WHERE sha256_hash_of_art_image_file = ?`
	insertFingerprintStatement       = `INSERT INTO image_hash_to_image_fingerprint_table(sha256_hash_of_art_image_file, path_to_art_image_file, new_model_image_fingerprint_vector, datetime_fingerprint_added_to_database, thumbnail_of_image, request_type, open_api_subset_id_string) VALUES(?,?,?,?,?,?,?)`
	getNumberOfFingerprintsStatement = `SELECT COUNT(*) as count FROM image_hash_to_image_fingerprint_table`
)

type dbFingerprints struct {
	Sha256HashOfArtImageFile           string `db:"sha256_hash_of_art_image_file,omitempty"`
	PathToArtImageFile                 string `db:"path_to_art_image_file,omitempty"`
	ImageFingerprintVector             []byte `db:"new_model_image_fingerprint_vector,omitempty"`
	DatetimeFingerprintAddedToDatabase string `db:"datetime_fingerprint_added_to_database,omitempty"`
	ImageThumbnailAsBase64             string `db:"thumbnail_of_image,omitempty"`
	RequestType                        string `db:"request_type,omitempty"`
	IDString                           string `db:"open_api_subset_id_string,omitempty"`
}

type dupeDetectionFingerprints struct {
	Sha256HashOfArtImageFile           string    `json:"sha256_hash_of_art_image_file,omitempty"`
	PathToArtImageFile                 string    `json:"path_to_art_image_file,omitempty"`
	ImageFingerprintVector             []float64 `json:"new_model_image_fingerprint_vector,omitempty"`
	DatetimeFingerprintAddedToDatabase string    `json:"datetime_fingerprint_added_to_database,omitempty"`
	ImageThumbnailAsBase64             string    `json:"thumbnail_of_image,omitempty"`
	RequestType                        string    `json:"request_type,omitempty"`
	IDString                           string    `json:"open_api_subset_id_string,omitempty"`
}

type service struct {
	config             *Config
	pastelClient       pastel.Client
	p2p                node.HermesP2PInterface
	sn                 node.SNClientInterface
	db                 *sqlx.DB
	isMasterNodeSynced bool
	latestBlockHeight  int

	currentNFTBlock    int
	currentActionBlock int
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
					log.WithContext(ctx).WithError(err).Warn("Failed to check synced status from master node")
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
				log.WithContext(ctx).WithError(err).Warn("Failed to check synced status from master node")
			} else {
				log.WithContext(ctx).Info("Done for waiting synchronization status")
				return nil
			}
		case <-timeoutCh:
			return errors.New("timeout expired")
		}
	}
}

func (r *dbFingerprints) toDomain() (*dupeDetectionFingerprints, error) {
	f := bytes.NewBuffer(r.ImageFingerprintVector)

	var fp []float64
	if err := npyio.Read(f, &fp); err != nil {
		return nil, errors.New("Failed to convert npy to float64")
	}

	return &dupeDetectionFingerprints{
		Sha256HashOfArtImageFile:           r.Sha256HashOfArtImageFile,
		PathToArtImageFile:                 r.PathToArtImageFile,
		ImageFingerprintVector:             fp,
		DatetimeFingerprintAddedToDatabase: r.DatetimeFingerprintAddedToDatabase,
		ImageThumbnailAsBase64:             r.ImageThumbnailAsBase64,
		RequestType:                        r.RequestType,
		IDString:                           r.IDString,
	}, nil
}

func (s *service) getLatestFingerprint(ctx context.Context) (*dupeDetectionFingerprints, error) {
	r := dbFingerprints{}
	err := s.db.Get(&r, getLatestFingerprintStatement)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("dd database is empty")
		}

		return nil, fmt.Errorf("failed to get record by key %w", err)
	}

	dd, err := r.toDomain()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("converting db data to dd-fingerprints struct failure")
	}

	return dd, nil
}

func (s *service) checkIfFingerprintExistsInDatabase(_ context.Context, hash string) (bool, error) {
	r := dbFingerprints{}
	err := s.db.Get(&r, getFingerprintFromHashStatement, hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, fmt.Errorf("failed to get record by key %w : key: %s", err, hash)
	}

	if r.Sha256HashOfArtImageFile == "" {
		return false, nil
	}

	return true, nil
}

func (s *service) storeFingerprint(ctx context.Context, input *dupeDetectionFingerprints) error {
	encodeFloat2Npy := func(v []float64) ([]byte, error) {
		// create numpy matrix Nx1
		m := mat.NewDense(len(v), 1, v)
		f := bytes.NewBuffer(nil)
		if err := npyio.Write(f, m); err != nil {
			return nil, errors.Errorf("encode to npy: %w", err)
		}
		return f.Bytes(), nil
	}

	fp, err := encodeFloat2Npy(input.ImageFingerprintVector)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`INSERT INTO image_hash_to_image_fingerprint_table(sha256_hash_of_art_image_file,
		 path_to_art_image_file, new_model_image_fingerprint_vector, datetime_fingerprint_added_to_database,
		  thumbnail_of_image, request_type, open_api_subset_id_string) VALUES(?,?,?,?,?,?,?)`, input.Sha256HashOfArtImageFile,
		input.PathToArtImageFile, fp, input.DatetimeFingerprintAddedToDatabase, input.ImageThumbnailAsBase64, input.RequestType, input.IDString)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to insert fingerprint record")
		return err
	}

	return nil
}

//Utility function to get dd and fp file from an id hash, where the file should be stored
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
	log.WithContext(ctx).Info("getting Reg tickets")

	nftRegTickets, err := s.pastelClient.RegTickets(ctx)
	if err != nil {
		return errors.Errorf("get registered ticket: %w", err)
	}

	if len(nftRegTickets) == 0 {
		return nil
	}
	log.WithContext(ctx).WithField("count", len(nftRegTickets)).Info("Reg tickets retrieved")

	//track latest block height, but don't set it until we check all the nft reg tickets and the sense tickets.
	latestBlockHeight := s.latestBlockHeight
	lastKnownGoodHeight := s.latestBlockHeight

	//loop through nft tickets and store newly found nft reg tickets
	for i := 0; i < len(nftRegTickets); i++ {
		if nftRegTickets[i].Height <= s.latestBlockHeight {
			continue
		}

		decTicket, err := pastel.DecodeNFTTicket(nftRegTickets[i].RegTicketData.NFTTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		nftRegTickets[i].RegTicketData.NFTTicketData = *decTicket

		ddFPIDs := nftRegTickets[i].RegTicketData.NFTTicketData.AppTicketData.DDAndFingerprintsIDs

		ddAndFpFromTicket := &pastel.DDAndFingerprints{}
		//Get the dd and fp file from the ticket
		for _, id := range ddFPIDs {
			ddAndFpFromTicket, err = s.tryToGetFingerprintFileFromHash(ctx, id)
			if err != nil {
				//probably too verbose even for debug.
				//log.WithContext(ctx).WithField("error", err).WithField("id", id).Debug("Could not get the fingerprint for this file hash")
				continue
			}
			break
		}
		if ddAndFpFromTicket == nil {
			log.WithContext(ctx).WithField("txid", nftRegTickets[i].TXID).
				WithField("NFTRegTicket", nftRegTickets[i].RegTicketData.NFTTicketData).Warnf("None of the dd and fp id files for this nft reg ticket could be properly unmarshalled")
			continue
		}
		if ddAndFpFromTicket.HashOfCandidateImageFile == "" {
			log.WithContext(ctx).WithField("NFTRegTicket", nftRegTickets[i].RegTicketData.NFTTicketData).Warnf("This NFT Reg ticket's DDAndFp struct has no HashOfCandidateImageFile, perhaps it's an older version.")
			continue
		}

		existsInDatabase, err := s.checkIfFingerprintExistsInDatabase(ctx, ddAndFpFromTicket.HashOfCandidateImageFile)
		if existsInDatabase {
			log.WithContext(ctx).WithField("hashOfCandidateImageFile", ddAndFpFromTicket.HashOfCandidateImageFile).Debug("Found hash of candidate image file already exists in database, not adding.")
			//can't directly update latest block height from here - if there's another ticket in this block we don't want to skip
			if nftRegTickets[i].RegTicketData.NFTTicketData.BlockNum > lastKnownGoodHeight {
				lastKnownGoodHeight = nftRegTickets[i].RegTicketData.NFTTicketData.BlockNum
			}
			continue
		}
		if err != nil {
			log.WithContext(ctx).WithField("hashOfCandidateImageFile", ddAndFpFromTicket.HashOfCandidateImageFile).Error("Could not properly query the dd database for this hash")
			continue
		}

		//make sure ImageFingerprintOfCnadidateImageFile exists.
		// this could fail if the ticket is an older version of the DDAndFingerprints struct, so we will continue to next fingerprint
		if ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile == nil {
			log.WithContext(ctx).WithField("NFTRegTicket", nftRegTickets[i].RegTicketData.NFTTicketData).Warnf("This NFT Reg ticket's DDAndFp struct has no ImageFingerprintOfCandidateImageFile, perhaps it's an older version.")
			continue
		}
		if len(ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile) < 1 {
			log.WithContext(ctx).WithField("NFTRegTicket", nftRegTickets[i].RegTicketData.NFTTicketData).Warnf("This NFT Reg ticket's DDAndFp struct's ImageFingerprintOfCandidateImageFile is zero length, perhaps it's an older version.")
			continue
		}

		// thumbnailHash := nftRegTickets[i].RegTicketData.NFTTicketData.AppTicketData.Thumbnail1Hash
		// thumbnail, err := s.p2pClient.Retrieve(ctx, string(thumbnailHash))
		// if err != nil {
		// 	log.WithContext(ctx).WithField("thumbnailHash", nftRegTickets[i].RegTicketData.NFTTicketData.AppTicketData.Thumbnail1Hash).Warnf("Could not get the thumbnail with this hash for nftticketdata")
		// 	continue
		// }

		if err := s.storeFingerprint(ctx, &dupeDetectionFingerprints{
			Sha256HashOfArtImageFile:           ddAndFpFromTicket.HashOfCandidateImageFile,
			ImageFingerprintVector:             toFloat64Array(ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile),
			DatetimeFingerprintAddedToDatabase: time.Now().Format("2006-01-02 15:04:05"),
			PathToArtImageFile:                 "",
			ImageThumbnailAsBase64:             "",
			RequestType:                        nftRegTickets[i].RegTicketData.Type,
			IDString:                           "",
		}); err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to store fingerprint")
		}
		if nftRegTickets[i].RegTicketData.NFTTicketData.BlockNum > latestBlockHeight {
			latestBlockHeight = nftRegTickets[i].RegTicketData.NFTTicketData.BlockNum
		}
	}
	//loop through action tickets and store newly found nft reg tickets

	senseRegTickets, err := s.pastelClient.ActionTickets(ctx)
	if err != nil {
		return errors.Errorf("get registered ticket: %w", err)
	}

	if len(senseRegTickets) == 0 {
		return nil
	}
	for i := 0; i < len(senseRegTickets); i++ {
		if senseRegTickets[i].ActionTicketData.CalledAt <= s.latestBlockHeight {
			continue
		}

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
			log.WithContext(ctx).WithField("senseRegTickets.ActionTicketData", senseRegTickets[i].ActionTicketData).Warnf("Could not get sense ticket for action ticket data")
			continue
		}

		ddAndFpFromTicket := &pastel.DDAndFingerprints{}
		//Get the dd and fp file from the ticket
		for _, id := range senseTicket.DDAndFingerprintsIDs {
			ddAndFpFromTicket, err = s.tryToGetFingerprintFileFromHash(ctx, id)
			if err != nil {
				//probably too verbose even for debug.
				//log.WithContext(ctx).WithField("error", err).Debug("Could not get the fingerprint for this file hash")
				continue
			}
			break
		}
		if ddAndFpFromTicket == nil {
			log.WithContext(ctx).WithField("senseTicket", senseTicket).Warnf("None of the dd and fp id files for this sense reg ticket could be properly unmarshalled")
			continue
		}
		if ddAndFpFromTicket.HashOfCandidateImageFile == "" {
			log.WithContext(ctx).WithField("senseTicket", nftRegTickets[i].RegTicketData.NFTTicketData).Warnf("This NFT sense ticket's DDAndFp struct has no HashOfCandidateImageFile, perhaps it's an older version.")
			continue
		}

		existsInDatabase, err := s.checkIfFingerprintExistsInDatabase(ctx, ddAndFpFromTicket.HashOfCandidateImageFile)
		if existsInDatabase {
			log.WithContext(ctx).WithField("hashOfCandidateImageFile", ddAndFpFromTicket.HashOfCandidateImageFile).Debug("Found hash of candidate image file already exists in database, not adding.")
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

		// imgFileHash := senseTicket.DataHash
		//Below code might be kept in case we want to add the ability to get thumbnails to sense
		// imgFile, err := s.p2pClient.Retrieve(ctx, string(imgFileHash))
		// if err != nil {
		// 	log.WithContext(ctx).WithField("imgFileHash", imgFileHash).Warnf("Could not get the image with this hash for senseTicket")
		// 	continue
		// }

		// imgBuffer := bytes.NewBuffer(imgFile)
		// src, err := imaging.Decode(imgBuffer)
		// if err != nil {
		// 	log.WithContext(ctx).WithField("imgFileHash", imgFileHash).Warnf("Could not open bytes as image")
		// 	continue
		// }

		// dstImage := imaging.Thumbnail(src, 200, 200, imaging.Lanczos)

		// var thumbByteArr bytes.Buffer

		// if err := webp.Encode(&thumbByteArr, dstImage, &webp.Options{Quality: 25}); err != nil {
		// 	log.WithContext(ctx).WithField("imgFileHash", imgFileHash).Warnf("Failed to encode thumbnail image")
		// 	continue
		// }

		// readableThumbnailStr := utils.B64Encode(thumbByteArr.Bytes())

		//make sure ImageFingerprintOfCnadidateImageFile exists.
		// this could fail if the ticket is an older version of the DDAndFingerprints struct, so we will continue to next fingerprint
		if ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile == nil {
			log.WithContext(ctx).WithField("senseTicket", senseRegTickets[i].ActionTicketData).Warnf("This sense ticket's DDAndFp struct has no ImageFingerprintOfCandidateImageFile, perhaps it's an older version.")
			continue
		}
		if len(ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile) < 1 {
			log.WithContext(ctx).WithField("senseTicket", senseRegTickets[i].ActionTicketData).Warnf("This sense reg ticket's DDAndFp struct's ImageFingerprintOfCandidateImageFile is zero length, perhaps it's an older version.")
			continue
		}

		if err := s.storeFingerprint(ctx, &dupeDetectionFingerprints{
			Sha256HashOfArtImageFile:           ddAndFpFromTicket.HashOfCandidateImageFile,
			ImageFingerprintVector:             toFloat64Array(ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile),
			DatetimeFingerprintAddedToDatabase: time.Now().Format("2006-01-02 15:04:05"),
			PathToArtImageFile:                 "",
			ImageThumbnailAsBase64:             "",
			RequestType:                        senseRegTickets[i].ActionTicketData.Type,
			IDString:                           "",
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

func (s *service) getRecordCount(_ context.Context) (int64, error) {
	var Data struct {
		Count int64 `db:"count"`
	}

	statement := getNumberOfFingerprintsStatement
	err := s.db.Get(&Data, statement)
	if err != nil {
		return 0, errors.Errorf("query: %w", err)
	}

	return Data.Count, nil
}

// Stats return status of dupe detection
func (s *service) Stats(ctx context.Context) (map[string]interface{}, error) {
	stats := map[string]interface{}{}

	// Get last inserted item
	lastItem, err := s.getLatestFingerprint(ctx)
	if err != nil {
		return nil, errors.Errorf("getLatestFingerprint: %w", err)
	}
	stats["last_insert_time"] = lastItem.DatetimeFingerprintAddedToDatabase

	// Get total of records
	recordCount, err := s.getRecordCount(ctx)
	if err != nil {
		return nil, errors.Errorf("getRecordCount: %w", err)
	}
	stats["record_count"] = recordCount

	return stats, nil
}

// NewService returns a new ddscan service
func NewService(config *Config, pastelClient pastel.Client, sn node.SNClientInterface) (svc.SvcInterface, error) {
	file := config.DataFile
	if os.Getenv("INTEGRATION_TEST_ENV") == "true" {
		tmpfile, err := ioutil.TempFile("", "registered_image_fingerprints_db.sqlite")
		if err != nil {
			panic(err.Error())
		}
		file = tmpfile.Name()
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, errors.Errorf("database dd service not found: %w", err)
	}

	db, err := sqlx.Connect("sqlite3", file)
	if err != nil {
		return nil, fmt.Errorf("cannot open dd-service database: %w", err)
	}

	return &service{
		config:             config,
		pastelClient:       pastelClient,
		sn:                 sn,
		db:                 db,
		currentNFTBlock:    1,
		currentActionBlock: 1,
	}, nil
}
