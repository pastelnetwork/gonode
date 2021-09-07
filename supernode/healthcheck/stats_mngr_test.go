package healthcheck

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ClientMock struct {
	stats map[string]interface{}
	err   error
}

func (cl *ClientMock) Stats(_ context.Context) (map[string]interface{}, error) {
	return cl.stats, cl.err
}

func TestStatsMngrNew(t *testing.T) {
	mngr := NewStatsMngr(1 * time.Second)
	clMock := &ClientMock{stats: nil, err: nil}
	mngr.Add("mock", clMock)

	assert.Equal(t, 1, len(mngr.clients))
	_, ok := mngr.clients["mock"]
	assert.True(t, ok)
}

func TestStatsMngrStats(t *testing.T) {
	mngr := NewStatsMngr(1 * time.Second)
	clMock := &ClientMock{stats: map[string]interface{}{"key": "value"}, err: nil}
	mngr.Add("mock", clMock)

	stats, err := mngr.Stats(context.Background())

	assert.Nil(t, err)
	mockStatsRaw, ok := stats["mock"]
	assert.True(t, ok)
	mockStats, ok := mockStatsRaw.(map[string]interface{})
	assert.True(t, ok)
	value, ok := mockStats["key"]
	assert.Equal(t, "value", value)
	assert.True(t, ok)
}
