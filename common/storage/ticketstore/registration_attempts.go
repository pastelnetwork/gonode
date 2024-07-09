package ticketstore

import (
	"github.com/pastelnetwork/gonode/common/types"
	"time"
)

type RegistrationAttemptsQueries interface {
	InsertRegistrationAttempt(attempt types.RegistrationAttempt) (int64, error)
	UpdateRegistrationAttempt(attempt types.RegistrationAttempt) (int64, error)
	GetRegistrationAttemptByID(id int) (*types.RegistrationAttempt, error)
	GetRegistrationAttemptsByFileIDAndBaseFileID(fileID, baseFileID string) ([]*types.RegistrationAttempt, error)
}

// InsertRegistrationAttempt insert a new registration attempt into the registration_attempts table
func (s *TicketStore) InsertRegistrationAttempt(attempt types.RegistrationAttempt) (int64, error) {
	const insertQuery = `
        INSERT INTO registration_attempts (
            file_id, base_file_id, reg_started_at, processor_sns, finished_at, is_successful, error_message, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        RETURNING id;`

	var id int64
	err := s.db.QueryRow(insertQuery,
		attempt.FileID, attempt.BaseFileID, attempt.RegStartedAt, attempt.ProcessorSNS,
		attempt.FinishedAt, attempt.IsSuccessful, attempt.ErrorMessage, time.Now().UTC(), time.Now().UTC()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateRegistrationAttempt update a new registration attempt into the registration_attempts table
func (s *TicketStore) UpdateRegistrationAttempt(attempt types.RegistrationAttempt) (int64, error) {
	const updateQuery = `
        UPDATE registration_attempts
        SET reg_started_at = ?,
            processor_sns = ?,
            finished_at = ?,
            is_successful = ?,
            error_message = ?,
            is_confirmed=?,
        	updated_at = ?
        WHERE id = ? AND file_id = ?
        RETURNING id;`

	var id int64
	err := s.db.QueryRow(updateQuery,
		attempt.RegStartedAt, attempt.ProcessorSNS, attempt.FinishedAt,
		attempt.IsSuccessful, attempt.ErrorMessage, attempt.IsConfirmed, time.Now().UTC(),
		attempt.ID, attempt.FileID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetRegistrationAttemptByID retrieves a registration attempt by its ID from the registration_attempts table
func (s *TicketStore) GetRegistrationAttemptByID(id int) (*types.RegistrationAttempt, error) {
	const selectQuery = `
        SELECT id, file_id, reg_started_at, processor_sns, finished_at, 
               is_successful, error_message
        FROM registration_attempts
        WHERE id = ?;`

	row := s.db.QueryRow(selectQuery, int64(id))

	var attempt types.RegistrationAttempt
	err := row.Scan(
		&attempt.ID, &attempt.FileID, &attempt.RegStartedAt, &attempt.ProcessorSNS,
		&attempt.FinishedAt, &attempt.IsSuccessful, &attempt.ErrorMessage)
	if err != nil {
		return nil, err
	}

	return &attempt, nil
}

// GetRegistrationAttemptsByFileIDAndBaseFileID retrieves registration attempts by file_id and baseFileID from the registration_attempts table
func (s *TicketStore) GetRegistrationAttemptsByFileIDAndBaseFileID(fileID, baseFileID string) ([]*types.RegistrationAttempt, error) {
	const selectQuery = `
        SELECT id, file_id, base_file_id, reg_started_at, processor_sns, finished_at, 
               is_successful, error_message
        FROM registration_attempts
        WHERE file_id = ? and base_file_id=?;`

	rows, err := s.db.Query(selectQuery, fileID, baseFileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*types.RegistrationAttempt
	for rows.Next() {
		var attempt types.RegistrationAttempt
		err := rows.Scan(
			&attempt.ID, &attempt.FileID, &attempt.BaseFileID, &attempt.RegStartedAt, &attempt.ProcessorSNS,
			&attempt.FinishedAt, &attempt.IsSuccessful, &attempt.ErrorMessage)
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
