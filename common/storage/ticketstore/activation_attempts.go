package ticketstore

import (
	"github.com/pastelnetwork/gonode/common/types"
)

type ActivationAttemptsQueries interface {
	InsertActivationAttempt(attempt types.ActivationAttempt) (int64, error)
	UpdateActivationAttempt(attempt types.ActivationAttempt) (int64, error)
	GetActivationAttemptByID(id int) (*types.ActivationAttempt, error)
	GetActivationAttemptsByFileID(fileID string) ([]*types.ActivationAttempt, error)
}

// InsertActivationAttempt insert a new activation attempt into the activation_attempts table
func (s *TicketStore) InsertActivationAttempt(attempt types.ActivationAttempt) (int64, error) {
	const insertQuery = `
        INSERT INTO activation_attempts (
            file_id, activation_attempt_at, is_successful, error_message
        ) VALUES (?, ?, ?, ?)
        RETURNING id;`

	var id int64
	err := s.db.QueryRow(insertQuery,
		attempt.FileID, attempt.ActivationAttemptAt,
		attempt.IsSuccessful, attempt.ErrorMessage).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateActivationAttempt update a new activation attempt into the activation_attempts table
func (s *TicketStore) UpdateActivationAttempt(attempt types.ActivationAttempt) (int64, error) {
	const updateQuery = `
        UPDATE activation_attempts
        SET activation_attempt_at = ?, is_successful = ?, error_message = ?
        WHERE id = ? AND file_id = ?
        RETURNING id`

	var id int64
	err := s.db.QueryRow(updateQuery,
		attempt.ActivationAttemptAt,
		attempt.IsSuccessful,
		attempt.ErrorMessage,
		attempt.ID,
		attempt.FileID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetActivationAttemptByID retrieves an activation attempt by its ID from the activation_attempts table
func (s *TicketStore) GetActivationAttemptByID(id int) (*types.ActivationAttempt, error) {
	const selectQuery = `
        SELECT id, file_id, activation_attempt_at, is_successful, error_message
        FROM activation_attempts
        WHERE id = ?;`

	row := s.db.QueryRow(selectQuery, int64(id))

	var attempt types.ActivationAttempt
	err := row.Scan(
		&attempt.ID, &attempt.FileID, &attempt.ActivationAttemptAt,
		&attempt.IsSuccessful, &attempt.ErrorMessage)
	if err != nil {
		return nil, err
	}

	return &attempt, nil
}

// GetActivationAttemptsByFileID retrieves activation attempts by file_id from the activation_attempts table
func (s *TicketStore) GetActivationAttemptsByFileID(fileID string) ([]*types.ActivationAttempt, error) {
	const selectQuery = `
        SELECT id, file_id, activation_attempt_at, is_successful, error_message
        FROM activation_attempts
        WHERE file_id = ?;`

	rows, err := s.db.Query(selectQuery, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*types.ActivationAttempt
	for rows.Next() {
		var attempt types.ActivationAttempt
		err := rows.Scan(
			&attempt.ID, &attempt.FileID, &attempt.ActivationAttemptAt,
			&attempt.IsSuccessful, &attempt.ErrorMessage)
		if err != nil {
			return nil, err
		}
		attempts = append(attempts, &attempt)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return attempts, nil
}
