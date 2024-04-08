package queries

import (
	"time"

	"github.com/pastelnetwork/gonode/common/types"
)

type AggregateScScoreQueries interface {
	GetAggregatedSCScoreByNodeID(nodeID string) (types.AggregatedSCScore, error)
	UpsertAggregatedScScore(aggregatedScScore types.AggregatedSCScore) error
}

func (s *SQLiteStore) GetAggregatedSCScoreByNodeID(nodeID string) (types.AggregatedSCScore, error) {
	const selectQuery = `
        SELECT node_id, ip_address, total_challenges_as_challengers, 
               total_challenges_as_recipients, total_challenges_as_observers,
               correct_challenger_evaluations, correct_recipient_evaluations, 
               correct_observer_evaluation, created_at, updated_at
        FROM aggregated_sc_scores
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
func (s *SQLiteStore) UpsertAggregatedScScore(aggregatedScScore types.AggregatedSCScore) error {
	const upsertQuery = `
        INSERT INTO aggregated_sc_scores (
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
