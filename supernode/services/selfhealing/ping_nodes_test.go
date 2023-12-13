package selfhealing_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/supernode/services/selfhealing"
	"github.com/stretchr/testify/assert"
)

func TestGetInfoToInsert(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()

	tests := []struct {
		testCase     string
		existedInfo  *types.PingInfo
		receivedInfo *types.PingInfo
		expectedInfo *types.PingInfo
	}{
		{
			testCase:     "when node is offline, should increment total ping only",
			existedInfo:  &types.PingInfo{TotalPings: 2, TotalSuccessfulPings: 2},
			receivedInfo: &types.PingInfo{TotalPings: 0, IsOnline: false},
			expectedInfo: &types.PingInfo{TotalPings: 3, TotalSuccessfulPings: 2, IsOnline: false},
		},
		{
			testCase:     "when node is online, should increment total ping & total successful pings",
			existedInfo:  &types.PingInfo{TotalPings: 2, TotalSuccessfulPings: 2, LastSeen: sql.NullTime{Time: now, Valid: true}},
			receivedInfo: &types.PingInfo{TotalPings: 0, IsOnline: true},
			expectedInfo: &types.PingInfo{TotalPings: 3, IsOnline: true, TotalSuccessfulPings: 3},
		},
		{
			testCase:     "when node is offline, last seen should be the same as existed record",
			existedInfo:  &types.PingInfo{TotalPings: 2, TotalSuccessfulPings: 2, LastSeen: sql.NullTime{Time: now, Valid: true}},
			receivedInfo: &types.PingInfo{TotalPings: 0, IsOnline: false, LastResponseTime: 0},
			expectedInfo: &types.PingInfo{TotalPings: 3, TotalSuccessfulPings: 2, IsOnline: false, LastSeen: sql.NullTime{Time: now, Valid: true}},
		},
		{
			testCase:     "when node is online, last seen should be the time node responded at",
			existedInfo:  &types.PingInfo{TotalPings: 2, TotalSuccessfulPings: 2, LastSeen: sql.NullTime{Time: now.Add(5 * time.Second), Valid: true}},
			receivedInfo: &types.PingInfo{TotalPings: 0, IsOnline: true, LastSeen: sql.NullTime{Time: now.Add(15 * time.Second), Valid: true}},
			expectedInfo: &types.PingInfo{TotalPings: 3, IsOnline: true, TotalSuccessfulPings: 3, LastSeen: sql.NullTime{Time: now.Add(15 * time.Second), Valid: true}},
		},
		{
			testCase:     "when node is offline, cumulative response time should be the same, avg response time should be recalculated ",
			existedInfo:  &types.PingInfo{TotalPings: 2, TotalSuccessfulPings: 2, LastSeen: sql.NullTime{Time: now, Valid: true}, CumulativeResponseTime: now.Sub(now.Add(-20 * time.Second)).Seconds(), AvgPingResponseTime: 24},
			receivedInfo: &types.PingInfo{TotalPings: 0, IsOnline: false, AvgPingResponseTime: 0.0},
			expectedInfo: &types.PingInfo{TotalPings: 3, TotalSuccessfulPings: 2, IsOnline: false, AvgPingResponseTime: 10, LastSeen: sql.NullTime{Time: now, Valid: true}, CumulativeResponseTime: now.Sub(now.Add(-20 * time.Second)).Seconds()},
		},
		{
			testCase:     "when node is online, cumulative response time should add the last ping time, and avg response time should be recalculated",
			existedInfo:  &types.PingInfo{TotalPings: 2, TotalSuccessfulPings: 2, LastSeen: sql.NullTime{Time: now, Valid: true}, CumulativeResponseTime: 20 * time.Second.Seconds(), AvgPingResponseTime: 24},
			receivedInfo: &types.PingInfo{TotalPings: 0, IsOnline: true, LastResponseTime: 4.0, LastSeen: sql.NullTime{Time: now, Valid: true}},
			expectedInfo: &types.PingInfo{TotalPings: 3, TotalSuccessfulPings: 3, IsOnline: true, AvgPingResponseTime: 8, LastSeen: sql.NullTime{Time: now, Valid: true}, CumulativeResponseTime: 24 * time.Second.Seconds(), LastResponseTime: 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testCase, func(t *testing.T) {
			assert.Equal(t, tt.expectedInfo, selfhealing.GetPingInfoToInsert(tt.existedInfo, tt.receivedInfo))
		})
	}
}
