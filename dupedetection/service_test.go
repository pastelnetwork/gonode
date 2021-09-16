package dupedetection

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	mrand "math/rand"
	"os"
	"testing"
	"time"

	fuzz "github.com/google/gofuzz"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/metadb/rqlite/db"
	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/assert"
)

const createTableStatement = `CREATE TABLE image_hash_to_image_fingerprint_table (sha256_hash_of_art_image_file text, path_to_art_image_file, model_1_image_fingerprint_vector array, model_2_image_fingerprint_vector array, model_3_image_fingerprint_vector array, model_4_image_fingerprint_vector array, datetime_fingerprint_added_to_database TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL, PRIMARY KEY (sha256_hash_of_art_image_file))`

func prepareService(t *testing.T) *service {
	tmpfile, err := ioutil.TempFile("", "dupe_detection_image_fingerprint_database.sqlite")
	if err != nil {
		panic(err.Error())
	}
	db, err := db.Open(tmpfile.Name())
	if err != nil {
		panic("failed to open database")
	}

	_, err = db.ExecuteStringStmt(createTableStatement)
	if err != nil {
		panic("failed to create table")
	}

	tmpfile.Close()

	s := &service{
		config: &Config{
			DataFile: tmpfile.Name(),
		},
	}
	s.db = db
	return s
}

func randFloats(n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = 0.0 + mrand.Float64()*(0.0-1.0)
	}
	return res
}

func generateFingerprint(t *testing.T) *dupeDetectionFingerprints {
	fp1 := randFloats(fingerprintSizeModel1)
	fp2 := randFloats(fingerprintSizeModel2)
	fp3 := randFloats(fingerprintSizeModel3)
	fp4 := randFloats(fingerprintSizeModel4)

	b := make([]byte, 10)
	rand.Read(b)
	return &dupeDetectionFingerprints{
		Sha256HashOfArtImageFile:           hex.EncodeToString(b),
		Model1ImageFingerprintVector:       fp1,
		Model2ImageFingerprintVector:       fp2,
		Model3ImageFingerprintVector:       fp3,
		Model4ImageFingerprintVector:       fp4,
		DatetimeFingerprintAddedToDatabase: time.Now(),
	}
}

func TestGetEmptyFingerprint(t *testing.T) {
	s := prepareService(t)
	defer s.db.Close()
	defer os.Remove(s.config.DataFile)

	emptyFp, err := s.getLatestFingerprint(context.Background())
	assert.True(t, emptyFp == nil)
	assert.Equal(t, err.Error(), "dd database is empty")
}

func TestSetFingerprint(t *testing.T) {
	s := prepareService(t)
	defer s.db.Close()
	defer os.Remove(s.config.DataFile)

	err := s.storeFingerprint(context.Background(), generateFingerprint(t))
	assert.True(t, err == nil)
}

func TestGetSetFingerprint(t *testing.T) {
	s := prepareService(t)
	defer s.db.Close()
	defer os.Remove(s.config.DataFile)

	setFp := generateFingerprint(t)
	ctx := context.Background()
	err := s.storeFingerprint(ctx, setFp)

	assert.True(t, err == nil)

	getFp, err := s.getLatestFingerprint(ctx)

	assert.True(t, err == nil)
	assert.True(t, getFp != nil)
	assert.Equal(t, setFp.Sha256HashOfArtImageFile, getFp.Sha256HashOfArtImageFile)
	assert.Equal(t, setFp.Model1ImageFingerprintVector, getFp.Model1ImageFingerprintVector)
	assert.Equal(t, setFp.Model2ImageFingerprintVector, getFp.Model2ImageFingerprintVector)
	assert.Equal(t, setFp.Model3ImageFingerprintVector, getFp.Model3ImageFingerprintVector)
	assert.Equal(t, setFp.Model4ImageFingerprintVector, getFp.Model4ImageFingerprintVector)
}

func TestWaitSynchronizationSuccessful(t *testing.T) {
	s := prepareService(t)
	defer s.db.Close()
	defer os.Remove(s.config.DataFile)

	status := &pastel.MasterNodeStatus{
		Status: "Masternode successfully started",
	}

	pMock := pastelMock.NewMockClient(t)
	pMock.ListenOnMasterNodeStatus(status, nil)
	s.pastelClient = pMock

	err := s.waitSynchronization(context.Background())
	assert.True(t, err == nil)
}

func TestWaitSynchronizationTimeout(t *testing.T) {
	s := prepareService(t)
	defer s.db.Close()
	defer os.Remove(s.config.DataFile)

	status := &pastel.MasterNodeStatus{
		Status: "hello",
	}

	pMock := pastelMock.NewMockClient(t)
	pMock.ListenOnMasterNodeStatus(status, nil)
	s.pastelClient = pMock

	err := s.waitSynchronization(context.Background())
	assert.Equal(t, err.Error(), "timeout")
}

func TestWaitSynchronizationError(t *testing.T) {
	s := prepareService(t)
	defer s.db.Close()
	defer os.Remove(s.config.DataFile)
	errMsg := "error in wait synchronization"

	pMock := pastelMock.NewMockClient(t)
	pMock.ListenOnMasterNodeStatus(nil, errors.New(errMsg))
	s.pastelClient = pMock

	err := s.waitSynchronization(context.Background())
	assert.Equal(t, err.Error(), errMsg)
}

func TestRunTaskSuccessful(t *testing.T) {
	s := prepareService(t)
	defer s.db.Close()
	defer os.Remove(s.config.DataFile)

	// Prepare reg art
	ticket := pastel.RegTicket{}
	f := fuzz.New()
	f.Fuzz(&ticket)
	ticket.Height = 2

	b, err := json.Marshal(ticket.RegTicketData.NFTTicketData.AppTicketData)
	if err != nil {
		t.Fatalf("faied to marshal, err: %s", err)
	}
	ticket.RegTicketData.NFTTicketData.AppTicket = b

	b, err = json.Marshal(ticket.RegTicketData.NFTTicketData)
	if err != nil {
		t.Fatalf("faied to marshal, err: %s", err)
	}
	ticket.RegTicketData.NFTTicket = b

	pMock := pastelMock.NewMockClient(t)
	pMock.ListenOnRegTickets(pastel.RegTickets{
		ticket,
	}, nil)
	s.pastelClient = pMock

	// Prepare latest fingerprint
	setFp := generateFingerprint(t)
	ctx := context.Background()
	err = s.storeFingerprint(ctx, setFp)
	assert.True(t, err == nil)
	time.Sleep(1 * time.Second)

	// Prepare p2p client
	fp := randFloats(fingerprintSizeModel)
	fpBuffer := new(bytes.Buffer)
	_ = binary.Write(fpBuffer, binary.LittleEndian, fp)

	p2pClient := p2pMock.NewMockClient(t)
	p2pClient.ListenOnRetrieve(fpBuffer.Bytes(), nil)
	s.p2pClient = p2pClient

	err = s.runTask(context.Background())
	assert.True(t, err == nil)

	getFp, err := s.getLatestFingerprint(ctx)
	fingerprintsHash := ticket.RegTicketData.NFTTicketData.AppTicketData.FingerprintsHash

	assert.True(t, err == nil)
	assert.True(t, getFp != nil)
	assert.Equal(t, hex.EncodeToString(fingerprintsHash), getFp.Sha256HashOfArtImageFile)
	// assert.Equal(t, ticket.Height, getFp.NumberOfBlock)

	start, end := 0, fingerprintSizeModel1
	assert.Equal(t, fp[start:end], getFp.Model1ImageFingerprintVector)

	start, end = start+fingerprintSizeModel1, end+fingerprintSizeModel2
	assert.Equal(t, fp[start:end], getFp.Model2ImageFingerprintVector)

	start, end = start+fingerprintSizeModel2, end+fingerprintSizeModel3
	assert.Equal(t, fp[start:end], getFp.Model3ImageFingerprintVector)

	start, end = start+fingerprintSizeModel3, end+fingerprintSizeModel4
	assert.Equal(t, fp[start:end], getFp.Model4ImageFingerprintVector)
}
