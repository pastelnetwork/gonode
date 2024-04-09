package queries

import (
	"github.com/pastelnetwork/gonode/common/types"
	"time"
)

type AggregateScScoreQueries interface {
	GetAccumulativeSCData(nodeID string) (types.AggregatedSCScore, error)
	GetAccumulativeHCData(nodeID string) (types.AggregatedSCScore, error)
	UpsertAccumulativeSCData(aggregatedScScore types.AggregatedSCScore) error
	UpsertAccumulativeHCData(aggregatedScScore types.AggregatedSCScore) error
}

func (s *SQLiteStore) GetAccumulativeSCData(nodeID string) (types.AggregatedSCScore, error) {
	const selectQuery = `
        SELECT node_id, ip_address, total_challenges_as_challengers, 
               total_challenges_as_recipients, total_challenges_as_observers,
               correct_challenger_evaluations, correct_recipient_evaluations, 
               correct_observer_evaluation, created_at, updated_at
        FROM accumulative_sc_data
        WHERE node_id = ?
    `

	var aggregatedScScore types.AggregatedSCScore
	err := s.db.QueryRow(selectQuery, nodeID).Scan(
		&aggregatedScScore.NodeID, &aggregatedScScore.IPAddress,
		&aggregatedScScore.TotalChallengesAsChallengers, &aggregatedScScore.TotalChallengesAsRecipients, &aggregatedScScore.TotalChallengesAsObservers,
		&aggregatedScScore.CorrectChallengerEvaluations, &aggregatedScScore.CorrectRecipientEvaluations, &aggregatedScScore.CorrectObserverEvaluations,
		&aggregatedScScore.CreatedAt, &aggregatedScScore.UpdatedAt,
	)
	if err != nil {
		return types.AggregatedSCScore{}, err
	}

	return aggregatedScScore, nil
}

func (s *SQLiteStore) UpsertAccumulativeSCData(aggregatedScScore types.AggregatedSCScore) error {
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

func (s *SQLiteStore) GetAccumulativeHCData(nodeID string) (types.AggregatedSCScore, error) {
	const selectQuery = `
        SELECT node_id, ip_address, total_challenges_as_challengers, 
               total_challenges_as_recipients, total_challenges_as_observers,
               correct_challenger_evaluations, correct_recipient_evaluations, 
               correct_observer_evaluation, created_at, updated_at
        FROM accumulative_hc_data
        WHERE node_id = ?
    `

	var aggregatedScScore types.AggregatedSCScore
	err := s.db.QueryRow(selectQuery, nodeID).Scan(
		&aggregatedScScore.NodeID, &aggregatedScScore.IPAddress,
		&aggregatedScScore.TotalChallengesAsChallengers, &aggregatedScScore.TotalChallengesAsRecipients, &aggregatedScScore.TotalChallengesAsObservers,
		&aggregatedScScore.CorrectChallengerEvaluations, &aggregatedScScore.CorrectRecipientEvaluations, &aggregatedScScore.CorrectObserverEvaluations,
		&aggregatedScScore.CreatedAt, &aggregatedScScore.UpdatedAt,
	)
	if err != nil {
		return types.AggregatedSCScore{}, err
	}

	return aggregatedScScore, nil
}

func (s *SQLiteStore) UpsertAccumulativeHCData(aggregatedScScore types.AggregatedSCScore) error {
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
