package ticketstore

import (
	"github.com/pastelnetwork/gonode/common/types"
)

type RegistrationAttemptsQueries interface {
	UpsertRegistrationAttempt(attempt types.RegistrationAttempt) error
	GetRegistrationAttemptByID(id int) (*types.RegistrationAttempt, error)
	GetRegistrationAttemptsByFileID(fileID string) ([]*types.RegistrationAttempt, error)
}

// UpsertRegistrationAttempt upsert a new registration attempt into the registration_attempts table
func (s *TicketStore) UpsertRegistrationAttempt(attempt types.RegistrationAttempt) error {
	const upsertQuery = `
        INSERT INTO registration_attempts (
            id, file_id, reg_started_at, processor_sns, finished_at, is_successful, error_message
        ) VALUES (?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) 
        DO UPDATE SET
            file_id = excluded.file_id,
            reg_started_at = excluded.reg_started_at,
            processor_sns = excluded.processor_sns,
            finished_at = excluded.finished_at,
            is_successful = excluded.is_successful,
            error_message = excluded.error_message;`

	_, err := s.db.Exec(upsertQuery,
		attempt.ID, attempt.FileID, attempt.RegStartedAt, attempt.ProcessorSNS,
		attempt.FinishedAt, attempt.IsSuccessful, attempt.ErrorMessage)
	if err != nil {
		return err
	}

	return nil
}

// GetRegistrationAttemptByID retrieves a registration attempt by its ID from the registration_attempts table
func (s *TicketStore) GetRegistrationAttemptByID(id int) (*types.RegistrationAttempt, error) {
	const selectQuery = `
        SELECT id, file_id, reg_started_at, processor_sns, finished_at, 
               is_successful, error_message
        FROM registration_attempts
        WHERE id = ?;`

	row := s.db.QueryRow(selectQuery, id)

	var attempt types.RegistrationAttempt
	err := row.Scan(
		&attempt.ID, &attempt.FileID, &attempt.RegStartedAt, &attempt.ProcessorSNS,
		&attempt.FinishedAt, &attempt.IsSuccessful, &attempt.ErrorMessage)
	if err != nil {
		return nil, err
	}

	return &attempt, nil
}

// GetRegistrationAttemptsByFileID retrieves registration attempts by file_id from the registration_attempts table
func (s *TicketStore) GetRegistrationAttemptsByFileID(fileID string) ([]*types.RegistrationAttempt, error) {
	const selectQuery = `
        SELECT id, file_id, reg_started_at, processor_sns, finished_at, 
               is_successful, error_message
        FROM registration_attempts
        WHERE file_id = ?;`

	rows, err := s.db.Query(selectQuery, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*types.RegistrationAttempt
	for rows.Next() {
		var attempt types.RegistrationAttempt
		err := rows.Scan(
			&attempt.ID, &attempt.FileID, &attempt.RegStartedAt, &attempt.ProcessorSNS,
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
