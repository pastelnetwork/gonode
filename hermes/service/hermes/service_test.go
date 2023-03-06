package hermes

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"math"
	"math/rand"

	"github.com/pastelnetwork/gonode/hermes/service/hermes/domain"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/store"

	"os"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/assert"
)

func prepareService(t *testing.T) *service {
	tmpfile, err := ioutil.TempFile("", "dupe_detection_image_fingerprint_database.sqlite")
	if err != nil {
		panic(err.Error())
	}
	store, err := store.NewSQLiteStore(tmpfile.Name())
	assert.Nil(t, err)

	tmpfile.Close()

	s := &service{
		config: &Config{
			DataFile: tmpfile.Name(),
		},
		store: store,
	}

	return s
}

func randFloats(n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = 0.0 + rand.Float64()*(math.MaxFloat64)
	}
	return res
}

func generateFingerprint(_ *testing.T) *domain.DDFingerprints {
	fps := randFloats(1500)

	b := make([]byte, 10)
	rand.Read(b)
	return &domain.DDFingerprints{
		Sha256HashOfArtImageFile:                   hex.EncodeToString(b),
		ImageFingerprintVector:                     fps,
		DatetimeFingerprintAddedToDatabase:         time.Now().Format("2006-01-02 15:04:05"),
		OpenAPIGroupIDString:                       "x",
		DoesNotImpactTheFollowingCollectionsString: "y,z",
	}
}

func TestGetEmptyFingerprint(t *testing.T) {
	s := prepareService(t)
	defer os.Remove(s.config.DataFile)

	emptyFp, err := s.store.GetLatestFingerprints(context.Background())
	assert.True(t, emptyFp == nil)
	assert.Equal(t, err.Error(), "dd database is empty")
}

func TestSetFingerprint(t *testing.T) {
	s := prepareService(t)
	defer os.Remove(s.config.DataFile)

	err := s.store.StoreFingerprint(context.Background(), generateFingerprint(t))
	assert.True(t, err == nil)
}

func TestGetSetFingerprint(t *testing.T) {
	s := prepareService(t)
	defer os.Remove(s.config.DataFile)

	setFp := generateFingerprint(t)
	ctx := context.Background()
	err := s.store.StoreFingerprint(ctx, setFp)

	assert.True(t, err == nil)

	getFp, err := s.store.GetLatestFingerprints(ctx)

	assert.True(t, err == nil)

	findFp, err := s.store.IfFingerprintExists(ctx, setFp.Sha256HashOfArtImageFile)

	assert.True(t, err == nil)
	assert.True(t, getFp != nil)
	assert.Equal(t, setFp.Sha256HashOfArtImageFile, getFp.Sha256HashOfArtImageFile)
	assert.Equal(t, setFp.ImageFingerprintVector, getFp.ImageFingerprintVector)
	assert.True(t, findFp)

	findFp, err = s.store.IfFingerprintExists(ctx, "blahblah")
	assert.True(t, err == nil)
	assert.False(t, findFp)
}

func TestWaitSynchronizationSuccessful(t *testing.T) {
	s := prepareService(t)
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
	defer os.Remove(s.config.DataFile)
	errMsg := "timeout expired"

	pMock := pastelMock.NewMockClient(t)
	pMock.ListenOnMasterNodeStatus(nil, errors.New(errMsg))
	s.pastelClient = pMock

	err := s.waitSynchronization(context.Background())
	assert.Equal(t, err.Error(), errMsg)
}
