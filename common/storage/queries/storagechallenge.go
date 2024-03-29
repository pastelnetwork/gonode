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

type StorageChallengeQueries interface {
	InsertStorageChallengeMessage(challenge types.StorageChallengeLogMessage) error
	InsertBroadcastMessage(challenge types.BroadcastLogMessage) error
	QueryStorageChallengeMessage(challengeID string, messageType int) (challenge types.StorageChallengeLogMessage, err error)
	CleanupStorageChallenges() (err error)
	GetStorageChallengeMetricsByChallengeID(challengeID string) ([]types.StorageChallengeLogMessage, error)

	BatchInsertSCMetrics(metrics []types.StorageChallengeLogMessage) error
	StorageChallengeMetrics(timestamp time.Time) ([]types.StorageChallengeLogMessage, error)
	InsertStorageChallengeMetric(metric types.StorageChallengeMetric) error
	GetSCSummaryStats(from time.Time) (scMetrics metrics.SCMetrics, err error)
	GetTotalSCGeneratedAndProcessedAndEvaluated(from time.Time) (metrics.SCMetrics, error)
	GetChallengerEvaluations(from time.Time) ([]types.StorageChallengeLogMessage, error)
	GetObserversEvaluations(from time.Time) ([]types.StorageChallengeLogMessage, error)
	GetMetricsDataByStorageChallengeID(ctx context.Context, challengeID string) ([]types.Message, error)
	GetLastNSCMetrics() ([]types.NScMetric, error)
}

// InsertStorageChallengeMessage inserts failed storage challenge to db
func (s *SQLiteStore) InsertStorageChallengeMessage(challenge types.StorageChallengeLogMessage) error {
	now := time.Now().UTC()
	const insertQuery = "INSERT INTO storage_challenge_messages(id, challenge_id, message_type, data, sender_id, sender_signature, created_at, updated_at) VALUES(NULL,?,?,?,?,?,?,?) ON CONFLICT DO NOTHING;"
	_, err := s.db.Exec(insertQuery, challenge.ChallengeID, challenge.MessageType, challenge.Data, challenge.Sender, challenge.SenderSignature, now, now)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStore) InsertStorageChallengeMetric(m types.StorageChallengeMetric) error {
	now := time.Now().UTC()

	const metricsQuery = "INSERT INTO storage_challenge_metrics(id, challenge_id, message_type, data, sender_id, created_at, updated_at) VALUES(NULL,?,?,?,?,?,?) ON CONFLICT DO NOTHING;"
	_, err := s.db.Exec(metricsQuery, m.ChallengeID, m.MessageType, m.Data, m.SenderID, now, now)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStore) BatchInsertSCMetrics(metrics []types.StorageChallengeLogMessage) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
        INSERT OR IGNORE INTO storage_challenge_metrics
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

func (s *SQLiteStore) GetMetricsDataByStorageChallengeID(ctx context.Context, challengeID string) (storageChallengeMessages []types.Message, err error) {
	scMetrics, err := s.GetStorageChallengeMetricsByChallengeID(challengeID)
	if err != nil {
		return storageChallengeMessages, err
	}
	log.WithContext(ctx).WithField("rows", len(scMetrics)).Info("storage-challenge metrics row count")

	for _, scMetric := range scMetrics {
		msg := types.MessageData{}
		if err := json.Unmarshal(scMetric.Data, &msg); err != nil {
			return storageChallengeMessages, fmt.Errorf("cannot unmarshal storage challenge data: %w", err)
		}

		storageChallengeMessages = append(storageChallengeMessages, types.Message{
			ChallengeID:     scMetric.ChallengeID,
			MessageType:     types.MessageType(scMetric.MessageType),
			Sender:          scMetric.Sender,
			SenderSignature: scMetric.SenderSignature,
			Data:            msg,
		})
	}

	return storageChallengeMessages, nil
}

func (s *SQLiteStore) GetTotalSCGeneratedAndProcessedAndEvaluated(from time.Time) (metrics.SCMetrics, error) {
	metrics := metrics.SCMetrics{}

	// Query for total number of challenges
	totalChallengeQuery := "SELECT COUNT(DISTINCT challenge_id) FROM storage_challenge_metrics WHERE message_type = 1 AND created_at > ?"
	err := s.db.QueryRow(totalChallengeQuery, from).Scan(&metrics.TotalChallenges)
	if err != nil {
		return metrics, err
	}

	// Query for total challenges responded
	totalChallengesProcessedQuery := "SELECT COUNT(DISTINCT challenge_id) FROM storage_challenge_metrics WHERE message_type = 2 AND created_at > ?"
	err = s.db.QueryRow(totalChallengesProcessedQuery, from).Scan(&metrics.TotalChallengesProcessed)
	if err != nil {
		return metrics, err
	}

	totalChallengesEvaluatedQuery := "SELECT COUNT(DISTINCT challenge_id) FROM storage_challenge_metrics WHERE message_type = 3 AND created_at > ?"
	err = s.db.QueryRow(totalChallengesEvaluatedQuery, from).Scan(&metrics.TotalChallengesEvaluatedByChallenger)
	if err != nil {
		return metrics, err
	}

	return metrics, nil
}

func (s *SQLiteStore) GetChallengerEvaluations(from time.Time) ([]types.StorageChallengeLogMessage, error) {
	var messages []types.StorageChallengeLogMessage

	query := "SELECT id, challenge_id, message_type, data, sender_id, created_at, updated_at FROM storage_challenge_metrics WHERE message_type = 3 and created_at > ?"
	rows, err := s.db.Query(query, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg types.StorageChallengeLogMessage
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

func (s *SQLiteStore) GetObserversEvaluations(from time.Time) ([]types.StorageChallengeLogMessage, error) {
	var messages []types.StorageChallengeLogMessage

	query := "SELECT id, challenge_id, message_type, data, sender_id, created_at, updated_at FROM storage_challenge_metrics WHERE message_type = 4 and created_at > ?"
	rows, err := s.db.Query(query, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg types.StorageChallengeLogMessage
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

func (s *SQLiteStore) GetSCSummaryStats(from time.Time) (scMetrics metrics.SCMetrics, err error) {
	scStats := metrics.SCMetrics{}
	scMetrics, err = s.GetTotalSCGeneratedAndProcessedAndEvaluated(from)
	if err != nil {
		return scMetrics, err
	}
	scStats.TotalChallenges = scMetrics.TotalChallenges
	scStats.TotalChallengesProcessed = scMetrics.TotalChallengesProcessed
	scStats.TotalChallengesEvaluatedByChallenger = scMetrics.TotalChallengesEvaluatedByChallenger

	observersEvaluations, err := s.GetObserversEvaluations(from)
	if err != nil {
		return scMetrics, err
	}
	log.WithField("observer_evaluations", len(observersEvaluations)).Info("observer evaluations retrieved")

	observerEvaluationMetrics := processObserverEvaluations(observersEvaluations)
	log.WithField("observer_evaluation_metrics", len(observerEvaluationMetrics)).Info("observer evaluation metrics retrieved")

	for _, obMetrics := range observerEvaluationMetrics {
		if obMetrics.ChallengesVerified > 2 {
			scMetrics.TotalChallengesVerified++
		} else {
			if obMetrics.FailedByInvalidTimestamps > 0 {
				scMetrics.SlowResponsesObservedByObservers++
			}
			if obMetrics.FailedByInvalidSignatures > 0 {
				scMetrics.InvalidSignaturesObservedByObservers++
			}
			if obMetrics.FailedByInvalidEvaluation > 0 {
				scMetrics.InvalidEvaluationObservedByObservers++
			}
		}
	}

	return scMetrics, nil
}

// InsertBroadcastMessage inserts broadcast storage challenge msg to db
func (s *SQLiteStore) InsertBroadcastMessage(challenge types.BroadcastLogMessage) error {
	now := time.Now().UTC()
	const insertQuery = "INSERT INTO broadcast_challenge_messages(id, challenge_id, data, challenger, recipient, observers, created_at, updated_at) VALUES(NULL,?,?,?,?,?,?,?);"
	_, err := s.db.Exec(insertQuery, challenge.ChallengeID, challenge.Data, challenge.Challenger, challenge.Recipient, challenge.Observers, now, now)
	if err != nil {
		return err
	}

	return nil
}

// StorageChallengeMetrics retrieves all the metrics needs to be broadcast
func (s *SQLiteStore) StorageChallengeMetrics(timestamp time.Time) ([]types.StorageChallengeLogMessage, error) {
	const query = `
    SELECT id, challenge_id, message_type, data, sender_id, created_at, updated_at
    FROM storage_challenge_metrics
    WHERE created_at > ?
    `

	rows, err := s.db.Query(query, timestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []types.StorageChallengeLogMessage
	for rows.Next() {
		var m types.StorageChallengeLogMessage
		err := rows.Scan(&m.ID, &m.ChallengeID, &m.MessageType, &m.Data, &m.Sender, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}

	return metrics, rows.Err()
}

// QueryStorageChallengeMessage retrieves storage challenge message against challengeID and messageType
func (s *SQLiteStore) QueryStorageChallengeMessage(challengeID string, messageType int) (challengeMessage types.StorageChallengeLogMessage, err error) {
	const selectQuery = "SELECT * FROM storage_challenge_messages WHERE challenge_id=? AND message_type=?"
	err = s.db.QueryRow(selectQuery, challengeID, messageType).Scan(
		&challengeMessage.ID, &challengeMessage.ChallengeID, &challengeMessage.MessageType, &challengeMessage.Data,
		&challengeMessage.Sender, &challengeMessage.SenderSignature, &challengeMessage.CreatedAt, &challengeMessage.UpdatedAt)

	if err != nil {
		return challengeMessage, err
	}

	return challengeMessage, nil
}

// CleanupStorageChallenges cleans up challenges stored in DB for self-healing
func (s *SQLiteStore) CleanupStorageChallenges() (err error) {
	const delQuery = "DELETE FROM storage_challenge_messages"
	_, err = s.db.Exec(delQuery)
	return err
}

// GetStorageChallengeMetricsByChallengeID retrieves all the metrics
func (s *SQLiteStore) GetStorageChallengeMetricsByChallengeID(challengeID string) ([]types.StorageChallengeLogMessage, error) {
	const query = `
    SELECT id, challenge_id, message_type, data, sender_id, created_at, updated_at
    FROM storage_challenge_metrics
    WHERE challenge_id = ?;`

	rows, err := s.db.Query(query, challengeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []types.StorageChallengeLogMessage
	for rows.Next() {
		var m types.StorageChallengeLogMessage
		err := rows.Scan(&m.ID, &m.ChallengeID, &m.MessageType, &m.Data, &m.Sender, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}

	return metrics, rows.Err()
}

func processObserverEvaluations(observersEvaluations []types.StorageChallengeLogMessage) map[string]ObserverEvaluationMetrics {
	evaluationMap := make(map[string]ObserverEvaluationMetrics)

	for _, observerEvaluation := range observersEvaluations {
		var oe types.MessageData
		if err := json.Unmarshal(observerEvaluation.Data, &oe); err != nil {
			continue
		}

		oem, exists := evaluationMap[observerEvaluation.ChallengeID]
		if !exists {
			oem = ObserverEvaluationMetrics{} // Initialize if not exists
		}

		if isObserverEvaluationVerified(oe.ObserverEvaluation) {
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

func isObserverEvaluationVerified(observerEvaluation types.ObserverEvaluationData) bool {
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
