package hermes

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/pastelnetwork/gonode/hermes/service/hermes/domain"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestIsChainReorgDetected(t *testing.T) {
	s := prepareService(t)
	defer os.Remove(s.config.DataFile)

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

func storeBlockHashes(s *service, blockCount int32) {
	for i := int32(1); i <= blockCount; i++ {
		s.store.StorePastelBlock(context.Background(), domain.PastelBlock{
			BlockHeight: i,
			BlockHash:   fmt.Sprintf("test-hash-%d", i),
		})
	}
}
