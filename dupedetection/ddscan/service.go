package ddscan

import (
	"bytes"
	"context"
	"encoding/json"
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

// Service interface
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
	fingerprintSizeModel1            = 4032
	fingerprintSizeModel2            = 2560
	fingerprintSizeModel3            = 1920
	fingerprintSizeModel4            = 1536
	fingerprintSizeModel             = 10048
	masterNodeSuccessfulStatus       = "Masternode successfully started"
	getLatestFingerprintStatement    = `SELECT * FROM image_hash_to_image_fingerprint_table ORDER BY datetime_fingerprint_added_to_database DESC LIMIT 1`
	insertFingerprintStatement       = `INSERT INTO image_hash_to_image_fingerprint_table(sha256_hash_of_art_image_file, model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector) VALUES(?,?,?,?,?)`
	getNumberOfFingerprintsStatement = `SELECT COUNT(*) FROM image_hash_to_image_fingerprint_table`
)

type dupeDetectionFingerprints struct {
	Sha256HashOfArtImageFile           string    `json:"sha256_hash_of_art_image_file,omitempty"`
	PathToArtImageFile                 string    `json:"path_to_art_image_file,omitempty"`
	NumberOfBlock                      int       `json:"number_of_block,omitempty"`
	DatetimeFingerprintAddedToDatabase time.Time `json:"datetime_fingerprint_added_to_database,omitempty"`
	Model1ImageFingerprintVector       []float64 `json:"model_1_image_fingerprint_vector,omitempty"`
	Model2ImageFingerprintVector       []float64 `json:"model_2_image_fingerprint_vector,omitempty"`
	Model3ImageFingerprintVector       []float64 `json:"model_3_image_fingerprint_vector,omitempty"`
	Model4ImageFingerprintVector       []float64 `json:"model_4_image_fingerprint_vector,omitempty"`
}

type service struct {
	config             *Config
	pastelClient       pastel.Client
	p2pClient          p2p.Client
	db                 *db.DB
	isMasterNodeSynced bool
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
		if err := s.run(ctx); err != nil {
			if utils.IsContextErr(err) {
				return err
			}
			log.DD().WithContext(ctx).WithError(err).Error("Failed to start dd-scan service, retrying.")
		} else {
			return nil
		}
	}
}

func (s *service) run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, "dd-scan")
	defer s.db.Close()

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
				} else {
					log.DD().WithContext(ctx).Debug("Done for waiting synchronization status")
					s.isMasterNodeSynced = true
				}
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
			// FIXME: these follow columns are nil, skip them
			if row[0].Columns[i] == "model_5_image_fingerprint_vector" ||
				row[0].Columns[i] == "model_6_image_fingerprint_vector" ||
				row[0].Columns[i] == "model_7_image_fingerprint_vector" {
				continue
			}

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

	fp1, err := encodeFloat2Npy(input.Model1ImageFingerprintVector)
	if err != nil {
		return err
	}

	fp2, err := encodeFloat2Npy(input.Model2ImageFingerprintVector)
	if err != nil {
		return err
	}

	fp3, err := encodeFloat2Npy(input.Model3ImageFingerprintVector)
	if err != nil {
		return err
	}

	fp4, err := encodeFloat2Npy(input.Model4ImageFingerprintVector)
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
						Value: &command.Parameter_Y{
							Y: fp1,
						},
					},
					{
						Value: &command.Parameter_Y{
							Y: fp2,
						},
					},
					{
						Value: &command.Parameter_Y{
							Y: fp3,
						},
					},
					{
						Value: &command.Parameter_Y{
							Y: fp4,
						},
					},
				},
			},
		},
	}

	_, err = s.db.Execute(req, false)
	if err != nil {
		log.DD().WithContext(ctx).WithError(err).Error("Failed to insert fingerprint record")
		return err
	}
	return nil
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

	lastestFp, err := s.getLatestFingerprint(ctx)
	if err != nil {
		return errors.Errorf("get lastest fingerprint: %w", err)
	}

	for i := 0; i < len(nftRegTickets); i++ {
		if nftRegTickets[i].Height <= lastestFp.NumberOfBlock {
			continue
		}

		_, err := pastel.DecodeNFTTicket(nftRegTickets[i].RegTicketData.NFTTicket)
		if err != nil {
			log.DD().WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}

		// fingerprintsHash := string(nftTicketData.AppTicketData.FingerprintsHash)
		// https://pastel-network.atlassian.net/browse/PSL-155 - update it
		fingerprintsHash := "TBD"
		compressedFingerprintBytes, err := s.p2pClient.Retrieve(ctx, fingerprintsHash)
		if err != nil {
			log.DD().WithContext(ctx).WithField("FingerprintsHash", fingerprintsHash).WithError(err).Error("Failed to retrieve fingerprint")
			continue
		}

		fingerprintFromBytes, err := zstd.Decompress(nil, compressedFingerprintBytes)
		if err != nil {
			log.DD().WithContext(ctx).WithField("FingerprintsHash", fingerprintsHash).WithError(err).Error("Failed to decompress fingerprint")
			continue
		}

		fingerprint, err := pastel.FingerprintFromBytes(fingerprintFromBytes)

		if err != nil {
			log.DD().WithContext(ctx).WithField("FingerprintsHash", fingerprintsHash).WithError(err).Error("Failed to convert fingerprint")
			continue
		}

		size := len(fingerprint)
		if size != fingerprintSizeModel {
			log.DD().WithContext(ctx).Errorf("invaild size fingerprint, size: %d", size)
			continue
		}

		fingerprint64 := toFloat64Array(fingerprint)

		model1ImageFingerprintVector := make([]float64, fingerprintSizeModel1)
		model2ImageFingerprintVector := make([]float64, fingerprintSizeModel2)
		model3ImageFingerprintVector := make([]float64, fingerprintSizeModel3)
		model4ImageFingerprintVector := make([]float64, fingerprintSizeModel4)

		start, end := 0, fingerprintSizeModel1
		copy(model1ImageFingerprintVector, fingerprint64[start:end])

		start, end = start+fingerprintSizeModel1, end+fingerprintSizeModel2
		copy(model2ImageFingerprintVector, fingerprint64[start:end])

		start, end = start+fingerprintSizeModel2, end+fingerprintSizeModel3
		copy(model3ImageFingerprintVector, fingerprint64[start:end])

		start, end = start+fingerprintSizeModel3, end+fingerprintSizeModel4
		copy(model4ImageFingerprintVector, fingerprint64[start:end])

		if err := s.storeFingerprint(ctx, &dupeDetectionFingerprints{
			Sha256HashOfArtImageFile:           fingerprintsHash,
			Model1ImageFingerprintVector:       model1ImageFingerprintVector,
			Model2ImageFingerprintVector:       model2ImageFingerprintVector,
			Model3ImageFingerprintVector:       model3ImageFingerprintVector,
			Model4ImageFingerprintVector:       model4ImageFingerprintVector,
			NumberOfBlock:                      nftRegTickets[i].Height,
			DatetimeFingerprintAddedToDatabase: time.Now(),
		}); err != nil {
			log.DD().WithContext(ctx).WithError(err).Error("Failed to store fingerprint")
		}
	}

	return nil
}

func (s *service) getRecordCount(_ context.Context) (int64, error) {
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

// Stats return status of dupde detection
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

// NewService return a new Service instance
func NewService(config *Config, pastelClient pastel.Client, p2pClient p2p.Client) (Service, error) {
	if _, err := os.Stat(config.DataFile); os.IsNotExist(err) {
		return nil, errors.Errorf("database dd service not found: %w", err)
	}

	db, err := db.Open(config.DataFile, true)
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
