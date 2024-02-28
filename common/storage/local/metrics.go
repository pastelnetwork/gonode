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
			ch.IsAck = true
			challenges[row.ChallengeID] = ch

		} else if row.MessageType == int(types.SelfHealingCompletionMessage) {
			ch := challenges[row.ChallengeID]
			ch.IsHealed = true
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

func (s *SQLiteStore) GetTotalSCGeneratedAndProcessed(from time.Time) (metrics.SCMetrics, error) {
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
	scMetrics, err = s.GetTotalSCGeneratedAndProcessed(from)
	if err != nil {
		return scMetrics, err
	}

	challengerEvaluations, err := s.GetChallengerEvaluations(from)
	if err != nil {
		return scMetrics, err
	}

	for i := 0; i < len(challengerEvaluations); i++ {
		challengerEvaluation := challengerEvaluations[i]

		var challengeMsg types.Message
		if err := json.Unmarshal(challengerEvaluation.Data, &challengeMsg); err != nil {
			continue
		}

		if challengeMsg.Data.ChallengerEvaluation.IsVerified {
			scMetrics.TotalChallengesVerifiedByChallenger++
		}
	}

	observersEvaluations, err := s.GetObserversEvaluations(from)
	if err != nil {
		return scMetrics, err
	}

	observerEvaluationMetrics := processObserverEvaluations(observersEvaluations)

	for _, obMetrics := range observerEvaluationMetrics {
		if obMetrics.ChallengesVerified > 2 {
			scMetrics.TotalChallengesVerifiedByObservers++
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

func processObserverEvaluations(observersEvaluations []types.StorageChallengeLogMessage) map[string]ObserverEvaluationMetrics {
	evaluationMap := make(map[string]ObserverEvaluationMetrics)

	for i := 0; i < len(observersEvaluations); i++ {
		observerEvaluation := observersEvaluations[i]

		var oe types.Message
		if err := json.Unmarshal(observerEvaluation.Data, &oe); err != nil {
			continue
		}

		oem := evaluationMap[oe.ChallengeID]

		if isObserverEvaluationVerified(oe.Data.ObserverEvaluation) {
			oem.ChallengesVerified++
		} else {
			if !oe.Data.ObserverEvaluation.IsChallengeTimestampOK ||
				!oe.Data.ObserverEvaluation.IsProcessTimestampOK ||
				!oe.Data.ObserverEvaluation.IsEvaluationTimestampOK {
				oem.FailedByInvalidTimestamps++
			}

			if !oe.Data.ObserverEvaluation.IsChallengerSignatureOK ||
				!oe.Data.ObserverEvaluation.IsRecipientSignatureOK {
				oem.FailedByInvalidSignatures++
			}

			if !oe.Data.ObserverEvaluation.IsEvaluationResultOK {
				oem.FailedByInvalidEvaluation++
			}
		}
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

// GetLastNSHChallenges ...
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

// GetSHChallengeReport ...
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
