package pastelblock

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pastelnetwork/gonode/hermes/domain"
	"github.com/pastelnetwork/gonode/hermes/store"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/assert"
)

func TestProcessBlocksWhenNoRecordExists(t *testing.T) {
	s, tmpFileName := prepareTestService(t)
	defer os.Remove(tmpFileName)

	pMock := pastelMock.NewMockClient(t)
	pMock.ListenOnGetBlockCount(1, nil)

	pMock.ListenOnGetBlockHash("test-hash", nil)

	s.pastelClient = pMock

	err := s.run(context.Background())
	assert.NoError(t, err)

	latestBlock, err := s.store.GetLatestPastelBlock(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, latestBlock.BlockHash, "test-hash")
	assert.Equal(t, latestBlock.BlockHeight, int32(1))
}

func TestProcessBlocksWhenRecordExistsWithEqualHeight(t *testing.T) {
	s, tmpFileName := prepareTestService(t)
	defer os.Remove(tmpFileName)

	pMock := pastelMock.NewMockClient(t)
	pMock.ListenOnGetBlockCount(1, nil)

	pMock.ListenOnGetBlockHash("test-hash", nil)

	s.pastelClient = pMock

	err := s.store.StorePastelBlock(context.Background(), domain.PastelBlock{BlockHeight: 1, BlockHash: "test-hash"})
	assert.NoError(t, err)

	err = s.run(context.Background())
	assert.NoError(t, err)

	_, err = s.store.GetLatestPastelBlock(context.Background())
	assert.NoError(t, err)
}

func TestProcessBlocksWhenRecordExistsWithLessHeight(t *testing.T) {
	s, tmpFileName := prepareTestService(t)
	defer os.Remove(tmpFileName)

	pMock := pastelMock.NewMockClient(t)
	pMock.ListenOnGetBlockCount(2, nil)

	pMock.ListenOnGetBlockHash("test-hash-2", nil)

	s.pastelClient = pMock

	err := s.store.StorePastelBlock(context.Background(), domain.PastelBlock{BlockHeight: 1, BlockHash: "test-hash"})
	assert.NoError(t, err)

	err = s.run(context.Background())
	assert.NoError(t, err)

	latestBlock, err := s.store.GetLatestPastelBlock(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, latestBlock.BlockHash, "test-hash-2")
	assert.Equal(t, latestBlock.BlockHeight, int32(2))
}

func prepareTestService(t *testing.T) (*pastelBlockService, string) {
	tmpFile, err := ioutil.TempFile("", "dupe_detection_image_fingerprint_database.sqlite")
	if err != nil {
		panic(err.Error())
	}
	store, err := store.NewSQLiteStore(tmpFile.Name())
	assert.Nil(t, err)

	tmpFile.Close()

	s := &pastelBlockService{
		store: store,
	}

	return s, tmpFile.Name()
}
