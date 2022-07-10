package hermes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DataDog/zstd"
	"github.com/jmoiron/sqlx"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/metadb/rqlite/command"
	"github.com/pastelnetwork/gonode/p2p"
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
	getNumberOfFingerprintsStatement = `SELECT COUNT(*) FROM image_hash_to_image_fingerprint_table`
)

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
	p2pClient          p2p.Client
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

	group.Go(func() error {
		return s.CleanupInactiveTickets(gctx)
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

			group.Go(func() error {
				return s.CleanupInactiveTickets(gctx)
			})

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

func (s *service) getLatestFingerprint(ctx context.Context) (*dupeDetectionFingerprints, error) {
	row, err := s.queryFromStatement(ctx, getLatestFingerprintStatement)
	if err != nil {
		return nil, errors.Errorf("query: %w", err)
	}

	if row == nil {
		return nil, errors.New("dd database is empty")
	}

	values := row.Values[0]
	resultStr := make(map[string]interface{})
	for i := 0; i < len(row.Columns); i++ {
		switch w := row.Values[0].GetParameters()[i].GetValue().(type) {
		case *command.Parameter_Y:
			val := values.GetParameters()[i].GetY()
			if val == nil {
				log.WithContext(ctx).Errorf("nil value at column: %s", row.Columns[i])
				continue
			}

			// even the shape if these should be Nx1, but for reading, we convert it into 1xN array
			f := bytes.NewBuffer(val)

			var fp []float64
			if err := npyio.Read(f, &fp); err != nil {
				log.WithContext(ctx).WithError(err).Error("Failed to convert npy to float64")
				continue
			}

			resultStr[row.Columns[i]] = fp
		case *command.Parameter_I:
			resultStr[row.Columns[i]] = w.I
		case *command.Parameter_D:
			resultStr[row.Columns[i]] = w.D
		case *command.Parameter_B:
			resultStr[row.Columns[i]] = w.B
		case *command.Parameter_S:
			resultStr[row.Columns[i]] = w.S
		case nil:
			resultStr[row.Columns[i]] = nil
		default:
			return nil, errors.Errorf("getLatestFingerprint unsupported type: %w", w)
		}
	}

	b, err := json.Marshal(resultStr)
	if err != nil {
		return nil, errors.Errorf("marshal output: %w", err)
	}

	ddFingerprint := &dupeDetectionFingerprints{}
	if err := json.Unmarshal(b, ddFingerprint); err != nil {
		return nil, errors.Errorf("unmarshal json: %w", err)
	}

	return ddFingerprint, nil
}

func (s *service) checkIfFingerprintExistsInDatabase(_ context.Context, hash string) (bool, error) {
	req := &command.Request{
		Statements: []*command.Statement{
			{
				Value: &command.Parameter_S{
					S: hash,
				},
			},
		},
	}

	conn, err := s.db.Conn(context.Background())
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to connect to database")
		return false, nil
	}
	defer conn.Close()

	parameters, err := parametersToValues(stmt.Parameters)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to insert fingerprint record")
	}

	queryResult, err := conn.QueryContext(context.Background(), stmt.Sql, parameters...)
	if err != nil {
		return false, errors.Errorf("querying fingerprint database for hash: %w", err)
	}
	if queryResult == nil {
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

	stmt := command.Statement{
		Sql: insertFingerprintStatement,
		Parameters: []*command.Parameter{
			{
				Value: &command.Parameter_S{
					S: input.Sha256HashOfArtImageFile,
				},
			},
			{
				Value: &command.Parameter_S{
					S: input.PathToArtImageFile,
				},
			},
			{
				Value: &command.Parameter_Y{
					Y: fp,
				},
			},
			{
				Value: &command.Parameter_S{
					S: input.DatetimeFingerprintAddedToDatabase,
				},
			},
			{
				Value: &command.Parameter_S{
					S: input.ImageThumbnailAsBase64,
				},
			},
			{
				Value: &command.Parameter_S{
					S: input.RequestType,
				},
			},
			{
				Value: &command.Parameter_S{
					S: input.IDString,
				},
			},
		},
	}

	conn, err := s.db.Conn(context.Background())
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to connect to database")
		return nil
	}
	defer conn.Close()

	parameters, err := parametersToValues(stmt.Parameters)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to insert fingerprint record")
	}

	executeResult, err := conn.ExecContext(context.Background(), stmt.Sql, parameters...)
	log.WithContext(ctx).WithField("execute result", executeResult).Debug("Execute result of adding new fp")
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
	/* // For debugging
	if cnt, err := s.getRecordCount(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to get count")
	} else {
		log.WithContext(ctx).WithField("cnt", cnt).Info("GetCount")
	}
	*/

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
	statement := getNumberOfFingerprintsStatement
	row, err := s.queryFromStatement(ctx, statement)
	if err != nil {
		return 0, errors.Errorf("query: %w", err)
	}

	if len(row.GetValues()) != 1 {
		return 0, errors.Errorf("invalid Values length: %d", len(row.Columns))
	}

	if len(row.GetValues()[0].GetParameters()) != 1 {
		return 0, errors.Errorf("invalid Values[0] length: %d", len(row.Columns))
	}

	switch w := row.GetValues()[0].GetParameters()[0].GetValue().(type) {
	case *command.Parameter_I:
		return row.GetValues()[0].GetParameters()[0].GetI(), nil
	default:
		return 0, errors.Errorf("invalid returned type : %v", w)
	}
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

func (s *service) queryFromStatement(ctx context.Context, queryString string) (*command.QueryRows, error) {
	conn, err := s.db.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows := &command.QueryRows{}

	rs, err := conn.QueryContext(context.Background(), queryString)
	if err != nil {
		log.WithError(err).Error("dd-scan failed to query database")
	}
	defer rs.Close()

	columns, err := rs.Columns()
	if err != nil {
		return nil, err
	}

	types, err := rs.ColumnTypes()
	if err != nil {
		return nil, err
	}
	xTypes := make([]string, len(types))
	for i := range types {
		xTypes[i] = strings.ToLower(types[i].DatabaseTypeName())
	}

	for rs.Next() {
		dest := make([]interface{}, len(columns))
		ptrs := make([]interface{}, len(dest))
		for i := range ptrs {
			ptrs[i] = &dest[i]
		}
		if err := rs.Scan(ptrs...); err != nil {
			return nil, err
		}
		rows.Values = append(rows.Values, &command.Values{
			Parameters: normalizeRowValues(dest, xTypes),
		})
	}

	// Check for errors from iterating over rows.
	if err := rs.Err(); err != nil {
		log.WithError(err).Error("dd-scan error after iterating over query")
		return nil, err
	}

	rows.Columns = columns
	rows.Types = xTypes

	return rows, nil
}

// parametersToValues maps values in the proto params to SQL driver values.
func parametersToValues(parameters []*command.Parameter) ([]interface{}, error) {
	if parameters == nil {
		return nil, nil
	}

	values := make([]interface{}, len(parameters))
	for i := range parameters {
		switch w := parameters[i].GetValue().(type) {
		case *command.Parameter_I:
			values[i] = w.I
		case *command.Parameter_D:
			values[i] = w.D
		case *command.Parameter_B:
			values[i] = w.B
		case *command.Parameter_Y:
			values[i] = w.Y
		case *command.Parameter_S:
			values[i] = w.S
		default:
			return nil, fmt.Errorf("unsupported type: %T", w)
		}
	}
	return values, nil
}

// normalizeRowValues performs some normalization of values in the returned rows.
// Text values come over (from sqlite-go) as []byte instead of strings
// for some reason, so we have explicitly convert (but only when type
// is "text" so we don't affect BLOB types)
func normalizeRowValues(row []interface{}, types []string) []*command.Parameter {
	values := make([]*command.Parameter, len(types))
	for i, v := range row {
		switch val := v.(type) {
		case int:
		case int64:
			values[i] = &command.Parameter{
				Value: &command.Parameter_I{
					I: val,
				},
			}
		case float64:
			values[i] = &command.Parameter{
				Value: &command.Parameter_D{
					D: val,
				},
			}
		case bool:
			values[i] = &command.Parameter{
				Value: &command.Parameter_B{
					B: val,
				},
			}
		case string:
			values[i] = &command.Parameter{
				Value: &command.Parameter_S{
					S: val,
				},
			}
		case []byte:
			if isTextType(types[i]) {
				values[i].Value = &command.Parameter_S{
					S: string(val),
				}
			} else {
				values[i] = &command.Parameter{
					Value: &command.Parameter_Y{
						Y: val,
					},
				}
			}
		}
	}
	return values
}

// isTextType returns whether the given type has a SQLite text affinity.
// http://www.sqlite.org/datatype3.html
func isTextType(t string) bool {
	return t == "text" ||
		t == "json" ||
		t == "" ||
		strings.HasPrefix(t, "varchar") ||
		strings.HasPrefix(t, "varying character") ||
		strings.HasPrefix(t, "nchar") ||
		strings.HasPrefix(t, "native character") ||
		strings.HasPrefix(t, "nvarchar") ||
		strings.HasPrefix(t, "clob")
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

	db, err := sqlx.Open("sqlite3", fmt.Sprintf("file:%s?_fk=%s", file, strconv.FormatBool(true)))
	if err != nil {
		return nil, errors.Errorf("open dd-service database: %w", err)
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
