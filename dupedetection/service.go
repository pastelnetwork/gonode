package dupedetection

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/metadb/rqlite/db"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/sbinet/npyio"
)

const (
	synchronizationIntervalSec    = 5
	synchronizationTimeoutSec     = 60
	runTaskIntervalMin            = 2
	fingerprintSizeModel1         = 4032
	fingerprintSizeModel2         = 2560
	fingerprintSizeModel3         = 1920
	fingerprintSizeModel4         = 1536
	fingerprintSizeModel          = 10048
	masterNodeSuccessfulStatus    = "Masternode successfully started"
	getLatestFingerprintStatement = `SELECT * FROM image_hash_to_image_fingerprint_table ORDER BY datetime_fingerprint_added_to_database DESC LIMIT 1`
	insertFingerprintStatement    = `INSERT INTO image_hash_to_image_fingerprint_table(sha256_hash_of_art_image_file, model_1_image_fingerprint_vector, model_2_image_fingerprint_vector, model_3_image_fingerprint_vector, model_4_image_fingerprint_vector) VALUES("%s", '%s', '%s', '%s', '%s')`
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
	config       *Config
	pastelClient pastel.Client
	p2pClient    p2p.Client
	db           *db.DB
}

// Run starts task
func (s *service) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, "dd-service")
	defer s.db.Close()

	if err := s.waitSynchronization(ctx); err != nil {
		return errors.Errorf("failed to wait synchronization, err: %w", err)
	}

	if err := s.runTask(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("First runTask() failed")
	}

	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(runTaskIntervalMin * time.Minute):
			if err := s.runTask(ctx); err != nil {
				log.WithContext(ctx).WithError(err).Errorf("runTask() failed")
			}
		}
	}
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
			st, err := s.pastelClient.MasterNodeStatus(ctx)
			if err != nil {
				log.WithContext(ctx).WithError(err).Warn("Failed to get status from master node")
			} else if st.Status == masterNodeSuccessfulStatus {
				log.WithContext(ctx).Debug("Done for waiting synchronization status")
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
		return nil, errors.Errorf("failed to query, err: %w", err)
	}

	if len(row) != 0 && len(row[0].Values) == 0 {
		return nil, errors.New("dd database is empty")
	}

	values := row[0].Values[0]
	resultStr := make(map[string]interface{})
	for i := 0; i < len(row[0].Columns); i++ {
		switch row[0].Types[i] {
		case "date", "datetime":
			// TODO
		case "array":
			str, ok := values[i].(string)
			if !ok {
				log.WithContext(ctx).Errorf("failed to get str from npy, columns: %s", row[0].Columns[i])
				continue
			}

			b, err := hex.DecodeString(str)
			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("Failed to hex decode npy str, columns: %s", row[0].Columns[i])
				continue
			}

			f := bytes.NewBuffer(b)

			var fp []float64
			if err := npyio.Read(f, &fp); err != nil {
				log.WithContext(ctx).WithError(err).Error("Failed to convert npy to float64")
				continue
			}
			resultStr[row[0].Columns[i]] = fp
		default:
			resultStr[row[0].Columns[i]] = values[i]
		}
	}

	b, err := json.Marshal(resultStr)
	if err != nil {
		return nil, errors.Errorf("failed to JSON marshal value, err: %w", err)
	}

	ddFingerprint := &dupeDetectionFingerprints{}
	if err := json.Unmarshal(b, ddFingerprint); err != nil {
		return nil, errors.Errorf("failed to JSON unmarshal, err: %w", err)
	}

	return ddFingerprint, nil
}

func (s *service) storeFingerprint(ctx context.Context, input *dupeDetectionFingerprints) error {
	encodeFloat2Npy := func(v []float64) (string, error) {
		f := bytes.NewBuffer(nil)
		if err := npyio.Write(f, v); err != nil {
			return "", errors.Errorf("failed to encode to npy, err: %w", err)
		}
		return hex.EncodeToString(f.Bytes()), nil
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

	statement := fmt.Sprintf(insertFingerprintStatement, input.Sha256HashOfArtImageFile, fp1, fp2, fp3, fp4)
	_, err = s.db.ExecuteStringStmt(statement)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to insert fingerprint record")
	}
	return nil
}

func (s *service) runTask(ctx context.Context) error {
	nftRegTickets, err := s.pastelClient.RegTickets(ctx)
	if err != nil {
		return errors.Errorf("failed to get registered ticket, err: %w", err)
	}

	if len(nftRegTickets) == 0 {
		return nil
	}

	lastestFp, err := s.getLatestFingerprint(ctx)
	if err != nil {
		return errors.Errorf("failed to get lastest fingerprint, err: %w", err)
	}

	for i := 0; i < len(nftRegTickets); i++ {
		if nftRegTickets[i].Height <= lastestFp.NumberOfBlock {
			continue
		}

		nftTicketData, err := pastel.DecodeNFTTicket(nftRegTickets[i].RegTicketData.NFTTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}

		fingerprintsHash := nftTicketData.AppTicketData.FingerprintsHash
		fingerprintByte, err := s.p2pClient.Retrieve(ctx, base58.Encode(fingerprintsHash))
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to retrieve fingerprint")
			continue
		}

		fingerprint, err := pastel.FingerprintFromBytes(fingerprintByte)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to convert fingerprint")
			continue
		}

		size := len(fingerprint)
		if size != fingerprintSizeModel {
			log.WithContext(ctx).Errorf("invaild size fingerprint, size: %d", size)
			continue
		}

		model1ImageFingerprintVector := make([]float64, fingerprintSizeModel1)
		model2ImageFingerprintVector := make([]float64, fingerprintSizeModel2)
		model3ImageFingerprintVector := make([]float64, fingerprintSizeModel3)
		model4ImageFingerprintVector := make([]float64, fingerprintSizeModel4)

		start, end := 0, fingerprintSizeModel1
		copy(model1ImageFingerprintVector, fingerprint[start:end])

		start, end = start+fingerprintSizeModel1, end+fingerprintSizeModel2
		copy(model2ImageFingerprintVector, fingerprint[start:end])

		start, end = start+fingerprintSizeModel2, end+fingerprintSizeModel3
		copy(model3ImageFingerprintVector, fingerprint[start:end])

		start, end = start+fingerprintSizeModel3, end+fingerprintSizeModel4
		copy(model4ImageFingerprintVector, fingerprint[start:end])

		if err := s.storeFingerprint(ctx, &dupeDetectionFingerprints{
			Sha256HashOfArtImageFile:           hex.EncodeToString(fingerprintsHash),
			Model1ImageFingerprintVector:       model1ImageFingerprintVector,
			Model2ImageFingerprintVector:       model2ImageFingerprintVector,
			Model3ImageFingerprintVector:       model3ImageFingerprintVector,
			Model4ImageFingerprintVector:       model4ImageFingerprintVector,
			NumberOfBlock:                      nftRegTickets[i].Height,
			DatetimeFingerprintAddedToDatabase: time.Now(),
		}); err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to store fingerprint")
		}
	}

	return nil
}

// NewService return a new Service instance
func NewService(config *Config, pastelClient pastel.Client, p2pClient p2p.Client) (Service, error) {
	if _, err := os.Stat(config.DataFile); os.IsNotExist(err) {
		return nil, errors.Errorf("database dd service not found, err: %w", err)
	}

	db, err := db.Open(config.DataFile)
	if err != nil {
		return nil, errors.Errorf("failed to open dd-service database, err: %w", err)
	}

	return &service{
		config:       config,
		pastelClient: pastelClient,
		p2pClient:    p2pClient,
		db:           db,
	}, nil
}
