package ddscan

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"math"
	"math/rand"

	"os"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/metadb/rqlite/db"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/assert"
)

const createTableStatement = `CREATE TABLE image_hash_to_image_fingerprint_table (sha256_hash_of_art_image_file text, path_to_art_image_file text, new_model_image_fingerprint_vector array, datetime_fingerprint_added_to_database text, thumbnail_of_image text, request_type text, open_api_subset_id_string text)`

func prepareService(_ *testing.T) *service {
	tmpfile, err := ioutil.TempFile("", "dupe_detection_image_fingerprint_database.sqlite")
	if err != nil {
		panic(err.Error())
	}
	db, err := db.Open(tmpfile.Name(), true)
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
		res[i] = 0.0 + rand.Float64()*(math.MaxFloat64)
	}
	return res
}

func generateFingerprint(_ *testing.T) *dupeDetectionFingerprints {
	fps := randFloats(1500)

	b := make([]byte, 10)
	rand.Read(b)
	return &dupeDetectionFingerprints{
		Sha256HashOfArtImageFile:           hex.EncodeToString(b),
		ImageFingerprintVector:             fps,
		DatetimeFingerprintAddedToDatabase: time.Now().Format("2006-01-02 15:04:05"),
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
	assert.Equal(t, setFp.ImageFingerprintVector, getFp.ImageFingerprintVector)
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
	assert.Equal(t, err.Error(), "timeout expired")
}

func TestWaitSynchronizationError(t *testing.T) {
	s := prepareService(t)
	defer s.db.Close()
	defer os.Remove(s.config.DataFile)
	errMsg := "timeout expired"

	pMock := pastelMock.NewMockClient(t)
	pMock.ListenOnMasterNodeStatus(nil, errors.New(errMsg))
	s.pastelClient = pMock

	err := s.waitSynchronization(context.Background())
	assert.Equal(t, err.Error(), errMsg)
}

// /* TODO: update later
// func TestRunTaskSuccessful(t *testing.T) {
// 	s := prepareService(t)
// 	defer s.db.Close()
// 	defer os.Remove(s.config.DataFile)

// 	// Prepare reg art
// 	ticket := pastel.RegTicket{}
// 	f := fuzz.New()
// 	f.Fuzz(&ticket)
// 	ticket.Height = 2

// 	b, err := json.Marshal(ticket.RegTicketData.NFTTicketData.AppTicketData)
// 	if err != nil {
// 		t.Fatalf("faied to marshal, err: %s", err)
// 	}
// 	ticket.RegTicketData.NFTTicketData.AppTicket = b85.Encode(b)

// 	b, err = json.Marshal(ticket.RegTicketData.NFTTicketData)
// 	if err != nil {
// 		t.Fatalf("faied to marshal, err: %s", err)
// 	}
// 	ticket.RegTicketData.NFTTicket = b

// 	pMock := pastelMock.NewMockClient(t)
// 	pMock.ListenOnRegTickets(pastel.RegTickets{
// 		ticket,
// 	}, nil)
// 	s.pastelClient = pMock

// 	// Prepare latest fingerprint
// 	setFp := generateFingerprint(t)
// 	ctx := context.Background()
// 	err = s.storeFingerprint(ctx, setFp)
// 	assert.True(t, err == nil)
// 	time.Sleep(1 * time.Second)

// 	// Prepare p2p client
// 	fp := randFloats(fingerprintSizeModel)

// 	fpBuffer := new(bytes.Buffer)
// 	_ = binary.Write(fpBuffer, binary.LittleEndian, fp)

// 	p2pClient := p2pMock.NewMockClient(t)
// 	p2pClient.ListenOnRetrieve(fpBuffer.Bytes(), nil)
// 	s.p2pClient = p2pClient

// 	err = s.runTask(context.Background())
// 	assert.True(t, err == nil)

// 	getFp, err := s.getLatestFingerprint(ctx)
// 	// fingerprintsHash := ticket.RegTicketData.NFTTicketData.AppTicketData.FingerprintsHash
// 	ddFPIDs := ticket.RegTicketData.NFTTicketData.AppTicketData.DDAndFingerprintsIDs
// 	// if len(ddFPIDs) < 1 {
// 	// 	return errors.Errorf("RegTicketData does not have properly registered DD and Fingerprint ID's")
// 	// }

// 	//just about everything will fail from this point onwards because the AppTicketData either doesn't exist or is randomly generated
// 	assert.True(t, len(ddFPIDs) > 0)
// 	check_for_differences := false
// 	for _, id := range ddFPIDs {
// 		if id != ddFPIDs[0] {
// 			check_for_differences = true
// 		}
// 	}
// 	assert.False(t, check_for_differences)
// 	fingerprintsHash := ddFPIDs[0]

// 	assert.True(t, err == nil)
// 	assert.True(t, getFp != nil)

// 	assert.Equal(t, hex.EncodeToString([]byte(fingerprintsHash)), getFp.Sha256HashOfArtImageFile)
// 	// assert.Equal(t, ticket.Height, getFp.NumberOfBlock)

// 	start, end := 0, fingerprintSizeModel1
// 	assert.Equal(t, fp[start:end], getFp.Model1ImageFingerprintVector)

// 	start, end = start+fingerprintSizeModel1, end+fingerprintSizeModel2
// 	assert.Equal(t, fp[start:end], getFp.Model2ImageFingerprintVector)

// 	start, end = start+fingerprintSizeModel2, end+fingerprintSizeModel3
// 	assert.Equal(t, fp[start:end], getFp.Model3ImageFingerprintVector)

// 	start, end = start+fingerprintSizeModel3, end+fingerprintSizeModel4
// 	assert.Equal(t, fp[start:end], getFp.Model4ImageFingerprintVector)
// }

// */
