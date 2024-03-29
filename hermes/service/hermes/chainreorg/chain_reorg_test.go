package chainreorg

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pastelnetwork/gonode/hermes/domain"
	"github.com/pastelnetwork/gonode/hermes/store"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestIsChainReorgDetected(t *testing.T) {
	s, tmpFileName := prepareTestService(t)
	defer os.Remove(tmpFileName)

	storeBlockHashes(s, int32(10))

	tests := []struct {
		testcase string
		expected bool
		setup    func(t *testing.T)
		expect   func(*testing.T, bool, bool, int32, error)
	}{
		{
			testcase: "when chain-reorg is not detected, it should return false with 0 as last good block count",
			expected: false,
			setup: func(t *testing.T) {
				pMock := pastelMock.NewMockClient(t)

				pMock.On(pastelMock.GetBlockCountMethod, mock.Anything).Return(int32(10), nil)

				var i int32
				for i = 10; i > 0; i-- {
					pMock.On(pastelMock.GetBlockHashMethod, mock.Anything, i).Return(fmt.Sprintf("test-hash-%d", i), nil)
				}

				s.pastelClient = pMock
			},
			expect: func(t *testing.T, expected, actual bool, lastGoodBlockCount int32, err error) {
				require.NoError(t, err)
				require.Equal(t, expected, actual)
				require.Equal(t, int32(0), lastGoodBlockCount)
			},
		},
		{
			testcase: "when chain-reorg detected, it should return true with last good block count",
			expected: true,
			setup: func(t *testing.T) {
				pMock := pastelMock.NewMockClient(t)

				pMock.On(pastelMock.GetBlockCountMethod, mock.Anything).Return(int32(10), nil)

				var i int32
				for i = 10; i > 5; i-- {
					pMock.On(pastelMock.GetBlockHashMethod, mock.Anything, i).Return(fmt.Sprintf("test-hash-%d", i), nil)
				}

				//hashes mismatched for 5,4 block, should return 3 as last good block count, with isDetected true
				pMock.On(pastelMock.GetBlockHashMethod, mock.Anything, int32(5)).Return("test", nil)
				pMock.On(pastelMock.GetBlockHashMethod, mock.Anything, int32(4)).Return("test", nil)

				for i = 3; i > 0; i-- {
					pMock.On(pastelMock.GetBlockHashMethod, mock.Anything, i).Return(fmt.Sprintf("test-hash-%d", i), nil)
				}

				s.pastelClient = pMock
			},
			expect: func(t *testing.T, expected, actual bool, lastGoodBlockCount int32, err error) {
				require.NoError(t, err)
				require.Equal(t, expected, actual)
				require.Equal(t, int32(3), lastGoodBlockCount)
			},
		},
	}

	for _, test := range tests {
		test := test // add this if there's subtest (t.Run)
		t.Run(test.testcase, func(t *testing.T) {
			test.setup(t)

			isDetected, lastGoodKnownBlockCount, err := s.IsChainReorgDetected(context.Background())

			test.expect(t, test.expected, isDetected, lastGoodKnownBlockCount, err)
		})

	}
}

func TestFixChainReorg(t *testing.T) {
	s, tmpFileName := prepareTestService(t)
	defer os.Remove(tmpFileName)

	storeBlockHashes(s, int32(10))

	tests := []struct {
		testcase                 string
		lastGoodKnownBlockHeight int32
		expected                 bool
		setup                    func(t *testing.T)
		expect                   func(*testing.T, int32, error)
	}{
		{
			testcase:                 "when chain-reorg is detected, it should fix from good block count up til the end",
			lastGoodKnownBlockHeight: 7,
			expected:                 false,
			setup: func(t *testing.T) {
				pMock := pastelMock.NewMockClient(t)

				pMock.On(pastelMock.GetBlockCountMethod, mock.Anything).Return(int32(10), nil)

				var i int32
				for i = 8; i <= 10; i++ {
					pMock.On(pastelMock.GetBlockHashMethod, mock.Anything, i).Return(fmt.Sprintf("test-hash-%d1", i), nil)
				}

				s.pastelClient = pMock
			},
			expect: func(t *testing.T, lastGoodBlockCount int32, err error) {
				require.NoError(t, err)

				pastelBlocks, err := getBlocks(s, lastGoodBlockCount)
				require.NoError(t, err)
				require.Equal(t, int32(8), pastelBlocks[0].BlockHeight)
				require.Equal(t, "test-hash-81", pastelBlocks[0].BlockHash)

				require.Equal(t, int32(9), pastelBlocks[1].BlockHeight)
				require.Equal(t, "test-hash-91", pastelBlocks[1].BlockHash)

				require.Equal(t, int32(10), pastelBlocks[2].BlockHeight)
				require.Equal(t, "test-hash-101", pastelBlocks[2].BlockHash)
			},
		},
	}

	for _, test := range tests {
		test := test // add this if there's subtest (t.Run)
		t.Run(test.testcase, func(t *testing.T) {
			test.setup(t)

			err := s.FixChainReorg(context.Background(), test.lastGoodKnownBlockHeight)

			test.expect(t, test.lastGoodKnownBlockHeight, err)
		})

	}
}

func storeBlockHashes(s *chainReorgService, blockCount int32) {
	for i := int32(1); i <= blockCount; i++ {
		s.store.StorePastelBlock(context.Background(), domain.PastelBlock{
			BlockHeight: i,
			BlockHash:   fmt.Sprintf("test-hash-%d", i),
		})
	}
}

func getBlocks(s *chainReorgService, blockCount int32) ([]domain.PastelBlock, error) {
	var pastelBlocks []domain.PastelBlock
	for i := blockCount + 1; i <= 10; i++ {
		pb, err := s.store.GetPastelBlockByHeight(context.Background(), i)
		if err != nil {
			return nil, err
		}

		pastelBlocks = append(pastelBlocks, pb)
	}

	return pastelBlocks, nil
}

func prepareTestService(t *testing.T) (*chainReorgService, string) {
	tmpfile, err := ioutil.TempFile("", "dupe_detection_image_fingerprint_database.sqlite")
	if err != nil {
		panic(err.Error())
	}
	store, err := store.NewSQLiteStore(tmpfile.Name())
	assert.Nil(t, err)

	tmpfile.Close()

	s := &chainReorgService{
		store: store,
	}

	return s, tmpfile.Name()
}
