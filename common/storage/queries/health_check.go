package queries

import (
	"context"
	"fmt"
	"time"

	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils/metrics"
)

type HealthCheckChallengeQueries interface {
	InsertHealthCheckChallengeMessage(challenge types.HealthCheckChallengeLogMessage) error
	InsertBroadcastHealthCheckMessage(challenge types.BroadcastHealthCheckLogMessage) error
	QueryHCChallengeMessage(challengeID string, messageType int) (challengeMessage types.HealthCheckChallengeLogMessage, err error)
	GetHealthCheckChallengeMetricsByChallengeID(challengeID string) ([]types.HealthCheckChallengeLogMessage, error)

	GetHCMetricsByChallengeIDAndMessageType(challengeID string, messageType types.MessageType) ([]types.HealthCheckChallengeLogMessage, error)
	BatchInsertHCMetrics(metrics []types.HealthCheckChallengeLogMessage) error
	HealthCheckChallengeMetrics(timestamp time.Time) ([]types.HealthCheckChallengeLogMessage, error)
	InsertHealthCheckChallengeMetric(metric types.HealthCheckChallengeMetric) error
	GetHCSummaryStats(from time.Time) (hcMetrics metrics.HCMetrics, err error)
	GetTotalHCGeneratedAndProcessedAndEvaluated(from time.Time) (metrics.HCMetrics, error)
	GetMetricsDataByHealthCheckChallengeID(ctx context.Context, challengeID string) ([]types.HealthCheckMessage, error)
	GetLastNHCMetrics() ([]types.NHcMetric, error)

	GetDistinctHCChallengeIDsCountForScoreAggregation(after, before time.Time) (int, error)
	GetDistinctHCChallengeIDs(after, before time.Time, batchNumber int) ([]string, error)
	BatchInsertHCScoreAggregationChallenges(challengeIDs []string, isAggregated bool) error
}

// GetTotalHCGeneratedAndProcessedAndEvaluated retrieves the total health-check challenges generated/processed/evaluated
func (s *SQLiteStore) GetTotalHCGeneratedAndProcessedAndEvaluated(from time.Time) (metrics.HCMetrics, error) {
	metrics := metrics.HCMetrics{}

	// Query for total number of challenges
	totalChallengeQuery := "SELECT COUNT(DISTINCT challenge_id) FROM healthcheck_challenge_metrics WHERE message_type = 1 AND created_at > ?"
	err := s.db.QueryRow(totalChallengeQuery, from).Scan(&metrics.TotalChallenges)
	if err != nil {
		return metrics, err
	}

	// Query for total challenges responded
	totalChallengesProcessedQuery := "SELECT COUNT(DISTINCT challenge_id) FROM healthcheck_challenge_metrics WHERE message_type = 2 AND created_at > ?"
	err = s.db.QueryRow(totalChallengesProcessedQuery, from).Scan(&metrics.TotalChallengesProcessed)
	if err != nil {
		return metrics, err
	}

	totalChallengesEvaluatedQuery := "SELECT COUNT(DISTINCT challenge_id) FROM healthcheck_challenge_metrics WHERE message_type = 3 AND created_at > ?"
	err = s.db.QueryRow(totalChallengesEvaluatedQuery, from).Scan(&metrics.TotalChallengesEvaluatedByChallenger)
	if err != nil {
		return metrics, err
	}

	return metrics, nil
}

// GetHCObserversEvaluations retrieves the observer's evaluations
func (s *SQLiteStore) GetHCObserversEvaluations(from time.Time) ([]types.HealthCheckChallengeLogMessage, error) {
	var messages []types.HealthCheckChallengeLogMessage

	query := "SELECT id, challenge_id, message_type, data, sender_id, created_at, updated_at FROM healthcheck_challenge_metrics WHERE message_type = 4 and created_at > ?"
	rows, err := s.db.Query(query, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg types.HealthCheckChallengeLogMessage
		err := rows.Scan(&msg.ID, &msg.ChallengeID, &msg.MessageType, &msg.Data, &msg.Sender, &msg.CreatedAt, &msg.UpdatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// GetHCSummaryStats get health-check summary stats
func (s *SQLiteStore) GetHCSummaryStats(from time.Time) (hcMetrics metrics.HCMetrics, err error) {
	hcStats := metrics.HCMetrics{}
	hcMetrics, err = s.GetTotalHCGeneratedAndProcessedAndEvaluated(from)
	if err != nil {
		return hcMetrics, err
	}
	hcStats.TotalChallenges = hcMetrics.TotalChallenges
	hcStats.TotalChallengesProcessed = hcMetrics.TotalChallengesProcessed
	hcStats.TotalChallengesEvaluatedByChallenger = hcMetrics.TotalChallengesEvaluatedByChallenger

	hcObserversEvaluations, err := s.GetHCObserversEvaluations(from)
	if err != nil {
		return hcMetrics, err
	}
	log.WithField("observer_evaluations", len(hcObserversEvaluations)).Info("observer evaluations retrieved")

	observerEvaluationMetrics := processHCObserverEvaluations(hcObserversEvaluations)
	log.WithField("observer_evaluation_metrics", len(observerEvaluationMetrics)).Info("observer evaluation metrics retrieved")

	for _, obMetrics := range observerEvaluationMetrics {
		if obMetrics.ChallengesVerified >= 3 {
			hcMetrics.TotalChallengesVerified++
		} else {
			if obMetrics.FailedByInvalidTimestamps > 0 {
				hcMetrics.SlowResponsesObservedByObservers++
			}
			if obMetrics.FailedByInvalidSignatures > 0 {
				hcMetrics.InvalidSignaturesObservedByObservers++
			}
			if obMetrics.FailedByInvalidEvaluation > 0 {
				hcMetrics.InvalidEvaluationObservedByObservers++
			}
		}
	}

	return hcMetrics, nil
}

// GetHealthCheckChallengeMetricsByChallengeID gets the health-check challenge by ID
func (s *SQLiteStore) GetHealthCheckChallengeMetricsByChallengeID(challengeID string) ([]types.HealthCheckChallengeLogMessage, error) {
	const query = `
    SELECT id, challenge_id, message_type, data, sender_id, created_at, updated_at
    FROM healthcheck_challenge_metrics
    WHERE challenge_id = ?;`

	rows, err := s.db.Query(query, challengeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []types.HealthCheckChallengeLogMessage
	for rows.Next() {
		var m types.HealthCheckChallengeLogMessage
		err := rows.Scan(&m.ID, &m.ChallengeID, &m.MessageType, &m.Data, &m.Sender, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}

	return metrics, rows.Err()
}

// GetMetricsDataByHealthCheckChallengeID gets the metrics data by health-check challenge id
func (s *SQLiteStore) GetMetricsDataByHealthCheckChallengeID(ctx context.Context, challengeID string) (healthCheckChallengeMessages []types.HealthCheckMessage, err error) {
	hcMetrics, err := s.GetHealthCheckChallengeMetricsByChallengeID(challengeID)
	if err != nil {
		return healthCheckChallengeMessages, err
	}
	log.WithContext(ctx).WithField("rows", len(hcMetrics)).Info("health-check-challenge metrics row count")

	for _, hcMetric := range hcMetrics {
		msg := types.HealthCheckMessageData{}
		if err := json.Unmarshal(hcMetric.Data, &msg); err != nil {
			return healthCheckChallengeMessages, fmt.Errorf("cannot unmarshal health check challenge data: %w", err)
		}

		healthCheckChallengeMessages = append(healthCheckChallengeMessages, types.HealthCheckMessage{
			ChallengeID:     hcMetric.ChallengeID,
			MessageType:     types.HealthCheckMessageType(hcMetric.MessageType),
			Sender:          hcMetric.Sender,
			SenderSignature: hcMetric.SenderSignature,
			Data:            msg,
		})
	}

	return healthCheckChallengeMessages, nil
}

// InsertHealthCheckChallengeMessage inserts failed healthcheck challenge to db
func (s *SQLiteStore) InsertHealthCheckChallengeMessage(challenge types.HealthCheckChallengeLogMessage) error {
	now := time.Now().UTC()
	const insertQuery = "INSERT INTO healthcheck_challenge_messages(id, challenge_id, message_type, data, sender_id, sender_signature, created_at, updated_at) VALUES(NULL,?,?,?,?,?,?,?);"
	_, err := s.db.Exec(insertQuery, challenge.ChallengeID, challenge.MessageType, challenge.Data, challenge.Sender, challenge.SenderSignature, now, now)

	if err != nil {
		return err
	}

	return nil
}

// InsertHealthCheckChallengeMetric inserts the health-check challenge metrics
func (s *SQLiteStore) InsertHealthCheckChallengeMetric(m types.HealthCheckChallengeMetric) error {
	now := time.Now().UTC()

	const metricsQuery = "INSERT INTO healthcheck_challenge_metrics(id, challenge_id, message_type, data, sender_id, created_at, updated_at) VALUES(NULL,?,?,?,?,?,?) ON CONFLICT DO NOTHING;"
	_, err := s.db.Exec(metricsQuery, m.ChallengeID, m.MessageType, m.Data, m.SenderID, now, now)
	if err != nil {
		return err
	}

	return nil
}

// BatchInsertHCMetrics inserts the health-check challenges in a batch
func (s *SQLiteStore) BatchInsertHCMetrics(metrics []types.HealthCheckChallengeLogMessage) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
        INSERT OR IGNORE INTO healthcheck_challenge_metrics
        (id, challenge_id, message_type, data, sender_id, created_at, updated_at)
        VALUES (NULL,?,?,?,?,?,?)
    `)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, metric := range metrics {
		now := time.Now().UTC()

		_, err = stmt.Exec(metric.ChallengeID, metric.MessageType, metric.Data, metric.Sender, now, now)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// HealthCheckChallengeMetrics retrieves all the metrics needs to be broadcast
func (s *SQLiteStore) HealthCheckChallengeMetrics(timestamp time.Time) ([]types.HealthCheckChallengeLogMessage, error) {
	const query = `
    SELECT id, challenge_id, message_type, data, sender_id, created_at, updated_at
    FROM healthcheck_challenge_metrics
    WHERE created_at > ?
    `

	rows, err := s.db.Query(query, timestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []types.HealthCheckChallengeLogMessage
	for rows.Next() {
		var m types.HealthCheckChallengeLogMessage
		err := rows.Scan(&m.ID, &m.ChallengeID, &m.MessageType, &m.Data, &m.Sender, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}

	return metrics, rows.Err()
}

// InsertBroadcastHealthCheckMessage inserts healthcheck healthcheck challenge msg to db
func (s *SQLiteStore) InsertBroadcastHealthCheckMessage(challenge types.BroadcastHealthCheckLogMessage) error {
	now := time.Now().UTC()
	const insertQuery = "INSERT INTO broadcast_healthcheck_challenge_messages(id, challenge_id, data, challenger, recipient, observers, created_at, updated_at) VALUES(NULL,?,?,?,?,?,?,?);"
	_, err := s.db.Exec(insertQuery, challenge.ChallengeID, challenge.Data, challenge.Challenger, challenge.Recipient, challenge.Observers, now, now)
	if err != nil {
		return err
	}

	return nil
}

// QueryHCChallengeMessage retrieves healthcheck challenge message against challengeID and messageType
func (s *SQLiteStore) QueryHCChallengeMessage(challengeID string, messageType int) (challengeMessage types.HealthCheckChallengeLogMessage, err error) {
	const selectQuery = "SELECT * FROM healthcheck_challenge_messages WHERE challenge_id=? AND message_type=?"
	err = s.db.QueryRow(selectQuery, challengeID, messageType).Scan(
		&challengeMessage.ID, &challengeMessage.ChallengeID, &challengeMessage.MessageType, &challengeMessage.Data,
		&challengeMessage.Sender, &challengeMessage.SenderSignature, &challengeMessage.CreatedAt, &challengeMessage.UpdatedAt)

	if err != nil {
		return challengeMessage, err
	}

	return challengeMessage, nil
}

// GetHCMetricsByChallengeIDAndMessageType retrieves all the metrics by challengeID and messageType
func (s *SQLiteStore) GetHCMetricsByChallengeIDAndMessageType(challengeID string, messageType types.MessageType) ([]types.HealthCheckChallengeLogMessage, error) {
	const query = `
    SELECT id, challenge_id, message_type, data, sender_id, created_at, updated_at
    FROM healthcheck_challenge_metrics
    WHERE challenge_id = ?
    AND message_type = ?;`

	rows, err := s.db.Query(query, challengeID, int(messageType))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []types.HealthCheckChallengeLogMessage
	for rows.Next() {
		var m types.HealthCheckChallengeLogMessage
		err := rows.Scan(&m.ID, &m.ChallengeID, &m.MessageType, &m.Data, &m.Sender, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}

	return metrics, rows.Err()
}

func processHCObserverEvaluations(observersEvaluations []types.HealthCheckChallengeLogMessage) map[string]HCObserverEvaluationMetrics {
	evaluationMap := make(map[string]HCObserverEvaluationMetrics)

	for _, observerEvaluation := range observersEvaluations {
		var oe types.HealthCheckMessageData
		if err := json.Unmarshal(observerEvaluation.Data, &oe); err != nil {
			continue
		}

		oem, exists := evaluationMap[observerEvaluation.ChallengeID]
		if !exists {
			oem = HCObserverEvaluationMetrics{} // Initialize if not exists
		}

		if isHCObserverEvaluationVerified(oe.ObserverEvaluation) {
			oem.ChallengesVerified++
		} else {
			if !oe.ObserverEvaluation.IsChallengeTimestampOK ||
				!oe.ObserverEvaluation.IsProcessTimestampOK ||
				!oe.ObserverEvaluation.IsEvaluationTimestampOK {
				oem.FailedByInvalidTimestamps++
			}

			if !oe.ObserverEvaluation.IsChallengerSignatureOK ||
				!oe.ObserverEvaluation.IsRecipientSignatureOK {
				oem.FailedByInvalidSignatures++
			}

			if !oe.ObserverEvaluation.IsEvaluationResultOK {
				oem.FailedByInvalidEvaluation++
			}
		}

		evaluationMap[observerEvaluation.ChallengeID] = oem
	}

	return evaluationMap
}

func isHCObserverEvaluationVerified(observerEvaluation types.HealthCheckObserverEvaluationData) bool {
	if !observerEvaluation.IsEvaluationResultOK {
		return false
	}

	if !observerEvaluation.IsChallengerSignatureOK {
		return false
	}

	if !observerEvaluation.IsRecipientSignatureOK {
		return false
	}

	if !observerEvaluation.IsChallengeTimestampOK {
		return false
	}

	if !observerEvaluation.IsProcessTimestampOK {
		return false
	}

	if !observerEvaluation.IsEvaluationTimestampOK {
		return false
	}

	return true
}

// GetDistinctHCChallengeIDsCountForScoreAggregation gets the count of distinct challenge ids for score aggregation
func (s *SQLiteStore) GetDistinctHCChallengeIDsCountForScoreAggregation(after, before time.Time) (int, error) {
	query := `
        SELECT COUNT(DISTINCT challenge_id)
        FROM healthcheck_challenge_metrics
        WHERE message_type = 4 AND created_at >= ? AND created_at < ?
    `

	var challengeIDsCount int
	err := s.db.QueryRow(query, after, before).Scan(&challengeIDsCount)
	if err != nil {
		return 0, err
	}

	return challengeIDsCount, nil
}

// GetDistinctHCChallengeIDs retrieves the distinct challenge ids for score aggregation
func (s *SQLiteStore) GetDistinctHCChallengeIDs(after, before time.Time, batchNumber int) ([]string, error) {
	offset := batchNumber * batchSizeForChallengeIDsRetrieval

	query := `
        SELECT DISTINCT challenge_id
        FROM healthcheck_challenge_metrics
        WHERE message_type = 4 AND created_at >= ? AND created_at < ?
        LIMIT ? OFFSET ?
    `

	rows, err := s.db.Query(query, after, before, batchSizeForChallengeIDsRetrieval, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challengeIDs []string
	for rows.Next() {
		var challengeID string
		if err := rows.Scan(&challengeID); err != nil {
			return nil, err
		}
		challengeIDs = append(challengeIDs, challengeID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return challengeIDs, nil
}

// BatchInsertHCScoreAggregationChallenges inserts the batch of challenge ids for score aggregation
func (s *SQLiteStore) BatchInsertHCScoreAggregationChallenges(challengeIDs []string, isAggregated bool) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
        INSERT OR IGNORE INTO sc_score_aggregation_queue
        (challenge_id, is_aggregated, created_at, updated_at)
        VALUES (?,?,?,?)
    `)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, id := range challengeIDs {
		now := time.Now().UTC()

		_, err = stmt.Exec(id, isAggregated, now, now)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}
