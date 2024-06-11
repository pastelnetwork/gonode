package ticketstore

import (
	"github.com/pastelnetwork/gonode/common/types"
)

type ActivationAttemptsQueries interface {
	UpsertActivationAttempt(attempt types.ActivationAttempt) error
	GetActivationAttemptByID(id string) (*types.ActivationAttempt, error)
}

// UpsertActivationAttempt upsert a new activation attempt into the activation_attempts table
func (s *TicketStore) UpsertActivationAttempt(attempt types.ActivationAttempt) error {
	const upsertQuery = `
        INSERT INTO activation_attempts (
            id, file_id, activation_attempt_at, is_successful, error_message
        ) VALUES (?, ?, ?, ?, ?)
        ON CONFLICT(id) 
        DO UPDATE SET
            file_id = excluded.file_id,
            activation_attempt_at = excluded.activation_attempt_at,
            is_successful = excluded.is_successful,
            error_message = excluded.error_message;`

	_, err := s.db.Exec(upsertQuery,
		attempt.ID, attempt.FileID, attempt.ActivationAttemptAt,
		attempt.IsSuccessful, attempt.ErrorMessage)
	if err != nil {
		return err
	}

	return nil
}

// GetActivationAttemptByID retrieves an activation attempt by its ID from the activation_attempts table
func (s *TicketStore) GetActivationAttemptByID(id string) (*types.ActivationAttempt, error) {
	const selectQuery = `
        SELECT id, file_id, activation_attempt_at, is_successful, error_message
        FROM activation_attempts
        WHERE id = ?;`

	row := s.db.QueryRow(selectQuery, id)

	var attempt types.ActivationAttempt
	err := row.Scan(
		&attempt.ID, &attempt.FileID, &attempt.ActivationAttemptAt,
		&attempt.IsSuccessful, &attempt.ErrorMessage)
	if err != nil {
		return nil, err
	}

	return &attempt, nil
}
