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
		TxIDTimestamp:                              1683877733,
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

/*
func TestFetchDDFpFileAndStoreFingerprints(t *testing.T) {
	type args struct {
		tType   string
		series  string
		actTXID string
		regTXID string
		ddFpIDs []string
	}

	tests := []struct {
		name               string
		args               args
		mockStore          func(*ddmock.DDStore)
		mockPastelClient   func(*pastelMock.Client)
		mockHermesP2P      func(*nodemock.HermesP2PInterface)
		expectedStored     bool
		expectedP2PRunning bool
	}{
		{
			name: "ticket already exists",
			args: args{
				tType:   "sense",
				series:  "series1",
				actTXID: "acttxid1",
				regTXID: "regtxid1",
				ddFpIDs: []string{"hash1"},
			},
			mockStore: func(store *ddmock.DDStore) {
				store.On("IfFingerprintExistsByRegTxid", mock.Anything, "regtxid1").Return(true, nil)
			},
			mockPastelClient:   func(client *pastelMock.Client) {},
			mockHermesP2P:      func(p2p *nodemock.HermesP2PInterface) {},
			expectedStored:     true,
			expectedP2PRunning: true,
		},
		{
			name: "p2p service not running",
			args: args{
				tType:   "sense",
				series:  "series1",
				actTXID: "acttxid1",
				regTXID: "regtxid1",
				ddFpIDs: []string{"hash1"},
			},
			mockStore: func(store *ddmock.DDStore) {
				store.On("IfFingerprintExistsByRegTxid", mock.Anything, "regtxid1").Return(false, nil)
			},
			mockPastelClient: func(client *pastelMock.Client) {},
			mockHermesP2P: func(p2p *nodemock.HermesP2PInterface) {
				p2p.On("Retrieve", mock.Anything, "hash1").Return(nil, errors.New("p2p service is not running"))
			},
			expectedStored:     false,
			expectedP2PRunning: false,
		},
		{
			name: "no dd and fp id files can be unmarshalled",
			args: args{
				tType:   "sense",
				series:  "series1",
				actTXID: "acttxid1",
				regTXID: "regtxid1",
				ddFpIDs: []string{"hash1"},
			},
			mockStore: func(store *ddmock.DDStore) {
				store.On("IfFingerprintExistsByRegTxid", mock.Anything, "regtxid1").Return(false, nil)
			},
			mockPastelClient: func(client *pastelMock.Client) {},
			mockHermesP2P: func(p2p *nodemock.HermesP2PInterface) {
				p2p.On("Retrieve", mock.Anything, "hash1").Return(nil, errors.New("error"))
			},
			expectedStored:     false,
			expectedP2PRunning: true,
		},
		{
			name: "hash of candidate image file missing",
			args: args{
				tType:   "sense",
				series:  "series1",
				actTXID: "acttxid1",
				regTXID: "regtxid1",
				ddFpIDs: []string{"hash1"},
			},
			mockStore: func(store *ddmock.DDStore) {
				store.On("IfFingerprintExistsByRegTxid", mock.Anything, "regtxid1").Return(false, nil)
			},
			mockPastelClient: func(client *pastelMock.Client) {},
			mockHermesP2P: func(p2p *nodemock.HermesP2PInterface) {
				p2p.On("Retrieve", mock.Anything, "hash1").Return([]byte(`{"HashOfCandidateImageFile": ""}`), nil)
			},
			expectedStored:     false,
			expectedP2PRunning: true,
		},
		{
			name: "failed to query the dd database",
			args: args{
				tType:   "sense",
				series:  "series1",
				actTXID: "acttxid1",
				regTXID: "regtxid1",
				ddFpIDs: []string{"hash1"},
			},
			mockStore: func(store *ddmock.DDStore) {
				store.On("IfFingerprintExistsByRegTxid", mock.Anything, "regtxid1").Return(false, nil)
				store.On("IfFingerprintExists", mock.Anything, "image_hash").Return(false, errors.New("query error"))
			},
			mockPastelClient: func(client *pastelMock.Client) {},
			mockHermesP2P: func(p2p *nodemock.HermesP2PInterface) {
				p2p.On("Retrieve", mock.Anything, "hash1").Return([]byte(`{"HashOfCandidateImageFile": "image_hash"}`), nil)
			},
			expectedStored:     false,
			expectedP2PRunning: true,
		},
		{
			name: "ImageFingerprintOfCandidateImageFile missing",
			args: args{
				tType:   "sense",
				series:  "series1",
				actTXID: "acttxid1",
				regTXID: "regtxid1",
				ddFpIDs: []string{"hash1"},
			},
			mockStore: func(store *ddmock.DDStore) {
				store.On("IfFingerprintExistsByRegTxid", mock.Anything, "regtxid1").Return(false, nil)
				store.On("IfFingerprintExists", mock.Anything, "image_hash").Return(false, nil)
			},
			mockPastelClient: func(client *pastelMock.Client) {},
			mockHermesP2P: func(p2p *nodemock.HermesP2PInterface) {
				p2p.On("Retrieve", mock.Anything, "hash1").Return([]byte(`{"HashOfCandidateImageFile": "image_hash"}`), nil)
			},
			expectedStored:     false,
			expectedP2PRunning: true,
		},
		{
			name: "ImageFingerprintOfCandidateImageFile has zero length",
			args: args{
				tType:   "sense",
				series:  "series1",
				actTXID: "acttxid1",
				regTXID: "regtxid1",
				ddFpIDs: []string{"hash1"},
			},
			mockStore: func(store *ddmock.DDStore) {
				store.On("IfFingerprintExistsByRegTxid", mock.Anything, "regtxid1").Return(false, nil)
				store.On("IfFingerprintExists", mock.Anything, "image_hash").Return(false, nil)
			},
			mockPastelClient: func(client *pastelMock.Client) {},
			mockHermesP2P: func(p2p *nodemock.HermesP2PInterface) {
				p2p.On("Retrieve", mock.Anything, "hash1").Return([]byte(`{"HashOfCandidateImageFile": "image_hash", "ImageFingerprintOfCandidateImageFile": []}`), nil)
			},
			expectedStored:     false,
			expectedP2PRunning: true,
		},
		{
			name: "store fingerprint failed",
			args: args{
				tType:   "sense",
				series:  "series1",
				actTXID: "acttxid1",
				regTXID: "regtxid1",
				ddFpIDs: []string{"hash1"},
			},
			mockStore: func(store *ddmock.DDStore) {
				store.On("IfFingerprintExistsByRegTxid", mock.Anything, "regtxid1").Return(false, nil)
				store.On("IfFingerprintExists", mock.Anything, "image_hash").Return(false, nil)
				store.On("StoreFingerprint", mock.Anything, mock.Anything).Return(errors.New("store error"))
			},
			mockPastelClient: func(client *pastelMock.Client) {},
			mockHermesP2P: func(p2p *nodemock.HermesP2PInterface) {
				p2p.On("Retrieve", mock.Anything, "hash1").Return([]byte(`{"HashOfCandidateImageFile": "image_hash", "ImageFingerprintOfCandidateImageFile": [0.1, 0.2]}`), nil)
			},
			expectedStored:     false,
			expectedP2PRunning: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			store := &ddmock.DDStore{}
			tt.mockStore(store)

			pastelClient := &pastelMock.Client{}
			tt.mockPastelClient(pastelClient)

			hermesP2P := &nodemock.HermesP2PInterface{}
			tt.mockHermesP2P(hermesP2P)

			sync := &synchronizer.Synchronizer{}
			s := &fingerprintService{
				pastelClient: pastelClient,
				store:        store,
				p2p:          hermesP2P,
				sync:         sync,
			}

			stored, p2pRunning := s.fetchDDFpFileAndStoreFingerprints(ctx, tt.args.tType, tt.args.series, tt.args.actTXID, tt.args.regTXID, tt.args.ddFpIDs)
			if stored != tt.expectedStored {
				t.Errorf("fetchDDFpFileAndStoreFingerprints() stored = %v, expectedStored = %v", stored, tt.expectedStored)
			}
			if p2pRunning != tt.expectedP2PRunning {
				t.Errorf("fetchDDFpFileAndStoreFingerprints() p2pRunning = %v, expectedP2PRunning = %v", p2pRunning, tt.expectedP2PRunning)
			}
		})
	}

}
*/
