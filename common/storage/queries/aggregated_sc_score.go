package queries

import (
	"database/sql"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/types"
	"time"
)

type AggregateScScoreQueries interface {
	GetAccumulativeSCDataForAllNodes() ([]types.AccumulativeChallengeData, error)
	GetAccumulativeHCDataForAllNodes() ([]types.AccumulativeChallengeData, error)
	GetAccumulativeSCData(nodeID string) (types.AccumulativeChallengeData, error)
	GetAccumulativeHCData(nodeID string) (types.AccumulativeChallengeData, error)
	UpsertAccumulativeSCData(aggregatedScScore types.AccumulativeChallengeData) error
	UpsertAccumulativeHCData(aggregatedScScore types.AccumulativeChallengeData) error
	UpsertAggregatedScore(score types.AggregatedScore) error
	GetAggregatedScore(nodeID string) (*types.AggregatedScore, error)
}

func (s *SQLiteStore) GetAccumulativeSCDataForAllNodes() ([]types.AccumulativeChallengeData, error) {
	rows, err := s.db.Query(`
        SELECT node_id, ip_address, total_challenges_as_challengers, 
               total_challenges_as_recipients, total_challenges_as_observers,
               correct_challenger_evaluations, correct_recipient_evaluations, 
               correct_observer_evaluations, created_at, updated_at
        FROM accumulative_sc_data
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []types.AccumulativeChallengeData
	for rows.Next() {
		var data types.AccumulativeChallengeData
		err = rows.Scan(&data.NodeID, &data.IPAddress, &data.TotalChallengesAsChallengers, &data.TotalChallengesAsRecipients,
			&data.TotalChallengesAsObservers, &data.CorrectChallengerEvaluations, &data.CorrectRecipientEvaluations,
			&data.CorrectObserverEvaluations, &data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			return nil, err
		}

		results = append(results, data)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *SQLiteStore) GetAccumulativeHCDataForAllNodes() ([]types.AccumulativeChallengeData, error) {
	rows, err := s.db.Query(`
        SELECT node_id, ip_address, total_challenges_as_challengers, 
               total_challenges_as_recipients, total_challenges_as_observers,
               correct_challenger_evaluations, correct_recipient_evaluations, 
               correct_observer_evaluations, created_at, updated_at
        FROM accumulative_hc_data
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []types.AccumulativeChallengeData
	for rows.Next() {
		var data types.AccumulativeChallengeData
		err = rows.Scan(&data.NodeID, &data.IPAddress, &data.TotalChallengesAsChallengers, &data.TotalChallengesAsRecipients,
			&data.TotalChallengesAsObservers, &data.CorrectChallengerEvaluations, &data.CorrectRecipientEvaluations,
			&data.CorrectObserverEvaluations, &data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			return nil, err
		}

		results = append(results, data)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *SQLiteStore) GetAccumulativeSCData(nodeID string) (types.AccumulativeChallengeData, error) {
	const selectQuery = `
        SELECT node_id, ip_address, total_challenges_as_challengers, 
               total_challenges_as_recipients, total_challenges_as_observers,
               correct_challenger_evaluations, correct_recipient_evaluations, 
               correct_observer_evaluation, created_at, updated_at
        FROM accumulative_sc_data
        WHERE node_id = ?
    `

	var aggregatedScScore types.AccumulativeChallengeData
	err := s.db.QueryRow(selectQuery, nodeID).Scan(
		&aggregatedScScore.NodeID, &aggregatedScScore.IPAddress,
		&aggregatedScScore.TotalChallengesAsChallengers, &aggregatedScScore.TotalChallengesAsRecipients, &aggregatedScScore.TotalChallengesAsObservers,
		&aggregatedScScore.CorrectChallengerEvaluations, &aggregatedScScore.CorrectRecipientEvaluations, &aggregatedScScore.CorrectObserverEvaluations,
		&aggregatedScScore.CreatedAt, &aggregatedScScore.UpdatedAt,
	)
	if err != nil {
		return types.AccumulativeChallengeData{}, err
	}

	return aggregatedScScore, nil
}

func (s *SQLiteStore) UpsertAccumulativeSCData(aggregatedScScore types.AccumulativeChallengeData) error {
	const upsertQuery = `
        INSERT INTO accumulative_sc_data (
            node_id, ip_address,
            total_challenges_as_challengers, total_challenges_as_recipients, total_challenges_as_observers,
            correct_challenger_evaluations, correct_recipient_evaluations, correct_observer_evaluation,
            created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(node_id, ip_address)
        DO UPDATE SET
            total_challenges_as_challengers = excluded.total_challenges_as_challengers,
            total_challenges_as_recipients = excluded.total_challenges_as_recipients,
            total_challenges_as_observers = excluded.total_challenges_as_observers,
            correct_challenger_evaluations = excluded.correct_challenger_evaluations,
            correct_recipient_evaluations = excluded.correct_recipient_evaluations,
            correct_observer_evaluation = excluded.correct_observer_evaluation,
            updated_at = excluded.updated_at;
    `

	_, err := s.db.Exec(
		upsertQuery,
		aggregatedScScore.NodeID, aggregatedScScore.IPAddress,
		aggregatedScScore.TotalChallengesAsChallengers, aggregatedScScore.TotalChallengesAsRecipients, aggregatedScScore.TotalChallengesAsObservers,
		aggregatedScScore.CorrectChallengerEvaluations, aggregatedScScore.CorrectRecipientEvaluations, aggregatedScScore.CorrectObserverEvaluations,
		time.Now().UTC(), time.Now().UTC(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStore) GetAccumulativeHCData(nodeID string) (types.AccumulativeChallengeData, error) {
	const selectQuery = `
        SELECT node_id, ip_address, total_challenges_as_challengers, 
               total_challenges_as_recipients, total_challenges_as_observers,
               correct_challenger_evaluations, correct_recipient_evaluations, 
               correct_observer_evaluation, created_at, updated_at
        FROM accumulative_hc_data
        WHERE node_id = ?
    `

	var aggregatedScScore types.AccumulativeChallengeData
	err := s.db.QueryRow(selectQuery, nodeID).Scan(
		&aggregatedScScore.NodeID, &aggregatedScScore.IPAddress,
		&aggregatedScScore.TotalChallengesAsChallengers, &aggregatedScScore.TotalChallengesAsRecipients, &aggregatedScScore.TotalChallengesAsObservers,
		&aggregatedScScore.CorrectChallengerEvaluations, &aggregatedScScore.CorrectRecipientEvaluations, &aggregatedScScore.CorrectObserverEvaluations,
		&aggregatedScScore.CreatedAt, &aggregatedScScore.UpdatedAt,
	)
	if err != nil {
		return types.AccumulativeChallengeData{}, err
	}

	return aggregatedScScore, nil
}

func (s *SQLiteStore) UpsertAccumulativeHCData(aggregatedScScore types.AccumulativeChallengeData) error {
	const upsertQuery = `
        INSERT INTO accumulative_hc_data (
            node_id, ip_address,
            total_challenges_as_challengers, total_challenges_as_recipients, total_challenges_as_observers,
            correct_challenger_evaluations, correct_recipient_evaluations, correct_observer_evaluation,
            created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(node_id, ip_address)
        DO UPDATE SET
            total_challenges_as_challengers = excluded.total_challenges_as_challengers,
            total_challenges_as_recipients = excluded.total_challenges_as_recipients,
            total_challenges_as_observers = excluded.total_challenges_as_observers,
            correct_challenger_evaluations = excluded.correct_challenger_evaluations,
            correct_recipient_evaluations = excluded.correct_recipient_evaluations,
            correct_observer_evaluation = excluded.correct_observer_evaluation,
            updated_at = excluded.updated_at;
    `

	_, err := s.db.Exec(
		upsertQuery,
		aggregatedScScore.NodeID, aggregatedScScore.IPAddress,
		aggregatedScScore.TotalChallengesAsChallengers, aggregatedScScore.TotalChallengesAsRecipients, aggregatedScScore.TotalChallengesAsObservers,
		aggregatedScScore.CorrectChallengerEvaluations, aggregatedScScore.CorrectRecipientEvaluations, aggregatedScScore.CorrectObserverEvaluations,
		time.Now().UTC(), time.Now().UTC(),
	)
	if err != nil {
		return err
	}

	return nil
}

// UpsertAggregatedScore inserts a new record or updates an existing record based on the node_id
func (s *SQLiteStore) UpsertAggregatedScore(score types.AggregatedScore) error {
	query := `
	INSERT INTO aggregated_challenge_scores(node_id, ip_address, storage_challenge_score, healthcheck_challenge_score, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT(node_id)
	DO UPDATE SET 
		storage_challenge_score = excluded.storage_challenge_score,
		healthcheck_challenge_score = excluded.healthcheck_challenge_score,
		updated_at = excluded.updated_at;
	`

	_, err := s.db.Exec(query, score.NodeID, score.IPAddress, score.StorageChallengeScore, score.HealthCheckChallengeScore, score.CreatedAt, score.UpdatedAt)
	return err
}

// GetAggregatedScore retrieves the aggregated score for a given node ID
func (s *SQLiteStore) GetAggregatedScore(nodeID string) (*types.AggregatedScore, error) {
	query := `
	SELECT node_id, ip_address, storage_challenge_score, healthcheck_challenge_score, created_at, updated_at
	FROM aggregated_challenge_scores
	WHERE node_id = ?
	`

	var score types.AggregatedScore
	row := s.db.QueryRow(query, nodeID)
	err := row.Scan(&score.NodeID, &score.IPAddress, &score.StorageChallengeScore, &score.HealthCheckChallengeScore, &score.CreatedAt, &score.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &score, nil
}
