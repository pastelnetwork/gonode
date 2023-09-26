package blocktracker

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakePastelClient struct {
	retBlockCnt int32
	retErr      error
}

func (fake *fakePastelClient) GetBlockCount(_ context.Context) (int32, error) {
	return fake.retBlockCnt, fake.retErr
}

func TestGetCountFirstTime(t *testing.T) {
	tests := []struct {
		name         string
		pastelClient *fakePastelClient
		expectErr    bool
	}{
		{
			name: "success",
			pastelClient: &fakePastelClient{
				retBlockCnt: 10,
				retErr:      nil,
			},
			expectErr: false,
		},
		{
			name: "fail",
			pastelClient: &fakePastelClient{
				retBlockCnt: 0,
				retErr:      errors.New("error"),
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := New(tt.pastelClient)
			tracker.retries = 1
			blkCnt, err := tracker.GetBlockCount()
			assert.Equal(t, tt.pastelClient.retBlockCnt, blkCnt)
			if tt.expectErr {
				assert.True(t, strings.Contains(err.Error(), tt.pastelClient.retErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetBlockCountNoRefresh(t *testing.T) {
	pastelClient := &fakePastelClient{
		retBlockCnt: 10,
		retErr:      errors.New("error"),
	}

	expectedBlk := int32(1)
	tracker := New(pastelClient)
	tracker.retries = 1
	tracker.curBlockCnt = expectedBlk
	tracker.lastRetried = time.Now().UTC()
	tracker.lastSuccess = time.Now().UTC()

	blkCnt, err := tracker.GetBlockCount()
	assert.Equal(t, expectedBlk, blkCnt)

	assert.Nil(t, err)
}

func TestGetBlockCountRefresh(t *testing.T) {
	expectedBlk := int32(10)
	pastelClient := &fakePastelClient{
		retBlockCnt: expectedBlk,
		retErr:      nil,
	}

	tracker := New(pastelClient)
	tracker.retries = 1
	tracker.curBlockCnt = 1
	tracker.lastRetried = time.Now().UTC().Add(-defaultSuccessUpdateDuration)
	tracker.lastSuccess = time.Now().UTC().Add(-defaultSuccessUpdateDuration)

	blkCnt, err := tracker.GetBlockCount()
	assert.Equal(t, expectedBlk, blkCnt)

	assert.Nil(t, err)
}
