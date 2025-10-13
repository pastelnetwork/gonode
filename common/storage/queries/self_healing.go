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

type SelfHealingQueries interface {
	BatchInsertSelfHealingChallengeEvents(ctx context.Context, event []types.SelfHealingChallengeEvent) error
	UpdateSHChallengeEventProcessed(challengeID string, isProcessed bool) error
	GetSelfHealingChallengeEvents() ([]types.SelfHealingChallengeEvent, error)
	CleanupSelfHealingChallenges() (err error)
	QuerySelfHealingChallenges() (challenges []types.SelfHealingChallenge, err error)

	QueryMetrics(ctx context.Context, from time.Time, to *time.Time) (m metrics.Metrics, err error)
	InsertSelfHealingGenerationMetrics(metrics types.SelfHealingGenerationMetric) error
	InsertSelfHealingExecutionMetrics(metrics types.SelfHealingExecutionMetric) error
	BatchInsertExecutionMetrics(metrics []types.SelfHealingExecutionMetric) error
	GetSelfHealingGenerationMetrics(timestamp time.Time) ([]types.SelfHealingGenerationMetric, error)
	GetSelfHealingExecutionMetrics(timestamp time.Time) ([]types.SelfHealingExecutionMetric, error)
	GetLastNSHChallenges(ctx context.Context, n int) (types.SelfHealingReports, error)
	GetSHChallengeReport(ctx context.Context, challengeID string) (types.SelfHealingReports, error)
	GetSHExecutionMetrics(ctx context.Context, from time.Time) (metrics.SHExecutionMetrics, error)
	RemoveSelfHealingStaleData(ctx context.Context, threshold string) error
}

var (
	oneYearAgo = time.Now().AddDate(-1, 0, 0)
)

// SHChallengeMetric represents the self-healing challenge metric
type SHChallengeMetric struct {
	ChallengeID string

	// healer node
	IsAck      bool
	IsAccepted bool
	IsRejected bool

	// verifier nodes
	HasMinVerifications                    bool
	IsVerified                             bool
	IsReconstructionRequiredVerified       bool
	IsReconstructionNotRequiredVerified    bool
	IsUnverified                           bool
	IsReconstructionRequiredNotVerified    bool
	IsReconstructionNotRequiredNotVerified bool
	IsReconstructionRequiredHashMismatch   bool

	IsHealed bool
}

type HCObserverEvaluationMetrics struct {
	ChallengesVerified        int
	FailedByInvalidTimestamps int
	FailedByInvalidSignatures int
	FailedByInvalidEvaluation int
}

type ObserverEvaluationMetrics struct {
	ChallengesVerified        int
	FailedByInvalidTimestamps int
	FailedByInvalidSignatures int
	FailedByInvalidEvaluation int
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

// BatchInsertExecutionMetrics inserts execution metrics in a batch
func (s *SQLiteStore) BatchInsertExecutionMetrics(metrics []types.SelfHealingExecutionMetric) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
        INSERT OR IGNORE INTO self_healing_execution_metrics
        (id, trigger_id, challenge_id, message_type, data, sender_id, sender_signature, created_at, updated_at) 
        VALUES (NULL,?,?,?,?,?,?,?,?)
    `)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, metric := range metrics {
		now := time.Now().UTC()

		_, err = stmt.Exec(metric.TriggerID, metric.ChallengeID, metric.MessageType, metric.Data, metric.SenderID, metric.SenderSignature, now, now)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
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

// GetLastNSCMetrics gets the N number of latest challenge IDs from the DB
func (s *SQLiteStore) GetLastNSCMetrics() ([]types.NScMetric, error) {
	const query = `
SELECT 
    count(*) AS count, 
    challenge_id, 
    MAX(created_at) AS most_recent
FROM 
    storage_challenge_metrics 
GROUP BY 
    challenge_id
HAVING 
    count(*) > 5
ORDER BY 
    most_recent DESC
LIMIT 20;`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []types.NScMetric
	for rows.Next() {
		var m types.NScMetric
		err := rows.Scan(&m.Count, &m.ChallengeID, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}

	return metrics, rows.Err()
}

// GetLastNHCMetrics gets the N number of latest health-check challenge IDs from the DB
func (s *SQLiteStore) GetLastNHCMetrics() ([]types.NHcMetric, error) {
	const query = `
SELECT 
    count(*) AS count, 
    challenge_id, 
    MAX(created_at) AS most_recent
FROM 
    healthcheck_challenge_metrics 
GROUP BY 
    challenge_id
HAVING 
    count(*) > 5
ORDER BY 
    most_recent DESC
LIMIT 20;`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []types.NHcMetric
	for rows.Next() {
		var m types.NHcMetric
		err := rows.Scan(&m.Count, &m.ChallengeID, &m.CreatedAt)
		if err != nil {
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

			if len(messages) >= minVerifications {
				ch := challenges[row.ChallengeID]
				ch.HasMinVerifications = true
				challenges[row.ChallengeID] = ch
			}

			reconReqVerified := 0
			reconNotReqVerified := 0
			reconReqUnverified := 0
			reconNotReqUnverified := 0
			reconReqHashMismatch := 0

			for _, message := range messages {
				if message.SelfHealingMessageData.Verification.VerifiedTicket.IsReconstructionRequired {
					if message.SelfHealingMessageData.Verification.VerifiedTicket.IsReconstructionRequiredByHealer {
						if message.SelfHealingMessageData.Verification.VerifiedTicket.IsVerified {
							reconReqVerified++
						} else {
							reconReqHashMismatch++
						}
					} else {
						reconNotReqUnverified++
					}
				} else {
					if message.SelfHealingMessageData.Verification.VerifiedTicket.IsReconstructionRequiredByHealer {
						reconReqUnverified++
					} else {
						reconNotReqVerified++
					}
				}
			}

			if reconReqVerified >= minVerifications {
				ch := challenges[row.ChallengeID]
				ch.IsVerified = true
				ch.IsReconstructionRequiredVerified = true
				challenges[row.ChallengeID] = ch
			} else if reconNotReqVerified >= minVerifications {
				ch := challenges[row.ChallengeID]
				ch.IsVerified = true
				ch.IsReconstructionNotRequiredVerified = true
				challenges[row.ChallengeID] = ch
			} else if reconReqUnverified >= minVerifications {
				ch := challenges[row.ChallengeID]
				ch.IsUnverified = true
				ch.IsReconstructionRequiredNotVerified = true
				challenges[row.ChallengeID] = ch
			} else if reconNotReqUnverified >= minVerifications {
				ch := challenges[row.ChallengeID]
				ch.IsUnverified = true
				ch.IsReconstructionNotRequiredNotVerified = true
				challenges[row.ChallengeID] = ch
			} else if reconReqHashMismatch >= minVerifications {
				ch := challenges[row.ChallengeID]
				ch.IsReconstructionRequiredHashMismatch = true
				challenges[row.ChallengeID] = ch
			}

		} else if row.MessageType == int(types.SelfHealingResponseMessage) {
			messages := types.SelfHealingMessages{}
			if err := json.Unmarshal(row.Data, &messages); err != nil {
				return m, fmt.Errorf("cannot unmarshal self healing execution message type 3: %w - row ID: %d", err, row.ID)
			}
			if len(messages) == 0 {
				return m, fmt.Errorf("len of selfhealing messages should not be 0 - problem with row ID %d", row.ID)
			}

			data := messages[0].SelfHealingMessageData

			ch := challenges[row.ChallengeID]
			if data.Response.RespondedTicket.IsReconstructionRequired {
				ch.IsAccepted = true
			} else {
				ch.IsRejected = true
			}
			challenges[row.ChallengeID] = ch

		} else if row.MessageType == int(types.SelfHealingCompletionMessage) {
			ch := challenges[row.ChallengeID]
			ch.IsHealed = true
			challenges[row.ChallengeID] = ch
		} else if row.MessageType == int(types.SelfHealingAcknowledgementMessage) {
			ch := challenges[row.ChallengeID]
			ch.IsAck = true
			challenges[row.ChallengeID] = ch
		}
	}

	log.WithContext(ctx).WithField("challenges", len(challenges)).Info("self-healing execution metrics challenges count")

	for _, challenge := range challenges {
		log.WithContext(ctx).WithField("challenge-id", challenge.ChallengeID).WithField("is-accepted", challenge.IsAccepted).
			WithField("is-verified", challenge.IsVerified).WithField("is-healed", challenge.IsHealed).
			Info("self-healing challenge metric")

		if challenge.IsAck {
			m.TotalChallengesAcknowledged++
		}

		if challenge.IsAccepted {
			m.TotalChallengesAccepted++
		}

		if challenge.IsRejected {
			m.TotalChallengesRejected++
		}

		if challenge.IsVerified {
			m.TotalChallengeEvaluationsVerified++
		}

		if challenge.IsReconstructionRequiredVerified {
			m.TotalReconstructionsApproved++
		}

		if challenge.IsReconstructionNotRequiredVerified {
			m.TotalReconstructionsNotRquiredApproved++
		}

		if challenge.IsUnverified {
			m.TotalChallengeEvaluationsUnverified++
		}

		if challenge.IsReconstructionRequiredNotVerified {
			m.TotalReconstructionsNotApproved++
		}

		if challenge.IsReconstructionNotRequiredNotVerified {
			m.TotalReconstructionsNotRequiredEvaluationNotApproved++
		}

		if challenge.IsReconstructionRequiredHashMismatch {
			m.TotalReconstructionRequiredHashMismatch++
		}

		if challenge.IsHealed {
			m.TotalFilesHealed++
		}
	}

	return m, nil
}

// QueryMetrics queries the self-healing metrics
func (s *SQLiteStore) QueryMetrics(ctx context.Context, from time.Time, _ *time.Time) (m metrics.Metrics, err error) {
	genMetric, err := s.GetSelfHealingGenerationMetrics(from)
	if err != nil {
		return metrics.Metrics{}, err
	}

	te := metrics.SHTriggerMetrics{}
	challengesIssued := 0
	for _, metric := range genMetric {
		t := metrics.SHTriggerMetric{}
		data := types.SelfHealingMessages{}
		if err := json.Unmarshal(metric.Data, &data); err != nil {
			return metrics.Metrics{}, fmt.Errorf("cannot unmarshal self healing generation message type 3: %w", err)
		}

		if len(data) < 1 {
			return metrics.Metrics{}, fmt.Errorf("len of selfhealing messages data JSON should not be 0")
		}

		t.TriggerID = metric.TriggerID
		t.ListOfNodes = data[0].SelfHealingMessageData.Challenge.NodesOnWatchlist
		t.TotalTicketsIdentified = len(data[0].SelfHealingMessageData.Challenge.ChallengeTickets)

		for _, ticket := range data[0].SelfHealingMessageData.Challenge.ChallengeTickets {
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
	em.TotalFileHealingFailed = em.TotalReconstructionsApproved - em.TotalFilesHealed

	m.SHTriggerMetrics = te

	m.SHExecutionMetrics = em

	return m, nil
}

// GetLastNSHChallenges retrieves the latest 'N' self-healing challenges
func (s *SQLiteStore) GetLastNSHChallenges(ctx context.Context, n int) (types.SelfHealingReports, error) {
	challenges := types.SelfHealingReports{}
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

			challenges[row.ChallengeID] = types.SelfHealingReport{}
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

// GetSHChallengeReport returns the self-healing report
func (s *SQLiteStore) GetSHChallengeReport(ctx context.Context, challengeID string) (types.SelfHealingReports, error) {
	challenges := types.SelfHealingReports{}
	rows, err := s.GetSelfHealingExecutionMetrics(oneYearAgo)
	if err != nil {
		return challenges, err
	}
	log.WithContext(ctx).WithField("rows", len(rows)).Info("self-healing execution metrics row count")

	for _, row := range rows {
		if row.ChallengeID == challengeID {
			if _, ok := challenges[row.ChallengeID]; !ok {
				challenges[row.ChallengeID] = types.SelfHealingReport{}
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

// QuerySelfHealingChallenges retrieves self-healing audit logs stored in DB for self-healing
func (s *SQLiteStore) QuerySelfHealingChallenges() (challenges []types.SelfHealingChallenge, err error) {
	const selectQuery = "SELECT * FROM self_healing_challenges"
	rows, err := s.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		challenge := types.SelfHealingChallenge{}
		err = rows.Scan(&challenge.ID, &challenge.ChallengeID, &challenge.MerkleRoot, &challenge.FileHash,
			&challenge.ChallengingNode, &challenge.RespondingNode, &challenge.VerifyingNode, &challenge.ReconstructedFileHash,
			&challenge.Status, &challenge.CreatedAt, &challenge.UpdatedAt)
		if err != nil {
			return nil, err
		}

		challenges = append(challenges, challenge)
	}

	return challenges, nil
}

// BatchInsertSelfHealingChallengeEvents inserts self-healing-challenge events in a batch
func (s *SQLiteStore) BatchInsertSelfHealingChallengeEvents(ctx context.Context, eventsBatch []types.SelfHealingChallengeEvent) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
        INSERT OR IGNORE INTO self_healing_challenge_events
        (trigger_id, ticket_id, challenge_id, data, sender_id, is_processed, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	stmt2, err := tx.Prepare(`
		INSERT OR IGNORE INTO self_healing_execution_metrics(id, trigger_id, challenge_id, message_type, data, sender_id, sender_signature, created_at, updated_at) 
		       VALUES(NULL,?,?,?,?,?,?,?,?);
    `)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt2.Close()

	for _, event := range eventsBatch {
		now := time.Now().UTC()

		_, err = stmt.Exec(event.TriggerID, event.TicketID, event.ChallengeID, event.Data, event.SenderID, false, now, now)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = stmt2.Exec(event.ExecMetric.TriggerID, event.ExecMetric.ChallengeID, event.ExecMetric.MessageType, event.ExecMetric.Data, event.ExecMetric.SenderID, event.ExecMetric.SenderSignature, now, now)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// GetSelfHealingChallengeEvents retrieves the challenge events from DB
func (s *SQLiteStore) GetSelfHealingChallengeEvents() ([]types.SelfHealingChallengeEvent, error) {
	const selectQuery = `
        SELECT trigger_id, ticket_id, challenge_id, data, sender_id, is_processed, created_at, updated_at
        FROM self_healing_challenge_events
        WHERE is_processed = false
    `
	rows, err := s.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []types.SelfHealingChallengeEvent

	for rows.Next() {
		var event types.SelfHealingChallengeEvent
		if err := rows.Scan(
			&event.TriggerID, &event.TicketID, &event.ChallengeID, &event.Data, &event.SenderID, &event.IsProcessed,
			&event.CreatedAt, &event.UpdatedAt,
		); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

// UpdateSHChallengeEventProcessed updates the is_processed flag of an event
func (s *SQLiteStore) UpdateSHChallengeEventProcessed(challengeID string, isProcessed bool) error {
	const updateQuery = `
        UPDATE self_healing_challenge_events
        SET is_processed = ?
        WHERE challenge_id = ?
    `
	_, err := s.db.Exec(updateQuery, isProcessed, challengeID)
	return err
}

// CleanupSelfHealingChallenges cleans up self-healing challenges stored in DB for inspection
func (s *SQLiteStore) CleanupSelfHealingChallenges() (err error) {
	const delQuery = "DELETE FROM self_healing_challenges"
	_, err = s.db.Exec(delQuery)
	return err
}

func (s *SQLiteStore) RemoveSelfHealingStaleData(ctx context.Context, threshold string) error {
	queries := []string{
		"DELETE FROM self_healing_execution_metrics WHERE created_at < $1",
		"DELETE FROM self_healing_generation_metrics WHERE created_at < $1",
		"DELETE from self_healing_challenge_events where is_processed = true and created_at < $1",
	}

	for _, query := range queries {
		if _, err := s.db.ExecContext(ctx, query, threshold); err != nil {
			return fmt.Errorf("failed to delete old metrics: %v", err)
		}
	}

	fmt.Println("Old metrics deleted successfully.")
	return nil
}
