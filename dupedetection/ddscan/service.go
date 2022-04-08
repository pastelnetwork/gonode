package ddscan

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/pastelnetwork/gonode/common/utils"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/metadb/rqlite/command"
	"github.com/pastelnetwork/gonode/metadb/rqlite/db"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/sbinet/npyio"
	"gonum.org/v1/gonum/mat"
)

// Service ...
type Service interface {
	// Run starts task
	Run(ctx context.Context) error

	// Stats returns current status of service
	Stats(ctx context.Context) (map[string]interface{}, error)
}

const (
	synchronizationIntervalSec       = 5
	synchronizationTimeoutSec        = 60
	runTaskInterval                  = 2 * time.Minute
	masterNodeSuccessfulStatus       = "Masternode successfully started"
	getLatestFingerprintStatement    = `SELECT * FROM image_hash_to_image_fingerprint_table ORDER BY datetime_fingerprint_added_to_database DESC LIMIT 1`
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
	db                 *db.DB
	isMasterNodeSynced bool
	latestBlockHeight  int
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
					return err
				}

				//log.DD().WithContext(ctx).WithError(err).Error("failed to run ddscan, retrying.")
			} else {
				return nil
			}
		}
	}
}

func (s *service) run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, "dd-scan")
	if _, err := os.Stat(s.config.DataFile); os.IsNotExist(err) {
		return errors.Errorf("dataFile dd service not found: %w", err)
	}

	if err := s.waitSynchronization(ctx); err != nil {
		log.DD().WithContext(ctx).WithError(err).Error("Failed to initial wait synchronization")
	} else {
		s.isMasterNodeSynced = true
	}

	if err := s.runTask(ctx); err != nil {
		log.DD().WithContext(ctx).WithError(err).Errorf("First runTask() failed")
	}

	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(runTaskInterval):
			// Check if node is synchronized or not
			if !s.isMasterNodeSynced {
				if err := s.checkSynchronized(ctx); err != nil {
					log.DD().WithContext(ctx).WithError(err).Warn("Failed to check synced status from master node")
					continue
				}

				log.DD().WithContext(ctx).Debug("Done for waiting synchronization status")
				s.isMasterNodeSynced = true
			}

			if err := s.runTask(ctx); err != nil {
				log.DD().WithContext(ctx).WithError(err).Errorf("runTask() failed")
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
				log.DD().WithContext(ctx).WithError(err).Warn("Failed to check synced status from master node")
			} else {
				log.DD().WithContext(ctx).Info("Done for waiting synchronization status")
				return nil
			}
		case <-timeoutCh:
			return errors.New("timeout expired")
		}
	}
}

func (s *service) getLatestFingerprint(ctx context.Context) (*dupeDetectionFingerprints, error) {
	row, err := s.db.QueryStringStmt(getLatestFingerprintStatement)
	if err != nil {
		return nil, errors.Errorf("query: %w", err)
	}

	if len(row) != 0 && len(row[0].Values) == 0 {
		return nil, errors.New("dd database is empty")
	}

	values := row[0].Values[0]
	resultStr := make(map[string]interface{})
	for i := 0; i < len(row[0].Columns); i++ {
		switch w := row[0].Values[0].GetParameters()[i].GetValue().(type) {
		case *command.Parameter_Y:
			val := values.GetParameters()[i].GetY()
			if val == nil {
				log.DD().WithContext(ctx).Errorf("nil value at column: %s", row[0].Columns[i])
				continue
			}

			// even the shape if these should be Nx1, but for reading, we convert it into 1xN array
			f := bytes.NewBuffer(val)

			var fp []float64
			if err := npyio.Read(f, &fp); err != nil {
				log.DD().WithContext(ctx).WithError(err).Error("Failed to convert npy to float64")
				continue
			}

			resultStr[row[0].Columns[i]] = fp
		case *command.Parameter_I:
			resultStr[row[0].Columns[i]] = w.I
		case *command.Parameter_D:
			resultStr[row[0].Columns[i]] = w.D
		case *command.Parameter_B:
			resultStr[row[0].Columns[i]] = w.B
		case *command.Parameter_S:
			resultStr[row[0].Columns[i]] = w.S
		case nil:
			resultStr[row[0].Columns[i]] = nil
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

// type dupeDetectionFingerprints struct {
// 	Sha256HashOfArtImageFile           string    `json:"sha256_hash_of_art_image_file,omitempty"`
// 	PathToArtImageFile                 string    `json:"path_to_art_image_file,omitempty"`
// 	ImageFingerprintVector             []float64 `json:"image_fingerprint_vector,omitempty"`
// 	DatetimeFingerprintAddedToDatabase string `json:"datetime_fingerprint_added_to_database,omitempty"`
// 	ImageThumbnailAsBase64             string    `json:"image_thumbnail_as_base64_string,omitempty"`
// 	RequestType                        string    `json:"request_type,omitempty"`
// 	IDString                           string    `json:"id_String,omitempty"`
// }

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

	req := &command.Request{
		Statements: []*command.Statement{
			{
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
			},
		},
	}

	executeResult, err := s.db.Execute(req, false)
	log.WithContext(ctx).WithField("execute result", executeResult).Debug("Execute result of adding new fp")
	if err != nil || executeResult[0].Error != "" {
		log.DD().WithContext(ctx).WithError(err).Error("Failed to insert fingerprint record")
		return err
	}
	return nil
}

//Utility function to get dd and fp file from an id hash, where the file should be stored
func (s *service) tryToGetFingerprintFileFromHash(ctx context.Context, hash string) (*pastel.DDAndFingerprints, error) {
	rawFile, err := s.p2pClient.Retrieve(ctx, hash)
	if err != nil {
		return nil, errors.Errorf("Error finding dd and fp file: %w", err)
	}

	dec, err := utils.B64Decode(rawFile)
	if err != nil {
		return nil, errors.Errorf("decode data: %w", err)
	}

	decData, err := zstd.Decompress(nil, dec)
	if err != nil {
		return nil, errors.Errorf("decompress: %w", err)
	}

	splits := bytes.Split(decData, []byte{pastel.SeparatorByte})
	file, err := utils.B64Decode(splits[0])
	if err != nil {
		return nil, errors.Errorf("decode file: %w", err)
	}

	ddFingerprint := &pastel.DDAndFingerprints{}
	if err := json.Unmarshal(file, ddFingerprint); err != nil {
		return nil, errors.Errorf("unmarshal json: %w", err)
	}
	return ddFingerprint, nil
}

func (s *service) runTask(ctx context.Context) error {
	/* // For debugging
	if cnt, err := s.getRecordCount(ctx); err != nil {
		log.DD().WithContext(ctx).WithError(err).Error("Failed to get count")
	} else {
		log.DD().WithContext(ctx).WithField("cnt", cnt).Info("GetCount")
	}
	*/

	nftRegTickets, err := s.pastelClient.RegTickets(ctx)
	if err != nil {
		return errors.Errorf("get registered ticket: %w", err)
	}

	if len(nftRegTickets) == 0 {
		return nil
	}

	//track latest block height, but don't set it until we check all the nft reg tickets and the sense tickets.
	latestBlockHeight := s.latestBlockHeight

	//loop through nft tickets and store newly found nft reg tickets
	for i := 0; i < len(nftRegTickets); i++ {
		if nftRegTickets[i].Height <= s.latestBlockHeight {
			continue
		}

		_, err := pastel.DecodeNFTTicket(nftRegTickets[i].RegTicketData.NFTTicket)
		if err != nil {
			log.DD().WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}

		ddFPIDs := nftRegTickets[i].RegTicketData.NFTTicketData.AppTicketData.DDAndFingerprintsIDs

		ddAndFpFromTicket := &pastel.DDAndFingerprints{}
		//Get the dd and fp file from the ticket
		for _, id := range ddFPIDs {
			ddAndFpFromTicket, err = s.tryToGetFingerprintFileFromHash(ctx, id)
			if err != nil {
				break
			}
		}
		if ddAndFpFromTicket == nil {
			log.WithContext(ctx).WithField("NFTRegTicket", nftRegTickets[i].RegTicketData.NFTTicketData).Warnf("None of the dd and fp id files for this nft reg ticket could be properly unmarshalled")
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
			log.DD().WithContext(ctx).WithError(err).Error("Failed to store fingerprint")
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
		if senseRegTickets[i].ActionTicketData.BlockNum <= s.latestBlockHeight {
			continue
		}
		if senseRegTickets[i].Type != "sense-reg" {
			continue
		}

		senseTicket, err := senseRegTickets[i].ActionTicketData.APISenseTicket()
		if err != nil {
			log.WithContext(ctx).WithField("senseRegTickets.ActionTicketData", senseRegTickets[i].ActionTicketData).Warnf("Could not get sense ticket for action ticket data")
			continue
		}

		ddAndFpFromTicket := &pastel.DDAndFingerprints{}
		//Get the dd and fp file from the ticket
		for _, id := range senseTicket.DDAndFingerprintsIDs {
			ddAndFpFromTicket, err = s.tryToGetFingerprintFileFromHash(ctx, id)
			if err != nil {
				break
			}
		}
		if ddAndFpFromTicket == nil {
			log.WithContext(ctx).WithField("senseTicket", senseTicket).Warnf("None of the dd and fp id files for thissense reg ticket could be properly unmarshalled")
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

		if err := s.storeFingerprint(ctx, &dupeDetectionFingerprints{
			Sha256HashOfArtImageFile:           ddAndFpFromTicket.HashOfCandidateImageFile,
			ImageFingerprintVector:             toFloat64Array(ddAndFpFromTicket.ImageFingerprintOfCandidateImageFile),
			DatetimeFingerprintAddedToDatabase: time.Now().Format("2006-01-02 15:04:05"),
			PathToArtImageFile:                 "",
			ImageThumbnailAsBase64:             "",
			RequestType:                        senseRegTickets[i].Type,
			IDString:                           "",
		}); err != nil {
			log.DD().WithContext(ctx).WithError(err).Error("Failed to store fingerprint")
		}
		if senseRegTickets[i].ActionTicketData.BlockNum > latestBlockHeight {
			latestBlockHeight = senseRegTickets[i].ActionTicketData.BlockNum
		}
	}

	s.latestBlockHeight = latestBlockHeight

	return nil
}

func (s *service) getRecordCount(ctx context.Context) (int64, error) {
	statement := getNumberOfFingerprintsStatement
	rows, err := s.db.QueryStringStmt(statement)
	if err != nil {
		return 0, errors.Errorf("query: %w", err)
	}
	row := rows[0]

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

// NewService returns a new ddscan service
func NewService(config *Config, pastelClient pastel.Client, p2pClient p2p.Client) (Service, error) {
	file := config.DataFile
	if os.Getenv("INTEGRATION_TEST_ENV") == "true" {
		tmpfile, err := ioutil.TempFile("", "dupe_detection_image_fingerprint_database.sqlite")
		if err != nil {
			panic(err.Error())
		}
		file = tmpfile.Name()
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, errors.Errorf("database dd service not found: %w", err)
	}

	db, err := db.Open(file, true)
	if err != nil {
		return nil, errors.Errorf("open dd-service database: %w", err)
	}

	return &service{
		config:       config,
		pastelClient: pastelClient,
		p2pClient:    p2pClient,
		db:           db,
	}, nil
}
