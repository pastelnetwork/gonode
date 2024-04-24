package queries

import (
	"time"

	"github.com/pastelnetwork/gonode/common/types"
)

type PingHistoryQueries interface {
	UpsertPingHistory(pingInfo types.PingInfo) error
	GetPingInfoBySupernodeID(supernodeID string) (*types.PingInfo, error)
	GetAllPingInfos() (types.PingInfos, error)
	GetWatchlistPingInfo() ([]types.PingInfo, error)
	GetAllPingInfoForOnlineNodes() (types.PingInfos, error)
	UpdatePingInfo(supernodeID string, isOnWatchlist, isAdjusted bool) error

	UpdateSCMetricsBroadcastTimestamp(nodeID string, broadcastAt time.Time) error
	UpdateMetricsBroadcastTimestamp(nodeID string) error
	UpdateGenerationMetricsBroadcastTimestamp(nodeID string) error
	UpdateExecutionMetricsBroadcastTimestamp(nodeID string) error
	UpdateHCMetricsBroadcastTimestamp(nodeID string, broadcastAt time.Time) error
}

// UpsertPingHistory inserts/update ping information into the ping_history table
func (s *SQLiteStore) UpsertPingHistory(pingInfo types.PingInfo) error {
	now := time.Now().UTC()

	const upsertQuery = `
		INSERT INTO ping_history (
			supernode_id, ip_address, total_pings, total_successful_pings, 
			avg_ping_response_time, is_online, is_on_watchlist, is_adjusted, last_seen, cumulative_response_time,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(supernode_id) 
		DO UPDATE SET
			total_pings = excluded.total_pings,
			total_successful_pings = excluded.total_successful_pings,
			avg_ping_response_time = excluded.avg_ping_response_time,
			is_online = excluded.is_online,
			is_on_watchlist = excluded.is_on_watchlist,
			is_adjusted = excluded.is_adjusted,
		    last_seen = excluded.last_seen,
		    cumulative_response_time = excluded.cumulative_response_time,
			updated_at = excluded.updated_at;`

	_, err := s.db.Exec(upsertQuery,
		pingInfo.SupernodeID, pingInfo.IPAddress, pingInfo.TotalPings,
		pingInfo.TotalSuccessfulPings, pingInfo.AvgPingResponseTime,
		pingInfo.IsOnline, pingInfo.IsOnWatchlist, pingInfo.IsAdjusted, pingInfo.LastSeen.Time, pingInfo.CumulativeResponseTime, now, now)
	if err != nil {
		return err
	}

	return nil
}

// GetPingInfoBySupernodeID retrieves a ping history record by supernode ID
func (s *SQLiteStore) GetPingInfoBySupernodeID(supernodeID string) (*types.PingInfo, error) {
	const selectQuery = `
        SELECT id, supernode_id, ip_address, total_pings, total_successful_pings,
               avg_ping_response_time, is_online, is_on_watchlist, is_adjusted, last_seen, cumulative_response_time,
               created_at, updated_at
        FROM ping_history
        WHERE supernode_id = ?;`

	var pingInfo types.PingInfo
	row := s.db.QueryRow(selectQuery, supernodeID)

	// Scan the row into the PingInfo struct
	err := row.Scan(
		&pingInfo.ID, &pingInfo.SupernodeID, &pingInfo.IPAddress, &pingInfo.TotalPings,
		&pingInfo.TotalSuccessfulPings, &pingInfo.AvgPingResponseTime,
		&pingInfo.IsOnline, &pingInfo.IsOnWatchlist, &pingInfo.IsAdjusted, &pingInfo.LastSeen, &pingInfo.CumulativeResponseTime,
		&pingInfo.CreatedAt, &pingInfo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &pingInfo, nil
}

// GetWatchlistPingInfo retrieves all the nodes that are on watchlist
func (s *SQLiteStore) GetWatchlistPingInfo() ([]types.PingInfo, error) {
	const selectQuery = `
        SELECT id, supernode_id, ip_address, total_pings, total_successful_pings,
               avg_ping_response_time, is_online, is_on_watchlist, is_adjusted, last_seen, cumulative_response_time,
               created_at, updated_at
        FROM ping_history
        WHERE is_on_watchlist = true AND is_adjusted = false;`

	rows, err := s.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pingInfos types.PingInfos
	for rows.Next() {
		var pingInfo types.PingInfo
		if err := rows.Scan(
			&pingInfo.ID, &pingInfo.SupernodeID, &pingInfo.IPAddress, &pingInfo.TotalPings,
			&pingInfo.TotalSuccessfulPings, &pingInfo.AvgPingResponseTime,
			&pingInfo.IsOnline, &pingInfo.IsOnWatchlist, &pingInfo.IsAdjusted, &pingInfo.LastSeen, &pingInfo.CumulativeResponseTime,
			&pingInfo.CreatedAt, &pingInfo.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pingInfos = append(pingInfos, pingInfo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pingInfos, nil
}

// UpdatePingInfo updates the ping info
func (s *SQLiteStore) UpdatePingInfo(supernodeID string, isOnWatchlist, isAdjusted bool) error {
	// Update query
	const updateQuery = `
UPDATE ping_history
SET is_adjusted = ?, is_on_watchlist = ?
WHERE supernode_id = ?;`

	// Execute the update query
	_, err := s.db.Exec(updateQuery, isAdjusted, isOnWatchlist, supernodeID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateMetricsBroadcastTimestamp updates the ping info metrics_last_broadcast_at
func (s *SQLiteStore) UpdateMetricsBroadcastTimestamp(nodeID string) error {
	// Update query
	const updateQuery = `
UPDATE ping_history
SET metrics_last_broadcast_at = ?
WHERE supernode_id = ?;`

	// Execute the update query
	_, err := s.db.Exec(updateQuery, time.Now().UTC(), nodeID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateGenerationMetricsBroadcastTimestamp updates the ping info generation_metrics_last_broadcast_at
func (s *SQLiteStore) UpdateGenerationMetricsBroadcastTimestamp(nodeID string) error {
	// Update query
	const updateQuery = `
UPDATE ping_history
SET generation_metrics_last_broadcast_at = ?
WHERE supernode_id = ?;`

	// Execute the update query
	_, err := s.db.Exec(updateQuery, time.Now().Add(-180*time.Minute).UTC(), nodeID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateExecutionMetricsBroadcastTimestamp updates the ping info execution_metrics_last_broadcast_at
func (s *SQLiteStore) UpdateExecutionMetricsBroadcastTimestamp(nodeID string) error {
	// Update query
	const updateQuery = `
UPDATE ping_history
SET execution_metrics_last_broadcast_at = ?
WHERE supernode_id = ?;`

	// Execute the update query
	_, err := s.db.Exec(updateQuery, time.Now().Add(-180*time.Minute).UTC(), nodeID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateSCMetricsBroadcastTimestamp updates the SC metrics last broadcast at timestamp
func (s *SQLiteStore) UpdateSCMetricsBroadcastTimestamp(nodeID string, updatedAt time.Time) error {
	// Update query
	const updateQuery = `
UPDATE ping_history
SET metrics_last_broadcast_at = ?
WHERE supernode_id = ?;`

	// Execute the update query
	_, err := s.db.Exec(updateQuery, time.Now().UTC().Add(-180*time.Minute), nodeID)
	if err != nil {
		return err
	}

	return nil
}

// GetAllPingInfos retrieves all ping infos
func (s *SQLiteStore) GetAllPingInfos() (types.PingInfos, error) {
	const selectQuery = `
        SELECT id, supernode_id, ip_address, total_pings, total_successful_pings,
               avg_ping_response_time, is_online, is_on_watchlist, is_adjusted, last_seen, cumulative_response_time,
               created_at, updated_at
        FROM ping_history
        `
	rows, err := s.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pingInfos types.PingInfos
	for rows.Next() {

		var pingInfo types.PingInfo
		if err := rows.Scan(
			&pingInfo.ID, &pingInfo.SupernodeID, &pingInfo.IPAddress, &pingInfo.TotalPings,
			&pingInfo.TotalSuccessfulPings, &pingInfo.AvgPingResponseTime,
			&pingInfo.IsOnline, &pingInfo.IsOnWatchlist, &pingInfo.IsAdjusted, &pingInfo.LastSeen, &pingInfo.CumulativeResponseTime,
			&pingInfo.CreatedAt, &pingInfo.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pingInfos = append(pingInfos, pingInfo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pingInfos, nil
}

// GetAllPingInfoForOnlineNodes retrieves all ping infos for nodes that are online
func (s *SQLiteStore) GetAllPingInfoForOnlineNodes() (types.PingInfos, error) {
	const selectQuery = `
        SELECT id, supernode_id, ip_address, total_pings, total_successful_pings,
               avg_ping_response_time, is_online, is_on_watchlist, is_adjusted, last_seen, cumulative_response_time, 
               metrics_last_broadcast_at, generation_metrics_last_broadcast_at, execution_metrics_last_broadcast_at,
               created_at, updated_at
        FROM ping_history
        WHERE is_online = true`
	rows, err := s.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pingInfos types.PingInfos
	for rows.Next() {

		var pingInfo types.PingInfo
		if err := rows.Scan(
			&pingInfo.ID, &pingInfo.SupernodeID, &pingInfo.IPAddress, &pingInfo.TotalPings,
			&pingInfo.TotalSuccessfulPings, &pingInfo.AvgPingResponseTime,
			&pingInfo.IsOnline, &pingInfo.IsOnWatchlist, &pingInfo.IsAdjusted, &pingInfo.LastSeen, &pingInfo.CumulativeResponseTime,
			&pingInfo.MetricsLastBroadcastAt, &pingInfo.GenerationMetricsLastBroadcastAt, &pingInfo.ExecutionMetricsLastBroadcastAt,
			&pingInfo.CreatedAt, &pingInfo.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pingInfos = append(pingInfos, pingInfo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pingInfos, nil
}

// UpdateHCMetricsBroadcastTimestamp updates health-check challenges last broadcast at
func (s *SQLiteStore) UpdateHCMetricsBroadcastTimestamp(nodeID string, updatedAt time.Time) error {
	// Update query
	const updateQuery = `
UPDATE ping_history
SET health_check_metrics_last_broadcast_at = ?
WHERE supernode_id = ?;`

	// Execute the update query
	_, err := s.db.Exec(updateQuery, time.Now().UTC().Add(-180*time.Minute), nodeID)
	if err != nil {
		return err
	}

	return nil
}
