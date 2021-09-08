package healthcheck

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ClientMock is an stub implementation of a StatsClient
type ClientMock struct {
	stats map[string]interface{}
	err   error
}

// Stats returns expected values
func (cl *ClientMock) Stats(_ context.Context) (map[string]interface{}, error) {
	return cl.stats, cl.err
}

//TestStatsMngrNew test NewStatsMngr()
func TestStatsMngrNew(t *testing.T) {
	mngr := NewStatsMngr(nil)
	clMock := &ClientMock{stats: nil, err: nil}
	mngr.Add("mock", clMock)

	assert.Equal(t, 1, len(mngr.clients))
	_, ok := mngr.clients["mock"]
	assert.True(t, ok)
}

// TestStatsMngrStats tests Stats() function
func TestStatsMngrStats(t *testing.T) {
	mngr := NewStatsMngr(nil)
	clMock := &ClientMock{stats: map[string]interface{}{"key": "value"}, err: nil}
	mngr.Add("mock", clMock)

	stats, err := mngr.Stats(context.Background())

	assert.Nil(t, err)
	timestamp, ok := stats["time_stamp"]
	assert.True(t, ok)
	assert.Equal(t, "", timestamp)
	assert.True(t, ok)
}
