package fingerprint

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/hermes/domain"
	"github.com/pastelnetwork/gonode/hermes/store"
	"github.com/stretchr/testify/assert"
)

func prepareFGService(t *testing.T) (store.DDStore, store.CollectionStore, string) {
	tmpfile, err := ioutil.TempFile("", "dupe_detection_image_fingerprint_database.sqlite")
	if err != nil {
		panic(err.Error())
	}
	store, err := store.NewSQLiteStore(tmpfile.Name())
	assert.Nil(t, err)

	tmpfile.Close()

	return store, store, tmpfile.Name()
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
	ddStore, _, tmpFileName := prepareFGService(t)
	defer os.Remove(tmpFileName)

	emptyFp, err := ddStore.GetLatestFingerprints(context.Background())
	assert.True(t, emptyFp == nil)
	assert.Equal(t, "dd database is empty", err.Error())
}

func TestSetFingerprint(t *testing.T) {
	ddStore, _, tmpFileName := prepareFGService(t)
	defer os.Remove(tmpFileName)

	err := ddStore.StoreFingerprint(context.Background(), generateFingerprint(t))
	assert.True(t, err == nil)
}

func TestDoesNotImpactCollections(t *testing.T) {
	ddStore, cStore, tmpFileName := prepareFGService(t)
	defer os.Remove(tmpFileName)

	setFp := generateFingerprint(t)
	ctx := context.Background()
	err := ddStore.StoreFingerprint(ctx, setFp)

	assert.NoError(t, err)

	nImpactedCollections, err := cStore.GetDoesNotImpactCollections(ctx, setFp.Sha256HashOfArtImageFile)
	assert.NoError(t, err)
	assert.Equal(t, len(nImpactedCollections), 2)

	collectionNames := strings.Split(setFp.DoesNotImpactTheFollowingCollectionsString, ",")
	assert.Equal(t, nImpactedCollections[0].CollectionName, collectionNames[0])
	assert.Equal(t, nImpactedCollections[1].CollectionName, collectionNames[1])
}
