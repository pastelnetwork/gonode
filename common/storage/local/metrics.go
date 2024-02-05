package local

import (
	"context"
	"fmt"
	"time"

	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils/metrics"
)

var (
	oneYearAgo = time.Now().AddDate(-1, 0, 0)
)

// SHChallengeMetric represents the self-healing challenge metric
type SHChallengeMetric struct {
	ChallengeID string
	IsAccepted  bool
	IsVerified  bool
	IsHealed    bool
}

// InsertSelfHealingGenerationMetrics inserts self-healing generation metrics
func (s *SQLiteStore) InsertSelfHealingGenerationMetrics(metrics types.SelfHealingGenerationMetric) error {
	now := time.Now().UTC()
	const insertQuery = "INSERT INTO self_healing_generation_metrics(id, trigger_id, message_type, data, sender_id, sender_signature, created_at, updated_at) VALUES(NULL,?,?,?,?,?,?,?) ON CONFLICT DO NOTHING;"
	_, err := s.db.Exec(insertQuery, metrics.TriggerID, metrics.MessageType, metrics.Data, metrics.SenderID, metrics.SenderSignature, now, now)
	if err != nil {
		return err
	}

	return nil
}

// InsertSelfHealingExecutionMetrics inserts self-healing execution metrics
func (s *SQLiteStore) InsertSelfHealingExecutionMetrics(metrics types.SelfHealingExecutionMetric) error {
	now := time.Now().UTC()
	const insertQuery = "INSERT INTO self_healing_execution_metrics(id, trigger_id, challenge_id, message_type, data, sender_id, sender_signature, created_at, updated_at) VALUES(NULL,?,?,?,?,?,?,?,?) ON CONFLICT DO NOTHING;"

	_, err := s.db.Exec(insertQuery, metrics.TriggerID, metrics.ChallengeID, metrics.MessageType, metrics.Data, metrics.SenderID, metrics.SenderSignature, now, now)
	if err != nil {
		return err
	}

	return nil
}

// GetSelfHealingExecutionMetrics retrieves all self_healing_execution_metrics records created after the specified timestamp.
func (s *SQLiteStore) GetSelfHealingExecutionMetrics(timestamp time.Time) ([]types.SelfHealingExecutionMetric, error) {
	const query = `
    SELECT id, trigger_id, challenge_id, message_type, data, sender_id, sender_signature, created_at, updated_at
    FROM self_healing_execution_metrics
    WHERE created_at > ?
    `

	rows, err := s.db.Query(query, timestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []types.SelfHealingExecutionMetric
	for rows.Next() {
		var m types.SelfHealingExecutionMetric
		if err := rows.Scan(&m.ID, &m.TriggerID, &m.ChallengeID, &m.MessageType, &m.Data, &m.SenderID, &m.SenderSignature, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}

	return metrics, rows.Err()
}

// GetSelfHealingGenerationMetrics retrieves all self_healing_generation_metrics records created after the specified timestamp.
func (s *SQLiteStore) GetSelfHealingGenerationMetrics(timestamp time.Time) ([]types.SelfHealingGenerationMetric, error) {
	const query = `
    SELECT id, trigger_id, message_type, data, sender_id, sender_signature, created_at, updated_at
    FROM self_healing_generation_metrics
    WHERE created_at > ?
    `

	rows, err := s.db.Query(query, timestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []types.SelfHealingGenerationMetric
	for rows.Next() {
		var m types.SelfHealingGenerationMetric
		if err := rows.Scan(&m.ID, &m.TriggerID, &m.MessageType, &m.Data, &m.SenderID, &m.SenderSignature, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}

	return metrics, rows.Err()
}

// GetSHExecutionMetrics retrieves self-healing execution metrics
func (s *SQLiteStore) GetSHExecutionMetrics(ctx context.Context, from time.Time) (metrics.SHExecutionMetrics, error) {
	m := metrics.SHExecutionMetrics{}
	rows, err := s.GetSelfHealingExecutionMetrics(from)
	if err != nil {
		return m, err
	}
	log.WithContext(ctx).WithField("rows", len(rows)).Info("self-healing execution metrics row count")

	challenges := make(map[string]SHChallengeMetric)
	for _, row := range rows {
		if _, ok := challenges[row.ChallengeID]; !ok {
			challenges[row.ChallengeID] = SHChallengeMetric{
				ChallengeID: row.ChallengeID,
			}
		}

		if row.MessageType == int(types.SelfHealingVerificationMessage) {
			messages := types.SelfHealingMessages{}
			if err := json.Unmarshal(row.Data, &messages); err != nil {
				return m, fmt.Errorf("cannot unmarshal self healing execution message type 3: %w - row ID: %d", err, row.ID)
			}

			verificationCount := 0
			for _, message := range messages {
				if message.SelfHealingMessageData.Verification.VerifiedTicket.IsVerified {
					verificationCount++
				}
			}

			if verificationCount >= minVerifications {
				ch := challenges[row.ChallengeID]
				ch.IsVerified = true
				challenges[row.ChallengeID] = ch
			}
		} else {
			messages := types.SelfHealingMessages{}
			if err := json.Unmarshal(row.Data, &messages); err != nil {
				return m, fmt.Errorf("cannot unmarshal self healing execution message type 3: %w - row ID: %d", err, row.ID)
			}
			if len(messages) == 0 {
				return m, fmt.Errorf("len of selfhealing messages should not be 0 - problem with row ID %d", row.ID)
			}

			data := messages[0].SelfHealingMessageData

			if row.MessageType == int(types.SelfHealingResponseMessage) {
				ch := challenges[row.ChallengeID]
				ch.IsAccepted = data.Response.RespondedTicket.IsReconstructionRequired
				challenges[row.ChallengeID] = ch
			}

			if row.MessageType == 4 {
				ch := challenges[row.ChallengeID]
				ch.IsHealed = true
				challenges[row.ChallengeID] = ch
			}
		}
	}

	log.WithContext(ctx).WithField("challenges", len(challenges)).Info("self-healing execution metrics challenges count")

	for _, challenge := range challenges {
		log.WithContext(ctx).WithField("challenge-id", challenge.ChallengeID).WithField("is-accepted", challenge.IsAccepted).
			WithField("is-verified", challenge.IsVerified).WithField("is-healed", challenge.IsHealed).
			Info("self-healing challenge metric")

		if challenge.IsAccepted {
			m.TotalChallengesAccepted++
		}

		if challenge.IsVerified {
			m.TotalChallengesSuccessful++
		}

		if challenge.IsHealed {
			m.TotalFilesHealed++
		}
	}

	return m, nil
}

// QueryMetrics queries metrics
func (s *SQLiteStore) QueryMetrics(ctx context.Context, from time.Time, _ *time.Time) (m metrics.Metrics, err error) {
	genMetric, err := s.GetSelfHealingGenerationMetrics(from)
	if err != nil {
		return metrics.Metrics{}, err
	}

	te := metrics.SHTriggerMetrics{}
	challengesIssued := 0
	for _, metric := range genMetric {
		t := metrics.SHTriggerMetric{}
		data := types.SelfHealingMessageData{}
		json.Unmarshal(metric.Data, &data)

		t.TriggerID = metric.TriggerID
		t.ListOfNodes = data.Challenge.NodesOnWatchlist
		t.TotalTicketsIdentified = len(data.Challenge.ChallengeTickets)

		for _, ticket := range data.Challenge.ChallengeTickets {
			t.TotalFilesIdentified += len(ticket.MissingKeys)
		}

		challengesIssued += t.TotalTicketsIdentified

		te = append(te, t)
	}

	em, err := s.GetSHExecutionMetrics(ctx, from)
	if err != nil {
		return metrics.Metrics{}, fmt.Errorf("cannot get self healing execution metrics: %w", err)
	}

	em.TotalChallengesIssued = challengesIssued
	em.TotalChallengesRejected = challengesIssued - em.TotalChallengesAccepted
	em.TotalChallengesFailed = em.TotalChallengesAccepted - em.TotalChallengesSuccessful
	em.TotalFileHealingFailed = em.TotalChallengesSuccessful - em.TotalFilesHealed

	m.SHTriggerMetrics = te

	m.SHExecutionMetrics = em

	return m, nil
}

// GetLastNSHChallenges
func (s *SQLiteStore) GetLastNSHChallenges(ctx context.Context, n int) (types.SelfHealingChallengeReports, error) {
	challenges := types.SelfHealingChallengeReports{}
	rows, err := s.GetSelfHealingExecutionMetrics(oneYearAgo)
	if err != nil {
		return challenges, err
	}
	log.WithContext(ctx).WithField("rows", len(rows)).Info("self-healing execution metrics row count")

	challengesInserted := 0
	for _, row := range rows {
		if _, ok := challenges[row.ChallengeID]; !ok {
			if challengesInserted == n {
				continue
			}

			challenges[row.ChallengeID] = types.SelfHealingChallengeReport{}
			challengesInserted++
		}

		messages := types.SelfHealingMessages{}
		if err := json.Unmarshal(row.Data, &messages); err != nil {
			return challenges, fmt.Errorf("cannot unmarshal self healing execution message type 3: %w", err)
		}

		msgType := types.SelfHealingMessageType(row.MessageType)
		challenges[row.ChallengeID][msgType.String()] = messages
	}

	return challenges, nil
}

// GetSHChallengeReport ...
func (s *SQLiteStore) GetSHChallengeReport(ctx context.Context, challengeID string) (types.SelfHealingChallengeReports, error) {
	challenges := types.SelfHealingChallengeReports{}
	rows, err := s.GetSelfHealingExecutionMetrics(oneYearAgo)
	if err != nil {
		return challenges, err
	}
	log.WithContext(ctx).WithField("rows", len(rows)).Info("self-healing execution metrics row count")

	for _, row := range rows {
		if row.ChallengeID == challengeID {
			if _, ok := challenges[row.ChallengeID]; !ok {
				challenges[row.ChallengeID] = types.SelfHealingChallengeReport{}
			}

			messages := types.SelfHealingMessages{}
			if err := json.Unmarshal(row.Data, &messages); err != nil {
				return challenges, fmt.Errorf("cannot unmarshal self healing execution message type 3: %w", err)
			}

			msgType := types.SelfHealingMessageType(row.MessageType)
			challenges[row.ChallengeID][msgType.String()] = messages
		}
	}

	return challenges, nil
}
