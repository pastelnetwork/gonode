package ticketstore

import (
	"github.com/pastelnetwork/gonode/common/types"
	"time"
)

type ActivationAttemptsQueries interface {
	InsertActivationAttempt(attempt types.ActivationAttempt) (int64, error)
	UpdateActivationAttempt(attempt types.ActivationAttempt) (int64, error)
	GetActivationAttemptByID(id int) (*types.ActivationAttempt, error)
	GetActivationAttemptsByFileIDAndBaseFileID(fileID, baseFileID string) ([]*types.ActivationAttempt, error)
}

// InsertActivationAttempt insert a new activation attempt into the activation_attempts table
func (s *TicketStore) InsertActivationAttempt(attempt types.ActivationAttempt) (int64, error) {
	const insertQuery = `
        INSERT INTO activation_attempts (
            file_id, base_file_id, activation_attempt_at, is_successful, error_message, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?)
        RETURNING id;`

	var id int64
	err := s.db.QueryRow(insertQuery,
		attempt.FileID, attempt.BaseFileID, attempt.ActivationAttemptAt,
		attempt.IsSuccessful, attempt.ErrorMessage, time.Now().UTC(), time.Now().UTC()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateActivationAttempt update a new activation attempt into the activation_attempts table
func (s *TicketStore) UpdateActivationAttempt(attempt types.ActivationAttempt) (int64, error) {
	const updateQuery = `
        UPDATE activation_attempts
        SET activation_attempt_at = ?, is_successful = ?, error_message = ?, updated_at = ?, is_confirmed=?
        WHERE id = ? AND file_id = ?
        RETURNING id`

	var id int64
	err := s.db.QueryRow(updateQuery,
		attempt.ActivationAttemptAt,
		attempt.IsSuccessful,
		attempt.ErrorMessage,
		time.Now().UTC(),
		attempt.IsConfirmed,
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

// GetActivationAttemptsByFileIDAndBaseFileID retrieves activation attempts by file_id from the activation_attempts table
func (s *TicketStore) GetActivationAttemptsByFileIDAndBaseFileID(fileID, baseFileID string) ([]*types.ActivationAttempt, error) {
	const selectQuery = `
        SELECT id, file_id, activation_attempt_at, is_successful, error_message
        FROM activation_attempts
        WHERE file_id = ? and base_file_id=?;`

	rows, err := s.db.Query(selectQuery, fileID, baseFileID)
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
